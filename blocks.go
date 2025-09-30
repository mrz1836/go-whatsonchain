package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetBlockByHash this endpoint retrieves block details with given hash.
//
// For more information: https://developers.whatsonchain.com/#get-by-hash
func (c *Client) GetBlockByHash(ctx context.Context, hash string) (blockInfo *BlockInfo, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/hash/<hash>
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/block/hash/%s", apiEndpointBase, c.Chain(), c.Network(), hash),
		http.MethodGet,
		nil,
	); err != nil {
		return blockInfo, err
	}

	if len(resp) == 0 {
		return nil, ErrBlockNotFound
	}
	err = json.Unmarshal([]byte(resp), &blockInfo)

	return blockInfo, err
}

// GetBlockByHeight this endpoint retrieves block details with given block height.
//
// For more information: https://developers.whatsonchain.com/#get-by-height
func (c *Client) GetBlockByHeight(ctx context.Context, height int64) (blockInfo *BlockInfo, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/height/<height>
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/block/height/%d", apiEndpointBase, c.Chain(), c.Network(), height),
		http.MethodGet, nil,
	); err != nil {
		return blockInfo, err
	}
	if len(resp) == 0 {
		return nil, ErrBlockNotFound
	}
	err = json.Unmarshal([]byte(resp), &blockInfo)
	return blockInfo, err
}

// GetBlockPages if the block has more than 1000 transactions the page URIs will
// be provided in the "pages element" when getting a block by hash or height.
//
// For more information: https://developers.whatsonchain.com/#get-block-pages
func (c *Client) GetBlockPages(ctx context.Context, hash string, page int) (txList BlockPagesInfo, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/hash/<hash>/page/1
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/block/hash/%s/page/%d", apiEndpointBase, c.Chain(), c.Network(), hash, page),
		http.MethodGet, nil,
	); err != nil {
		return txList, err
	}
	if len(resp) == 0 {
		return nil, ErrBlockNotFound
	}
	err = json.Unmarshal([]byte(resp), &txList)
	return txList, err
}

// GetHeaderByHash this endpoint retrieves block header details with given hash.
//
// For more information: https://developers.whatsonchain.com/#get-header-by-hash
func (c *Client) GetHeaderByHash(ctx context.Context, hash string) (headerInfo *BlockInfo, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/<hash>/header
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/block/%s/header", apiEndpointBase, c.Chain(), c.Network(), hash),
		http.MethodGet, nil,
	); err != nil {
		return headerInfo, err
	}
	if len(resp) == 0 {
		return nil, ErrBlockNotFound
	}
	err = json.Unmarshal([]byte(resp), &headerInfo)
	return headerInfo, err
}

// GetHeaders this endpoint retrieves last 10 block headers.
//
// For more information: https://developers.whatsonchain.com/#get-headers
func (c *Client) GetHeaders(ctx context.Context) (blockHeaders []*BlockInfo, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/headers
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/block/headers", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodGet, nil,
	); err != nil {
		return blockHeaders, err
	}
	if len(resp) == 0 {
		return nil, ErrHeadersNotFound
	}
	err = json.Unmarshal([]byte(resp), &blockHeaders)
	return blockHeaders, err
}

// GetHeaderBytesFileLinks this endpoint retrieves header bytes file links.
//
// For more information: https://developers.whatsonchain.com/#get-header-bytes
func (c *Client) GetHeaderBytesFileLinks(ctx context.Context) (resource *HeaderBytesResource, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/headers/resources
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/block/headers/resources", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodGet, nil,
	); err != nil {
		return resource, err
	}
	if len(resp) == 0 {
		return nil, ErrHeadersNotFound
	}
	err = json.Unmarshal([]byte(resp), &resource)
	return resource, err
}

// GetLatestHeaderBytes this endpoint retrieves latest header bytes.
//
// For more information: https://developers.whatsonchain.com/#get-latest-headers
func (c *Client) GetLatestHeaderBytes(ctx context.Context, count int) (headerBytes string, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/headers/latest?count=<count>
	url := fmt.Sprintf("%s%s/%s/block/headers/latest", apiEndpointBase, c.Chain(), c.Network())
	if count > 0 {
		url = fmt.Sprintf("%s?count=%d", url, count)
	}

	if resp, err = c.request(
		ctx,
		url,
		http.MethodGet, nil,
	); err != nil {
		return headerBytes, err
	}
	if len(resp) == 0 {
		return "", ErrHeadersNotFound
	}
	return resp, nil
}
