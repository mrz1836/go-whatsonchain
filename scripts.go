package whatsonchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetScriptHistory this endpoint retrieves confirmed and unconfirmed script transactions
//
// For more information: https://developers.whatsonchain.com/#get-script-history
func (c *Client) GetScriptHistory(scriptHash string) (history ScriptList, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/script/<scriptHash>/history
	if resp, err = c.Request(fmt.Sprintf("%s%s/script/%s/history", apiEndpoint, c.Network, scriptHash), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &history)
	return
}

// GetScriptUnspentTransactions this endpoint retrieves ordered list of UTXOs
//
// For more information: https://developers.whatsonchain.com/#get-script-unspent-transactions
func (c *Client) GetScriptUnspentTransactions(scriptHash string) (scriptList ScriptList, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/script/<scriptHash>/unspent
	if resp, err = c.Request(fmt.Sprintf("%s%s/script/%s/unspent", apiEndpoint, c.Network, scriptHash), http.MethodGet, nil); err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &scriptList)

	return
}
