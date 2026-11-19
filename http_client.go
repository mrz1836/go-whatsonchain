package whatsonchain

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ExponentialBackoff provides exponential backoff functionality
type ExponentialBackoff struct {
	initialTimeout    time.Duration
	maxTimeout        time.Duration
	exponentFactor    float64
	maxJitterInterval time.Duration
}

// NewExponentialBackoff creates a new exponential backoff instance
func NewExponentialBackoff(initialTimeout, maxTimeout time.Duration, exponentFactor float64, maxJitterInterval time.Duration) *ExponentialBackoff {
	return &ExponentialBackoff{
		initialTimeout:    initialTimeout,
		maxTimeout:        maxTimeout,
		exponentFactor:    exponentFactor,
		maxJitterInterval: maxJitterInterval,
	}
}

// newDefaultExponentialBackoff returns the backoff used when none is configured
func newDefaultExponentialBackoff() *ExponentialBackoff {
	return NewExponentialBackoff(
		defaultBackOffInitialTimeout,
		defaultBackOffMaxTimeout,
		defaultBackOffExponentFactor,
		defaultBackOffMaxJitterInterval,
	)
}

// NextInterval calculates the next backoff interval for the given attempt
func (eb *ExponentialBackoff) NextInterval(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	// Calculate exponential backoff: initialTimeout * (exponentFactor ^ attempt)
	backoffDuration := float64(eb.initialTimeout) * math.Pow(eb.exponentFactor, float64(attempt))

	// Cap at maxTimeout
	if backoffDuration > float64(eb.maxTimeout) {
		backoffDuration = float64(eb.maxTimeout)
	}

	// Add jitter to prevent thundering herd
	var jitter time.Duration
	if eb.maxJitterInterval > 0 {
		// Use crypto/rand for better security
		if maxJitter := big.NewInt(int64(eb.maxJitterInterval)); maxJitter.Int64() > 0 {
			if n, err := rand.Int(rand.Reader, maxJitter); err == nil {
				jitter = time.Duration(n.Int64())
			}
		}
	}

	return time.Duration(backoffDuration) + jitter
}

// RetryableHTTPClient is a native Go HTTP client with retry capability
type RetryableHTTPClient struct {
	client     *http.Client
	retryCount int
	backoff    *ExponentialBackoff
}

// NewRetryableHTTPClient creates a new retryable HTTP client.
//
// A nil httpClient uses a plain http.Client, and a nil backoff uses the SDK's
// default backoff (250 ms initial, 2 s maximum, factor 2.0, 100 ms jitter).
func NewRetryableHTTPClient(httpClient *http.Client, retryCount int, backoff *ExponentialBackoff) *RetryableHTTPClient {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if backoff == nil {
		backoff = newDefaultExponentialBackoff()
	}

	return &RetryableHTTPClient{
		client:     httpClient,
		retryCount: retryCount,
		backoff:    backoff,
	}
}

// Do will execute an HTTP request with retry logic.
//
// Network errors and 429, 500, 502, 503 and 504 responses are retried up to the
// configured retry count, waiting for the backoff interval between attempts.
// On a 429 or 503 with a Retry-After header the wait is at least the requested
// time. When Retry-After asks for longer than the maximum backoff, or the wait
// would outlast the request context's deadline, Do returns that response (body
// unread) straight away so the caller sees the real status.
func (r *RetryableHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var lastResp *http.Response
	var lastErr error

	// If no retries configured, just execute once
	if r.retryCount <= 0 {
		return r.client.Do(req)
	}

	// Read and store the request body once so we can reuse it for retries
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		_ = req.Body.Close()
	}

	maxAttempts := r.retryCount + 1 // retryCount doesn't include the initial attempt

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Create a new request for each attempt
		var reqForAttempt *http.Request
		var err error

		if bodyBytes != nil {
			// Create new request with fresh body
			reqForAttempt, err = http.NewRequestWithContext( //nolint:gosec // G704: URL is controlled by this library, not user input
				req.Context(),
				req.Method,
				req.URL.String(),
				bytes.NewReader(bodyBytes),
			)
		} else {
			// There is no "body", just clone the request
			reqForAttempt, err = http.NewRequestWithContext( //nolint:gosec // G704: URL is controlled by this library, not user input
				req.Context(),
				req.Method,
				req.URL.String(),
				nil,
			)
		}

		if err != nil {
			return nil, err
		}

		// Copy headers from original request
		for key, values := range req.Header {
			for _, value := range values {
				reqForAttempt.Header.Add(key, value)
			}
		}

		// Execute the request
		var resp *http.Response
		resp, err = r.client.Do(reqForAttempt) //nolint:gosec // G704: URL is controlled by this library, not user input

		// If this is the last attempt, return whatever we got
		if attempt == maxAttempts-1 {
			return resp, err
		}

		// Check if we should retry
		if !r.shouldRetry(resp, err) {
			return resp, err
		}

		// Work out the wait, and give up now if waiting cannot produce an answer
		delay, ok := r.retryDelay(req.Context(), resp, attempt)
		if !ok {
			return resp, err
		}

		// Store the response/error for potential return
		lastResp = resp
		lastErr = err

		// Close the response body if it exists (to free up the connection)
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}

		// Wait before retrying, respecting context cancellation
		timer := time.NewTimer(delay)
		select {
		case <-req.Context().Done():
			timer.Stop()
			return lastResp, req.Context().Err()
		case <-timer.C:
			// Continue to next attempt
		}
	}

	return lastResp, lastErr
}

