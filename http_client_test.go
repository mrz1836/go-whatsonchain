package whatsonchain

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errNetworkError = errors.New("network error")

func TestNewExponentialBackoff(t *testing.T) {
	t.Parallel()

	initialTimeout := 100 * time.Millisecond
	maxTimeout := 5 * time.Second
	exponentFactor := 2.0
	maxJitterInterval := 50 * time.Millisecond

	backoff := NewExponentialBackoff(initialTimeout, maxTimeout, exponentFactor, maxJitterInterval)

	assert.Equal(t, initialTimeout, backoff.initialTimeout)
	assert.Equal(t, maxTimeout, backoff.maxTimeout)
	assert.InDelta(t, exponentFactor, backoff.exponentFactor, 0.01)
	assert.Equal(t, maxJitterInterval, backoff.maxJitterInterval)
}

func TestExponentialBackoff_NextInterval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		initialTimeout    time.Duration
		maxTimeout        time.Duration
		exponentFactor    float64
		maxJitterInterval time.Duration
		attempt           int
		expectedMin       time.Duration
		expectedMax       time.Duration
	}{
		{
			name:              "first attempt",
			initialTimeout:    100 * time.Millisecond,
			maxTimeout:        5 * time.Second,
			exponentFactor:    2.0,
			maxJitterInterval: 0,
			attempt:           0,
			expectedMin:       100 * time.Millisecond,
			expectedMax:       100 * time.Millisecond,
		},
		{
			name:              "second attempt",
			initialTimeout:    100 * time.Millisecond,
			maxTimeout:        5 * time.Second,
			exponentFactor:    2.0,
			maxJitterInterval: 0,
			attempt:           1,
			expectedMin:       200 * time.Millisecond,
			expectedMax:       200 * time.Millisecond,
		},
		{
			name:              "third attempt",
			initialTimeout:    100 * time.Millisecond,
			maxTimeout:        5 * time.Second,
			exponentFactor:    2.0,
			maxJitterInterval: 0,
			attempt:           2,
			expectedMin:       400 * time.Millisecond,
			expectedMax:       400 * time.Millisecond,
		},
		{
			name:              "with jitter",
			initialTimeout:    100 * time.Millisecond,
			maxTimeout:        5 * time.Second,
			exponentFactor:    2.0,
			maxJitterInterval: 50 * time.Millisecond,
			attempt:           1,
			expectedMin:       200 * time.Millisecond,
			expectedMax:       250 * time.Millisecond,
		},
		{
			name:              "capped at max timeout",
			initialTimeout:    1 * time.Second,
			maxTimeout:        2 * time.Second,
			exponentFactor:    2.0,
			maxJitterInterval: 0,
			attempt:           5, // Would be 32 seconds without cap
			expectedMin:       2 * time.Second,
			expectedMax:       2 * time.Second,
		},
		{
			name:              "negative attempt",
			initialTimeout:    100 * time.Millisecond,
			maxTimeout:        5 * time.Second,
			exponentFactor:    2.0,
			maxJitterInterval: 0,
			attempt:           -1,
			expectedMin:       100 * time.Millisecond,
			expectedMax:       100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			backoff := NewExponentialBackoff(tt.initialTimeout, tt.maxTimeout, tt.exponentFactor, tt.maxJitterInterval)
			interval := backoff.NextInterval(tt.attempt)

			assert.GreaterOrEqual(t, interval, tt.expectedMin)
			assert.LessOrEqual(t, interval, tt.expectedMax)
		})
	}
}

func TestRetryableHTTPClient_shouldRetry(t *testing.T) {
	t.Parallel()

	client := NewRetryableHTTPClient(nil, 3, nil)

	tests := []struct {
		name        string
		resp        *http.Response
		err         error
		shouldRetry bool
	}{
		{
			name:        "network error",
			resp:        nil,
			err:         errNetworkError,
			shouldRetry: true,
		},
		{
			name: "500 internal server error",
			resp: &http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			err:         nil,
			shouldRetry: true,
		},
		{
			name: "502 bad gateway",
			resp: &http.Response{
				StatusCode: http.StatusBadGateway,
			},
			err:         nil,
			shouldRetry: true,
		},
		{
			name: "503 service unavailable",
			resp: &http.Response{
				StatusCode: http.StatusServiceUnavailable,
			},
			err:         nil,
			shouldRetry: true,
		},
		{
			name: "504 gateway timeout",
			resp: &http.Response{
				StatusCode: http.StatusGatewayTimeout,
			},
			err:         nil,
			shouldRetry: true,
		},
		{
			name: "429 too many requests",
			resp: &http.Response{
				StatusCode: http.StatusTooManyRequests,
			},
			err:         nil,
			shouldRetry: true,
		},
		{
			name: "200 ok",
			resp: &http.Response{
				StatusCode: http.StatusOK,
			},
			err:         nil,
			shouldRetry: false,
		},
		{
			name: "400 bad request",
			resp: &http.Response{
				StatusCode: http.StatusBadRequest,
			},
			err:         nil,
			shouldRetry: false,
		},
		{
			name: "404 not found",
			resp: &http.Response{
				StatusCode: http.StatusNotFound,
			},
			err:         nil,
			shouldRetry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := client.shouldRetry(tt.resp, tt.err)
			assert.Equal(t, tt.shouldRetry, result)
		})
	}
}

