package whatsonchain

import (
	"context"
	"testing"
)

// TestClientWithChain tests the new constructor with chain parameter
func TestClientWithChain(t *testing.T) {
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

			client, err := NewClient(
				context.Background(),
				WithChain(tt.chain),
				WithNetwork(tt.network),
			)
			if err != nil {
				t.Fatal(err)
			}

			if client.Chain() != tt.expectedChain {
				t.Errorf("Chain() = %v, want %v", client.Chain(), tt.expectedChain)
			}

			if client.Network() != tt.expectedNetwork {
				t.Errorf("Network() = %v, want %v", client.Network(), tt.expectedNetwork)
			}
		})
	}
}

// TestClient_DefaultChain tests that the new constructor defaults to BSV
func TestClient_DefaultChain(t *testing.T) {
	t.Parallel()

	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Should default to BSV for backward compatibility
	if client.Chain() != ChainBSV {
		t.Errorf("Chain() = %v, want %v", client.Chain(), ChainBSV)
	}

	if client.Network() != NetworkMain {
		t.Errorf("Network() = %v, want %v", client.Network(), NetworkMain)
	}
}

// TestChainSupport_BSV tests BSV-specific functionality
func TestChainSupport_BSV(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
	)
	if err != nil {
		t.Fatal(err)
	}

	if client.Chain() != ChainBSV {
		t.Errorf("Expected BSV chain, got %v", client.Chain())
	}

	// Verify that BSV client has BSV-specific interface methods
	// This is just a compile-time check that BSV methods exist
	var _ BSVService = client
}

// TestChainSupport_BTC tests BTC-specific functionality
func TestChainSupport_BTC(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		context.Background(),
		WithChain(ChainBTC),
		WithNetwork(NetworkMain),
	)
	if err != nil {
		t.Fatal(err)
	}

	if client.Chain() != ChainBTC {
		t.Errorf("Expected BTC chain, got %v", client.Chain())
	}

	// Verify that BTC client has BTC-specific interface methods
	// This is just a compile-time check that BTC methods exist
	var _ BTCService = client
}

// TestChain_NetworkSupport tests all network combinations
func TestChain_NetworkSupport(t *testing.T) {
	t.Parallel()

	networks := []NetworkType{NetworkMain, NetworkTest, NetworkStn}
	chains := []ChainType{ChainBSV, ChainBTC}

	for _, chain := range chains {
		for _, network := range networks {
			t.Run(string(chain)+"-"+string(network), func(t *testing.T) {
				client, err := NewClient(
					context.Background(),
					WithChain(chain),
					WithNetwork(network),
				)
				if err != nil {
					t.Fatal(err)
				}

				if client.Chain() != chain {
					t.Errorf("Chain() = %v, want %v", client.Chain(), chain)
				}

				if client.Network() != network {
					t.Errorf("Network() = %v, want %v", client.Network(), network)
				}
			})
		}
	}
}

// TestClient_ChainSetter tests the SetChain method
func TestClient_ChainSetter(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		context.Background(),
		WithChain(ChainBSV),
	)
	if err != nil {
		t.Fatal(err)
	}

	if client.Chain() != ChainBSV {
		t.Errorf("Initial Chain() = %v, want %v", client.Chain(), ChainBSV)
	}

	// Change to BTC
	if err = client.SetChain(ChainBTC); err != nil {
		t.Fatal(err)
	}

	if client.Chain() != ChainBTC {
		t.Errorf("After SetChain, Chain() = %v, want %v", client.Chain(), ChainBTC)
	}
}

// TestClient_NetworkSetter tests the SetNetwork method
func TestClient_NetworkSetter(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		context.Background(),
		WithNetwork(NetworkMain),
	)
	if err != nil {
		t.Fatal(err)
	}

	if client.Network() != NetworkMain {
		t.Errorf("Initial Network() = %v, want %v", client.Network(), NetworkMain)
	}

	// Change to test network
	if err = client.SetNetwork(NetworkTest); err != nil {
		t.Fatal(err)
	}

	if client.Network() != NetworkTest {
		t.Errorf("After SetNetwork, Network() = %v, want %v", client.Network(), NetworkTest)
	}
}
