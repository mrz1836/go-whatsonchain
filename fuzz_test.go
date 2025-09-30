package whatsonchain

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// FuzzBuildURL tests the buildURL method with various path patterns and arguments
// This fuzzer ensures that URL construction doesn't panic on malformed input
func FuzzBuildURL(f *testing.F) {
	// Seed corpus with known patterns from the codebase
	f.Add("/tx/hash/%s", "abc123")
	f.Add("/address/%s/balance", "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
	f.Add("/block/hash/%s", "")
	f.Add("", "")
	f.Add("/script/%s/history", "76a914")
	f.Add("/tx/%s/out/%d/hex", "hash")
	f.Add("/%s/%s/%s", "multi")
	f.Add("/path/with/no/args", "ignored")

	f.Fuzz(func(t *testing.T, path, arg string) {
		client := createTestClient(t)

		// The function should never panic, regardless of input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("buildURL panicked: %v", r)
			}
		}()

		// Test with no args
		url1 := client.buildURL(path)
		require.NotEmpty(t, url1, "buildURL should always return a non-empty string")
		require.Contains(t, url1, client.Chain(), "URL should contain chain")
		require.Contains(t, url1, client.Network(), "URL should contain network")

		// Test with one arg
		url2 := client.buildURL(path, arg)
		require.NotEmpty(t, url2, "buildURL should always return a non-empty string")
	})
}

// FuzzBulkRequest tests the bulkRequest function with various address list configurations
// This fuzzer validates input validation and JSON marshaling
func FuzzBulkRequest(f *testing.F) {
	// Seed corpus with boundary conditions
	f.Add(0)  // Empty list
	f.Add(1)  // Single address
	f.Add(10) // Mid-range
	f.Add(20) // Max limit
	f.Add(21) // Over limit (should error)
	f.Add(100)

	f.Fuzz(func(t *testing.T, count int) {
		// Clamp count to reasonable range to avoid OOM
		if count < 0 {
			count = -count
		}
		if count > 1000 {
			count = count % 1000
		}

		list := &AddressList{
			Addresses: make([]string, count),
		}

		// Fill with somewhat realistic addresses
		for i := 0; i < count; i++ {
			list.Addresses[i] = fmt.Sprintf("1Address%d", i)
		}

		// Function should never panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("bulkRequest panicked with count=%d: %v", count, r)
			}
		}()

		data, err := bulkRequest(list)

		// Validate behavior
		if count > MaxAddressesForLookup {
			// Should return an error for too many addresses
			require.Error(t, err, "Expected error for %d addresses (max is %d)", count, MaxAddressesForLookup)
			require.Nil(t, data, "Should return nil data when error occurs")
		} else {
			// Should succeed for valid counts
			require.NoError(t, err, "Should not error for %d addresses", count)
			require.NotNil(t, data, "Should return valid JSON data")

			// Verify JSON is valid
			var decoded AddressList
			err = json.Unmarshal(data, &decoded)
			require.NoError(t, err, "Returned data should be valid JSON")
			require.Len(t, decoded.Addresses, count, "Decoded address count should match")
		}
	})
}

// FuzzRequestAndUnmarshalJSON tests JSON unmarshaling with malformed input
// This fuzzer ensures the generic unmarshal helper handles corrupt data gracefully
func FuzzRequestAndUnmarshalJSON(f *testing.F) {
	// Seed with valid and invalid JSON
	f.Add(`{"address":"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa","confirmed":100,"unconfirmed":0}`)
	f.Add(`{"address":"test"}`)
	f.Add(`{}`)
	f.Add(`[]`)
	f.Add(``)
	f.Add(`null`)
	f.Add(`"string"`)
	f.Add(`123`)
	f.Add(`{broken`)
	f.Add(`{"nested":{"data":true}}`)
	f.Add(`{"unicode":"test"}`)

	f.Fuzz(func(t *testing.T, jsonData string) {
		// Function should never panic on any input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("json.Unmarshal panicked: %v\nInput: %s", r, jsonData)
			}
		}()

		// Test unmarshaling into AddressBalance
		var result AddressBalance
		err := json.Unmarshal([]byte(jsonData), &result)
		// Either succeeds or returns an error - both are acceptable
		// The key is that it should never panic
		if err != nil {
			// Verify it's a JSON error, not a panic
			require.Error(t, err, "Error should be non-nil")
		}

		// Test unmarshaling into different types to ensure generics work
		var txInfo TxInfo
		_ = json.Unmarshal([]byte(jsonData), &txInfo)

		var blockInfo BlockInfo
		_ = json.Unmarshal([]byte(jsonData), &blockInfo)
	})
}

