package whatsonchain

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPChainInfoBenchmark provides mock HTTP responses for chain info benchmarks
type mockHTTPChainInfoBenchmark struct{}

func (m *mockHTTPChainInfoBenchmark) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusOK

	// Chain info
	if strings.Contains(req.URL.String(), "/chain/info") {
		resp.Body = io.NopCloser(strings.NewReader(`{"chain":"main","blocks":700000,"headers":700000,"bestblockhash":"00000000000000000373d17e4f4f8e0f0f3f3e3d3c3b3a39383736353433323","difficulty":1234567.89,"mediantime":1609456000,"verificationprogress":0.999999,"chainwork":"00000000000000000000000000000000000000000000000000000000000000ff","pruned":false}`))
		return resp, nil
	}

	// Chain tips
	if strings.Contains(req.URL.String(), "/chain/tips") {
		resp.Body = io.NopCloser(strings.NewReader(`[{"height":700000,"hash":"00000000000000000373d17e4f4f8e0f0f3f3e3d3c3b3a39383736353433323","branchlen":0,"status":"active"}]`))
		return resp, nil
	}

	// Circulating supply
	if strings.Contains(req.URL.String(), "/circulatingsupply") {
		resp.Body = io.NopCloser(strings.NewReader(`18750000.00000000`))
		return resp, nil
	}

	// Exchange rate
	if strings.Contains(req.URL.String(), "/exchangerate") && !strings.Contains(req.URL.String(), "/historical") {
		resp.Body = io.NopCloser(strings.NewReader(`{"rate":50000.50,"currency":"USD","time":1609459200}`))
		return resp, nil
	}

	// Historical exchange rate
	if strings.Contains(req.URL.String(), "/historical") {
		resp.Body = io.NopCloser(strings.NewReader(`[{"rate":50000.50,"currency":"USD","time":1609459200},{"rate":49500.25,"currency":"USD","time":1609459100}]`))
		return resp, nil
	}

	// Peer info
	if strings.Contains(req.URL.String(), "/peer/info") {
		resp.Body = io.NopCloser(strings.NewReader(`[{"id":1,"addr":"192.168.1.1:8333","addrlocal":"192.168.1.100:54321","services":"000000000000040d","relaytxes":true,"lastsend":1609459200,"lastrecv":1609459200,"bytessent":123456,"bytesrecv":654321,"conntime":1609450000,"timeoffset":0,"pingtime":0.05,"minping":0.03,"version":70015,"subver":"/Bitcoin SV:1.0.0/","inbound":false,"addnode":false,"startingheight":700000,"txninvsize":0,"banscore":0,"synced_headers":700000,"synced_blocks":700000,"whitelisted":false}]`))
		return resp, nil
	}

	// Mempool info - mempoolminfee should be int64 (satoshis)
	if strings.Contains(req.URL.String(), "/mempool/info") {
		resp.Body = io.NopCloser(strings.NewReader(`{"size":5000,"bytes":2500000,"usage":3000000,"maxmempool":300000000,"mempoolminfee":1000}`))
		return resp, nil
	}

	// Mempool transactions
	if strings.Contains(req.URL.String(), "/mempool/raw") {
		resp.Body = io.NopCloser(strings.NewReader(`["tx1","tx2","tx3"]`))
		return resp, nil
	}

	resp.Body = io.NopCloser(strings.NewReader(`{}`))
	return resp, nil
}

// BenchmarkGetChainInfo benchmarks getting chain information
func BenchmarkGetChainInfo(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPChainInfoBenchmark{}),
	)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetChainInfo(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetChainTips benchmarks getting chain tips
func BenchmarkGetChainTips(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPChainInfoBenchmark{}),
	)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetChainTips(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetCirculatingSupply benchmarks getting circulating supply
func BenchmarkGetCirculatingSupply(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPChainInfoBenchmark{}),
	)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetCirculatingSupply(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetExchangeRate benchmarks getting current exchange rate
func BenchmarkGetExchangeRate(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPChainInfoBenchmark{}),
	)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetExchangeRate(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetHistoricalExchangeRate benchmarks getting historical exchange rates
func BenchmarkGetHistoricalExchangeRate(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPChainInfoBenchmark{}),
	)

	ctx := context.Background()

	tests := []struct {
		name     string
		fromTime int64
		toTime   int64
	}{
		{"1Hour", 1609455600, 1609459200},
		{"1Day", 1609372800, 1609459200},
		{"1Week", 1608854400, 1609459200},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := client.GetHistoricalExchangeRate(ctx, tt.fromTime, tt.toTime)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkGetPeerInfo benchmarks getting peer information
func BenchmarkGetPeerInfo(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPChainInfoBenchmark{}),
	)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetPeerInfo(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetMempoolInfo benchmarks getting mempool information
func BenchmarkGetMempoolInfo(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPChainInfoBenchmark{}),
	)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetMempoolInfo(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetMempoolTransactions benchmarks getting mempool transactions
func BenchmarkGetMempoolTransactions(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPChainInfoBenchmark{}),
	)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetMempoolTransactions(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}
