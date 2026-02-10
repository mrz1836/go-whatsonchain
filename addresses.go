package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	netURL "net/url"
	"time"
)

// AddressInfo this endpoint retrieves various address info.
//
// For more information: https://docs.whatsonchain.com/#address
func (c *Client) AddressInfo(ctx context.Context, address string) (*AddressInfo, error) {
	url := c.buildURL("/address/%s/info", address)
	return requestAndUnmarshal[AddressInfo](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// AddressBalance retrieves the balance for an address.
//
// Deprecated: AddressBalance uses a combined balance endpoint that is no longer in the API.
// Use AddressConfirmedBalance and AddressUnconfirmedBalance instead.
//
// For more information: https://docs.whatsonchain.com/#get-balance
func (c *Client) AddressBalance(ctx context.Context, address string) (*AddressBalance, error) {
	url := c.buildURL("/address/%s/balance", address)
	return requestAndUnmarshal[AddressBalance](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// AddressHistory retrieves the transaction history for an address.
//
// Deprecated: AddressHistory uses a combined history endpoint that is no longer in the API.
// Use AddressConfirmedHistory and AddressUnconfirmedHistory instead.
//
// For more information: https://docs.whatsonchain.com/#get-history
func (c *Client) AddressHistory(ctx context.Context, address string) (AddressHistory, error) {
	url := c.buildURL("/address/%s/history", address)
	return requestAndUnmarshalSlice[*HistoryRecord](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// AddressUnspentTransactions this endpoint retrieves ordered list of UTXOs.
//
// For more information: https://docs.whatsonchain.com/#get-unspent-transactions
func (c *Client) AddressUnspentTransactions(ctx context.Context, address string) (AddressHistory, error) {
	url := c.buildURL("/address/%s/unspent/all", address)
	resp, err := requestAndUnmarshal[addressUnspentAllResponse](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, resp.Error)
	}
	return resp.Result, nil
}

// AddressUnspentTransactionDetails this endpoint retrieves transaction details for a given address
// Use max transactions to filter if there are more UTXOs returned than needed by the user
//
// For more information: (custom request for this go package)
func (c *Client) AddressUnspentTransactionDetails(ctx context.Context, address string, maxTransactions int) (history AddressHistory, err error) {
	// Get the address UTXO history
	var utxos AddressHistory
	if utxos, err = c.AddressUnspentTransactions(ctx, address); err != nil {
		return history, err
	} else if len(utxos) == 0 {
		return history, err
	}

	// Do we have a "custom max" amount?
	if maxTransactions > 0 {
		total := len(utxos)
		if total > maxTransactions {
			utxos = utxos[:total-(total-maxTransactions)]
		}
	}

	// Break up the UTXOs into batches
	batches := chunkSlice(utxos, MaxTransactionsUTXO)

	// Loop Batches - and get each batch (multiple batches of MaxTransactionsUTXO)
	txHashes := &TxHashes{}
	for _, batch := range batches {
		// Check for context cancellation before processing
		select {
		case <-ctx.Done():
			return history, ctx.Err()
		default:
		}

		// Reuse the TxHashes struct
		txHashes.TxIDs = txHashes.TxIDs[:0]

		// Loop the batch (max MaxTransactionsUTXO)
		for _, utxo := range batch {
			txHashes.TxIDs = append(txHashes.TxIDs, utxo.TxHash)
			history = append(history, utxo)
		}

		// Get the tx details (max of MaxTransactionsUTXO)
		var txList TxList
		if txList, err = c.BulkTransactionDetails(ctx, txHashes); err != nil {
			return history, err
		}

		// Build a map from tx hash to tx info for O(1) lookup
		txInfoMap := make(map[string]*TxInfo, len(txList))
		for _, tx := range txList {
			txInfoMap[tx.TxID] = tx
		}
		// Attach tx info to all UTXOs in this batch with matching tx hash
		for _, utxo := range batch {
			if info, ok := txInfoMap[utxo.TxHash]; ok {
				utxo.Info = info
			}
		}
	}

	return history, err
}

// DownloadStatement this endpoint downloads an address statement (PDF)
// The contents will be returned in plain-text and need to be converted to a file.pdf
//
// For more information: https://docs.whatsonchain.com/#download-statement
func (c *Client) DownloadStatement(ctx context.Context, address string) (string, error) {
	// This endpoint does not follow the convention of the WOC API v1
	url := fmt.Sprintf("https://%s.whatsonchain.com/statement/%s", c.Network(), netURL.PathEscape(address))
	return requestString(ctx, c, url)
}

// bulkRequest is the common parts of the bulk requests
func bulkRequest(list *AddressList) ([]byte, error) {
	if list == nil {
		return nil, ErrMissingRequest
	}

	// The max limit by WOC
	if len(list.Addresses) > MaxAddressesForLookup {
		return nil, fmt.Errorf("%w: %d addresses requested, max is %d", ErrMaxAddressesExceeded, len(list.Addresses), MaxAddressesForLookup)
	}

	// Convert to JSON
	return json.Marshal(list)
}

// BulkBalance retrieves balances for multiple addresses.
// Max of 20 addresses at a time.
//
// Deprecated: BulkBalance uses a combined balance endpoint that is no longer in the API.
// Use BulkAddressConfirmedBalance and BulkAddressUnconfirmedBalance instead.
//
// For more information: https://docs.whatsonchain.com/#bulk-balance
func (c *Client) BulkBalance(ctx context.Context, list *AddressList) (AddressBalances, error) {
	postData, err := bulkRequest(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/addresses/balance")
	return requestAndUnmarshalSlice[*AddressBalanceRecord](ctx, c, url, http.MethodPost, postData, ErrAddressNotFound)
}

// BulkUnspentTransactionsProcessor processes bulk unspent transactions in batches.
// Max of 20 addresses at a time.
//
// Deprecated: BulkUnspentTransactionsProcessor wraps BulkUnspentTransactions which uses an endpoint
// no longer in the API. Use BulkAddressConfirmedUTXOs and BulkAddressUnconfirmedUTXOs instead.
//
// For more information: https://docs.whatsonchain.com/#bulk-unspent-transactions
func (c *Client) BulkUnspentTransactionsProcessor(ctx context.Context, list *AddressList) (responseList BulkUnspentResponse, err error) {
	if list == nil {
		return nil, ErrMissingRequest
	}
	batches := chunkSlice(list.Addresses, MaxTransactionsUTXO)
	// Set up rate limiting with a ticker
	ticker := time.NewTicker(time.Second / time.Duration(c.RateLimit()))
	defer ticker.Stop()

	addressList := &AddressList{}
	for _, batch := range batches {
		// Check for context cancellation before processing
		select {
		case <-ctx.Done():
			return responseList, ctx.Err()
		default:
		}

		// Wait for rate limit tick
		select {
		case <-ctx.Done():
			return responseList, ctx.Err()
		case <-ticker.C:
		}

		addressList.Addresses = addressList.Addresses[:0]
		addressList.Addresses = append(addressList.Addresses, batch...)
		var returnedList BulkUnspentResponse
		if returnedList, err = c.BulkUnspentTransactions(ctx, addressList); err != nil {
			return responseList, err
		}
		responseList = append(responseList, returnedList...)
	}
	return responseList, err
}

// BulkUnspentTransactions retrieves unspent transactions for multiple addresses.
// Max of 20 addresses at a time.
//
// Deprecated: BulkUnspentTransactions uses a combined unspent endpoint that is no longer in the API.
// Use BulkAddressConfirmedUTXOs and BulkAddressUnconfirmedUTXOs instead.
//
// For more information: https://docs.whatsonchain.com/#bulk-unspent-transactions
func (c *Client) BulkUnspentTransactions(ctx context.Context, list *AddressList) (BulkUnspentResponse, error) {
	postData, err := bulkRequest(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/addresses/unspent/all")
	return requestAndUnmarshalSlice[*BulkResponseRecord](ctx, c, url, http.MethodPost, postData, ErrAddressNotFound)
}

// AddressUnconfirmedUTXOs retrieves unconfirmed UTXOs for an address
//
// For more information: https://docs.whatsonchain.com/#get-unconfirmed-utxos
func (c *Client) AddressUnconfirmedUTXOs(ctx context.Context, address string) (AddressHistory, error) {
	url := c.buildURL("/address/%s/unconfirmed/unspent", address)
	return requestAndUnmarshalSlice[*HistoryRecord](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// BulkAddressUnconfirmedUTXOs retrieves unconfirmed UTXOs for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/#bulk-unconfirmed-utxos
func (c *Client) BulkAddressUnconfirmedUTXOs(ctx context.Context, list *AddressList) (BulkUnspentResponse, error) {
	postData, err := bulkRequest(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/addresses/unconfirmed/unspent")
	return requestAndUnmarshalSlice[*BulkResponseRecord](ctx, c, url, http.MethodPost, postData, ErrAddressNotFound)
}

// AddressConfirmedUTXOs retrieves confirmed UTXOs for an address
//
// For more information: https://docs.whatsonchain.com/#get-confirmed-utxos
func (c *Client) AddressConfirmedUTXOs(ctx context.Context, address string) (AddressHistory, error) {
	url := c.buildURL("/address/%s/confirmed/unspent", address)
	return requestAndUnmarshalSlice[*HistoryRecord](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// BulkAddressConfirmedUTXOs retrieves confirmed UTXOs for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/#bulk-confirmed-utxos
func (c *Client) BulkAddressConfirmedUTXOs(ctx context.Context, list *AddressList) (BulkUnspentResponse, error) {
	postData, err := bulkRequest(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/addresses/confirmed/unspent")
	return requestAndUnmarshalSlice[*BulkResponseRecord](ctx, c, url, http.MethodPost, postData, ErrAddressNotFound)
}

// AddressUsed retrieves whether an address has been used
//
// For more information: https://docs.whatsonchain.com/api/address#get-address-usage
func (c *Client) AddressUsed(ctx context.Context, address string) (*AddressUsed, error) {
	url := c.buildURL("/address/%s/used", address)
	return requestAndUnmarshal[AddressUsed](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// AddressScripts retrieves associated scripthashes for an address
//
// For more information: https://docs.whatsonchain.com/api/address#get-associated-scripthashes
func (c *Client) AddressScripts(ctx context.Context, address string) (*AddressScripts, error) {
	url := c.buildURL("/address/%s/scripts", address)
	return requestAndUnmarshal[AddressScripts](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// AddressUnconfirmedBalance retrieves unconfirmed balance for an address
//
// For more information: https://docs.whatsonchain.com/api/address#get-unconfirmed-balance
func (c *Client) AddressUnconfirmedBalance(ctx context.Context, address string) (*AddressUnconfirmedBalance, error) {
	url := c.buildURL("/address/%s/unconfirmed/balance", address)
	return requestAndUnmarshal[AddressUnconfirmedBalance](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// AddressConfirmedBalance retrieves confirmed balance for an address
//
// For more information: https://docs.whatsonchain.com/api/address#get-confirmed-balance
func (c *Client) AddressConfirmedBalance(ctx context.Context, address string) (*AddressConfirmedBalance, error) {
	url := c.buildURL("/address/%s/confirmed/balance", address)
	return requestAndUnmarshal[AddressConfirmedBalance](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// AddressUnconfirmedHistory retrieves unconfirmed transaction history for an address
//
// For more information: https://docs.whatsonchain.com/api/address#get-unconfirmed-history
func (c *Client) AddressUnconfirmedHistory(ctx context.Context, address string) (AddressHistory, error) {
	url := c.buildURL("/address/%s/unconfirmed/history", address)
	return requestAndUnmarshalSlice[*HistoryRecord](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// AddressConfirmedHistory retrieves confirmed transaction history for an address
//
// For more information: https://docs.whatsonchain.com/api/address#get-confirmed-history
func (c *Client) AddressConfirmedHistory(ctx context.Context, address string) (AddressHistory, error) {
	url := c.buildURL("/address/%s/confirmed/history", address)
	return requestAndUnmarshalSlice[*HistoryRecord](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// BulkAddressUnconfirmedBalance retrieves unconfirmed balances for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/api/address#bulk-unconfirmed-balance
func (c *Client) BulkAddressUnconfirmedBalance(ctx context.Context, list *AddressList) (AddressBalances, error) {
	postData, err := bulkRequest(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/addresses/unconfirmed/balance")
	return requestAndUnmarshalSlice[*AddressBalanceRecord](ctx, c, url, http.MethodPost, postData, ErrAddressNotFound)
}

// BulkAddressConfirmedBalance retrieves confirmed balances for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/api/address#bulk-confirmed-balance
func (c *Client) BulkAddressConfirmedBalance(ctx context.Context, list *AddressList) (AddressBalances, error) {
	postData, err := bulkRequest(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/addresses/confirmed/balance")
	return requestAndUnmarshalSlice[*AddressBalanceRecord](ctx, c, url, http.MethodPost, postData, ErrAddressNotFound)
}

// BulkAddressUnconfirmedHistory retrieves unconfirmed transaction history for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/api/address#bulk-unconfirmed-history
func (c *Client) BulkAddressUnconfirmedHistory(ctx context.Context, list *AddressList) (BulkAddressHistoryResponse, error) {
	postData, err := bulkRequest(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/addresses/unconfirmed/history")
	return requestAndUnmarshalSlice[*BulkAddressHistoryRecord](ctx, c, url, http.MethodPost, postData, ErrAddressNotFound)
}

// BulkAddressConfirmedHistory retrieves confirmed transaction history for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/api/address#bulk-confirmed-history
func (c *Client) BulkAddressConfirmedHistory(ctx context.Context, list *AddressList) (BulkAddressHistoryResponse, error) {
	postData, err := bulkRequest(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/addresses/confirmed/history")
	return requestAndUnmarshalSlice[*BulkAddressHistoryRecord](ctx, c, url, http.MethodPost, postData, ErrAddressNotFound)
}

// BulkAddressHistory retrieves all transaction history for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/api/address#bulk-history
func (c *Client) BulkAddressHistory(ctx context.Context, list *AddressList) (BulkAddressHistoryResponse, error) {
	postData, err := bulkRequest(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/addresses/history/all")
	return requestAndUnmarshalSlice[*BulkAddressHistoryRecord](ctx, c, url, http.MethodPost, postData, ErrAddressNotFound)
}
