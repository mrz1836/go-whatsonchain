package whatsonchain

import (
	"fmt"
	"net/url"
)

// buildURL constructs a URL with the chain and network prefix
// This centralizes URL construction to avoid repetition across all API methods
func (c *Client) buildURL(path string, args ...any) string {
	// Read both chain and network under a single lock to prevent
	// mismatched values if SetChain/SetNetwork is called concurrently
	c.optionsMu.RLock()
	chain := c.options.chain
	network := c.options.network
	c.optionsMu.RUnlock()

	// Build the base URL with chain and network
	baseURL := fmt.Sprintf("%s%s/%s", apiEndpointBase, chain, network)

	// If args are provided, escape string arguments and format the path
	if len(args) > 0 {
		for i, arg := range args {
			if s, ok := arg.(string); ok {
				args[i] = url.PathEscape(s)
			}
		}
		path = fmt.Sprintf(path, args...)
	}

	// Combine base and path
	return baseURL + path
}
