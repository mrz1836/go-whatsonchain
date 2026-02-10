package whatsonchain

import (
	"context"
	"testing"
	"time"
)

// BenchmarkClientCreation benchmarks client creation with different option combinations
func BenchmarkClientCreation(b *testing.B) {
	tests := []struct {
		name string
		opts []ClientOption
	}{
		{
			name: "Minimal",
			opts: []ClientOption{},
		},
		{
			name: "WithChainAndNetwork",
			opts: []ClientOption{
				WithChain(ChainBSV),
				WithNetwork(NetworkMain),
			},
		},
		{
			name: "FullyConfigured",
			opts: []ClientOption{
				WithChain(ChainBSV),
				WithNetwork(NetworkMain),
				WithAPIKey("test-api-key"),
				WithUserAgent("test-agent/1.0"),
				WithRateLimit(10),
				WithRequestTimeout(30 * time.Second),
				WithRequestRetryCount(2),
				WithBackoff(2*time.Millisecond, 10*time.Millisecond, 2.0, 2*time.Millisecond),
				WithDialer(20*time.Second, 5*time.Second),
				WithTransport(20*time.Second, 5*time.Second, 3*time.Second, 10),
			},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			ctx := context.Background()
			for i := 0; i < b.N; i++ {
				client, err := NewClient(ctx, tt.opts...)
				if err != nil {
					b.Fatal(err)
				}
				_ = client
			}
		})
	}
}

// BenchmarkClientGetters benchmarks the getter methods
func BenchmarkClientGetters(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithAPIKey("test-key"),
		WithUserAgent("test-agent"),
		WithRateLimit(5),
	)

	tests := []struct {
		name string
		fn   func()
	}{
		{"Chain", func() { _ = client.Chain() }},
		{"Network", func() { _ = client.Network() }},
		{"UserAgent", func() { _ = client.UserAgent() }},
		{"RateLimit", func() { _ = client.RateLimit() }},
		{"APIKey", func() { _ = client.APIKey() }},
		{"RequestTimeout", func() { _ = client.RequestTimeout() }},
		{"RequestRetryCount", func() { _ = client.RequestRetryCount() }},
		{"LastRequest", func() { _ = client.LastRequest() }},
		{"HTTPClient", func() { _ = client.HTTPClient() }},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				tt.fn()
			}
		})
	}
}

// BenchmarkClientSetters benchmarks the setter methods
func BenchmarkClientSetters(b *testing.B) {
	tests := []struct {
		name string
		fn   func(*Client)
	}{
		{"SetChain", func(c *Client) { _ = c.SetChain(ChainBSV) }},
		{"SetNetwork", func(c *Client) { _ = c.SetNetwork(NetworkMain) }},
		{"SetAPIKey", func(c *Client) { c.SetAPIKey("test-key") }},
		{"SetUserAgent", func(c *Client) { c.SetUserAgent("agent") }},
		{"SetRateLimit", func(c *Client) { c.SetRateLimit(5) }},
		{"SetRequestTimeout", func(c *Client) { c.SetRequestTimeout(30 * time.Second) }},
		{"SetRequestRetryCount", func(c *Client) { c.SetRequestRetryCount(3) }},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			client, _ := NewClient(context.Background())
			c := client.(*Client)
			for i := 0; i < b.N; i++ {
				tt.fn(c)
			}
		})
	}
}

// BenchmarkBuildURL benchmarks the URL building (hot path)
func BenchmarkBuildURL(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
	)
	c := client.(*Client)

	tests := []struct {
		name string
		path string
		args []any
	}{
		{
			name: "Simple",
			path: "/chain/info",
			args: nil,
		},
		{
			name: "WithOneArg",
			path: "/tx/hash/%s",
			args: []any{"abc123def456"},
		},
		{
			name: "WithTwoArgs",
			path: "/block/hash/%s/page/%d",
			args: []any{"abc123def456", 1},
		},
		{
			name: "LongAddress",
			path: "/address/%s/balance",
			args: []any{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = c.buildURL(tt.path, tt.args...)
			}
		})
	}
}

// BenchmarkBackoffConfig benchmarks backoff configuration operations
func BenchmarkBackoffConfig(b *testing.B) {
	client, _ := NewClient(context.Background())
	c := client.(*Client)

	b.Run("Get", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _, _ = c.BackoffConfig()
		}
	})

	b.Run("Set", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			c.SetBackoffConfig(2*time.Millisecond, 10*time.Millisecond, 2.0, 2*time.Millisecond)
		}
	})
}

// BenchmarkDialerConfig benchmarks dialer configuration operations
func BenchmarkDialerConfig(b *testing.B) {
	client, _ := NewClient(context.Background())
	c := client.(*Client)

	b.Run("Get", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = c.DialerConfig()
		}
	})

	b.Run("Set", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			c.SetDialerConfig(20*time.Second, 5*time.Second)
		}
	})
}

// BenchmarkTransportConfig benchmarks transport configuration operations
func BenchmarkTransportConfig(b *testing.B) {
	client, _ := NewClient(context.Background())
	c := client.(*Client)

	b.Run("Get", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _, _ = c.TransportConfig()
		}
	})

	b.Run("Set", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			c.SetTransportConfig(20*time.Second, 5*time.Second, 3*time.Second, 10)
		}
	})
}
