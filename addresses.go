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

// AddressUnspentTransactionDetails this endpoint retrieves transaction details for a given address
// Max returned is the limit in the definitions.go
//
// For more information: (custom request for this go wrapper)
func (c *Client) AddressUnspentTransactionDetails(address string, maxTransactions int) (history AddressHistory, err error) {

	// Get the address UTXO history
	var tempHistory AddressHistory
	if tempHistory, err = c.AddressUnspentTransactions(address); err != nil {
		return
	} else if len(tempHistory) == 0 {
		return
	}

	// Set the max to return
	if maxTransactions < 0 || maxTransactions > MaxTransactionsUTXO {
		maxTransactions = MaxTransactionsUTXO
	}

	// Get the hashes
	txHashes := new(TxHashes)
	foundTxs := 0
	for index, tx := range tempHistory {
		if foundTxs >= maxTransactions {
			break
		}
		txHashes.TxIDs = append(txHashes.TxIDs, tx.TxHash)
		history = append(history, tempHistory[index])
		foundTxs = foundTxs + 1
	}

	// Hashes into json
	var postData []byte
	if postData, err = json.Marshal(txHashes); err != nil {
		return
	}

	// Fire the request
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/txs
	if resp, err = c.Request(fmt.Sprintf("%s%s/txs", apiEndpoint, c.Parameters.Network), http.MethodPost, postData); err != nil {
		return
	}

	// Unmarshal the response
	var txList TxList
	if err = json.Unmarshal([]byte(resp), &txList); err != nil {
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

	return
}
