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

// GetTxsByHashes this endpoint retrieves transaction details with given transaction hashes
//
// For more information: (undocumented) // todo: update url when documented
func (c *Client) GetTxsByHashes(hashes *TxHashes) (txList TxList, err error) {

	// Max limit by WOC
	if len(hashes.TxIDs) == 0 {
		err = fmt.Errorf("missing hashes")
		return
	}

	// Testing turning off the limit
	// todo: turn back on when limit is known
	/*else if len(hashes.TxIDs) > MaxTransactionsUTXO {
		err = fmt.Errorf("max limit of utxos is %d and you sent %d", MaxTransactionsUTXO, len(hashes.TxIDs))
		return
	}*/

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
func (c *Client) GetMerkleProof(hash string) (merkleInfo *MerkleInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tx/<hash>/merkleproof
	if resp, err = c.Request(fmt.Sprintf("%s%s/tx/%s/merkleproof", apiEndpoint, c.Parameters.Network, hash), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &merkleInfo)
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
