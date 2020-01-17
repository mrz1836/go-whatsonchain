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
