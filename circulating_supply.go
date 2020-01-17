package whatsonchain

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// GetCirculatingSupply this endpoint retrieves the current circulating supply
//
// For more information: (undocumented) //todo: add link once in documentation
func (c *Client) GetCirculatingSupply() (supply float64, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/circulatingsupply
	if resp, err = c.Request(fmt.Sprintf("%s%s/circulatingsupply", apiEndpoint, c.Parameters.Network), http.MethodGet, nil); err != nil {
		return
	}

	supply, err = strconv.ParseFloat(strings.TrimSpace(resp), 64)
	return
}
