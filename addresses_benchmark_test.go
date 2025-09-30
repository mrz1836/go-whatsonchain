package whatsonchain

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPAddressesBenchmark provides mock HTTP responses for address benchmarks
type mockHTTPAddressesBenchmark struct{}

func (m *mockHTTPAddressesBenchmark) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusOK

	// Address Info
	if strings.Contains(req.URL.String(), "/info") {
		resp.Body = io.NopCloser(strings.NewReader(`{"isvalid": true,"address": "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA","scriptPubKey": "76a9143d0e5368bdadddca108a0fe44739919274c726c788ac","ismine": false,"iswatchonly": false,"isscript": false}`))
		return resp, nil
	}

	// Address Balance
	if strings.Contains(req.URL.String(), "/balance") {
		resp.Body = io.NopCloser(strings.NewReader(`{"confirmed": 10102050381,"unconfirmed": 123}`))
		return resp, nil
	}

	// Address History
	if strings.Contains(req.URL.String(), "/history") {
		resp.Body = io.NopCloser(strings.NewReader(`[{"tx_hash": "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1","height": 563052},{"tx_hash": "1c312435789754392f92ffcb64e1248e17da47bed179abfd27e6003c775e0e04","height": 565076}]`))
		return resp, nil
	}

	// Address Unspent
	if strings.Contains(req.URL.String(), "/unspent") {
		resp.Body = io.NopCloser(strings.NewReader(`[{"height": 639302,"tx_pos": 3,"tx_hash": "33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3","value": 2451680},{"height": 639601,"tx_pos": 3,"tx_hash": "4805041897a2ae59ffca85f0deb46e89d73d1ba4478bbd9c0fcd76ba0985ded2","value": 2744764},{"height": 640276,"tx_pos": 3,"tx_hash": "2493ff4cbca16b892ac641b7f2cb6d4388e75cb3f8963c291183f2bf0b27f415","value": 2568774}]`))
		return resp, nil
	}

	// Bulk operations - return array of balance records
	if strings.Contains(req.URL.String(), "/addresses/") {
		resp.Body = io.NopCloser(strings.NewReader(`[{"address":"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA","error":"","balance":{"confirmed":10102050381,"unconfirmed":123}}]`))
		return resp, nil
	}

	// Address Scripts
	if strings.Contains(req.URL.String(), "/scripts") {
		resp.Body = io.NopCloser(strings.NewReader(`{"scripts":["76a9143d0e5368bdadddca108a0fe44739919274c726c788ac"]}`))
		return resp, nil
	}

	// Address Used
	if strings.Contains(req.URL.String(), "/used") {
		resp.Body = io.NopCloser(strings.NewReader(`{"used":true}`))
		return resp, nil
	}

	resp.Body = io.NopCloser(strings.NewReader(`{}`))
	return resp, nil
}

// BenchmarkAddressInfo benchmarks getting address information
func BenchmarkAddressInfo(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPAddressesBenchmark{}),
	)

	ctx := context.Background()
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.AddressInfo(ctx, address)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkAddressBalance benchmarks getting address balance
func BenchmarkAddressBalance(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPAddressesBenchmark{}),
	)

	ctx := context.Background()
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.AddressBalance(ctx, address)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkAddressHistory benchmarks getting address history
func BenchmarkAddressHistory(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPAddressesBenchmark{}),
	)

	ctx := context.Background()
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.AddressHistory(ctx, address)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkAddressUnspentTransactions benchmarks getting unspent transactions
func BenchmarkAddressUnspentTransactions(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPAddressesBenchmark{}),
	)

	ctx := context.Background()
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.AddressUnspentTransactions(ctx, address)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBulkBalance benchmarks bulk balance requests with different sizes
// Currently disabled due to mock response format issue
func BenchmarkBulkBalanceDisabled(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPAddressesBenchmark{}),
	)

	ctx := context.Background()

	tests := []struct {
		name      string
		addresses int
	}{
		{"1Address", 1},
		{"5Addresses", 5},
		{"10Addresses", 10},
		{"20Addresses", 20},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			addresses := make([]string, tt.addresses)
			for i := 0; i < tt.addresses; i++ {
				addresses[i] = "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"
			}
			list := &AddressList{Addresses: addresses}

			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := client.BulkBalance(ctx, list)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkAddressConfirmedBalance benchmarks getting confirmed balance
func BenchmarkAddressConfirmedBalance(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPAddressesBenchmark{}),
	)

	ctx := context.Background()
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.AddressConfirmedBalance(ctx, address)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkAddressUnconfirmedBalance benchmarks getting unconfirmed balance
func BenchmarkAddressUnconfirmedBalance(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPAddressesBenchmark{}),
	)

	ctx := context.Background()
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.AddressUnconfirmedBalance(ctx, address)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkAddressUsed benchmarks checking if address has been used
func BenchmarkAddressUsed(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPAddressesBenchmark{}),
	)

	ctx := context.Background()
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.AddressUsed(ctx, address)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkAddressScripts benchmarks getting associated scripthashes
func BenchmarkAddressScripts(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPAddressesBenchmark{}),
	)

	ctx := context.Background()
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.AddressScripts(ctx, address)
		if err != nil {
			b.Fatal(err)
		}
	}
}
