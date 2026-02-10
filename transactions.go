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
// For more information: https://docs.whatsonchain.com/#get-by-tx-hash
func (c *Client) GetTxByHash(ctx context.Context, hash string) (*TxInfo, error) {
	url := c.buildURL("/tx/hash/%s", hash)
	return requestAndUnmarshal[TxInfo](ctx, c, url, http.MethodGet, nil, ErrTransactionNotFound)
}

// BulkTransactionDetails this fetches details for multiple transactions in single request
// Max 20 transactions per request
//
// For more information: https://docs.whatsonchain.com/#bulk-transaction-details
func (c *Client) BulkTransactionDetails(ctx context.Context, hashes *TxHashes) (TxList, error) {
	if len(hashes.TxIDs) > MaxTransactionsUTXO {
		return nil, fmt.Errorf("%w: %d UTXOs requested, max is %d", ErrMaxUTXOsExceeded, len(hashes.TxIDs), MaxTransactionsUTXO)
	}

	postData, err := json.Marshal(hashes)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/txs")
	return requestAndUnmarshalSlice[*TxInfo](ctx, c, url, http.MethodPost, postData, nil)
}

// BulkTransactionDetailsProcessor will get the details for ALL transactions in batches
// Processes 20 transactions per request
// See: BulkTransactionDetails()
func (c *Client) BulkTransactionDetailsProcessor(ctx context.Context, hashes *TxHashes) (txList TxList, err error) {
	// Break up the transactions into batches
	chunkSize := MaxTransactionsUTXO
	numBatches := (len(hashes.TxIDs) + chunkSize - 1) / chunkSize
	batches := make([][]string, 0, numBatches)

	for i := 0; i < len(hashes.TxIDs); i += chunkSize {
		end := i + chunkSize

		if end > len(hashes.TxIDs) {
			end = len(hashes.TxIDs)
		}

		batches = append(batches, hashes.TxIDs[i:end])
	}

	// Set up rate limiting with a ticker
	ticker := time.NewTicker(time.Second / time.Duration(c.RateLimit()))
	defer ticker.Stop()

	txHashes := &TxHashes{}

	// Loop Batches - and get each batch (multiple batches of MaxTransactionsUTXO)
	for _, batch := range batches {
		// Check for context cancellation before processing
		select {
		case <-ctx.Done():
			return txList, ctx.Err()
		default:
		}

		// Wait for rate limit tick
		select {
		case <-ctx.Done():
			return txList, ctx.Err()
		case <-ticker.C:
		}

		// Reuse the TxHashes struct
		txHashes.TxIDs = txHashes.TxIDs[:0]
		txHashes.TxIDs = append(txHashes.TxIDs, batch...)

		// Get the tx details (max of MaxTransactionsUTXO)
		var returnedList TxList
		if returnedList, err = c.BulkTransactionDetails(
			ctx, txHashes,
		); err != nil {
			return txList, err
		}

		// Add to the list
		txList = append(txList, returnedList...)
	}

	return txList, err
}

// GetMerkleProof retrieves the merkle proof for a transaction.
//
// Deprecated: GetMerkleProof uses a non-TSC proof endpoint that is no longer in the API.
// Use GetMerkleProofTSC instead.
//
// For more information: https://docs.whatsonchain.com/#get-merkle-proof
func (c *Client) GetMerkleProof(ctx context.Context, hash string) (MerkleResults, error) {
	url := c.buildURL("/tx/%s/proof", hash)
	return requestAndUnmarshalSlice[*MerkleInfo](ctx, c, url, http.MethodGet, nil, ErrTransactionNotFound)
}

// GetMerkleProofTSC this endpoint returns TSC compliant proof to a confirmed transaction
//
// For more information: TODO! No link today
func (c *Client) GetMerkleProofTSC(ctx context.Context, hash string) (MerkleTSCResults, error) {
	url := c.buildURL("/tx/%s/proof/tsc", hash)
	return requestAndUnmarshalSlice[*MerkleTSCInfo](ctx, c, url, http.MethodGet, nil, ErrTransactionNotFound)
}

