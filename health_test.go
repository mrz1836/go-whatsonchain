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

// TestClient_GetHealthWithNetworks tests the GetHealth() method across all BSV networks
func TestClient_GetHealthWithNetworks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		network NetworkType
	}{
		{name: "main network", network: NetworkMain},
		{name: "test network", network: NetworkTest},
		{name: "stn network", network: NetworkStn},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create client with a specific network
			client, err := NewClient(context.Background(), WithNetwork(tt.network), WithHTTPClient(&mockHTTPHealthBSV{}))
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()

			// Test GetHealth
			result, err := client.GetHealth(ctx)
			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
			}
			if result != "Whats On Chain" {
				t.Errorf("Expected 'Whats On Chain', got '%s'", result)
			}

			// Verify the client has the correct network
			if client.Network() != tt.network {
				t.Errorf("Expected network %s, got %s", tt.network, client.Network())
			}
		})
	}
}

// TestClient_GetHealthURLConstruction tests that the correct BSV URLs are constructed per network
func TestClient_GetHealthURLConstruction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		network     NetworkType
		expectedURL string
	}{
		{
			name:        "main network URL",
			network:     NetworkMain,
			expectedURL: "https://api.whatsonchain.com/v1/bsv/main/woc",
		},
		{
			name:        "test network URL",
			network:     NetworkTest,
			expectedURL: "https://api.whatsonchain.com/v1/bsv/test/woc",
		},
		{
			name:        "stn network URL",
			network:     NetworkStn,
			expectedURL: "https://api.whatsonchain.com/v1/bsv/stn/woc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(context.Background(), WithNetwork(tt.network), WithHTTPClient(&mockHTTPHealthBSV{}))
			if err != nil {
				t.Fatal(err)
			}

			c, ok := client.(*Client)
			if !ok {
				t.Fatal("expected *Client")
			}
			if got := c.buildURL("/woc"); got != tt.expectedURL {
				t.Errorf("Expected URL %s, got %s", tt.expectedURL, got)
			}

			// The request should succeed against the BSV mock
			if _, err = client.GetHealth(context.Background()); err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
			}
		})
	}
}
