package whatsonchain

import (
	"context"
)

// BSVService is the interface for BSV-specific endpoints
type BSVService interface {
	GetOpReturnData(ctx context.Context, txHash string) (string, error)
	TokenService
}

// GetOpReturnData gets OP_RETURN data by transaction hash (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/#get-op_return-data-by-tx-hash
func (c *Client) GetOpReturnData(ctx context.Context, txHash string) (string, error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return "", ErrBSVChainRequired
	}

	url := c.buildURL("/tx/%s/opreturn", txHash)
	return requestString(ctx, c, url)
}