func TestNewRetryableHTTPClient(t *testing.T) {
	t.Parallel()

	// Test with nil client
	client1 := NewRetryableHTTPClient(nil, 3, nil)
	assert.NotNil(t, client1.client)
	assert.Equal(t, 3, client1.retryCount)

	// Test with custom client
	customClient := &http.Client{Timeout: 10 * time.Second}
	backoff := NewExponentialBackoff(100*time.Millisecond, 5*time.Second, 2.0, 50*time.Millisecond)
	client2 := NewRetryableHTTPClient(customClient, 5, backoff)
	assert.Equal(t, customClient, client2.client)
	assert.Equal(t, 5, client2.retryCount)
	assert.Equal(t, backoff, client2.backoff)
}

func TestRetryableHTTPClient_Do_NoRetries(t *testing.T) {
	t.Parallel()

	// Create a test server that always returns 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "success")
	}))
	defer server.Close()

	// Create client with 0 retries
	client := NewRetryableHTTPClient(nil, 0, nil)

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "success")
}

func TestRetryableHTTPClient_Do_SuccessFirstTry(t *testing.T) {
	t.Parallel()

	// Create a test server that always returns 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "success")
	}))
	defer server.Close()

	backoff := NewExponentialBackoff(10*time.Millisecond, 1*time.Second, 2.0, 0)
	client := NewRetryableHTTPClient(nil, 3, backoff)

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "success")
}

func TestRetryableHTTPClient_Do_RetryOnServerError(t *testing.T) {
	t.Parallel()

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requestCount++
		if requestCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, "server error")
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, "success")
		}
	}))
	defer server.Close()

	backoff := NewExponentialBackoff(1*time.Millisecond, 100*time.Millisecond, 2.0, 0)
	client := NewRetryableHTTPClient(nil, 3, backoff)

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 3, requestCount) // Should have been called 3 times

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "success")
}

func TestRetryableHTTPClient_Do_MaxRetriesExceeded(t *testing.T) {
	t.Parallel()

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, "server error")
	}))
	defer server.Close()

	backoff := NewExponentialBackoff(1*time.Millisecond, 100*time.Millisecond, 2.0, 0)
	client := NewRetryableHTTPClient(nil, 2, backoff) // Only 2 retries

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, 3, requestCount) // Initial attempt + 2 retries = 3 total attempts

	defer func() { _ = resp.Body.Close() }()
}

func TestRetryableHTTPClient_Do_ContextCancellation(t *testing.T) {
	t.Parallel()

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, "server error")
	}))
	defer server.Close()

	backoff := NewExponentialBackoff(100*time.Millisecond, 1*time.Second, 2.0, 0)
	client := NewRetryableHTTPClient(nil, 5, backoff)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	require.NoError(t, err)

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	// Should return quickly due to context cancellation
	assert.Less(t, duration, 200*time.Millisecond)

	if err != nil {
		assert.Equal(t, context.DeadlineExceeded, err)
	} else {
		// If we got a response, it should be the first error response
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		_ = resp.Body.Close()
	}

	// Should not have made many requests due to quick cancellation
	assert.LessOrEqual(t, requestCount, 2)
}

func TestRetryableHTTPClient_Do_NoRetryOn4xx(t *testing.T) {
	t.Parallel()

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "bad request")
	}))
	defer server.Close()

	backoff := NewExponentialBackoff(1*time.Millisecond, 100*time.Millisecond, 2.0, 0)
	client := NewRetryableHTTPClient(nil, 3, backoff)

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, 1, requestCount) // Should only be called once (no retries)

	defer func() { _ = resp.Body.Close() }()
}

