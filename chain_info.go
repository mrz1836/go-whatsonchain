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
		fmt.Sprintf("%s%s/%s/chain/info", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodGet,
		nil,
	); err != nil {
		return chainInfo, err
	}

	if len(resp) == 0 {
		return nil, ErrChainInfoNotFound
	}
	err = json.Unmarshal([]byte(resp), &chainInfo)
	return chainInfo, err
}

// GetCirculatingSupply this endpoint retrieves the current circulating supply
//
// For more information: https://developers.whatsonchain.com/#get-circulating-supply
func (c *Client) GetCirculatingSupply(ctx context.Context) (supply float64, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/circulatingsupply
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/circulatingsupply", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodGet, nil,
	); err != nil {
		return supply, err
	}

	return strconv.ParseFloat(strings.TrimSpace(resp), 64)
}

// GetChainTips this endpoint retrieves the chain tips
//
// For more information: https://developers.whatsonchain.com/#get-chain-tips
func (c *Client) GetChainTips(ctx context.Context) (chainTips []*ChainTip, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/chain/tips
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/chain/tips", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodGet,
		nil,
	); err != nil {
		return chainTips, err
	}

	if len(resp) == 0 {
		return nil, ErrChainTipsNotFound
	}
	err = json.Unmarshal([]byte(resp), &chainTips)
	return chainTips, err
}

// GetPeerInfo this endpoint retrieves information about peers connected to the node
//
// For more information: https://developers.whatsonchain.com/#get-peer-info
func (c *Client) GetPeerInfo(ctx context.Context) (peerInfo []*PeerInfo, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/peer/info
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/peer/info", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodGet,
		nil,
	); err != nil {
		return peerInfo, err
	}

	if len(resp) == 0 {
		return nil, ErrPeerInfoNotFound
	}
	err = json.Unmarshal([]byte(resp), &peerInfo)
	return peerInfo, err
}