// retryDelay returns how long to wait before retrying after the given attempt
// (counting from 0). The wait is the backoff interval, raised to the server's
// Retry-After value when there is one.
//
// It returns false when retrying is pointless: Retry-After asks for longer than
// the maximum backoff, or the wait would end after the context's deadline.
func (r *RetryableHTTPClient) retryDelay(ctx context.Context, resp *http.Response, attempt int) (time.Duration, bool) {
	delay := r.backoff.NextInterval(attempt)

	if wait, ok := retryAfter(resp, time.Now()); ok {
		if wait > r.backoff.maxTimeout {
			return 0, false
		}
		delay = max(delay, wait)
	}

	if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) < delay {
		return 0, false
	}

	return delay, true
}

// maxRetryAfterSeconds is the largest Retry-After delta-seconds value that fits
// in a time.Duration; larger values are clamped to the maximum duration.
const maxRetryAfterSeconds = math.MaxInt64 / int64(time.Second)

// retryAfter returns the wait requested by the Retry-After header of a 429 or
// 503 response, relative to now. The header holds either a non-negative number
// of seconds or an HTTP date; a date in the past means no wait. It returns false
// for other status codes and for a missing or malformed header.
func retryAfter(resp *http.Response, now time.Time) (time.Duration, bool) {
	if resp == nil ||
		(resp.StatusCode != http.StatusTooManyRequests && resp.StatusCode != http.StatusServiceUnavailable) {
		return 0, false
	}

	value := strings.TrimSpace(resp.Header.Get("Retry-After"))
	if value == "" {
		return 0, false
	}

	// delta-seconds: one or more digits, no sign
	if strings.Trim(value, "0123456789") == "" {
		seconds, err := strconv.ParseInt(value, 10, 64)
		if err != nil || seconds > maxRetryAfterSeconds {
			// Only an overflow fails here; it means "far longer than we would wait"
			return time.Duration(math.MaxInt64), true
		}
		return time.Duration(seconds) * time.Second, true
	}

	// HTTP-date
	date, err := http.ParseTime(value)
	if err != nil {
		return 0, false
	}
	return max(date.Sub(now), 0), true
}

// shouldRetry determines if a request should be retried based on the response
func (r *RetryableHTTPClient) shouldRetry(resp *http.Response, err error) bool {
	// Retry on network errors
	if err != nil {
		return true
	}

	// Retry on server errors (5xx) and specific client errors
	if resp != nil {
		switch resp.StatusCode {
		case http.StatusInternalServerError, // 500
			http.StatusBadGateway,         // 502
			http.StatusServiceUnavailable, // 503
			http.StatusGatewayTimeout,     // 504
			http.StatusTooManyRequests:    // 429
			return true
		}
	}

	return false
}

// SimpleHTTPClient is a simple wrapper around http.Client for non-retry scenarios
type SimpleHTTPClient struct {
	client *http.Client
}

// NewSimpleHTTPClient creates a new simple HTTP client wrapper
func NewSimpleHTTPClient(httpClient *http.Client) *SimpleHTTPClient {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &SimpleHTTPClient{
		client: httpClient,
	}
}

// Do executes an HTTP request without retry logic
func (s *SimpleHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return s.client.Do(req)
}