// GetRawTransactionData this endpoint returns raw hex for the transaction with given hash
//
// For more information: https://docs.whatsonchain.com/#get-raw-transaction-data
func (c *Client) GetRawTransactionData(ctx context.Context, hash string) (string, error) {
	url := c.buildURL("/tx/%s/hex", hash)
	return requestString(ctx, c, url)
}

// BulkRawTransactionData this fetches raw hex data for multiple
// transactions in single request
// Max 20 transactions per request
//
// For more information: https://docs.whatsonchain.com/#bulk-raw-transaction-data
func (c *Client) BulkRawTransactionData(ctx context.Context, hashes *TxHashes) (TxList, error) {
	if len(hashes.TxIDs) > MaxTransactionsRaw {
		return nil, fmt.Errorf("%w: %d transactions requested, max is %d", ErrMaxRawTransactionsExceeded, len(hashes.TxIDs), MaxTransactionsRaw)
	}

	postData, err := json.Marshal(hashes)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/txs/hex")
	return requestAndUnmarshalSlice[*TxInfo](ctx, c, url, http.MethodPost, postData, nil)
}

// BulkRawTransactionDataProcessor this fetches raw hex data for
// multiple transactions in single request and handles chunking
// Max 20 transactions per request
//
// For more information: https://docs.whatsonchain.com/#bulk-raw-transaction-data
func (c *Client) BulkRawTransactionDataProcessor(ctx context.Context, hashes *TxHashes) (txList TxList, err error) {
	// Break up the transactions into batches
	chunkSize := MaxTransactionsRaw
	numBatches := (len(hashes.TxIDs) + chunkSize - 1) / chunkSize
	batches := make([][]string, 0, numBatches)

	for i := 0; i < len(hashes.TxIDs); i += chunkSize {
		end := i + chunkSize

		if end > len(hashes.TxIDs) {
			end = len(hashes.TxIDs)
		}

		batches = append(batches, hashes.TxIDs[i:end])
	}

	// Set up rate limiting with a ticker
	ticker := time.NewTicker(time.Second / time.Duration(c.RateLimit()))
	defer ticker.Stop()

	txHashes := &TxHashes{}

	// Loop Batches - and get each batch (multiple batches of MaxTransactionsRaw)
	for _, batch := range batches {
		// Check for context cancellation before processing
		select {
		case <-ctx.Done():
			return txList, ctx.Err()
		default:
		}

		// Wait for rate limit tick
		select {
		case <-ctx.Done():
			return txList, ctx.Err()
		case <-ticker.C:
		}

		// Reuse the TxHashes struct
		txHashes.TxIDs = txHashes.TxIDs[:0]
		txHashes.TxIDs = append(txHashes.TxIDs, batch...)

		// Get the tx details (max of MaxTransactionsRaw)
		var returnedList TxList
		if returnedList, err = c.BulkRawTransactionData(
			ctx, txHashes,
		); err != nil {
			return txList, err
		}

		// Add to the list
		txList = append(txList, returnedList...)
	}

	return txList, err
}

// GetRawTransactionOutputData this endpoint returns raw hex for the transaction output with given hash and index
//
// For more information: https://docs.whatsonchain.com/#get-raw-transaction-output-data
func (c *Client) GetRawTransactionOutputData(ctx context.Context, hash string, vOutIndex int) (string, error) {
	url := c.buildURL("/tx/%s/out/%d/hex", hash, vOutIndex)
	return requestString(ctx, c, url)
}

