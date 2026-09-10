package whatsonchain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeprecated_DefaultChain verifies the client always reports the BSV chain.
func TestDeprecated_DefaultChain(t *testing.T) {
	t.Parallel()

	client, err := NewClient(context.Background())
	require.NoError(t, err)
	assert.Equal(t, ChainBSV, client.Chain())
}

// TestDeprecated_WithChainIsNoOp verifies WithChain no longer changes behavior;
// BSV is always the configured chain and requests always target the /bsv/ path.
func TestDeprecated_WithChainIsNoOp(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		context.Background(),
		WithChain(ChainBSV),
		WithNetwork(NetworkMain),
	)
	require.NoError(t, err)

	assert.Equal(t, ChainBSV, client.Chain())

	c, ok := client.(*Client)
	require.True(t, ok)
	assert.Equal(t, "https://api.whatsonchain.com/v1/bsv/main/health", c.buildURL("/health"))
}

// TestDeprecated_SetChainIsNoOp verifies SetChain is a no-op that returns nil
// and never alters the resolved chain or request URL.
func TestDeprecated_SetChainIsNoOp(t *testing.T) {
	t.Parallel()

	client, err := NewClient(context.Background(), WithNetwork(NetworkTest))
	require.NoError(t, err)

	require.NoError(t, client.SetChain(ChainBSV))
	assert.Equal(t, ChainBSV, client.Chain())

	c, ok := client.(*Client)
	require.True(t, ok)
	assert.Equal(t, "https://api.whatsonchain.com/v1/bsv/test/health", c.buildURL("/health"))
}

// TestDeprecated_AllNetworksTargetBSV verifies every supported network resolves
// to a /bsv/ base URL.
func TestDeprecated_AllNetworksTargetBSV(t *testing.T) {
	t.Parallel()

	tests := []struct {
		network  NetworkType
		expected string
	}{
		{NetworkMain, "https://api.whatsonchain.com/v1/bsv/main/health"},
		{NetworkTest, "https://api.whatsonchain.com/v1/bsv/test/health"},
		{NetworkStn, "https://api.whatsonchain.com/v1/bsv/stn/health"},
	}

	for _, tt := range tests {
		t.Run(string(tt.network), func(t *testing.T) {
			t.Parallel()

			client := newMockClient(&mockHTTPEmpty{})
			require.NoError(t, client.SetNetwork(tt.network))

			c, ok := client.(*Client)
			require.True(t, ok)
			assert.Equal(t, tt.expected, c.buildURL("/health"))
		})
	}
}
