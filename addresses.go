package whatsonchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddressInfo this endpoint retrieves various address info.
//
// For more information: https://developers.whatsonchain.com/#address
func (c *Client) AddressInfo(address string) (addressInfo *AddressInfo, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/info
	if resp, err = c.Request(fmt.Sprintf("%s%s/address/%s/info", apiEndpoint, c.Network, address), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &addressInfo)
	return
}

// AddressBalance this endpoint retrieves confirmed and unconfirmed address balance.
//
// For more information: https://developers.whatsonchain.com/#get-balance
func (c *Client) AddressBalance(address string) (balance *AddressBalance, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/balance
	if resp, err = c.Request(fmt.Sprintf("%s%s/address/%s/balance", apiEndpoint, c.Network, address), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &balance)
	return
}

// AddressHistory this endpoint retrieves confirmed and unconfirmed address transactions.
//
// For more information: https://developers.whatsonchain.com/#get-history
func (c *Client) AddressHistory(address string) (history AddressHistory, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/history
	if resp, err = c.Request(fmt.Sprintf("%s%s/address/%s/history", apiEndpoint, c.Network, address), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &history)
	return
}

// AddressUnspentTransactions this endpoint retrieves ordered list of UTXOs.
//
// For more information: https://developers.whatsonchain.com/#get-unspent-transactions
func (c *Client) AddressUnspentTransactions(address string) (history AddressHistory, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/unspent
	if resp, err = c.Request(fmt.Sprintf("%s%s/address/%s/unspent", apiEndpoint, c.Network, address), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &history)
	return
}

// AddressUnspentTransactionDetails this endpoint retrieves transaction details for a given address
// Use max transactions to filter if there are more UTXOs returned than needed by the user
//
// For more information: (custom request for this go wrapper)
func (c *Client) AddressUnspentTransactionDetails(address string, maxTransactions int) (history AddressHistory, err error) {

	// Get the address UTXO history
	var utxos AddressHistory
	if utxos, err = c.AddressUnspentTransactions(address); err != nil {
		return
	} else if len(utxos) == 0 {
		return
	}

	// Do we have a "custom max" amount?
	if maxTransactions > 0 {
		total := len(utxos)
		if total > maxTransactions {
			utxos = utxos[:total-(total-maxTransactions)]
		}
	}

	// Loop as long as we still have utxos to get
	for len(utxos) > 0 {

		// Get the hashes
		txHashes := new(TxHashes)
		foundTxs := 0
		for index, tx := range utxos {

			// Only grab the max that can be sent
			if foundTxs >= MaxTransactionsUTXO {
				break
			}

			// Append to the list to send and return
			txHashes.TxIDs = append(txHashes.TxIDs, tx.TxHash)
			history = append(history, tx)

			// Removing from our list to fetch (reducing the list each pass)
			if len(utxos) >= MaxTransactionsUTXO {
				utxos = append(utxos[:index], utxos[index+1:]...)
			} else {
				utxos = AddressHistory{}
			}
			foundTxs = foundTxs + 1
		}

		// Get the tx details
		var txList TxList
		if txList, err = c.BulkTransactionDetails(txHashes); err != nil {
			return
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

	return
}

// DownloadStatement this endpoint downloads an address statement (PDF)
// The contents will be returned in plain-text and need to be converted to a file.pdf
//
// For more information: https://developers.whatsonchain.com/#download-statement
func (c *Client) DownloadStatement(address string) (string, error) {

	// https://<network>.whatsonchain.com/statement/<hash>
	// todo: this endpoint does not follow the convention of the WOC API v1
	return c.Request(fmt.Sprintf("https://%s.whatsonchain.com/statement/%s", c.Network, address), http.MethodGet, nil)
}
