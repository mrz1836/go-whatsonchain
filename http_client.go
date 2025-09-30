package whatsonchain

import (
	"bytes"
	"crypto/rand"
	"io"
	"math"
	"math/big"
	"net/http"
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

// NewRetryableHTTPClient creates a new retryable HTTP client
func NewRetryableHTTPClient(httpClient *http.Client, retryCount int, backoff *ExponentialBackoff) *RetryableHTTPClient {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &RetryableHTTPClient{
		client:     httpClient,
		retryCount: retryCount,
		backoff:    backoff,
	}
}

// Do will execute an HTTP request with retry logic
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
			reqForAttempt, err = http.NewRequestWithContext(
				req.Context(),
				req.Method,
				req.URL.String(),
				bytes.NewReader(bodyBytes),
			)
		} else {
			// There is no "body", just clone the request
			reqForAttempt, err = http.NewRequestWithContext(
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
		resp, err = r.client.Do(reqForAttempt)

		// If this is the last attempt, return whatever we got
		if attempt == maxAttempts-1 {
			return resp, err
		}

		// Check if we should retry
		if !r.shouldRetry(resp, err) {
			return resp, err
		}

		// Store the response/error for potential return
		lastResp = resp
		lastErr = err

		// Close the response body if it exists (to free up the connection)
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}

		// Calculate backoff duration
		var backoffDuration time.Duration
		if r.backoff != nil {
			backoffDuration = r.backoff.NextInterval(attempt)
		} else {
			// Default exponential backoff if none provided
			backoffDuration = time.Duration(math.Pow(2, float64(attempt))) * time.Millisecond * 100
		}

		// Wait before retrying, respecting context cancellation
		select {
		case <-req.Context().Done():
			return lastResp, req.Context().Err()
		case <-time.After(backoffDuration):
			// Continue to next attempt
		}
	}

	return lastResp, lastErr
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
