/*
Package whatsonchain is the unofficial golang implementation for the whatsonchain.com API

Example:

// Create a client
client, _ := whatsonchain.NewClient(whatsonchain.NetworkMain, nil)

// Get a balance for an address
balance, _ := client.AddressBalance("1JSSSgcyufLgbXFw6WAXyXgBrmgFpnqXWh")
*/
package whatsonchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"unsafe"
)

// NewClient creates a new client to submit requests
func NewClient(network NetworkType, clientOptions *Options) (c *Client, err error) {

	// Create a client using the given options
	c = createClient(clientOptions)

	// Set the network
	c.Parameters.Network = network

	return
}

// Request is a generic request wrapper that can be used without constraints
func (c *Client) Request(url string, method string, payload []byte) (response string, err error) {

	// Set reader
	var bodyReader io.Reader

	// Switch on Method
	switch method {
	case http.MethodPost, http.MethodPut:
		{
			bodyReader = bytes.NewBuffer(payload)
		}
	}

	// Store for debugging purposes
	c.LastRequest.Method = method
	c.LastRequest.URL = url

	// Start the request
	var request *http.Request
	if request, err = http.NewRequest(method, url, bodyReader); err != nil {
		return
	}

	// Change the header (user agent is in case they block default Go user agents)
	request.Header.Set("User-Agent", c.Parameters.UserAgent)

	// Set the content type on Method
	if method == http.MethodPost || method == http.MethodPut {
		request.Header.Set("Content-Type", "application/json")
	}

	// Fire the http request
	var resp *http.Response
	if resp, err = c.httpClient.Do(request); err != nil {
		return
	}

	// Close the response body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Save the status
	c.LastRequest.StatusCode = resp.StatusCode

	// Read the body
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	// Parse the response
	response = string(body)
	return
}

// GetHealth simple endpoint to show API server is up and running
//
// For more information: https://developers.whatsonchain.com/#health
func (c *Client) GetHealth() (status string, err error) {

	// https://api.whatsonchain.com/v1/bsv/<network>/woc
	url := fmt.Sprintf("%s%s/woc", apiEndpoint, c.Parameters.Network)
	return c.Request(url, http.MethodGet, nil)
}

// GetChainInfo this endpoint retrieves various state info of the chain for the selected network.
//
// For more information: https://developers.whatsonchain.com/#chain-info
func (c *Client) GetChainInfo() (chainInfo *ChainInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/chain/info
	url := fmt.Sprintf("%s%s/chain/info", apiEndpoint, c.Parameters.Network)
	if resp, err = c.Request(url, http.MethodGet, nil); err != nil {
		return
	}

	chainInfo = new(ChainInfo)
	err = json.Unmarshal([]byte(resp), chainInfo)
	return
}

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

// GetTxByHash this endpoint retrieves transaction details with given transaction hash
//
// For more information: https://developers.whatsonchain.com/#get-by-tx-hash
func (c *Client) GetTxByHash(hash string) (txInfo *TxInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tx/hash/<hash>
	url := fmt.Sprintf("%s%s/tx/hash/%s", apiEndpoint, c.Parameters.Network, hash)
	if resp, err = c.Request(url, http.MethodGet, nil); err != nil {
		return
	}

	txInfo = new(TxInfo)
	err = json.Unmarshal([]byte(resp), txInfo)
	return
}

// GetMerkleProof this endpoint returns merkle branch to a confirmed transaction
//
// For more information: https://developers.whatsonchain.com/#get-merkle-proof
func (c *Client) GetMerkleProof(hash string) (merkleInfo *MerkleInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tx/<hash>/merkleproof
	url := fmt.Sprintf("%s%s/tx/%s/merkleproof", apiEndpoint, c.Parameters.Network, hash)
	if resp, err = c.Request(url, http.MethodGet, nil); err != nil {
		return
	}

	merkleInfo = new(MerkleInfo)
	err = json.Unmarshal([]byte(resp), merkleInfo)
	return
}

// BroadcastTx will broadcast transaction using this endpoint.
// Get tx_id in response or error msg from node.
//
// For more information: https://developers.whatsonchain.com/#broadcast-transaction
func (c *Client) BroadcastTx(txHex string) (txID string, err error) {

	// Start the post data
	stringVal := fmt.Sprintf(`{"txhex":"%s"}`, txHex)
	postData := []byte(stringVal)

	// https://api.whatsonchain.com/v1/bsv/<network>/tx/raw
	url := fmt.Sprintf("%s%s/tx/raw", apiEndpoint, c.Parameters.Network)
	if txID, err = c.Request(url, http.MethodPost, postData); err != nil {
		return
	}

	// Got an error
	if c.LastRequest.StatusCode > 200 {
		err = fmt.Errorf("error broadcasting: %s", txID)
		txID = "" // remove the error message
	} else {
		// Remove quotes or spaces
		txID = strings.TrimSpace(strings.Replace(txID, `"`, "", -1))
	}

	return
}