// BroadcastTx will broadcast transaction using this endpoint.
// Get tx_id in response or error msg from node.
//
// For more information: https://docs.whatsonchain.com/#broadcast-transaction
func (c *Client) BroadcastTx(ctx context.Context, txHex string) (txID string, err error) {
	// Start the post data
	postData, err := json.Marshal(map[string]string{"txhex": txHex})
	if err != nil {
		return "", err
	}

	// https://api.whatsonchain.com/v1/bsv/<network>/tx/raw
	var resp []byte
	var statusCode int
	if resp, statusCode, err = c.request(
		ctx,
		c.buildURL("/tx/raw"),
		http.MethodPost, postData,
	); err != nil {
		return "", fmt.Errorf("%w: %w", ErrBroadcastFailed, err)
	}

	// Check for non-OK status codes using the standard helper
	if err = checkStatusCode(statusCode, resp); err != nil {
		return "", fmt.Errorf("%w: %w", ErrBroadcastFailed, err)
	}

	// Remove quotes or spaces from successful response
	txID = strings.TrimSpace(strings.ReplaceAll(string(resp), `"`, ""))
	return txID, nil
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
// Deprecated: BulkBroadcastTx uses a bulk broadcast endpoint that is no longer in the API.
// Use BroadcastTx for individual transactions instead.
//
// submitted transactions, eg 'QUEUED' or 'PROCESSED', with response data from node. You can poll the provided unique
// url until all transactions are marked as 'PROCESSED'. Progress of the transactions are tracked on this unique url
// for up to 5 hours.
//
// For more information: https://docs.whatsonchain.com/#bulk-broadcast
func (c *Client) BulkBroadcastTx(ctx context.Context, rawTxs []string,
	feedback bool,
) (response *BulkBroadcastResponse, err error) {
	// Set a max (from WOC)
	if len(rawTxs) > MaxBroadcastTransactions {
		err = fmt.Errorf("%w: %d transactions, max is %d", ErrMaxTransactionsExceeded, len(rawTxs), MaxBroadcastTransactions)
		return response, err
	}

	// Set a total max
	if len(strings.Join(rawTxs[:], ",")) > MaxCombinedTransactionSize {
		err = fmt.Errorf("%w: payload size %.0f bytes, max is %.0f bytes", ErrMaxPayloadSizeExceeded, float64(len(strings.Join(rawTxs[:], ","))), MaxCombinedTransactionSize)
		return response, err
	}

	// Check size of each tx
	for _, tx := range rawTxs {
		if len(tx) > MaxSingleTransactionSize {
			err = fmt.Errorf("%w: transaction size %d bytes, max is %d bytes", ErrMaxTransactionSizeExceeded, len(tx), MaxSingleTransactionSize)
			return response, err
		}
	}

	// Start the post data
	var postData []byte
	if postData, err = json.Marshal(rawTxs); err != nil {
		return nil, err
	}

	var resp []byte
	var statusCode int

	// https://api.whatsonchain.com/v1/bsv/tx/broadcast?feedback=<feedback>
	if resp, statusCode, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/tx/broadcast?feedback=%t", apiEndpointBase, c.Chain(), c.Network(), feedback),
		http.MethodPost, postData,
	); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBroadcastFailed, err)
	}

	// Check for non-OK status codes using the standard helper
	if err = checkStatusCode(statusCode, resp); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBroadcastFailed, err)
	}

	response = &BulkBroadcastResponse{Feedback: feedback}
	if feedback {
		if err = json.Unmarshal(resp, response); err != nil {
			return response, err
		}
	}

	return response, nil
}

// DecodeTransaction this endpoint decodes raw transaction
//
// For more information: https://docs.whatsonchain.com/#decode-transaction
func (c *Client) DecodeTransaction(ctx context.Context, txHex string) (*TxInfo, error) {
	postData, err := json.Marshal(map[string]string{"txhex": txHex})
	if err != nil {
		return nil, err
	}
	url := c.buildURL("/tx/decode")
	return requestAndUnmarshal[TxInfo](ctx, c, url, http.MethodPost, postData, ErrTransactionNotFound)
}

// DownloadReceipt this endpoint downloads a transaction receipt (PDF)
// The contents will be returned in plain-text and need to be converted to a file.pdf
//
// For more information: https://docs.whatsonchain.com/#download-receipt
func (c *Client) DownloadReceipt(ctx context.Context, hash string) (string, error) {
	// This endpoint does not follow the convention of the WOC API v1
	url := fmt.Sprintf("https://%s.whatsonchain.com/receipt/%s", c.Network(), hash)
	return requestString(ctx, c, url)
}

// GetTransactionPropagationStatus this endpoint retrieves transaction propagation status (BSV only)
//
// For more information: https://docs.whatsonchain.com/#get-tx-propagation
func (c *Client) GetTransactionPropagationStatus(ctx context.Context, hash string) (*PropagationStatus, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/tx/hash/%s/propagation", hash)
	return requestAndUnmarshal[PropagationStatus](ctx, c, url, http.MethodGet, nil, ErrTransactionNotFound)
}

