package whatsonchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetBlockByHash this endpoint retrieves block details with given hash.
//
// For more information: https://developers.whatsonchain.com/#get-by-hash
func (c *Client) GetBlockByHash(hash string) (blockInfo *BlockInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/hash/<hash>
	url := fmt.Sprintf("%s%s/block/hash/%s", apiEndpoint, c.Parameters.Network, hash)
	if resp, err = c.Request(url, http.MethodGet, nil); err != nil {
		return
	}

	blockInfo = new(BlockInfo)
	err = json.Unmarshal([]byte(resp), blockInfo)
	return
}

// GetBlockByHeight this endpoint retrieves block details with given block height.
//
// For more information: https://developers.whatsonchain.com/#get-by-height
func (c *Client) GetBlockByHeight(height int64) (blockInfo *BlockInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/height/<height>
	url := fmt.Sprintf("%s%s/block/height/%d", apiEndpoint, c.Parameters.Network, height)
	if resp, err = c.Request(url, http.MethodGet, nil); err != nil {
		return
	}

	blockInfo = new(BlockInfo)
	err = json.Unmarshal([]byte(resp), blockInfo)
	return
}

// GetBlockPages If the block has more that 1000 transactions the page URIs will
// be provided in the pages element when getting a block by hash or height.
//
// For more information: https://developers.whatsonchain.com/#get-block-pages
func (c *Client) GetBlockPages(hash string, page int) (txList *BlockPagesInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/block/hash/<hash>/page/1
	url := fmt.Sprintf("%s%s/block/hash/%s/page/%d", apiEndpoint, c.Parameters.Network, hash, page)
	if resp, err = c.Request(url, http.MethodGet, nil); err != nil {
		return
	}

	txList = new(BlockPagesInfo)
	err = json.Unmarshal([]byte(resp), txList)
	return
}
