package whatsonchain

import (
	"context"
	"net/http"
)

// GetMempoolInfo this endpoint retrieves various info about the node's mempool for the selected network
//
// For more information: https://docs/#get-mempool-info
func (c *Client) GetMempoolInfo(ctx context.Context) (*MempoolInfo, error) {
	url := c.buildURL("/mempool/info")
	return requestAndUnmarshal[MempoolInfo](ctx, c, url, http.MethodGet, nil, ErrMempoolInfoNotFound)
}

// GetMempoolTransactions this endpoint will retrieve a list of transaction ids from the node's mempool
// for the selected network
//
// For more information: https://docs/#get-mempool-transactions
func (c *Client) GetMempoolTransactions(ctx context.Context) ([]string, error) {
	url := c.buildURL("/mempool/raw")
	return requestAndUnmarshalSlice[string](ctx, c, url, http.MethodGet, nil, ErrMempoolInfoNotFound)
}
