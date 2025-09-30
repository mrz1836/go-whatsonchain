package whatsonchain

import "fmt"

// buildURL constructs a URL with the chain and network prefix
// This centralizes URL construction to avoid repetition across all API methods
func (c *Client) buildURL(path string, args ...interface{}) string {
	// Build the base URL with chain and network
	baseURL := fmt.Sprintf("%s%s/%s", apiEndpointBase, c.Chain(), c.Network())

	// If args are provided, format the path with them
	if len(args) > 0 {
		path = fmt.Sprintf(path, args...)
	}

	// Combine base and path
	return baseURL + path
}