func TestRetryableHTTPClient_Do_PostRequest(t *testing.T) {
	t.Parallel()

	requestCount := 0
	var receivedBodies []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		body, _ := io.ReadAll(r.Body)
		receivedBodies = append(receivedBodies, string(body))

		if requestCount < 2 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, "success")
		}
	}))
	defer server.Close()

	backoff := NewExponentialBackoff(1*time.Millisecond, 100*time.Millisecond, 2.0, 0)
	client := NewRetryableHTTPClient(nil, 3, backoff)

	reqBody := "test post data"
	req, err := http.NewRequestWithContext(context.Background(), "POST", server.URL, strings.NewReader(reqBody))
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, requestCount)

	// Verify that the POST body was correctly sent in both attempts
	assert.Len(t, receivedBodies, 2)
	assert.Equal(t, reqBody, receivedBodies[0])
	assert.Equal(t, reqBody, receivedBodies[1])

	defer func() { _ = resp.Body.Close() }()
}

func TestSimpleHTTPClient(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "simple success")
	}))
	defer server.Close()

	// Test with nil client
	client1 := NewSimpleHTTPClient(nil)
	assert.NotNil(t, client1.client)

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	require.NoError(t, err)

	resp, err := client1.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "simple success")

	// Test with custom client
	customClient := &http.Client{Timeout: 10 * time.Second}
	client2 := NewSimpleHTTPClient(customClient)
	assert.Equal(t, customClient, client2.client)
}

// errReadError is a static error for testing
var errReadError = errors.New("read error")

// errorReader always returns an error when read
type errorReader struct{}

func (e *errorReader) Read(_ []byte) (n int, err error) {
	return 0, errReadError
}

// TestRetryableHTTPClient_Do_BodyReadError tests that body read errors are handled properly
func TestRetryableHTTPClient_Do_BodyReadError(t *testing.T) {
	t.Parallel()

	backoff := NewExponentialBackoff(1*time.Millisecond, 100*time.Millisecond, 2.0, 0)
	client := NewRetryableHTTPClient(nil, 2, backoff)

	// Create a request with an error reader body
	req, err := http.NewRequestWithContext(context.Background(), "POST", "http://example.com", &errorReader{})
	require.NoError(t, err)

	// This should fail when trying to read the body for retries
	resp, err := client.Do(req)
	if resp != nil && resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read error")
}

// TestRetryableHTTPClient_Do_DefaultBackoff tests that default backoff works when no backoff is provided
func TestRetryableHTTPClient_Do_DefaultBackoff(t *testing.T) {
	t.Parallel()

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requestCount++
		if requestCount < 2 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, "success")
		}
	}))
	defer server.Close()

	// Create client without backoff (will use default)
	client := NewRetryableHTTPClient(nil, 3, nil)

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	require.NoError(t, err)

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, requestCount)

	// Should have some delay due to default backoff
	assert.Greater(t, duration, 50*time.Millisecond)

	defer func() { _ = resp.Body.Close() }()
}

// TestNewRetryableHTTPClient_DefaultBackoff checks that a nil backoff falls back
// to the SDK's default backoff rather than an unbounded built-in schedule.
func TestNewRetryableHTTPClient_DefaultBackoff(t *testing.T) {
	t.Parallel()

	client := NewRetryableHTTPClient(nil, 2, nil)
	require.NotNil(t, client.backoff)
	assert.Equal(t, defaultBackOffInitialTimeout, client.backoff.initialTimeout)
	assert.Equal(t, defaultBackOffMaxTimeout, client.backoff.maxTimeout)
	assert.InDelta(t, defaultBackOffExponentFactor, client.backoff.exponentFactor, 1e-9)
	assert.Equal(t, defaultBackOffMaxJitterInterval, client.backoff.maxJitterInterval)

	// The defaults must be long enough to ride out a brief rate limit (issue #223)
	assert.GreaterOrEqual(t, defaultBackOffInitialTimeout, 100*time.Millisecond)
	assert.GreaterOrEqual(t, defaultBackOffMaxTimeout, time.Second)

	// The default schedule is capped: a late attempt never waits longer than max + jitter
	assert.LessOrEqual(t, client.backoff.NextInterval(30), defaultBackOffMaxTimeout+defaultBackOffMaxJitterInterval)
}

// retryAfterHeader builds a header carrying the given Retry-After value
func retryAfterHeader(value string) http.Header {
	header := http.Header{}
	header.Set("Retry-After", value)
	return header
}