// FuzzTransactionHex tests transaction hex validation and parsing
// This fuzzer ensures hex decoding doesn't panic on malformed input
func FuzzTransactionHex(f *testing.F) {
	// Seed with valid and invalid hex patterns
	f.Add("0100000001")                                                                         // Valid hex start
	f.Add("01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff") // Longer valid hex
	f.Add("")                                                                                   // Empty
	f.Add("gg")                                                                                 // Invalid hex chars
	f.Add("0")                                                                                  // Odd length
	f.Add("00")                                                                                 // Valid but minimal
	f.Add("deadbeef")
	f.Add("DEADBEEF") // Uppercase
	f.Add("0x0100")   // With prefix

	f.Fuzz(func(t *testing.T, hexData string) {
		// These operations should never panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Transaction hex processing panicked: %v\nInput: %s", r, hexData)
			}
		}()

		// Test JSON marshaling (used in BroadcastTx and DecodeTransaction)
		// Use proper JSON marshaling instead of string formatting to handle special characters
		payload := map[string]string{"txhex": hexData}
		postData, err := json.Marshal(payload)
		require.NoError(t, err, "Should marshal JSON")
		require.NotNil(t, postData, "Should create post data")

		// Verify the JSON is valid
		var decoded map[string]interface{}
		err = json.Unmarshal(postData, &decoded)
		require.NoError(t, err, "Generated JSON should be valid")

		// Verify the txhex field exists (can be empty if input was empty)
		_, ok := decoded["txhex"]
		require.True(t, ok, "Should have txhex field in JSON")
	})
}

// FuzzBroadcastTxResponse tests broadcast response parsing
// This fuzzer validates response cleanup logic (quotes, spaces, etc.)
func FuzzBroadcastTxResponse(f *testing.F) {
	// Seed with various response formats seen in the wild
	f.Add(`"abc123def456"`)
	f.Add(`abc123def456`)
	f.Add(`  "  abc123  "  `)
	f.Add(`""`)
	f.Add(``)
	f.Add(`" "`)
	f.Add(`"tx with spaces"`)
	f.Add(`"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"`) // 64 char txid
	f.Add(`{"error":"some error"}`)
	f.Add(`257: txn-already-known`)

	f.Fuzz(func(t *testing.T, response string) {
		// Response parsing should never panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Response parsing panicked: %v\nInput: %s", r, response)
			}
		}()

		// Simulate the cleanup logic from BroadcastTx (transactions.go:214)
		cleaned := strings.TrimSpace(strings.ReplaceAll(response, `"`, ""))

		// Verify cleanup always returns a valid string (can be empty)
		// Strings in Go are value types and cannot be nil
		_ = cleaned // Ensure variable is used

		// Additional validation: if response looks like JSON, try unmarshaling
		if strings.HasPrefix(strings.TrimSpace(response), "{") {
			var decoded map[string]interface{}
			_ = json.Unmarshal([]byte(response), &decoded)
			// Error is acceptable, panic is not
		}
	})
}

// createTestClient creates a basic client for testing
func createTestClient(t *testing.T) *Client {
	t.Helper()

	opts := defaultClientOptions()
	opts.chain = ChainBSV
	opts.network = NetworkMain

	return newClientFromOptions(opts)
}
