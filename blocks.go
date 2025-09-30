package whatsonchain

import (
	"context"
	"fmt"
	"net/http"
)

// GetBlockByHash this endpoint retrieves block details with given hash.
//
// For more information: https://docs.whatsonchain.com/#get-by-hash
func (c *Client) GetBlockByHash(ctx context.Context, hash string) (*BlockInfo, error) {
	url := c.buildURL("/block/hash/%s", hash)
	return requestAndUnmarshal[BlockInfo](ctx, c, url, http.MethodGet, nil, ErrBlockNotFound)
}

// GetBlockByHeight this endpoint retrieves block details with given block height.
//
// For more information: https://docs.whatsonchain.com/#get-by-height
func (c *Client) GetBlockByHeight(ctx context.Context, height int64) (*BlockInfo, error) {
	url := c.buildURL("/block/height/%d", height)
	return requestAndUnmarshal[BlockInfo](ctx, c, url, http.MethodGet, nil, ErrBlockNotFound)
}

// GetBlockPages if the block has more than 1000 transactions the page URIs will
// be provided in the "pages element" when getting a block by hash or height.
//
// For more information: https://docs.whatsonchain.com/#get-block-pages
func (c *Client) GetBlockPages(ctx context.Context, hash string, page int) (BlockPagesInfo, error) {
	url := c.buildURL("/block/hash/%s/page/%d", hash, page)
	return requestAndUnmarshalSlice[string](ctx, c, url, http.MethodGet, nil, ErrBlockNotFound)
}

// GetHeaderByHash this endpoint retrieves block header details with given hash.
//
// For more information: https://docs.whatsonchain.com/#get-header-by-hash
func (c *Client) GetHeaderByHash(ctx context.Context, hash string) (*BlockInfo, error) {
	url := c.buildURL("/block/%s/header", hash)
	return requestAndUnmarshal[BlockInfo](ctx, c, url, http.MethodGet, nil, ErrBlockNotFound)
}

// GetHeaders this endpoint retrieves last 10 block headers.
//
// For more information: https://docs.whatsonchain.com/#get-headers
func (c *Client) GetHeaders(ctx context.Context) ([]*BlockInfo, error) {
	url := c.buildURL("/block/headers")
	return requestAndUnmarshalSlice[*BlockInfo](ctx, c, url, http.MethodGet, nil, ErrHeadersNotFound)
}

// GetHeaderBytesFileLinks this endpoint retrieves header bytes file links.
//
// For more information: https://docs.whatsonchain.com/#get-header-bytes
func (c *Client) GetHeaderBytesFileLinks(ctx context.Context) (*HeaderBytesResource, error) {
	url := c.buildURL("/block/headers/resources")
	return requestAndUnmarshal[HeaderBytesResource](ctx, c, url, http.MethodGet, nil, ErrHeadersNotFound)
}

// GetLatestHeaderBytes this endpoint retrieves latest header bytes.
//
// For more information: https://docs.whatsonchain.com/#get-latest-headers
func (c *Client) GetLatestHeaderBytes(ctx context.Context, count int) (string, error) {
	path := "/block/headers/latest"
	if count > 0 {
		path = fmt.Sprintf("%s?count=%d", path, count)
	}
	url := c.buildURL(path)
	resp, err := requestString(ctx, c, url)
	if err != nil {
		return "", err
	}
	if len(resp) == 0 {
		return "", ErrHeadersNotFound
	}
	return resp, nil
}
