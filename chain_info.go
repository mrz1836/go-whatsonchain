package whatsonchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetChainInfo this endpoint retrieves various state info of the chain for the selected network.
//
// For more information: https://developers.whatsonchain.com/#chain-info
func (c *Client) GetChainInfo() (chainInfo *ChainInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/chain/info
	if resp, err = c.Request(fmt.Sprintf("%s%s/chain/info", apiEndpoint, c.Parameters.Network), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &chainInfo)
	return
}
