package whatsonchain

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewClient tests the new NewClient function with default options
func TestNewClient(t *testing.T) {
	t.Setenv(EnvAPIKey, "")

	client, err := NewClient(context.Background())
	require.NoError(t, err)
	require.NotNil(t, client)

	// Check defaults
	assert.Equal(t, NetworkMain, client.Network())
	assert.Equal(t, defaultUserAgent, client.UserAgent())
	assert.Equal(t, defaultRateLimit, client.RateLimit())
	assert.Empty(t, client.APIKey())
}

// TestWithNetwork tests the WithNetwork option
func TestWithNetwork(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		network NetworkType
	}{
		{"main net", NetworkMain},
		{"test net", NetworkTest},
		{"stn net", NetworkStn},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(context.Background(), WithNetwork(tt.network))
			require.NoError(t, err)
			assert.Equal(t, tt.network, client.Network())
		})
	}
}

// TestWithAPIKey tests the WithAPIKey option
func TestWithAPIKey(t *testing.T) {
	t.Parallel()

	apiKey := "12345"
	client, err := NewClient(context.Background(), WithAPIKey(apiKey))
	require.NoError(t, err)
	assert.Equal(t, apiKey, client.APIKey())
}

// TestWithUserAgent tests the WithUserAgent option
func TestWithUserAgent(t *testing.T) {
	t.Parallel()

	userAgent := "custom-user-agent/1.0"
	client, err := NewClient(context.Background(), WithUserAgent(userAgent))
	require.NoError(t, err)
	assert.Equal(t, userAgent, client.UserAgent())
}

// TestWithRateLimit tests the WithRateLimit option
func TestWithRateLimit(t *testing.T) {
	t.Parallel()

	rateLimit := 10
	client, err := NewClient(context.Background(), WithRateLimit(rateLimit))
	require.NoError(t, err)
	assert.Equal(t, rateLimit, client.RateLimit())
}

// TestWithHTTPClient tests the WithHTTPClient option
func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	customClient := http.DefaultClient
	client, err := NewClient(context.Background(), WithHTTPClient(customClient))
	require.NoError(t, err)
	assert.Equal(t, customClient, client.HTTPClient())
}

// TestWithRequestTimeout tests the WithRequestTimeout option
func TestWithRequestTimeout(t *testing.T) {
	t.Parallel()

	timeout := 60 * time.Second
	client, err := NewClient(context.Background(), WithRequestTimeout(timeout))
	require.NoError(t, err)
	assert.Equal(t, timeout, client.RequestTimeout())
}

// TestWithRequestRetryCount tests the WithRequestRetryCount option
func TestWithRequestRetryCount(t *testing.T) {
	t.Parallel()

	retryCount := 5
	client, err := NewClient(context.Background(), WithRequestRetryCount(retryCount))
	require.NoError(t, err)
	assert.Equal(t, retryCount, client.RequestRetryCount())
}

// TestWithRequestRetryCount_NoRetry tests zero retry count
func TestWithRequestRetryCount_NoRetry(t *testing.T) {
	t.Parallel()

	client, err := NewClient(context.Background(), WithRequestRetryCount(0))
	require.NoError(t, err)
	assert.Equal(t, 0, client.RequestRetryCount())
}

// TestWithBackoff tests the WithBackoff option
func TestWithBackoff(t *testing.T) {
	t.Parallel()

	initialTimeout := 5 * time.Millisecond
	maxTimeout := 50 * time.Millisecond
	exponentFactor := 3.0
	maxJitter := 5 * time.Millisecond

	client, err := NewClient(
		context.Background(),
		WithBackoff(initialTimeout, maxTimeout, exponentFactor, maxJitter),
	)
	require.NoError(t, err)

	initial, maxBackoff, factor, jitter := client.BackoffConfig()
	assert.Equal(t, initialTimeout, initial)
	assert.Equal(t, maxTimeout, maxBackoff)
	assert.InEpsilon(t, exponentFactor, factor, 0.0001)
	assert.Equal(t, maxJitter, jitter)
}

// TestWithDialer tests the WithDialer option
func TestWithDialer(t *testing.T) {
	t.Parallel()

	keepAlive := 30 * time.Second
	timeout := 10 * time.Second

	client, err := NewClient(
		context.Background(),
		WithDialer(keepAlive, timeout),
	)
	require.NoError(t, err)

	ka, to := client.DialerConfig()
	assert.Equal(t, keepAlive, ka)
	assert.Equal(t, timeout, to)
}

