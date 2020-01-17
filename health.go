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
	url := fmt.Sprintf("%s%s/woc", apiEndpoint, c.Parameters.Network)
	return c.Request(url, http.MethodGet, nil)
}
