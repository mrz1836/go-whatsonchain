package whatsonchain

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPBlocksBenchmark provides mock HTTP responses for block benchmarks
type mockHTTPBlocksBenchmark struct{}

func (m *mockHTTPBlocksBenchmark) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusOK

	blockResponse := `{"hash":"00000000000000000373d17e4f4f8e0f0f3f3e3d3c3b3a39383736353433323","confirmations":100,"size":1234567,"height":700000,"version":536870912,"versionHex":"20000000","merkleroot":"abc123def456","tx":["tx1","tx2","tx3"],"txcount":3,"time":1609459200,"mediantime":1609456000,"nonce":123456789,"bits":"1a012345","difficulty":1234567.89,"chainwork":"00000000000000000000000000000000000000000000000000000000000000ff","previousblockhash":"00000000000000000373d17e4f4f8e0f0f3f3e3d3c3b3a39383736353433322","nextblockhash":"00000000000000000373d17e4f4f8e0f0f3f3e3d3c3b3a39383736353433324","coinbaseTx":{"txid":"coinbase123","hash":"coinbase123","version":1,"size":100,"locktime":0,"vin":[],"vout":[]},"totalFees":1.5,"miner":"Unknown","pages":{"size":3,"uri":["/v1/bsv/main/block/hash/00000000000000000373d17e4f4f8e0f0f3f3e3d3c3b3a39383736353433323/page/0"]}}`

	// Block by hash
	if strings.Contains(req.URL.String(), "/block/hash/") && !strings.Contains(req.URL.String(), "/page/") && !strings.Contains(req.URL.String(), "/header") {
		resp.Body = io.NopCloser(strings.NewReader(blockResponse))
		return resp, nil
	}

	// Block by height
	if strings.Contains(req.URL.String(), "/block/height/") {
		resp.Body = io.NopCloser(strings.NewReader(blockResponse))
		return resp, nil
	}

	// Block pages
	if strings.Contains(req.URL.String(), "/page/") {
		resp.Body = io.NopCloser(strings.NewReader(`["tx1","tx2","tx3"]`))
		return resp, nil
	}

	// Headers (last 10) - return array of block headers (check before /header to avoid false match)
	if strings.Contains(req.URL.String(), "/block/headers") && !strings.Contains(req.URL.String(), "/resources") && !strings.Contains(req.URL.String(), "/latest") {
		// Return a properly formatted block header array matching the API structure
		headersResponse := `[{"hash":"0000000000000000008605a63392a85ebc7e055af19334b2a2f3952e1fdeb3b2","confirmations":1,"height":700000,"version":536870912,"versionHex":"20000000","merkleroot":"abc123def456","time":1609459200,"mediantime":1609456000,"nonce":123456789,"bits":"1a012345","difficulty":1234567.89,"chainwork":"00000000000000000000000000000000000000000000000000000000000000ff","previousblockhash":"00000000000000000373d17e4f4f8e0f0f3f3e3d3c3b3a39383736353433322","nextblockhash":""}]`
		resp.Body = io.NopCloser(strings.NewReader(headersResponse))
		return resp, nil
	}

	// Header by hash
	if strings.Contains(req.URL.String(), "/header") {
		resp.Body = io.NopCloser(strings.NewReader(blockResponse))
		return resp, nil
	}

	// Header bytes resources
	if strings.Contains(req.URL.String(), "/resources") {
		resp.Body = io.NopCloser(strings.NewReader(`{"description":"Header bytes file links","links":[{"format":"hex","uri":"https://example.com/headers.hex"}]}`))
		return resp, nil
	}

	// Latest header bytes
	if strings.Contains(req.URL.String(), "/latest") {
		resp.Body = io.NopCloser(strings.NewReader(`abcdef0123456789`))
		return resp, nil
	}

	resp.Body = io.NopCloser(strings.NewReader(`{}`))
	return resp, nil
}

// BenchmarkGetBlockByHash benchmarks getting a block by hash
func BenchmarkGetBlockByHash(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPBlocksBenchmark{}),
	)

	ctx := context.Background()
	blockHash := "00000000000000000373d17e4f4f8e0f0f3f3e3d3c3b3a39383736353433323"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetBlockByHash(ctx, blockHash)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetBlockByHeight benchmarks getting a block by height
func BenchmarkGetBlockByHeight(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPBlocksBenchmark{}),
	)

	ctx := context.Background()

	tests := []struct {
		name   string
		height int64
	}{
		{"LowHeight", 100},
		{"MediumHeight", 500000},
		{"HighHeight", 700000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := client.GetBlockByHeight(ctx, tt.height)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkGetBlockPages benchmarks getting block pages
func BenchmarkGetBlockPages(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPBlocksBenchmark{}),
	)

	ctx := context.Background()
	blockHash := "00000000000000000373d17e4f4f8e0f0f3f3e3d3c3b3a39383736353433323"

	tests := []struct {
		name string
		page int
	}{
		{"Page0", 0},
		{"Page1", 1},
		{"Page10", 10},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := client.GetBlockPages(ctx, blockHash, tt.page)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkGetHeaderByHash benchmarks getting a block header by hash
func BenchmarkGetHeaderByHash(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPBlocksBenchmark{}),
	)

	ctx := context.Background()
	blockHash := "00000000000000000373d17e4f4f8e0f0f3f3e3d3c3b3a39383736353433323"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetHeaderByHash(ctx, blockHash)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetHeaders benchmarks getting the last 10 block headers
func BenchmarkGetHeaders(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPBlocksBenchmark{}),
	)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetHeaders(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetHeaderBytesFileLinks benchmarks getting header bytes file links
func BenchmarkGetHeaderBytesFileLinks(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPBlocksBenchmark{}),
	)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetHeaderBytesFileLinks(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetLatestHeaderBytes benchmarks getting latest header bytes
func BenchmarkGetLatestHeaderBytes(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPBlocksBenchmark{}),
	)

	ctx := context.Background()

	tests := []struct {
		name  string
		count int
	}{
		{"1Header", 1},
		{"10Headers", 10},
		{"100Headers", 100},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := client.GetLatestHeaderBytes(ctx, tt.count)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