// BulkBroadcastTx will broadcast many transactions at once
// You can bulk broadcast transactions using this endpoint.
// 		Size per transaction should be less than 100KB
//		Overall payload per request should be less than 10MB
//		Max 100 transactions per request
//		Only available for mainnet
//
// For more information: https://developers.whatsonchain.com/#bulk-broadcast
func (c *Client) BulkBroadcastTx(rawTxs []string, feedback bool) (response *BulkBroadcastResponse, err error) {

	// Set a max (from Whats on Chain)
	if len(rawTxs) > 100 {
		err = fmt.Errorf("max transactions are 100")
		return
	}

	// Set a total max
	if size := unsafe.Sizeof(rawTxs); size > 1e+7 {
		err = fmt.Errorf("max overall payload of 10MB (1e+7 bytes)")
		return
	}

	// Check size of each tx
	for _, tx := range rawTxs {
		if size := unsafe.Sizeof(tx); size > 102400 {
			err = fmt.Errorf("max tx size of 100kb (102400 bytes)")
			return
		}
	}

	// Start the post data
	postData, _ := json.Marshal(rawTxs)

	var resp string

	// https://api.whatsonchain.com/v1/bsv/tx/broadcast?feedback=<feedback>
	url := fmt.Sprintf("%stx/broadcast?feedback=%t", apiEndpoint, feedback)
	if resp, err = c.Request(url, http.MethodPost, postData); err != nil {
		return
	}

	response = new(BulkBroadcastResponse)
	response.Feedback = feedback
	if feedback {
		if err = json.Unmarshal([]byte(resp), response); err != nil {
			return
		}
	}

	// Got an error
	if c.LastRequest.StatusCode > 200 {
		err = fmt.Errorf("error broadcasting: %s", resp)
	}

	return
}

// AddressInfo this endpoint retrieves various address info.
//
// For more information: https://developers.whatsonchain.com/#address
func (c *Client) AddressInfo(address string) (addressInfo *AddressInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/info
	url := fmt.Sprintf("%s%s/address/%s/info", apiEndpoint, c.Parameters.Network, address)
	if resp, err = c.Request(url, http.MethodGet, nil); err != nil {
		return
	}

	addressInfo = new(AddressInfo)
	err = json.Unmarshal([]byte(resp), addressInfo)
	return
}

// AddressBalance this endpoint retrieves confirmed and unconfirmed address balance.
//
// For more information: https://developers.whatsonchain.com/#get-balance
func (c *Client) AddressBalance(address string) (balance *AddressBalance, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/balance
	url := fmt.Sprintf("%s%s/address/%s/balance", apiEndpoint, c.Parameters.Network, address)
	if resp, err = c.Request(url, http.MethodGet, nil); err != nil {
		return
	}

	balance = new(AddressBalance)
	err = json.Unmarshal([]byte(resp), balance)
	return
}

// AddressHistory this endpoint retrieves confirmed and unconfirmed address transactions.
//
// For more information: https://developers.whatsonchain.com/#get-history
func (c *Client) AddressHistory(address string) (history AddressHistory, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/history
	url := fmt.Sprintf("%s%s/address/%s/history", apiEndpoint, c.Parameters.Network, address)
	if resp, err = c.Request(url, http.MethodGet, nil); err != nil {
		return
	}

	history = *new(AddressHistory)
	err = json.Unmarshal([]byte(resp), &history)
	return
}

// AddressUnspentTransactions this endpoint retrieves ordered list of UTXOs.
//
// For more information: https://developers.whatsonchain.com/#get-unspent-transactions
func (c *Client) AddressUnspentTransactions(address string) (history AddressHistory, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/unspent
	url := fmt.Sprintf("%s%s/address/%s/unspent", apiEndpoint, c.Parameters.Network, address)
	if resp, err = c.Request(url, http.MethodGet, nil); err != nil {
		return
	}

	history = *new(AddressHistory)
	err = json.Unmarshal([]byte(resp), &history)
	return
}
