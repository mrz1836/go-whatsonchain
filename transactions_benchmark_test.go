package whatsonchain

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPTransactionsBenchmark provides mock HTTP responses for transaction benchmarks
type mockHTTPTransactionsBenchmark struct{}

func (m *mockHTTPTransactionsBenchmark) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusOK

	txResponse := `{"txid":"293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946","hash":"293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946","version":1,"size":225,"locktime":0,"vin":[{"coinbase":"03c2c70e","txid":"0000000000000000000000000000000000000000000000000000000000000000","vout":4294967295,"scriptSig":{"asm":"03c2c70e","hex":"03c2c70e"},"sequence":4294967295}],"vout":[{"value":50.00000000,"n":0,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 a94a8f3b09b432f1e5c2c1c8b0b4d7e0b6e3e2f7 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914a94a8f3b09b432f1e5c2c1c8b0b4d7e0b6e3e2f788ac","reqSigs":1,"type":"pubkeyhash","addresses":["1GSEjCJaEzPKrJWWqhPVPuUGpfzU2c2tz1"]}}],"blockhash":"00000000000000000373d17e4f4f8e0f0f3f3e3d3c3b3a39383736353433323","confirmations":100,"time":1609459200,"blocktime":1609459200}`

	// Single transaction
	if strings.Contains(req.URL.String(), "/tx/hash/") {
		resp.Body = io.NopCloser(strings.NewReader(txResponse))
		return resp, nil
	}

	// Bulk transactions
	if strings.Contains(req.URL.String(), "/txs") && req.Method == http.MethodPost {
		bulkResponse := `[` + txResponse + `,` + txResponse + `]`
		resp.Body = io.NopCloser(strings.NewReader(bulkResponse))
		return resp, nil
	}

	// Transaction status
	if strings.Contains(req.URL.String(), "/status") {
		resp.Body = io.NopCloser(strings.NewReader(`[{"txid":"293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946","valid":true,"height":100}]`))
		return resp, nil
	}

	// Raw transaction
	if strings.Contains(req.URL.String(), "/hex") {
		resp.Body = io.NopCloser(strings.NewReader(`01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0403c2c70effffffff0100f2052a01000000434104a94a8f3b09b432f1e5c2c1c8b0b4d7e0b6e3e2f7ac00000000`))
		return resp, nil
	}

	// Broadcast
	if strings.Contains(req.URL.String(), "/tx/raw") {
		resp.Body = io.NopCloser(strings.NewReader(`293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946`))
		return resp, nil
	}

	// Decode
	if strings.Contains(req.URL.String(), "/decode") {
		resp.Body = io.NopCloser(strings.NewReader(txResponse))
		return resp, nil
	}

	// Merkle proof
	if strings.Contains(req.URL.String(), "/proof") {
		resp.Body = io.NopCloser(strings.NewReader(`[{"blockHash":"00000000000000000373d17e4f4f8e0f0f3f3e3d3c3b3a39383736353433323","branches":[{"hash":"abc123","pos":"R"}],"hash":"293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946","merkleRoot":"def456"}]`))
		return resp, nil
	}

	// Spent output
	if strings.Contains(req.URL.String(), "/spent") {
		resp.Body = io.NopCloser(strings.NewReader(`{"txid":"abc123def456","vin":0}`))
		return resp, nil
	}

	resp.Body = io.NopCloser(strings.NewReader(`{}`))
	return resp, nil
}

// BenchmarkGetTxByHash benchmarks getting a transaction by hash
func BenchmarkGetTxByHash(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPTransactionsBenchmark{}),
	)

	ctx := context.Background()
	txHash := "293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetTxByHash(ctx, txHash)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBulkTransactionDetails benchmarks bulk transaction details
func BenchmarkBulkTransactionDetails(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPTransactionsBenchmark{}),
	)

	ctx := context.Background()

	tests := []struct {
		name  string
		count int
	}{
		{"1Tx", 1},
		{"5Txs", 5},
		{"10Txs", 10},
		{"20Txs", 20},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			txIDs := make([]string, tt.count)
			for i := 0; i < tt.count; i++ {
				txIDs[i] = "293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946"
			}
			hashes := &TxHashes{TxIDs: txIDs}

			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := client.BulkTransactionDetails(ctx, hashes)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkBulkTransactionStatus benchmarks bulk transaction status
func BenchmarkBulkTransactionStatus(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPTransactionsBenchmark{}),
	)

	ctx := context.Background()
	hashes := &TxHashes{
		TxIDs: []string{
			"293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946",
			"293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946",
		},
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.BulkTransactionStatus(ctx, hashes)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetRawTransactionData benchmarks getting raw transaction data
func BenchmarkGetRawTransactionData(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPTransactionsBenchmark{}),
	)

	ctx := context.Background()
	txHash := "293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetRawTransactionData(ctx, txHash)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBroadcastTx benchmarks broadcasting a transaction
func BenchmarkBroadcastTx(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPTransactionsBenchmark{}),
	)

	ctx := context.Background()
	txHex := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0403c2c70effffffff0100f2052a01000000434104a94a8f3b09b432f1e5c2c1c8b0b4d7e0b6e3e2f7ac00000000"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.BroadcastTx(ctx, txHex)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDecodeTransaction benchmarks decoding a transaction
func BenchmarkDecodeTransaction(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPTransactionsBenchmark{}),
	)

	ctx := context.Background()
	txHex := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0403c2c70effffffff0100f2052a01000000434104a94a8f3b09b432f1e5c2c1c8b0b4d7e0b6e3e2f7ac00000000"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.DecodeTransaction(ctx, txHex)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetMerkleProof benchmarks getting merkle proof
func BenchmarkGetMerkleProof(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPTransactionsBenchmark{}),
	)

	ctx := context.Background()
	txHash := "293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetMerkleProof(ctx, txHash)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetSpentOutput benchmarks getting spent output information
func BenchmarkGetSpentOutput(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPTransactionsBenchmark{}),
	)

	ctx := context.Background()
	txHash := "293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := client.GetSpentOutput(ctx, txHash, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBulkRawTransactionData benchmarks bulk raw transaction data
func BenchmarkBulkRawTransactionData(b *testing.B) {
	client, _ := NewClient(context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
		WithHTTPClient(&mockHTTPTransactionsBenchmark{}),
	)

	ctx := context.Background()

	tests := []struct {
		name  string
		count int
	}{
		{"1Tx", 1},
		{"10Txs", 10},
		{"20Txs", 20},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			txIDs := make([]string, tt.count)
			for i := 0; i < tt.count; i++ {
				txIDs[i] = "293cd46be8e436099e183a5bc145e19e7bb1ccc7a7f9c0bc842b3c6a992c7946"
			}
			hashes := &TxHashes{TxIDs: txIDs}

			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := client.BulkRawTransactionData(ctx, hashes)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
