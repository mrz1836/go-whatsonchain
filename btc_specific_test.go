package whatsonchain

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockHTTPEmpty for simple mocking
type mockHTTPEmpty struct{}

// Do is a mock http request
func (m *mockHTTPEmpty) Do(_ *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusOK
	return resp, nil
}

// TestBTCService_Interface tests that Client implements BTCService interface
func TestBTCService_Interface(t *testing.T) {
	t.Parallel()

	// Test that Client implements BTCService interface
	var _ BTCService = (*Client)(nil)

	// Test that interface can be used with client
	client := newMockClientBTC(&mockHTTPEmpty{})
	assert.NotNil(t, client)

	// Verify client is of the correct type that implements the interface
	btcService, ok := interface{}(client).(BTCService)
	assert.True(t, ok, "Client should implement BTCService interface")
	assert.NotNil(t, btcService)
}

// TestBTCService_ChainSpecific tests BTC chain-specific behavior
func TestBTCService_ChainSpecific(t *testing.T) {
	t.Parallel()

	// Test BTC chain creation
	client := newMockClientBTC(&mockHTTPEmpty{})
	assert.NotNil(t, client)
	assert.Equal(t, ChainBTC, client.Chain())

	// Test that BTC service interface is still satisfied
	var _ BTCService = client
}

// TestBTCService_InterfaceCompletion tests interface is properly defined
func TestBTCService_InterfaceCompletion(t *testing.T) {
	t.Parallel()

	// Test that the interface type is defined correctly
	client := newMockClientBTC(&mockHTTPEmpty{})

	// Test casting to BTCService
	btcService := BTCService(client)
	assert.NotNil(t, btcService)
}

// TestBTCService_ChainComparison tests BSV vs BTC chain differences
func TestBTCService_ChainComparison(t *testing.T) {
	t.Parallel()

	// Test BSV vs BTC client differences
	bsvClient := newMockClientBSV(&mockHTTPEmpty{})
	btcClient := newMockClientBTC(&mockHTTPEmpty{})

	assert.Equal(t, ChainBSV, bsvClient.Chain())
	assert.Equal(t, ChainBTC, btcClient.Chain())
}
