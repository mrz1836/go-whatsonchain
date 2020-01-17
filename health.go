package whatsonchain

import (
	"fmt"
	"net/http"
)

// GetHealth simple endpoint to show API server is up and running
//
// For more information: https://developers.whatsonchain.com/#health
func (c *Client) GetHealth() (status string, err error) {

	// https://api.whatsonchain.com/v1/bsv/<network>/woc
	return c.Request(fmt.Sprintf("%s%s/woc", apiEndpoint, c.Parameters.Network), http.MethodGet, nil)
}