// BulkTransactionStatus this endpoint fetches status for multiple transactions in single request
// Max 20 transactions per request
//
// For more information: https://docs.whatsonchain.com/#bulk-transaction-status
func (c *Client) BulkTransactionStatus(ctx context.Context, hashes *TxHashes) (TxStatusList, error) {
	if len(hashes.TxIDs) > MaxTransactionsUTXO {
		return nil, fmt.Errorf("%w: %d transactions requested, max is %d", ErrMaxUTXOsExceeded, len(hashes.TxIDs), MaxTransactionsUTXO)
	}

	postData, err := json.Marshal(hashes)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/txs/status")
	return requestAndUnmarshalSlice[*TxStatus](ctx, c, url, http.MethodPost, postData, nil)
}

// GetTransactionAsBinary this endpoint retrieves transaction data as binary
//
// For more information: https://docs.whatsonchain.com/#get-tx-binary
func (c *Client) GetTransactionAsBinary(ctx context.Context, hash string) ([]byte, error) {
	url := c.buildURL("/tx/%s/bin", hash)
	resp, err := requestString(ctx, c, url)
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, ErrTransactionNotFound
	}
	return []byte(resp), nil
}

// BulkRawTransactionOutputData this endpoint fetches raw output data for multiple transactions in single request
// Max 20 transactions per request
//
// For more information: https://docs.whatsonchain.com/#bulk-raw-tx-output
func (c *Client) BulkRawTransactionOutputData(ctx context.Context, request *BulkRawOutputRequest) ([]*BulkRawOutputResponse, error) {
	if len(request.TxIDs) > MaxTransactionsUTXO {
		return nil, fmt.Errorf("%w: %d transactions requested, max is %d", ErrMaxUTXOsExceeded, len(request.TxIDs), MaxTransactionsUTXO)
	}

	postData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/txs/vouts/hex")
	return requestAndUnmarshalSlice[*BulkRawOutputResponse](ctx, c, url, http.MethodPost, postData, nil)
}

// GetUnconfirmedSpentOutput this endpoint retrieves unconfirmed spent transaction output details
//
// For more information: https://docs.whatsonchain.com/#get-unconfirmed-spent
func (c *Client) GetUnconfirmedSpentOutput(ctx context.Context, txHash string, index int) (*SpentOutput, error) {
	url := c.buildURL("/tx/%s/%d/unconfirmed/spent", txHash, index)
	return requestAndUnmarshal[SpentOutput](ctx, c, url, http.MethodGet, nil, ErrTransactionNotFound)
}

// GetConfirmedSpentOutput this endpoint retrieves confirmed spent transaction output details
//
// For more information: https://docs.whatsonchain.com/#get-confirmed-spent
func (c *Client) GetConfirmedSpentOutput(ctx context.Context, txHash string, index int) (*SpentOutput, error) {
	url := c.buildURL("/tx/%s/%d/confirmed/spent", txHash, index)
	return requestAndUnmarshal[SpentOutput](ctx, c, url, http.MethodGet, nil, ErrTransactionNotFound)
}

// GetSpentOutput this endpoint retrieves spent transaction output details (both confirmed and unconfirmed)
//
// For more information: https://docs.whatsonchain.com/#get-spent-output
func (c *Client) GetSpentOutput(ctx context.Context, txHash string, index int) (*SpentOutput, error) {
	url := c.buildURL("/tx/%s/%d/spent", txHash, index)
	return requestAndUnmarshal[SpentOutput](ctx, c, url, http.MethodGet, nil, ErrTransactionNotFound)
}

// BulkSpentOutputs this endpoint retrieves spent output details for multiple UTXOs
// Max 20 UTXOs per request
//
// For more information: https://docs.whatsonchain.com/#bulk-spent-outputs
func (c *Client) BulkSpentOutputs(ctx context.Context, request *BulkSpentOutputRequest) (BulkSpentOutputResponse, error) {
	if len(request.UTXOs) > MaxTransactionsUTXO {
		return nil, fmt.Errorf("%w: %d UTXOs requested, max is %d", ErrMaxUTXOsExceeded, len(request.UTXOs), MaxTransactionsUTXO)
	}

	postData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/utxos/spent")
	return requestAndUnmarshalSlice[BulkSpentOutputResult](ctx, c, url, http.MethodPost, postData, nil)
}
