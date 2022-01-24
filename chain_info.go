package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// GetChainInfo this endpoint retrieves various state info of the chain for the selected network.
//
// For more information: https://developers.whatsonchain.com/#chain-info
func (c *Client) GetChainInfo(ctx context.Context) (chainInfo *ChainInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/chain/info
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/chain/info", apiEndpoint, c.Network()),
		http.MethodGet, nil,
	); err != nil {
		return
	}
	if len(resp) == 0 {
		return nil, ErrChainInfoNotFound
	}
	err = json.Unmarshal([]byte(resp), &chainInfo)
	return
}

// GetCirculatingSupply this endpoint retrieves the current circulating supply
//
// For more information: https://developers.whatsonchain.com/#get-circulating-supply
func (c *Client) GetCirculatingSupply(ctx context.Context) (supply float64, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/circulatingsupply
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/circulatingsupply", apiEndpoint, c.Network()),
		http.MethodGet, nil,
	); err != nil {
		return
	}

	return strconv.ParseFloat(strings.TrimSpace(resp), 64)
}