// TestRetryAfter covers parsing of the Retry-After header (RFC 9110 section 10.2.3)
func TestRetryAfter(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.November, 19, 12, 0, 0, 0, time.UTC)
	maxSeconds := strconv.FormatInt(maxRetryAfterSeconds, 10)
	overSeconds := strconv.FormatInt(maxRetryAfterSeconds+1, 10)

	tests := []struct {
		name   string
		resp   *http.Response
		want   time.Duration
		wantOK bool
	}{
		{"nil response", nil, 0, false},
		{"429 without header", &http.Response{StatusCode: http.StatusTooManyRequests, Header: http.Header{}}, 0, false},
		{"429 with empty header", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("")}, 0, false},
		{"429 seconds", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("5")}, 5 * time.Second, true},
		{"503 seconds", &http.Response{StatusCode: http.StatusServiceUnavailable, Header: retryAfterHeader("120")}, 120 * time.Second, true},
		{"seconds with surrounding spaces", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader(" 7 ")}, 7 * time.Second, true},
		{"zero seconds", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("0")}, 0, true},
		{"leading zeros", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("007")}, 7 * time.Second, true},
		{"negative seconds", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("-5")}, 0, false},
		{"explicit plus sign", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("+5")}, 0, false},
		{"fractional seconds", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("1.5")}, 0, false},
		{"seconds with unit", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("5s")}, 0, false},
		{"garbage", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("soon")}, 0, false},
		{"largest representable seconds", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader(maxSeconds)}, time.Duration(maxRetryAfterSeconds) * time.Second, true},
		{"seconds just past time.Duration", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader(overSeconds)}, time.Duration(math.MaxInt64), true},
		{"seconds past int64", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("99999999999999999999999")}, time.Duration(math.MaxInt64), true},
		{"future HTTP date", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader(now.Add(30 * time.Second).Format(http.TimeFormat))}, 30 * time.Second, true},
		{"future RFC 850 date", &http.Response{StatusCode: http.StatusServiceUnavailable, Header: retryAfterHeader(now.Add(10 * time.Second).Format(time.RFC850))}, 10 * time.Second, true},
		{"future ANSI C date", &http.Response{StatusCode: http.StatusServiceUnavailable, Header: retryAfterHeader(now.Add(time.Minute).Format(time.ANSIC))}, time.Minute, true},
		{"HTTP date equal to now", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader(now.Format(http.TimeFormat))}, 0, true},
		{"past HTTP date", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader(now.Add(-time.Hour).Format(http.TimeFormat))}, 0, true},
		{"malformed date", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("Thu, 99 Nov 2026 12:00:00 GMT")}, 0, false},
		{"ignored on 500", &http.Response{StatusCode: http.StatusInternalServerError, Header: retryAfterHeader("5")}, 0, false},
		{"ignored on 502", &http.Response{StatusCode: http.StatusBadGateway, Header: retryAfterHeader("5")}, 0, false},
		{"ignored on 200", &http.Response{StatusCode: http.StatusOK, Header: retryAfterHeader("5")}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := retryAfter(tt.resp, now)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestRetryableHTTPClient_retryDelay covers how the backoff, Retry-After and the
// context deadline combine into the wait before a retry.
func TestRetryableHTTPClient_retryDelay(t *testing.T) {
	t.Parallel()

	// Deterministic backoff: 100ms, 200ms, 400ms, ... capped at 1s, no jitter
	backoff := NewExponentialBackoff(100*time.Millisecond, time.Second, 2.0, 0)
	client := NewRetryableHTTPClient(nil, 3, backoff)

	tests := []struct {
		name     string
		resp     *http.Response
		attempt  int
		deadline time.Duration // 0: no deadline
		want     time.Duration
		wantOK   bool
	}{
		{"network error uses the backoff", nil, 0, 0, 100 * time.Millisecond, true},
		{"later attempt grows the backoff", nil, 2, 0, 400 * time.Millisecond, true},
		{"backoff is capped at max", nil, 10, 0, time.Second, true},
		{"429 without Retry-After uses the backoff", &http.Response{StatusCode: http.StatusTooManyRequests, Header: http.Header{}}, 0, 0, 100 * time.Millisecond, true},
		{"Retry-After shorter than the backoff keeps the backoff", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("0")}, 1, 0, 200 * time.Millisecond, true},
		{"Retry-After longer than the backoff wins", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("1")}, 0, 0, time.Second, true},
		{"Retry-After equal to max is honored", &http.Response{StatusCode: http.StatusServiceUnavailable, Header: retryAfterHeader("1")}, 0, 0, time.Second, true},
		{"Retry-After above max gives up", &http.Response{StatusCode: http.StatusServiceUnavailable, Header: retryAfterHeader("2")}, 0, 0, 0, false},
		{"huge Retry-After gives up", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("99999999999999999999")}, 0, 0, 0, false},
		{"malformed Retry-After falls back to the backoff", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("later")}, 0, 0, 100 * time.Millisecond, true},
		{"Retry-After on 500 is ignored", &http.Response{StatusCode: http.StatusInternalServerError, Header: retryAfterHeader("120")}, 0, 0, 100 * time.Millisecond, true},
		{"deadline after the delay allows the retry", nil, 0, 10 * time.Second, 100 * time.Millisecond, true},
		{"deadline before the backoff gives up", nil, 0, 50 * time.Millisecond, 0, false},
		{"deadline before Retry-After gives up", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("1")}, 0, 500 * time.Millisecond, 0, false},
		{"deadline after Retry-After allows the retry", &http.Response{StatusCode: http.StatusTooManyRequests, Header: retryAfterHeader("1")}, 0, 10 * time.Second, time.Second, true},
		{"expired deadline gives up", nil, 0, -time.Second, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			if tt.deadline != 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.deadline)
				defer cancel()
			}

			got, ok := client.retryDelay(ctx, tt.resp, tt.attempt)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

