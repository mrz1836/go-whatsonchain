package whatsonchain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unsafe"
)

// GetTxByHash this endpoint retrieves transaction details with given transaction hash
//
// For more information: https://developers.whatsonchain.com/#get-by-tx-hash
func (c *Client) GetTxByHash(hash string) (txInfo *TxInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tx/hash/<hash>
	if resp, err = c.Request(fmt.Sprintf("%s%s/tx/hash/%s", apiEndpoint, c.Parameters.Network, hash), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &txInfo)
	return
}

// BulkTransactionDetails this fetches details for multiple transactions in single request
// Max 20 transactions per request
//
// For more information: https://developers.whatsonchain.com/#bulk-transaction-details
func (c *Client) BulkTransactionDetails(hashes *TxHashes) (txList TxList, err error) {

	// Max limit by WOC
	if len(hashes.TxIDs) == 0 {
		err = fmt.Errorf("missing hashes")
		return
	} else if len(hashes.TxIDs) > MaxTransactionsUTXO {
		err = fmt.Errorf("max limit of utxos is %d and you sent %d", MaxTransactionsUTXO, len(hashes.TxIDs))
		return
	}

	// Hashes into json
	var postData []byte
	if postData, err = json.Marshal(hashes); err != nil {
		return
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/txs
	if resp, err = c.Request(fmt.Sprintf("%s%s/txs", apiEndpoint, c.Parameters.Network), http.MethodPost, postData); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &txList)
	return
}

// GetMerkleProof this endpoint returns merkle branch to a confirmed transaction
//
// For more information: https://developers.whatsonchain.com/#get-merkle-proof
func (c *Client) GetMerkleProof(hash string) (merkleResults MerkleResults, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tx/<hash>/proof
	if resp, err = c.Request(fmt.Sprintf("%s%s/tx/%s/proof", apiEndpoint, c.Parameters.Network, hash), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &merkleResults)
	return
}

// GetRawTransactionData this endpoint returns raw hex for the transaction with given hash
//
// For more information: https://developers.whatsonchain.com/#get-raw-transaction-data
func (c *Client) GetRawTransactionData(hash string) (hex string, err error) {

	// https://api.whatsonchain.com/v1/bsv/<network>/tx/<hash>/hex
	hex, err = c.Request(fmt.Sprintf("%s%s/tx/%s/hex", apiEndpoint, c.Parameters.Network, hash), http.MethodGet, nil)

	return
}

// GetRawTransactionOutputData this endpoint returns raw hex for the transaction output with given hash and index
//
// For more information: https://developers.whatsonchain.com/#get-raw-transaction-output-data
func (c *Client) GetRawTransactionOutputData(hash string, vOutIndex int) (hex string, err error) {

	// https://api.whatsonchain.com/v1/bsv/<network>/tx/<hash>/out/<index>/hex
	hex, err = c.Request(fmt.Sprintf("%s%s/tx/%s/out/%d/hex", apiEndpoint, c.Parameters.Network, hash, vOutIndex), http.MethodGet, nil)

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
	if txID, err = c.Request(fmt.Sprintf("%s%s/tx/raw", apiEndpoint, c.Parameters.Network), http.MethodPost, postData); err != nil {
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
// Tip: First transaction in the list should have an output to WOC tip address '16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA'
//
// Feedback: true/false: true if response from the node is required for each transaction, otherwise, set it to false.
// (For stress testing set it to false). When set to true a unique url is provided to check the progress of the
// submitted transactions, eg 'QUEUED' or 'PROCESSED', with response data from node. You can poll the provided unique
// url until all transactions are marked as 'PROCESSED'. Progress of the transactions are tracked on this unique url
// for up to 5 hours.
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
	if resp, err = c.Request(fmt.Sprintf("%stx/broadcast?feedback=%t", apiEndpoint, feedback), http.MethodPost, postData); err != nil {
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

// DecodeTransaction this endpoint decodes raw transaction
//
// For more information: https://developers.whatsonchain.com/#decode-transaction
func (c *Client) DecodeTransaction(txHex string) (txInfo *TxInfo, err error) {

	// Start the post data
	stringVal := fmt.Sprintf(`{"txhex":"%s"}`, txHex)
	postData := []byte(stringVal)

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tx/decode
	if resp, err = c.Request(fmt.Sprintf("%s%s/tx/decode", apiEndpoint, c.Parameters.Network), http.MethodPost, postData); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &txInfo)
	return
}

// DownloadReceipt this endpoint downloads a transaction receipt (PDF)
// The contents will be returned in plain-text and need to be converted to a file.pdf
//
// For more information: https://developers.whatsonchain.com/#download-receipt
func (c *Client) DownloadReceipt(hash string) (pdfRawContent string, err error) {

	// https://<network>.whatsonchain.com/receipt/<hash>
	// todo: this endpoint does not follow the convention of the WOC API v1
	pdfRawContent, err = c.Request(fmt.Sprintf("https://%s.whatsonchain.com/receipt/%s", c.Parameters.Network, hash), http.MethodGet, nil)

	return
}
