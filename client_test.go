package whatsonchain

import (
	"context"
	"net/http"
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
	assert.Equal(t, ChainBSV, client.Chain())
	assert.Equal(t, NetworkMain, client.Network())
	assert.Equal(t, defaultUserAgent, client.UserAgent())
	assert.Equal(t, defaultRateLimit, client.RateLimit())
	assert.Empty(t, client.APIKey())
}

// TestWithChain tests the WithChain option
func TestWithChain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		chain ChainType
	}{
		{"BSV", ChainBSV},
		{"BTC", ChainBTC},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(context.Background(), WithChain(tt.chain))
			require.NoError(t, err)
			assert.Equal(t, tt.chain, client.Chain())
		})
	}
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
	chain := ChainBTC
	network := NetworkTest
	rateLimit := 5

	client, err := NewClient(
		context.Background(),
		WithAPIKey(apiKey),
		WithUserAgent(userAgent),
		WithChain(chain),
		WithNetwork(network),
		WithRateLimit(rateLimit),
	)
	require.NoError(t, err)

	assert.Equal(t, apiKey, client.APIKey())
	assert.Equal(t, userAgent, client.UserAgent())
	assert.Equal(t, chain, client.Chain())
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

// TestSetChain tests the SetChain method
func TestSetChain(t *testing.T) {
	t.Parallel()

	client, err := NewClient(context.Background())
	require.NoError(t, err)

	require.NoError(t, client.SetChain(ChainBTC))
	assert.Equal(t, ChainBTC, client.Chain())
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
			WithChain(ChainBTC),
			WithNetwork(NetworkTest),
			WithAPIKey("test-key"),
			WithRateLimit(10),
		)
	}
}
