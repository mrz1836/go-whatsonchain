package whatsonchain

import (
	"context"
	"fmt"
	"net/http"
)

// GetExplorerLinks this endpoint identifies whether the posted query text is a block hash, txid or address and
// responds with WoC links. Ideal for extending customized search in apps.
//
// For more information: https://docs.whatsonchain.com/#get-history
func (c *Client) GetExplorerLinks(ctx context.Context, query string) (SearchResults, error) {
	postData := []byte(fmt.Sprintf(`{"query":"%s"}`, query))
	url := c.buildURL("/search/links")
	result, err := requestAndUnmarshal[SearchResults](ctx, c, url, http.MethodPost, postData, ErrChainInfoNotFound)
	if err != nil {
		return SearchResults{}, err
	}
	return *result, nil
}
