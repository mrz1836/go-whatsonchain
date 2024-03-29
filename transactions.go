package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GetTxByHash this endpoint retrieves transaction details with given transaction hash
//
// For more information: https://developers.whatsonchain.com/#get-by-tx-hash
func (c *Client) GetTxByHash(ctx context.Context, hash string) (txInfo *TxInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tx/hash/<hash>
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/tx/hash/%s", apiEndpoint, c.Network(), hash),
		http.MethodGet, nil,
	); err != nil {
		return
	}
	if len(resp) == 0 {
		return nil, ErrTransactionNotFound
	}
	err = json.Unmarshal([]byte(resp), &txInfo)
	return
}

// BulkTransactionDetails this fetches details for multiple transactions in single request
// Max 20 transactions per request
//
// For more information: https://developers.whatsonchain.com/#bulk-transaction-details
func (c *Client) BulkTransactionDetails(ctx context.Context, hashes *TxHashes) (txList TxList, err error) {

	// The max limit by WOC
	if len(hashes.TxIDs) > MaxTransactionsUTXO {
		err = fmt.Errorf(
			"max limit of utxos is %d and you sent %d",
			MaxTransactionsUTXO, len(hashes.TxIDs),
		)
		return
	}

	// Convert to JSON
	var postData []byte
	if postData, err = json.Marshal(hashes); err != nil {
		return
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/txs
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/txs", apiEndpoint, c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return
	}

	if len(resp) > 0 {
		err = json.Unmarshal([]byte(resp), &txList)
	}
	return
}

// BulkTransactionDetailsProcessor will get the details for ALL transactions in batches
// Processes 20 transactions per request
// See: BulkTransactionDetails()
func (c *Client) BulkTransactionDetailsProcessor(ctx context.Context, hashes *TxHashes) (txList TxList, err error) {

	// Break up the transactions into batches
	var batches [][]string
	chunkSize := MaxTransactionsUTXO

	for i := 0; i < len(hashes.TxIDs); i += chunkSize {
		end := i + chunkSize

		if end > len(hashes.TxIDs) {
			end = len(hashes.TxIDs)
		}

		batches = append(batches, hashes.TxIDs[i:end])
	}

	var currentRateLimit int

	// Loop Batches - and get each batch (multiple batches of MaxTransactionsUTXO)
	for _, batch := range batches {

		txHashes := new(TxHashes)

		// Loop the batch (max MaxTransactionsUTXO)
		txHashes.TxIDs = append(txHashes.TxIDs, batch...)

		// Get the tx details (max of MaxTransactionsUTXO)
		var returnedList TxList
		if returnedList, err = c.BulkTransactionDetails(
			ctx, txHashes,
		); err != nil {
			return
		}

		// Add to the list
		txList = append(txList, returnedList...)

		// Accumulate / sleep to prevent rate limiting
		currentRateLimit++
		if currentRateLimit >= c.RateLimit() {
			time.Sleep(1 * time.Second)
			currentRateLimit = 0
		}
	}

	return
}

// GetMerkleProof this endpoint returns merkle branch to a confirmed transaction
//
// For more information: https://developers.whatsonchain.com/#get-merkle-proof
func (c *Client) GetMerkleProof(ctx context.Context, hash string) (merkleResults MerkleResults, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tx/<hash>/proof
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/tx/%s/proof", apiEndpoint, c.Network(), hash),
		http.MethodGet, nil,
	); err != nil {
		return
	}
	if len(resp) == 0 {
		return nil, ErrTransactionNotFound
	}
	err = json.Unmarshal([]byte(resp), &merkleResults)
	return
}

// GetMerkleProofTSC this endpoint returns TSC compliant proof to a confirmed transaction
//
// For more information: TODO! No link today
func (c *Client) GetMerkleProofTSC(ctx context.Context, hash string) (merkleResults MerkleTSCResults, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tx/<hash>/proof/tsc
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/tx/%s/proof/tsc", apiEndpoint, c.Network(), hash),
		http.MethodGet, nil,
	); err != nil {
		return
	}
	if len(resp) == 0 {
		return nil, ErrTransactionNotFound
	}
	err = json.Unmarshal([]byte(resp), &merkleResults)
	return
}

// GetRawTransactionData this endpoint returns raw hex for the transaction with given hash
//
// For more information: https://developers.whatsonchain.com/#get-raw-transaction-data
func (c *Client) GetRawTransactionData(ctx context.Context, hash string) (string, error) {

	// https://api.whatsonchain.com/v1/bsv/<network>/tx/<hash>/hex
	return c.request(
		ctx,
		fmt.Sprintf("%s%s/tx/%s/hex", apiEndpoint, c.Network(), hash),
		http.MethodGet, nil,
	)
}

// BulkRawTransactionData this fetches raw hex data for multiple
// transactions in single request
// Max 20 transactions per request
//
// For more information: https://developers.whatsonchain.com/#bulk-raw-transaction-data
func (c *Client) BulkRawTransactionData(ctx context.Context, hashes *TxHashes) (txList TxList, err error) {

	// The max limit by WOC
	if len(hashes.TxIDs) > MaxTransactionsRaw {
		err = fmt.Errorf(
			"max limit of transactions is %d and you sent %d",
			MaxTransactionsRaw, len(hashes.TxIDs),
		)
		return
	}

	// Convert to JSON
	var postData []byte
	if postData, err = json.Marshal(hashes); err != nil {
		return
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/txs/hex

	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/txs/hex", apiEndpoint, c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return
	}

	if len(resp) > 0 {
		err = json.Unmarshal([]byte(resp), &txList)
	}
	return
}

