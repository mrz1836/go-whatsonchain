package whatsonchain

import (
	"context"
	"fmt"
	"net/http"
)

// GetHealth simple endpoint to show API server is up and running
//
// For more information: https://developers.whatsonchain.com/#health
func (c *Client) GetHealth(ctx context.Context) (string, error) {

	// https://api.whatsonchain.com/v1/bsv/<network>/woc
	return c.request(
		ctx,
		fmt.Sprintf("%s%s/woc", apiEndpoint, c.Network),
		http.MethodGet, nil,
	)
}
