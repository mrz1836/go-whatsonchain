package whatsonchain

import (
	"context"
)

// GetHealth simple endpoint to show API server is up and running
//
// For more information: https://docs/#health
func (c *Client) GetHealth(ctx context.Context) (string, error) {
	url := c.buildURL("/woc")
	return requestString(ctx, c, url)
}