// retryAfterServer answers the first request with status and the given
// Retry-After header (computed per request), then 200 "ok". It counts requests.
func retryAfterServer(t *testing.T, status int, retryAfterValue func() string) (*httptest.Server, *atomic.Int32) {
	t.Helper()

	var count atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if count.Add(1) == 1 {
			if retryAfterValue != nil {
				w.Header().Set("Retry-After", retryAfterValue())
			}
			w.WriteHeader(status)
			_, _ = fmt.Fprint(w, "slow down")
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "ok")
	}))
	t.Cleanup(server.Close)

	return server, &count
}

// doGet runs a GET through client and returns the status, body and elapsed time
func doGet(ctx context.Context, t *testing.T, client *RetryableHTTPClient, url string) (int, string, time.Duration) {
	t.Helper()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	require.NoError(t, err)

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "the returned body must still be open")

	return resp.StatusCode, string(body), elapsed
}

// TestRetryableHTTPClient_Do_HonorsRetryAfter checks that a 429 or 503 retry waits
// for the server's Retry-After even when the backoff alone would retry sooner.
func TestRetryableHTTPClient_Do_HonorsRetryAfter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		status  int
		value   func() string
		minWait time.Duration
	}{
		{"429 with seconds", http.StatusTooManyRequests, func() string { return "1" }, time.Second},
		{"503 with seconds", http.StatusServiceUnavailable, func() string { return "1" }, time.Second},
		{"429 with HTTP date", http.StatusTooManyRequests, func() string {
			// Dates have one-second resolution, so this asks for a wait of 1 to 2 seconds
			return time.Now().Add(2 * time.Second).UTC().Format(http.TimeFormat)
		}, time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server, count := retryAfterServer(t, tt.status, tt.value)

			// A 1ms backoff would retry almost at once without Retry-After
			client := NewRetryableHTTPClient(nil, 2, NewExponentialBackoff(time.Millisecond, 2*time.Second, 2.0, 0))
			status, body, elapsed := doGet(context.Background(), t, client, server.URL)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, "ok", body)
			assert.Equal(t, int32(2), count.Load())
			assert.GreaterOrEqual(t, elapsed, tt.minWait, "the retry must wait for Retry-After")
		})
	}
}