// TestWithTransport tests the WithTransport option
func TestWithTransport(t *testing.T) {
	t.Parallel()

	idleTimeout := 30 * time.Second
	tlsTimeout := 10 * time.Second
	expectContinueTimeout := 5 * time.Second
	maxIdleConnections := 20

	client, err := NewClient(
		context.Background(),
		WithTransport(idleTimeout, tlsTimeout, expectContinueTimeout, maxIdleConnections),
	)
	require.NoError(t, err)

	idle, tls, expect, maxIdle := client.TransportConfig()
	assert.Equal(t, idleTimeout, idle)
	assert.Equal(t, tlsTimeout, tls)
	assert.Equal(t, expectContinueTimeout, expect)
	assert.Equal(t, maxIdleConnections, maxIdle)
}

// TestMultipleOptions tests combining multiple options
func TestMultipleOptions(t *testing.T) {
	t.Parallel()

	apiKey := "test-key"
	userAgent := "test-agent"
	network := NetworkTest
	rateLimit := 5

	client, err := NewClient(
		context.Background(),
		WithAPIKey(apiKey),
		WithUserAgent(userAgent),
		WithNetwork(network),
		WithRateLimit(rateLimit),
	)
	require.NoError(t, err)

	assert.Equal(t, apiKey, client.APIKey())
	assert.Equal(t, userAgent, client.UserAgent())
	assert.Equal(t, network, client.Network())
	assert.Equal(t, rateLimit, client.RateLimit())
}

// TestSetters tests all setter methods
func TestSetAPIKey(t *testing.T) {
	t.Parallel()

	client, err := NewClient(context.Background())
	require.NoError(t, err)

	newKey := "new-api-key"
	client.SetAPIKey(newKey)
	assert.Equal(t, newKey, client.APIKey())
}

// TestSetUserAgent tests the SetUserAgent method
func TestSetUserAgent(t *testing.T) {
	t.Parallel()

	client, err := NewClient(context.Background())
	require.NoError(t, err)

	newAgent := "new-user-agent"
	client.SetUserAgent(newAgent)
	assert.Equal(t, newAgent, client.UserAgent())
}

// TestSetRateLimit tests the SetRateLimit method
func TestSetRateLimit(t *testing.T) {
	t.Parallel()

	client, err := NewClient(context.Background())
	require.NoError(t, err)

	newLimit := 15
	client.SetRateLimit(newLimit)
	assert.Equal(t, newLimit, client.RateLimit())
}

// TestSetNetwork tests the SetNetwork method
func TestSetNetwork(t *testing.T) {
	t.Parallel()

	client, err := NewClient(context.Background())
	require.NoError(t, err)

	require.NoError(t, client.SetNetwork(NetworkTest))
	assert.Equal(t, NetworkTest, client.Network())
}

// TestNewClient_EnvAPIKey tests that WHATS_ON_CHAIN_API_KEY env var is auto-loaded
func TestNewClient_EnvAPIKey(t *testing.T) {
	t.Setenv(EnvAPIKey, "env-api-key-123")

	client, err := NewClient(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "env-api-key-123", client.APIKey())
}

// TestNewClient_EnvAPIKey_ExplicitOverrides tests that WithAPIKey takes precedence over env var
func TestNewClient_EnvAPIKey_ExplicitOverrides(t *testing.T) {
	t.Setenv(EnvAPIKey, "env-api-key-123")

	client, err := NewClient(context.Background(), WithAPIKey("explicit-key"))
	require.NoError(t, err)
	assert.Equal(t, "explicit-key", client.APIKey())
}

// TestNewClient_EnvAPIKey_NotSet tests that no API key is set when env var is absent
func TestNewClient_EnvAPIKey_NotSet(t *testing.T) {
	t.Setenv(EnvAPIKey, "")

	client, err := NewClient(context.Background())
	require.NoError(t, err)
	assert.Empty(t, client.APIKey())
}

// TestDefaultBackoff checks the default retry backoff (issue #223: the old 2ms to
// 10ms default sent both retries of a 429 within about 10ms) and that the default
// HTTP client is built with it.
func TestDefaultBackoff(t *testing.T) {
	t.Parallel()

	client, err := NewClient(context.Background())
	require.NoError(t, err)

	initial, maxBackoff, factor, jitter := client.BackoffConfig()
	assert.Equal(t, 250*time.Millisecond, initial)
	assert.Equal(t, 2*time.Second, maxBackoff)
	assert.InDelta(t, 2.0, factor, 1e-9)
	assert.Equal(t, 100*time.Millisecond, jitter)
	assert.Equal(t, 2, client.RequestRetryCount())

	retryClient, ok := client.HTTPClient().(*RetryableHTTPClient)
	require.True(t, ok, "the default client must retry")
	assert.Equal(t, initial, retryClient.backoff.initialTimeout)
	assert.Equal(t, maxBackoff, retryClient.backoff.maxTimeout)
	assert.InDelta(t, factor, retryClient.backoff.exponentFactor, 1e-9)
	assert.Equal(t, jitter, retryClient.backoff.maxJitterInterval)
}

