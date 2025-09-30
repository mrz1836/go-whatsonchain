package whatsonchain

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPHealthValid for mocking requests
type mockHTTPHealthValid struct{}

// Do is a mock http request
func (m *mockHTTPHealthValid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	// Valid
	if strings.Contains(req.URL.String(), "/woc") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`Whats On Chain`))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPHealthBSV for mocking BSV chain requests
type mockHTTPHealthBSV struct{}

// Do is a mock http request that validates BSV chain in URL
func (m *mockHTTPHealthBSV) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	// Valid BSV health endpoint
	if strings.Contains(req.URL.String(), "/bsv/") && strings.Contains(req.URL.String(), "/woc") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`Whats On Chain`))
	} else {
		// Return empty body for non-matching requests to avoid nil pointer
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, ErrBadRequest
	}

	return resp, nil
}

// mockHTTPHealthBTC for mocking BTC chain requests
type mockHTTPHealthBTC struct{}

// Do is a mock http request that validates BTC chain in URL
func (m *mockHTTPHealthBTC) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	// Valid BTC health endpoint
	if strings.Contains(req.URL.String(), "/btc/") && strings.Contains(req.URL.String(), "/woc") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`Whats On Chain`))
	} else {
		// Return empty body for non-matching requests to avoid nil pointer
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, ErrBadRequest
	}

	return resp, nil
}

// mockHTTPHealthInvalid for mocking requests
type mockHTTPHealthInvalid struct{}

// Do is a mock http request
func (m *mockHTTPHealthInvalid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	// Invalid
	if strings.Contains(req.URL.String(), "/woc") {
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, ErrBadRequest
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetHealth tests the GetHealth()
func TestClient_GetHealth(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPHealthValid{})
	ctx := context.Background()

	// Test the valid response
	info, err := client.GetHealth(ctx)
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if info != "Whats On Chain" {
		t.Errorf("%s Failed: response was [%s] expected [%s]", t.Name(), info, "Whats On Chain")
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPHealthInvalid{})

	// Test invalid response
	_, err = client.GetHealth(ctx)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}

// TestClient_GetHealthWithChains tests the GetHealth() method with both BSV and BTC chains
func TestClient_GetHealthWithChains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		chain      ChainType
		network    NetworkType
		mockClient HTTPInterface
		expectErr  bool
	}{
		{
			name:       "BSV main network",
			chain:      ChainBSV,
			network:    NetworkMain,
			mockClient: &mockHTTPHealthBSV{},
			expectErr:  false,
		},
		{
			name:       "BSV test network",
			chain:      ChainBSV,
			network:    NetworkTest,
			mockClient: &mockHTTPHealthBSV{},
			expectErr:  false,
		},
		{
			name:       "BSV stn network",
			chain:      ChainBSV,
			network:    NetworkStn,
			mockClient: &mockHTTPHealthBSV{},
			expectErr:  false,
		},
		{
			name:       "BTC main network",
			chain:      ChainBTC,
			network:    NetworkMain,
			mockClient: &mockHTTPHealthBTC{},
			expectErr:  false,
		},
		{
			name:       "BTC test network",
			chain:      ChainBTC,
			network:    NetworkTest,
			mockClient: &mockHTTPHealthBTC{},
			expectErr:  false,
		},
		{
			name:       "BTC stn network",
			chain:      ChainBTC,
			network:    NetworkStn,
			mockClient: &mockHTTPHealthBTC{},
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create client with specific chain and network
			client := NewClientWithChain(tt.chain, tt.network, nil, tt.mockClient)
			ctx := context.Background()

			// Test GetHealth
			result, err := client.GetHealth(ctx)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tt.name, err)
				}
				if result != "Whats On Chain" {
					t.Errorf("Expected 'Whats On Chain', got '%s'", result)
				}
			}

			// Verify the client has the correct chain and network
			if client.Chain() != tt.chain {
				t.Errorf("Expected chain %s, got %s", tt.chain, client.Chain())
			}
			if client.Network() != tt.network {
				t.Errorf("Expected network %s, got %s", tt.network, client.Network())
			}
		})
	}
}

// TestClient_GetHealthURLConstruction tests that the correct URLs are constructed for different chains
func TestClient_GetHealthURLConstruction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		chain       ChainType
		network     NetworkType
		expectedURL string
		shouldFail  bool
	}{
		{
			name:        "BSV main network URL",
			chain:       ChainBSV,
			network:     NetworkMain,
			expectedURL: "https://api.whatsonchain.com/v1/bsv/main/woc",
			shouldFail:  false,
		},
		{
			name:        "BTC main network URL",
			chain:       ChainBTC,
			network:     NetworkMain,
			expectedURL: "https://api.whatsonchain.com/v1/btc/main/woc",
			shouldFail:  false,
		},
		{
			name:        "BSV test network URL",
			chain:       ChainBSV,
			network:     NetworkTest,
			expectedURL: "https://api.whatsonchain.com/v1/bsv/test/woc",
			shouldFail:  false,
		},
		{
			name:        "BTC test network URL",
			chain:       ChainBTC,
			network:     NetworkTest,
			expectedURL: "https://api.whatsonchain.com/v1/btc/test/woc",
			shouldFail:  false,
		},
		{
			name:        "Wrong chain in BSV mock",
			chain:       ChainBTC,
			network:     NetworkMain,
			expectedURL: "https://api.whatsonchain.com/v1/btc/main/woc",
			shouldFail:  true, // Will use BSV mock which expects /bsv/ in URL
		},
		{
			name:        "Wrong chain in BTC mock",
			chain:       ChainBSV,
			network:     NetworkMain,
			expectedURL: "https://api.whatsonchain.com/v1/bsv/main/woc",
			shouldFail:  true, // Will use BTC mock which expects /btc/ in URL
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Use wrong mock client for the last two test cases to verify URL validation
			var mockClient HTTPInterface
			if tt.shouldFail {
				if tt.chain == ChainBTC {
					mockClient = &mockHTTPHealthBSV{} // Wrong mock for BTC
				} else {
					mockClient = &mockHTTPHealthBTC{} // Wrong mock for BSV
				}
			} else {
				if tt.chain == ChainBSV {
					mockClient = &mockHTTPHealthBSV{}
				} else {
					mockClient = &mockHTTPHealthBTC{}
				}
			}

			client := NewClientWithChain(tt.chain, tt.network, nil, mockClient)
			ctx := context.Background()

			_, err := client.GetHealth(ctx)

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected request to fail for %s due to URL mismatch, but it succeeded", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tt.name, err)
				}
			}
		})
	}
}
