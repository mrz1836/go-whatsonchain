package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AddressInfo this endpoint retrieves various address info.
//
// For more information: https://docs.whatsonchain.com/#address
func (c *Client) AddressInfo(ctx context.Context, address string) (*AddressInfo, error) {
	url := c.buildURL("/address/%s/info", address)
	return requestAndUnmarshal[AddressInfo](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// AddressBalance this endpoint retrieves confirmed and unconfirmed address balance.
//
// For more information: https://docs.whatsonchain.com/#get-balance
func (c *Client) AddressBalance(ctx context.Context, address string) (*AddressBalance, error) {
	url := c.buildURL("/address/%s/balance", address)
	return requestAndUnmarshal[AddressBalance](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
}

// AddressHistory this endpoint retrieves confirmed and unconfirmed address transactions.
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
	return requestAndUnmarshalSlice[*HistoryRecord](ctx, c, url, http.MethodGet, nil, ErrAddressNotFound)
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
	var batches []AddressHistory
	chunkSize := MaxTransactionsUTXO

	for i := 0; i < len(utxos); i += chunkSize {
		end := i + chunkSize

		if end > len(utxos) {
			end = len(utxos)
		}

		batches = append(batches, utxos[i:end])
	}

	// todo: use channels/wait group to fire all requests at the same time (rate limiting)

	// Loop Batches - and get each batch (multiple batches of MaxTransactionsUTXO)
	for _, batch := range batches {

		txHashes := new(TxHashes)

		// Loop the batch (max MaxTransactionsUTXO)
		for _, utxo := range batch {

			// Append to the list to send and return
			txHashes.TxIDs = append(txHashes.TxIDs, utxo.TxHash)
			history = append(history, utxo)
		}

		// Get the tx details (max of MaxTransactionsUTXO)
		var txList TxList
		if txList, err = c.BulkTransactionDetails(ctx, txHashes); err != nil {
			return history, err
		}

		// Add to the history list
		for index, tx := range txList {
			for _, utxo := range history {
				if utxo.TxHash == tx.TxID {
					utxo.Info = txList[index]
					continue
				}
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
	url := fmt.Sprintf("https://%s.whatsonchain.com/statement/%s", c.Network(), address)
	return requestString(ctx, c, url)
}

// bulkRequest is the common parts of the bulk requests
func bulkRequest(list *AddressList) ([]byte, error) {
	// The max limit by WOC
	if len(list.Addresses) > MaxAddressesForLookup {
		return nil, fmt.Errorf("%w: %d addresses requested, max is %d", ErrMaxAddressesExceeded, len(list.Addresses), MaxAddressesForLookup)
	}

	// Convert to JSON
	return json.Marshal(list)
}

// BulkBalance this endpoint retrieves confirmed and unconfirmed address balances
// Max of 20 addresses at a time
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

// BulkUnspentTransactionsProcessor will fetch UTXOs for multiple addresses in a single request while automatically batching
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/#bulk-unspent-transactions
func (c *Client) BulkUnspentTransactionsProcessor(ctx context.Context, list *AddressList) (responseList BulkUnspentResponse, err error) {
	var batches [][]string
	chunkSize := MaxTransactionsUTXO
	for i := 0; i < len(list.Addresses); i += chunkSize {
		end := i + chunkSize
		if end > len(list.Addresses) {
			end = len(list.Addresses)
		}
		batches = append(batches, list.Addresses[i:end])
	}
	var currentRateLimit int
	for _, batch := range batches {
		addressList := new(AddressList)
		addressList.Addresses = append(addressList.Addresses, batch...)
		var returnedList BulkUnspentResponse
		if returnedList, err = c.BulkUnspentTransactions(ctx, addressList); err != nil {
			return responseList, err
		}
		responseList = append(responseList, returnedList...)
		currentRateLimit++
		if currentRateLimit >= c.RateLimit() {
			time.Sleep(1 * time.Second)
			currentRateLimit = 0
		}
	}
	return responseList, err
}

// BulkUnspentTransactions will fetch UTXOs for multiple addresses in a single request
// Max of 20 addresses at a time
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