// countingRoundTripper answers each request with the next status from statuses
// (repeating the last one), using the status text as the body. It counts calls
// and records the last User-Agent it saw.
type countingRoundTripper struct {
	statuses  []int
	calls     atomic.Int32
	userAgent atomic.Value
}

// RoundTrip implements http.RoundTripper
func (rt *countingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	n := int(rt.calls.Add(1))
	rt.userAgent.Store(req.Header.Get("User-Agent"))
	status := rt.statuses[min(n, len(rt.statuses))-1]

	return &http.Response{
		StatusCode: status,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(http.StatusText(status))),
		Request:    req,
	}, nil
}

// TestWithRoundTripper checks that a custom round tripper sends the requests
// while the SDK's retries, backoff and timeout still apply (issue #223).
func TestWithRoundTripper(t *testing.T) {
	t.Parallel()

	fastBackoff := WithBackoff(time.Millisecond, 10*time.Millisecond, 2.0, 0)

	t.Run("retries run through the round tripper", func(t *testing.T) {
		t.Parallel()

		rt := &countingRoundTripper{statuses: []int{http.StatusServiceUnavailable, http.StatusOK}}
		client, err := NewClient(context.Background(), WithRoundTripper(rt), fastBackoff, WithUserAgent("rt-test/1.0"))
		require.NoError(t, err)

		health, err := client.GetHealth(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "OK", health)
		assert.Equal(t, int32(2), rt.calls.Load(), "the 503 must be retried")
		assert.Equal(t, "rt-test/1.0", rt.userAgent.Load(), "requests must carry the SDK headers")
	})

	t.Run("retry count is honored", func(t *testing.T) {
		t.Parallel()

		rt := &countingRoundTripper{statuses: []int{http.StatusBadGateway}}
		client, err := NewClient(context.Background(), WithRoundTripper(rt), fastBackoff, WithRequestRetryCount(3))
		require.NoError(t, err)

		_, err = client.GetHealth(context.Background())
		require.ErrorIs(t, err, ErrRequestFailed)
		assert.Equal(t, int32(4), rt.calls.Load(), "one attempt plus three retries")
	})

	t.Run("no retries when disabled", func(t *testing.T) {
		t.Parallel()

		rt := &countingRoundTripper{statuses: []int{http.StatusServiceUnavailable, http.StatusOK}}
		client, err := NewClient(context.Background(), WithRoundTripper(rt), WithRequestRetryCount(0))
		require.NoError(t, err)

		simple, ok := client.HTTPClient().(*SimpleHTTPClient)
		require.True(t, ok)
		assert.Same(t, rt, simple.client.Transport)

		_, err = client.GetHealth(context.Background())
		require.EqualError(t, err, "API request failed: HTTP 503: Service Unavailable")
		assert.Equal(t, int32(1), rt.calls.Load())
	})

	t.Run("wrapped with the request timeout", func(t *testing.T) {
		t.Parallel()

		rt := &countingRoundTripper{statuses: []int{http.StatusOK}}
		client, err := NewClient(context.Background(), WithRoundTripper(rt), WithRequestTimeout(7*time.Second))
		require.NoError(t, err)

		retryClient, ok := client.HTTPClient().(*RetryableHTTPClient)
		require.True(t, ok)
		assert.Same(t, rt, retryClient.client.Transport)
		assert.Equal(t, 7*time.Second, retryClient.client.Timeout)
	})

	t.Run("nil keeps the default transport and its options", func(t *testing.T) {
		t.Parallel()

		client, err := NewClient(
			context.Background(),
			WithRoundTripper(nil),
			WithTransport(31*time.Second, 6*time.Second, 4*time.Second, 12),
		)
		require.NoError(t, err)

		retryClient, ok := client.HTTPClient().(*RetryableHTTPClient)
		require.True(t, ok)
		transport, ok := retryClient.client.Transport.(*http.Transport)
		require.True(t, ok, "expected the default *http.Transport, got %T", retryClient.client.Transport)
		assert.Equal(t, 31*time.Second, transport.IdleConnTimeout)
		assert.Equal(t, 6*time.Second, transport.TLSHandshakeTimeout)
		assert.Equal(t, 4*time.Second, transport.ExpectContinueTimeout)
		assert.Equal(t, 12, transport.MaxIdleConns)
		assert.NotNil(t, transport.DialContext)
		assert.NotNil(t, transport.Proxy)
	})
}

