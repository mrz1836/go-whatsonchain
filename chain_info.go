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
	url := fmt.Sprintf("%s%s/chain/info", apiEndpoint, c.Parameters.Network)
	if resp, err = c.Request(url, http.MethodGet, nil); err != nil {
		return
	}

	chainInfo = new(ChainInfo)
	err = json.Unmarshal([]byte(resp), chainInfo)
	return
}
