package whatsonchain

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
