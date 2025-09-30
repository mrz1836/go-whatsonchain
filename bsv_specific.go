package whatsonchain

import (
	"context"
	"fmt"
	"net/http"
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

	// https://api.whatsonchain.com/v1/bsv/<network>/tx/<txHash>/opreturn
	return c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/tx/%s/opreturn", apiEndpointBase, c.Chain(), c.Network(), txHash),
		http.MethodGet, nil,
	)
}
