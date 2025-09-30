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
// For more information: https://developers.whatsonchain.com/#address
func (c *Client) AddressInfo(ctx context.Context, address string) (addressInfo *AddressInfo, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/address/<address>/info
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/info", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet,
		nil,
	); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &addressInfo)

	return addressInfo, err
}

// AddressBalance this endpoint retrieves confirmed and unconfirmed address balance.
//
// For more information: https://developers.whatsonchain.com/#get-balance
func (c *Client) AddressBalance(ctx context.Context, address string) (balance *AddressBalance, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/address/<address>/balance
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/balance", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return balance, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &balance)

	return balance, err
}

// AddressHistory this endpoint retrieves confirmed and unconfirmed address transactions.
//
// For more information: https://developers.whatsonchain.com/#get-history
func (c *Client) AddressHistory(ctx context.Context, address string) (history AddressHistory, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/history
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/history", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return history, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &history)
	return history, err
}

// AddressUnspentTransactions this endpoint retrieves ordered list of UTXOs.
//
// For more information: https://developers.whatsonchain.com/#get-unspent-transactions
func (c *Client) AddressUnspentTransactions(ctx context.Context, address string) (history AddressHistory, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/unspent/all
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/unspent/all", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return history, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &history)
	return history, err
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
// For more information: https://developers.whatsonchain.com/#download-statement
func (c *Client) DownloadStatement(ctx context.Context, address string) (string, error) {
	// https://<network>.whatsonchain.com/statement/<hash>
	// todo: this endpoint does not follow the convention of the WOC API v1
	return c.request(
		ctx,
		fmt.Sprintf("https://%s.whatsonchain.com/statement/%s", c.Network(), address),
		http.MethodGet, nil,
	)
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
// For more information: https://developers.whatsonchain.com/#bulk-balance
func (c *Client) BulkBalance(ctx context.Context, list *AddressList) (balances AddressBalances, err error) {
	// Get the JSON
	var postData []byte
	if postData, err = bulkRequest(list); err != nil {
		return balances, err
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/addresses/balance
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/addresses/balance", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return balances, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &balances)
	return balances, err
}

// BulkUnspentTransactionsProcessor will fetch UTXOs for multiple addresses in a single request while automatically batching
// Max of 20 addresses at a time
//
// For more information: https://developers.whatsonchain.com/#bulk-unspent-transactions
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
// For more information: https://developers.whatsonchain.com/#bulk-unspent-transactions
func (c *Client) BulkUnspentTransactions(ctx context.Context, list *AddressList) (response BulkUnspentResponse, err error) {
	// Get the JSON
	var postData []byte
	if postData, err = bulkRequest(list); err != nil {
		return response, err
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/addresses/unspent/all
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/addresses/unspent/all", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return response, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &response)
	return response, err
}

// AddressUnconfirmedUTXOs retrieves unconfirmed UTXOs for an address
//
// For more information: https://developers.whatsonchain.com/#get-unconfirmed-utxos
func (c *Client) AddressUnconfirmedUTXOs(ctx context.Context, address string) (history AddressHistory, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/address/<address>/unconfirmed/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/unconfirmed/unspent", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return history, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &history)
	return history, err
}

// BulkAddressUnconfirmedUTXOs retrieves unconfirmed UTXOs for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://developers.whatsonchain.com/#bulk-unconfirmed-utxos
func (c *Client) BulkAddressUnconfirmedUTXOs(ctx context.Context, list *AddressList) (response BulkUnspentResponse, err error) {
	// Get the JSON
	var postData []byte
	if postData, err = bulkRequest(list); err != nil {
		return response, err
	}

	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/addresses/unconfirmed/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/addresses/unconfirmed/unspent", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return response, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &response)
	return response, err
}

// AddressConfirmedUTXOs retrieves confirmed UTXOs for an address
//
// For more information: https://developers.whatsonchain.com/#get-confirmed-utxos
func (c *Client) AddressConfirmedUTXOs(ctx context.Context, address string) (history AddressHistory, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/address/<address>/confirmed/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/confirmed/unspent", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return history, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &history)
	return history, err
}

// BulkAddressConfirmedUTXOs retrieves confirmed UTXOs for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://developers.whatsonchain.com/#bulk-confirmed-utxos
func (c *Client) BulkAddressConfirmedUTXOs(ctx context.Context, list *AddressList) (response BulkUnspentResponse, err error) {
	// Get the JSON
	var postData []byte
	if postData, err = bulkRequest(list); err != nil {
		return response, err
	}

	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/addresses/confirmed/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/addresses/confirmed/unspent", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return response, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &response)
	return response, err
}

// AddressUsed retrieves whether an address has been used
//
// For more information: https://docs.whatsonchain.com/api/address#get-address-usage
func (c *Client) AddressUsed(ctx context.Context, address string) (used *AddressUsed, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/address/<address>/used
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/used", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return used, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &used)
	return used, err
}

// AddressScripts retrieves associated scripthashes for an address
//
// For more information: https://docs.whatsonchain.com/api/address#get-associated-scripthashes
func (c *Client) AddressScripts(ctx context.Context, address string) (scripts *AddressScripts, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/address/<address>/scripts
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/scripts", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return scripts, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &scripts)
	return scripts, err
}

// AddressUnconfirmedBalance retrieves unconfirmed balance for an address
//
// For more information: https://docs.whatsonchain.com/api/address#get-unconfirmed-balance
func (c *Client) AddressUnconfirmedBalance(ctx context.Context, address string) (balance *AddressUnconfirmedBalance, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/address/<address>/unconfirmed/balance
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/unconfirmed/balance", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return balance, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &balance)
	return balance, err
}

// AddressConfirmedBalance retrieves confirmed balance for an address
//
// For more information: https://docs.whatsonchain.com/api/address#get-confirmed-balance
func (c *Client) AddressConfirmedBalance(ctx context.Context, address string) (balance *AddressConfirmedBalance, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/address/<address>/confirmed/balance
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/confirmed/balance", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return balance, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &balance)
	return balance, err
}

// AddressUnconfirmedHistory retrieves unconfirmed transaction history for an address
//
// For more information: https://docs.whatsonchain.com/api/address#get-unconfirmed-history
func (c *Client) AddressUnconfirmedHistory(ctx context.Context, address string) (history AddressHistory, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/address/<address>/unconfirmed/history
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/unconfirmed/history", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return history, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &history)
	return history, err
}

// AddressConfirmedHistory retrieves confirmed transaction history for an address
//
// For more information: https://docs.whatsonchain.com/api/address#get-confirmed-history
func (c *Client) AddressConfirmedHistory(ctx context.Context, address string) (history AddressHistory, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/address/<address>/confirmed/history
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/confirmed/history", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return history, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &history)
	return history, err
}

// BulkAddressUnconfirmedBalance retrieves unconfirmed balances for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/api/address#bulk-unconfirmed-balance
func (c *Client) BulkAddressUnconfirmedBalance(ctx context.Context, list *AddressList) (balances AddressBalances, err error) {
	// Get the JSON
	var postData []byte
	if postData, err = bulkRequest(list); err != nil {
		return balances, err
	}

	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/addresses/unconfirmed/balance
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/addresses/unconfirmed/balance", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return balances, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &balances)
	return balances, err
}

// BulkAddressConfirmedBalance retrieves confirmed balances for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/api/address#bulk-confirmed-balance
func (c *Client) BulkAddressConfirmedBalance(ctx context.Context, list *AddressList) (balances AddressBalances, err error) {
	// Get the JSON
	var postData []byte
	if postData, err = bulkRequest(list); err != nil {
		return balances, err
	}

	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/addresses/confirmed/balance
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/addresses/confirmed/balance", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return balances, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &balances)
	return balances, err
}

// BulkAddressUnconfirmedHistory retrieves unconfirmed transaction history for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/api/address#bulk-unconfirmed-history
func (c *Client) BulkAddressUnconfirmedHistory(ctx context.Context, list *AddressList) (history BulkAddressHistoryResponse, err error) {
	// Get the JSON
	var postData []byte
	if postData, err = bulkRequest(list); err != nil {
		return history, err
	}

	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/addresses/unconfirmed/history
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/addresses/unconfirmed/history", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return history, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &history)
	return history, err
}

// BulkAddressConfirmedHistory retrieves confirmed transaction history for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/api/address#bulk-confirmed-history
func (c *Client) BulkAddressConfirmedHistory(ctx context.Context, list *AddressList) (history BulkAddressHistoryResponse, err error) {
	// Get the JSON
	var postData []byte
	if postData, err = bulkRequest(list); err != nil {
		return history, err
	}

	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/addresses/confirmed/history
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/addresses/confirmed/history", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return history, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &history)
	return history, err
}

// BulkAddressHistory retrieves all transaction history for multiple addresses
// Max of 20 addresses at a time
//
// For more information: https://docs.whatsonchain.com/api/address#bulk-history
func (c *Client) BulkAddressHistory(ctx context.Context, list *AddressList) (history BulkAddressHistoryResponse, err error) {
	// Get the JSON
	var postData []byte
	if postData, err = bulkRequest(list); err != nil {
		return history, err
	}

	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/addresses/history/all
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/addresses/history/all", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return history, err
	}
	if len(resp) == 0 {
		return nil, ErrAddressNotFound
	}
	err = json.Unmarshal([]byte(resp), &history)
	return history, err
}
