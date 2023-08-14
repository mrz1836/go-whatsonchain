package whatsonchain

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewClient test new client
func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("main net", func(t *testing.T) {
		client := NewClient(NetworkMain, nil, nil)
		require.NotNil(t, client)
		assert.NotEqual(t, 0, len(client.UserAgent()))
	})

	t.Run("test net", func(t *testing.T) {
		client := NewClient(NetworkTest, nil, nil)
		require.NotNil(t, client)
		assert.NotEqual(t, 0, len(client.UserAgent()))
	})

	t.Run("stn", func(t *testing.T) {
		client := NewClient(NetworkStn, nil, nil)
		require.NotNil(t, client)
		assert.NotEqual(t, 0, len(client.UserAgent()))
	})

	t.Run("with API key", func(t *testing.T) {
		opts := ClientDefaultOptions()
		opts.APIKey = "test1234567"
		client := NewClient(NetworkStn, opts, nil)
		require.NotNil(t, client)
		assert.NotEqual(t, 0, len(client.UserAgent()))
	})
}

// TestNewClient_CustomHTTPClient test new client
func TestNewClient_CustomHTTPClient(t *testing.T) {
	t.Parallel()

	h := http.DefaultClient
	client := NewClient(NetworkTest, nil, h)
	require.NotNil(t, client)
	assert.Equal(t, h, client.HTTPClient())
}

// ExampleNewClient example using NewClient()
func ExampleNewClient() {
	client := NewClient(NetworkTest, nil, nil)
	fmt.Println(client.UserAgent())
	// Output:go-whatsonchain: v0.13.0
}

// BenchmarkNewClient benchmarks the NewClient method
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewClient(NetworkTest, nil, nil)
	}
}

// TestClientDefaultOptions tests the method ClientDefaultOptions()
func TestClientDefaultOptions(t *testing.T) {
	t.Parallel()

	options := ClientDefaultOptions()
	require.NotNil(t, options)
	assert.Equal(t, defaultUserAgent, options.UserAgent)
	assert.Equal(t, defaultRateLimit, options.RateLimit)
	assert.Equal(t, 2.0, options.BackOffExponentFactor)
	assert.Equal(t, 2*time.Millisecond, options.BackOffInitialTimeout)
	assert.Equal(t, 2*time.Millisecond, options.BackOffMaximumJitterInterval)
	assert.Equal(t, 10*time.Millisecond, options.BackOffMaxTimeout)
	assert.Equal(t, 20*time.Second, options.DialerKeepAlive)
	assert.Equal(t, 5*time.Second, options.DialerTimeout)
	assert.Equal(t, 2, options.RequestRetryCount)
	assert.Equal(t, 30*time.Second, options.RequestTimeout)
	assert.Equal(t, 3*time.Second, options.TransportExpectContinueTimeout)
	assert.Equal(t, 20*time.Second, options.TransportIdleTimeout)
	assert.Equal(t, 10, options.TransportMaxIdleConnections)
	assert.Equal(t, 5*time.Second, options.TransportTLSHandshakeTimeout)
	assert.Equal(t, "", options.APIKey)
}

// TestClientDefaultOptions_NoRetry will set 0 retry counts
func TestClientDefaultOptions_NoRetry(t *testing.T) {
	options := ClientDefaultOptions()
	require.NotNil(t, options)

	options.RequestRetryCount = 0

	client := NewClient(NetworkTest, options, nil)
	require.NotNil(t, client)
	assert.Equal(t, defaultUserAgent, client.UserAgent())
}
