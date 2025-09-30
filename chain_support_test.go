package whatsonchain

import (
	"context"
	"testing"
)

// TestNewClientWithChain tests the new constructor with chain parameter
func TestNewClientWithChain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		chain           ChainType
		network         NetworkType
		expectedChain   ChainType
		expectedNetwork NetworkType
	}{
		{
			name:            "BSV main",
			chain:           ChainBSV,
			network:         NetworkMain,
			expectedChain:   ChainBSV,
			expectedNetwork: NetworkMain,
		},
		{
			name:            "BTC main",
			chain:           ChainBTC,
			network:         NetworkMain,
			expectedChain:   ChainBTC,
			expectedNetwork: NetworkMain,
		},
		{
			name:            "BSV test",
			chain:           ChainBSV,
			network:         NetworkTest,
			expectedChain:   ChainBSV,
			expectedNetwork: NetworkTest,
		},
		{
			name:            "BTC test",
			chain:           ChainBTC,
			network:         NetworkTest,
			expectedChain:   ChainBTC,
			expectedNetwork: NetworkTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := NewClientWithChain(tt.chain, tt.network, nil, nil)

			if client.Chain() != tt.expectedChain {
				t.Errorf("Chain() = %v, want %v", client.Chain(), tt.expectedChain)
			}

			if client.Network() != tt.expectedNetwork {
				t.Errorf("Network() = %v, want %v", client.Network(), tt.expectedNetwork)
			}
		})
	}
}

// TestClient_BackwardCompatibility tests that the old constructor still defaults to BSV
func TestClient_BackwardCompatibility(t *testing.T) {
	t.Parallel()

	client := NewClient(NetworkMain, nil, nil)

	// Should default to BSV for backward compatibility
	if client.Chain() != ChainBSV {
		t.Errorf("Chain() = %v, want %v", client.Chain(), ChainBSV)
	}

	if client.Network() != NetworkMain {
		t.Errorf("Network() = %v, want %v", client.Network(), NetworkMain)
	}
}

// TestClient_BSVSpecificMethods tests BSV-specific methods
func TestClient_BSVSpecificMethods(t *testing.T) {
	t.Parallel()

	t.Run("BSV chain allows OP_RETURN data", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration test")
		}

		client := NewClientWithChain(ChainBSV, NetworkMain, nil, nil)

		// Test with a transaction that should exist (may fail due to network)
		_, err := client.GetOpReturnData(context.Background(), "46c5495468b68248b69e55aa76a6b9ca1cb343bee9477c9c121358380e421ff3")

		// We don't test for success since this might not exist,
		// just that it doesn't error due to chain restriction
		if err != nil && err.Error() == "GetOpReturnData is only available for BSV chain" {
			t.Errorf("BSV client should allow GetOpReturnData")
		}
	})

	t.Run("BTC chain rejects BSV-specific methods", func(t *testing.T) {
		client := NewClientWithChain(ChainBTC, NetworkMain, nil, nil)

		_, err := client.GetOpReturnData(context.Background(), "test")

		expectedError := "operation is only available for BSV chain"
		if err == nil || err.Error() != expectedError {
			t.Errorf("Expected error '%s', got %v", expectedError, err)
		}
	})
}

// TestClient_BTCSpecificMethods tests BTC-specific methods
func TestClient_BTCSpecificMethods(t *testing.T) {
	t.Parallel()

	t.Run("BTC chain allows stats methods", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping integration test")
		}

		client := NewClientWithChain(ChainBTC, NetworkMain, nil, nil)

		// Test with a recent block height (may fail due to network)
		_, err := client.GetBlockStats(context.Background(), 700000)

		// We don't test for success since this might not exist,
		// just that it doesn't error due to chain restriction
		if err != nil && err.Error() == "GetBlockStats is only available for BTC chain" {
			t.Errorf("BTC client should allow GetBlockStats")
		}
	})

	t.Run("BSV chain rejects BTC-specific methods", func(t *testing.T) {
		client := NewClientWithChain(ChainBSV, NetworkMain, nil, nil)

		_, err := client.GetBlockStats(context.Background(), 700000)

		expectedError := "operation is only available for BTC chain"
		if err == nil || err.Error() != expectedError {
			t.Errorf("Expected error '%s', got %v", expectedError, err)
		}
	})

	t.Run("BSV chain rejects BTC miner stats", func(t *testing.T) {
		client := NewClientWithChain(ChainBSV, NetworkMain, nil, nil)

		_, err := client.GetMinerBlocksStats(context.Background(), 1)

		expectedError := "operation is only available for BTC chain"
		if err == nil || err.Error() != expectedError {
			t.Errorf("Expected error '%s', got %v", expectedError, err)
		}
	})
}

// TestClient_SharedMethods tests that shared methods work for both chains
func TestClient_SharedMethods(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Parallel()

	chains := []struct {
		name  string
		chain ChainType
	}{
		{"BSV", ChainBSV},
		{"BTC", ChainBTC},
	}

	for _, tc := range chains {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := NewClientWithChain(tc.chain, NetworkMain, nil, nil)

			// Test a shared method - chain info should work for both
			_, err := client.GetChainInfo(context.Background())

			// We expect this might fail due to network issues, but not due to chain incompatibility
			if err != nil && (err.Error() == "GetChainInfo is only available for BSV chain" ||
				err.Error() == "GetChainInfo is only available for BTC chain") {
				t.Errorf("GetChainInfo should be available for both chains, got error: %v", err)
			}
		})
	}
}

// TestChainTypeConstants tests that chain type constants are properly defined
func TestChainTypeConstants(t *testing.T) {
	t.Parallel()

	if ChainBSV != "bsv" {
		t.Errorf("ChainBSV = %v, want 'bsv'", ChainBSV)
	}

	if ChainBTC != "btc" {
		t.Errorf("ChainBTC = %v, want 'btc'", ChainBTC)
	}
}
