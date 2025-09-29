package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetMempoolInfo this endpoint retrieves various info about the node's mempool for the selected network
//
// For more information: https://developers.whatsonchain.com/#get-mempool-info
func (c *Client) GetMempoolInfo(ctx context.Context) (info *MempoolInfo, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/mempool/info
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/mempool/info", apiEndpoint, c.Network()),
		http.MethodGet, nil,
	); err != nil {
		return info, err
	}
	if len(resp) == 0 {
		return nil, ErrMempoolInfoNotFound
	}
	err = json.Unmarshal([]byte(resp), &info)
	return info, err
}

// GetMempoolTransactions this endpoint will retrieve a list of transaction ids from the node's mempool
// for the selected network
//
// For more information: https://developers.whatsonchain.com/#get-mempool-transactions
func (c *Client) GetMempoolTransactions(ctx context.Context) (transactions []string, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/mempool/raw
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/mempool/raw", apiEndpoint, c.Network()),
		http.MethodGet, nil,
	); err != nil {
		return transactions, err
	}
	if len(resp) == 0 {
		return nil, ErrMempoolInfoNotFound
	}
	err = json.Unmarshal([]byte(resp), &transactions)
	return transactions, err
}