// TestWithHTTPClient_TakesPrecedence checks that WithHTTPClient wins over
// WithRoundTripper, whatever the order of the options.
func TestWithHTTPClient_TakesPrecedence(t *testing.T) {
	t.Parallel()

	for _, httpClientFirst := range []bool{true, false} {
		t.Run(fmt.Sprintf("http client first: %t", httpClientFirst), func(t *testing.T) {
			t.Parallel()

			custom := &mockHTTPStatusCode{statusCode: http.StatusOK, body: "custom"}
			rt := &countingRoundTripper{statuses: []int{http.StatusOK}}
			opts := []ClientOption{WithRoundTripper(rt), WithHTTPClient(custom)}
			if httpClientFirst {
				opts = []ClientOption{WithHTTPClient(custom), WithRoundTripper(rt)}
			}

			client, err := NewClient(context.Background(), opts...)
			require.NoError(t, err)
			assert.Same(t, custom, client.HTTPClient())

			health, err := client.GetHealth(context.Background())
			require.NoError(t, err)
			assert.Equal(t, "custom", health)
			assert.Equal(t, int32(0), rt.calls.Load(), "the round tripper must not be used")
		})
	}
}

// TestWithHTTPClient_SkipsRetries documents the trade-off of WithHTTPClient: the
// client is used as is, so the SDK's retry options do not apply to it.
func TestWithHTTPClient_SkipsRetries(t *testing.T) {
	t.Parallel()

	rt := &countingRoundTripper{statuses: []int{http.StatusServiceUnavailable, http.StatusOK}}
	client, err := NewClient(
		context.Background(),
		WithHTTPClient(&http.Client{Transport: rt}),
		WithRequestRetryCount(5),
		WithBackoff(time.Millisecond, 10*time.Millisecond, 2.0, 0),
	)
	require.NoError(t, err)

	_, err = client.GetHealth(context.Background())
	require.ErrorIs(t, err, ErrRequestFailed)
	assert.Equal(t, int32(1), rt.calls.Load(), "a custom HTTP client is never retried by the SDK")
}

// TestWithoutEnvAPIKey checks that the WHATS_ON_CHAIN_API_KEY lookup can be turned
// off (issue #223), and that an explicit key still wins in every combination.
func TestWithoutEnvAPIKey(t *testing.T) {
	t.Setenv(EnvAPIKey, "env-api-key-123")

	tests := []struct {
		name string
		opts []ClientOption
		want string
	}{
		{"env key is used by default", nil, "env-api-key-123"},
		{"empty explicit key falls back to the env key", []ClientOption{WithAPIKey("")}, "env-api-key-123"},
		{"opt-out ignores the env key", []ClientOption{WithoutEnvAPIKey()}, ""},
		{"opt-out with an empty explicit key stays empty", []ClientOption{WithAPIKey(""), WithoutEnvAPIKey()}, ""},
		{"opt-out keeps an explicit key", []ClientOption{WithoutEnvAPIKey(), WithAPIKey("explicit-key")}, "explicit-key"},
		{"explicit key before the opt-out is kept", []ClientOption{WithAPIKey("explicit-key"), WithoutEnvAPIKey()}, "explicit-key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(context.Background(), tt.opts...)
			require.NoError(t, err)
			assert.Equal(t, tt.want, client.APIKey())
		})
	}
}

// TestWithoutEnvAPIKey_NoHeaderSent checks on the wire that the opt-out keeps
// the environment key out of requests.
func TestWithoutEnvAPIKey_NoHeaderSent(t *testing.T) {
	t.Setenv(EnvAPIKey, "env-api-key-123")

	var gotKey atomic.Value
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey.Store(r.Header.Get(apiHeaderKey))
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	tests := []struct {
		name       string
		optOut     bool
		wantHeader string
	}{
		{"env key is sent by default", false, "env-api-key-123"},
		{"opt-out sends no key", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []ClientOption{WithBaseURL(server.URL + "/v1/")}
			if tt.optOut {
				opts = append(opts, WithoutEnvAPIKey())
			}

			client, err := NewClient(context.Background(), opts...)
			require.NoError(t, err)

			_, err = client.GetHealth(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tt.wantHeader, gotKey.Load())
		})
	}
}

// BenchmarkNewClient benchmarks the NewClient method
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewClient(context.Background())
	}
}

// BenchmarkNewClientWithOptions benchmarks NewClient with options
func BenchmarkNewClientWithOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewClient(
			context.Background(),
			WithNetwork(NetworkTest),
			WithAPIKey("test-key"),
			WithRateLimit(10),
		)
	}
}
