package whatsonchain

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// GetChainInfo this endpoint retrieves various state info of the chain for the selected network.
//
// For more information: https://docs.whatsonchain.com/#chain-info
func (c *Client) GetChainInfo(ctx context.Context) (*ChainInfo, error) {
	url := c.buildURL("/chain/info")
	return requestAndUnmarshal[ChainInfo](ctx, c, url, http.MethodGet, nil, ErrChainInfoNotFound)
}

// GetCirculatingSupply this endpoint retrieves the current circulating supply
//
// For more information: https://docs.whatsonchain.com/#get-circulating-supply
func (c *Client) GetCirculatingSupply(ctx context.Context) (float64, error) {
	url := c.buildURL("/circulatingsupply")
	resp, err := requestString(ctx, c, url)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(strings.TrimSpace(resp), 64)
}

// GetChainTips this endpoint retrieves the chain tips
//
// For more information: https://docs.whatsonchain.com/#get-chain-tips
func (c *Client) GetChainTips(ctx context.Context) ([]*ChainTip, error) {
	url := c.buildURL("/chain/tips")
	return requestAndUnmarshalSlice[*ChainTip](ctx, c, url, http.MethodGet, nil, ErrChainTipsNotFound)
}

// GetPeerInfo this endpoint retrieves information about peers connected to the node
//
// For more information: https://docs.whatsonchain.com/#get-peer-info
func (c *Client) GetPeerInfo(ctx context.Context) ([]*PeerInfo, error) {
	url := c.buildURL("/peer/info")
	return requestAndUnmarshalSlice[*PeerInfo](ctx, c, url, http.MethodGet, nil, ErrPeerInfoNotFound)
}
