package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetExplorerLinks this endpoint identifies whether the posted query text is a block hash, txid or address and
// responds with WoC links. Ideal for extending customized search in apps.
//
// For more information: https://developers.whatsonchain.com/#get-history
func (c *Client) GetExplorerLinks(ctx context.Context, query string) (results SearchResults, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/search/links
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/search/links", apiEndpoint, c.Network()),
		http.MethodPost, []byte(fmt.Sprintf(`{"query":"%s"}`, query)),
	); err != nil {
		return
	}
	if len(resp) == 0 {
		return results, ErrChainInfoNotFound
	}
	err = json.Unmarshal([]byte(resp), &results)
	return
}
