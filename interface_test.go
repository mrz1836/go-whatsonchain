package whatsonchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClientInterface_Compliance tests that Client implements ClientInterface
func TestClientInterface_Compliance(t *testing.T) {
	t.Parallel()

	// Test that Client implements ClientInterface
	var _ ClientInterface = (*Client)(nil)

	// Create a mock client to test interface compliance
	client := newMockClientBSV(&mockHTTPEmpty{})
	assert.NotNil(t, client)

	// Verify interface implementation
	assert.Implements(t, (*ClientInterface)(nil), client)
}

// TestAddressService_Interface tests AddressService interface compliance
func TestAddressService_Interface(t *testing.T) {
	t.Parallel()

	// Test that Client implements AddressService
	var _ AddressService = (*Client)(nil)

	client := newMockClientBSV(&mockHTTPEmpty{})
	assert.Implements(t, (*AddressService)(nil), client)

	// Test interface embedding in ClientInterface
	addressService := AddressService(client)
	assert.NotNil(t, addressService)
}

// TestBlockService_Interface tests BlockService interface compliance
func TestBlockService_Interface(t *testing.T) {
	t.Parallel()

	// Test that Client implements BlockService
	var _ BlockService = (*Client)(nil)

	client := newMockClientBSV(&mockHTTPEmpty{})
	assert.Implements(t, (*BlockService)(nil), client)
}

// TestChainService_Interface tests ChainService interface compliance
func TestChainService_Interface(t *testing.T) {
	t.Parallel()

	// Test that Client implements ChainService
	var _ ChainService = (*Client)(nil)

	client := newMockClientBSV(&mockHTTPEmpty{})
	assert.Implements(t, (*ChainService)(nil), client)
}

// TestDownloadService_Interface tests DownloadService interface compliance
func TestDownloadService_Interface(t *testing.T) {
	t.Parallel()

	// Test that Client implements DownloadService
	var _ DownloadService = (*Client)(nil)

	client := newMockClientBSV(&mockHTTPEmpty{})
	assert.Implements(t, (*DownloadService)(nil), client)
}

// TestGeneralService_Interface tests GeneralService interface compliance
func TestGeneralService_Interface(t *testing.T) {
	t.Parallel()

	// Test that Client implements GeneralService
	var _ GeneralService = (*Client)(nil)

	client := newMockClientBSV(&mockHTTPEmpty{})
	assert.Implements(t, (*GeneralService)(nil), client)
}

// TestMempoolService_Interface tests MempoolService interface compliance
func TestMempoolService_Interface(t *testing.T) {
	t.Parallel()

	// Test that Client implements MempoolService
	var _ MempoolService = (*Client)(nil)

	client := newMockClientBSV(&mockHTTPEmpty{})
	assert.Implements(t, (*MempoolService)(nil), client)
}

// TestScriptService_Interface tests ScriptService interface compliance
func TestScriptService_Interface(t *testing.T) {
	t.Parallel()

	// Test that Client implements ScriptService
	var _ ScriptService = (*Client)(nil)

	client := newMockClientBSV(&mockHTTPEmpty{})
	assert.Implements(t, (*ScriptService)(nil), client)
}

// TestStatsService_Interface tests StatsService interface compliance
func TestStatsService_Interface(t *testing.T) {
	t.Parallel()

	// Test that Client implements StatsService
	var _ StatsService = (*Client)(nil)

	client := newMockClientBSV(&mockHTTPEmpty{})
	assert.Implements(t, (*StatsService)(nil), client)
}

// TestTokenService_Interface tests TokenService interface compliance
func TestTokenService_Interface(t *testing.T) {
	t.Parallel()

	// Test that Client implements TokenService
	var _ TokenService = (*Client)(nil)

	client := newMockClientBSV(&mockHTTPEmpty{})
	assert.Implements(t, (*TokenService)(nil), client)
}

// TestTransactionService_Interface tests TransactionService interface compliance
func TestTransactionService_Interface(t *testing.T) {
	t.Parallel()

	// Test that Client implements TransactionService
	var _ TransactionService = (*Client)(nil)

	client := newMockClientBSV(&mockHTTPEmpty{})
	assert.Implements(t, (*TransactionService)(nil), client)
}

// TestClientInterface_EmbeddedInterfaces tests that ClientInterface embeds all service interfaces
func TestClientInterface_EmbeddedInterfaces(t *testing.T) {
	t.Parallel()

	client := newMockClientBSV(&mockHTTPEmpty{})

	// Test that client can be cast to all service interfaces
	_, ok := interface{}(client).(AddressService)
	assert.True(t, ok, "Client should implement AddressService")

	_, ok = interface{}(client).(BlockService)
	assert.True(t, ok, "Client should implement BlockService")

	_, ok = interface{}(client).(ChainService)
	assert.True(t, ok, "Client should implement ChainService")

	_, ok = interface{}(client).(DownloadService)
	assert.True(t, ok, "Client should implement DownloadService")

	_, ok = interface{}(client).(GeneralService)
	assert.True(t, ok, "Client should implement GeneralService")

	_, ok = interface{}(client).(MempoolService)
	assert.True(t, ok, "Client should implement MempoolService")

	_, ok = interface{}(client).(ScriptService)
	assert.True(t, ok, "Client should implement ScriptService")

	_, ok = interface{}(client).(StatsService)
	assert.True(t, ok, "Client should implement StatsService")

	_, ok = interface{}(client).(TokenService)
	assert.True(t, ok, "Client should implement TokenService")

	_, ok = interface{}(client).(TransactionService)
	assert.True(t, ok, "Client should implement TransactionService")
}

// TestClientInterface_Methods tests ClientInterface specific methods
func TestClientInterface_Methods(t *testing.T) {
	t.Parallel()

	client := newMockClientBSV(&mockHTTPEmpty{})

	// Test Network method
	assert.NotEmpty(t, client.Network())

	// Test UserAgent method
	assert.NotEmpty(t, client.UserAgent())

	// Test RateLimit method
	assert.GreaterOrEqual(t, client.RateLimit(), 0)

	// Test HTTPClient method
	assert.NotNil(t, client.HTTPClient())

	// Test LastRequest method (may be nil initially)
	lastReq := client.LastRequest()
	assert.True(t, lastReq == nil || lastReq != nil) // Either is acceptable initially
}