// TestRetryableHTTPClient_Do_GivesUpEarly checks that Do returns the server's
// response at once, with its body open, when waiting cannot pay off.
func TestRetryableHTTPClient_Do_GivesUpEarly(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		status   int
		value    func() string
		backoff  *ExponentialBackoff
		deadline time.Duration // 0: no deadline
	}{
		{
			name:    "429 Retry-After above max backoff",
			status:  http.StatusTooManyRequests,
			value:   func() string { return "120" },
			backoff: NewExponentialBackoff(time.Millisecond, 2*time.Second, 2.0, 0),
		},
		{
			name:    "503 Retry-After date above max backoff",
			status:  http.StatusServiceUnavailable,
			value:   func() string { return time.Now().Add(time.Hour).UTC().Format(http.TimeFormat) },
			backoff: NewExponentialBackoff(time.Millisecond, 2*time.Second, 2.0, 0),
		},
		{
			name:    "default backoff with Retry-After above its max",
			status:  http.StatusTooManyRequests,
			value:   func() string { return "120" },
			backoff: nil,
		},
		{
			name:     "context deadline before the backoff",
			status:   http.StatusServiceUnavailable,
			backoff:  NewExponentialBackoff(2*time.Second, 5*time.Second, 2.0, 0),
			deadline: 500 * time.Millisecond,
		},
		{
			name:     "context deadline before Retry-After",
			status:   http.StatusTooManyRequests,
			value:    func() string { return "2" },
			backoff:  NewExponentialBackoff(time.Millisecond, 5*time.Second, 2.0, 0),
			deadline: 500 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server, count := retryAfterServer(t, tt.status, tt.value)
			client := NewRetryableHTTPClient(nil, 3, tt.backoff)

			ctx := context.Background()
			if tt.deadline > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.deadline)
				defer cancel()
			}

			status, body, elapsed := doGet(ctx, t, client, server.URL)

			assert.Equal(t, tt.status, status, "the caller must see the real status")
			assert.Equal(t, "slow down", body)
			assert.Equal(t, int32(1), count.Load(), "no retry may be sent")
			assert.Less(t, elapsed, 250*time.Millisecond, "Do must not wait before giving up")
		})
	}
}

// TestRetryableHTTPClient_Do_RetryAfterIgnoredOn500 checks that Retry-After on a
// status other than 429 or 503 neither delays the retry nor stops it.
func TestRetryableHTTPClient_Do_RetryAfterIgnoredOn500(t *testing.T) {
	t.Parallel()

	server, count := retryAfterServer(t, http.StatusInternalServerError, func() string { return "120" })
	client := NewRetryableHTTPClient(nil, 2, NewExponentialBackoff(time.Millisecond, 2*time.Second, 2.0, 0))

	status, body, elapsed := doGet(context.Background(), t, client, server.URL)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "ok", body)
	assert.Equal(t, int32(2), count.Load())
	assert.Less(t, elapsed, time.Second)
}

// TestClient_RetryAfter_EndToEnd checks Retry-After handling through NewClient:
// a short Retry-After is waited out and retried, while a long one surfaces the
// real 429 to the caller as ErrRequestFailed straight away.
func TestClient_RetryAfter_EndToEnd(t *testing.T) {
	t.Parallel()

	t.Run("short Retry-After is retried", func(t *testing.T) {
		t.Parallel()

		server, count := retryAfterServer(t, http.StatusTooManyRequests, func() string { return "1" })
		client, err := NewClient(context.Background(), WithBaseURL(server.URL+"/v1/"))
		require.NoError(t, err)

		start := time.Now()
		health, err := client.GetHealth(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "ok", health)
		assert.Equal(t, int32(2), count.Load())
		assert.GreaterOrEqual(t, time.Since(start), time.Second)
	})

	t.Run("long Retry-After fails fast", func(t *testing.T) {
		t.Parallel()

		server, count := retryAfterServer(t, http.StatusTooManyRequests, func() string { return "3600" })
		client, err := NewClient(context.Background(), WithBaseURL(server.URL+"/v1/"))
		require.NoError(t, err)

		start := time.Now()
		health, err := client.GetHealth(context.Background())
		require.ErrorIs(t, err, ErrRequestFailed)
		require.EqualError(t, err, "API request failed: HTTP 429: slow down")
		assert.Empty(t, health)
		assert.Equal(t, int32(1), count.Load())
		assert.Less(t, time.Since(start), 250*time.Millisecond)
	})
}

// TestRetryableHTTPClient_Do_CanceledDuringWait checks that canceling the context
// while Do waits between attempts stops the wait at once and returns the
// cancellation error, without sending another request.
func TestRetryableHTTPClient_Do_CanceledDuringWait(t *testing.T) {
	t.Parallel()

	server, count := retryAfterServer(t, http.StatusServiceUnavailable, nil)
	// One 3s wait: a regression that ignores cancellation fails in about 3s instead of hanging
	client := NewRetryableHTTPClient(nil, 1, NewExponentialBackoff(3*time.Second, 3*time.Second, 2.0, 0))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	timer := time.AfterFunc(100*time.Millisecond, cancel)
	defer timer.Stop()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)
	if resp != nil {
		_ = resp.Body.Close()
	}

	require.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, int32(1), count.Load(), "no retry may be sent after cancellation")
	assert.Less(t, elapsed, 2*time.Second, "the wait must stop when the context is canceled")
}