// BulkRawTransactionDataProcessor this fetches raw hex data for
// multiple transactions in single request and handles chunking
// Max 20 transactions per request
//
// For more information: https://developers.whatsonchain.com/#bulk-raw-transaction-data
func (c *Client) BulkRawTransactionDataProcessor(ctx context.Context, hashes *TxHashes) (txList TxList, err error) {

	// Break up the transactions into batches
	var batches [][]string
	chunkSize := MaxTransactionsRaw

	for i := 0; i < len(hashes.TxIDs); i += chunkSize {
		end := i + chunkSize

		if end > len(hashes.TxIDs) {
			end = len(hashes.TxIDs)
		}

		batches = append(batches, hashes.TxIDs[i:end])
	}

	var currentRateLimit int

	// Loop Batches - and get each batch (multiple batches of MaxTransactionsRaw)
	for _, batch := range batches {

		txHashes := new(TxHashes)

		// Loop the batch (max MaxTransactionsRaw)
		txHashes.TxIDs = append(txHashes.TxIDs, batch...)

		// Get the tx details (max of MaxTransactionsUTXO)
		var returnedList TxList
		if returnedList, err = c.BulkRawTransactionData(
			ctx, txHashes,
		); err != nil {
			return
		}

		// Add to the list
		txList = append(txList, returnedList...)

		// Accumulate / sleep to prevent rate limiting
		currentRateLimit++
		if currentRateLimit >= c.RateLimit() {
			time.Sleep(1 * time.Second)
			currentRateLimit = 0
		}
	}

	return
}

// GetRawTransactionOutputData this endpoint returns raw hex for the transaction output with given hash and index
//
// For more information: https://developers.whatsonchain.com/#get-raw-transaction-output-data
func (c *Client) GetRawTransactionOutputData(ctx context.Context, hash string, vOutIndex int) (string, error) {

	// https://api.whatsonchain.com/v1/bsv/<network>/tx/<hash>/out/<index>/hex
	return c.request(
		ctx,
		fmt.Sprintf("%s%s/tx/%s/out/%d/hex", apiEndpoint, c.Network(), hash, vOutIndex),
		http.MethodGet, nil,
	)
}

// BroadcastTx will broadcast transaction using this endpoint.
// Get tx_id in response or error msg from node.
//
// For more information: https://developers.whatsonchain.com/#broadcast-transaction
func (c *Client) BroadcastTx(ctx context.Context, txHex string) (txID string, err error) {

	// Start the post data
	postData := []byte(fmt.Sprintf(`{"txhex":"%s"}`, txHex))

	// https://api.whatsonchain.com/v1/bsv/<network>/tx/raw
	if txID, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/tx/raw", apiEndpoint, c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return
	}

	// Got an error
	if c.lastRequest.StatusCode > http.StatusOK {
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
//
//	Size per transaction should be less than 100KB
//	Overall payload per request should be less than 10MB
//	Max 100 transactions per request
//	Only available for mainnet
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
func (c *Client) BulkBroadcastTx(ctx context.Context, rawTxs []string,
	feedback bool) (response *BulkBroadcastResponse, err error) {

	// Set a max (from WOC)
	if len(rawTxs) > MaxBroadcastTransactions {
		err = fmt.Errorf("max transactions are %d", MaxBroadcastTransactions)
		return
	}

	// Set a total max
	if len(strings.Join(rawTxs[:], ",")) > MaxCombinedTransactionSize {
		err = fmt.Errorf("max overall payload of 10MB (%f bytes)", MaxCombinedTransactionSize)
		return
	}

	// Check size of each tx
	for _, tx := range rawTxs {
		if len(tx) > MaxSingleTransactionSize {
			err = fmt.Errorf("max tx size of 100kb (%d bytes)", MaxSingleTransactionSize)
			return
		}
	}

	// Start the post data
	var postData []byte
	if postData, err = json.Marshal(rawTxs); err != nil {
		return nil, err
	}

	var resp string

	// https://api.whatsonchain.com/v1/bsv/tx/broadcast?feedback=<feedback>
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%stx/broadcast?feedback=%t", apiEndpoint, feedback),
		http.MethodPost, postData,
	); err != nil {
		return
	}

	response = &BulkBroadcastResponse{Feedback: feedback}
	if feedback {
		if err = json.Unmarshal([]byte(resp), response); err != nil {
			return
		}
	}

	// Got an error
	if c.lastRequest.StatusCode > http.StatusOK {
		err = fmt.Errorf("error broadcasting: %s", resp)
	}

	return
}

// DecodeTransaction this endpoint decodes raw transaction
//
// For more information: https://developers.whatsonchain.com/#decode-transaction
func (c *Client) DecodeTransaction(ctx context.Context, txHex string) (txInfo *TxInfo, err error) {

	// Start the post data
	postData := []byte(fmt.Sprintf(`{"txhex":"%s"}`, txHex))

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tx/decode
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/tx/decode", apiEndpoint, c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return
	}
	if len(resp) == 0 {
		return nil, ErrTransactionNotFound
	}
	err = json.Unmarshal([]byte(resp), &txInfo)
	return
}

// DownloadReceipt this endpoint downloads a transaction receipt (PDF)
// The contents will be returned in plain-text and need to be converted to a file.pdf
//
// For more information: https://developers.whatsonchain.com/#download-receipt
func (c *Client) DownloadReceipt(ctx context.Context, hash string) (string, error) {

	// https://<network>.whatsonchain.com/receipt/<hash>
	// todo: this endpoint does not follow the convention of the WOC API v1
	return c.request(
		ctx,
		fmt.Sprintf("https://%s.whatsonchain.com/receipt/%s", c.Network(), hash),
		http.MethodGet, nil,
	)
}
