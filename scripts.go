package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetScriptHistory this endpoint retrieves confirmed and unconfirmed script transactions
//
// For more information: https://developers.whatsonchain.com/#get-script-history
func (c *Client) GetScriptHistory(ctx context.Context, scriptHash string) (history ScriptList, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/script/<scriptHash>/history
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/script/%s/history", apiEndpoint, c.Network(), scriptHash),
		http.MethodGet, nil,
	); err != nil {
		return
	}
	if len(resp) == 0 {
		return nil, ErrScriptNotFound
	}
	err = json.Unmarshal([]byte(resp), &history)
	return
}

// GetScriptUnspentTransactions this endpoint retrieves ordered list of UTXOs
//
// For more information: https://developers.whatsonchain.com/#get-script-unspent-transactions
func (c *Client) GetScriptUnspentTransactions(ctx context.Context,
	scriptHash string) (scriptList ScriptList, err error) {

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/script/<scriptHash>/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/script/%s/unspent", apiEndpoint, c.Network(), scriptHash),
		http.MethodGet, nil,
	); err != nil {
		return
	}
	if len(resp) == 0 {
		return nil, ErrScriptNotFound
	}
	err = json.Unmarshal([]byte(resp), &scriptList)

	return
}

// BulkScriptUnspentTransactions will fetch UTXOs for multiple scripts in a single request
// Max of 20 scripts at a time
//
// For more information: https://developers.whatsonchain.com/#bulk-script-unspent-transactions
func (c *Client) BulkScriptUnspentTransactions(ctx context.Context,
	list *ScriptsList) (response BulkScriptUnspentResponse, err error) {

	// The max limit by WOC
	if len(list.Scripts) > MaxScriptsForLookup {
		return nil, fmt.Errorf(
			"max limit of scripts is %d and you sent %d",
			MaxScriptsForLookup, len(list.Scripts),
		)
	}

	// Get the JSON
	var postData []byte
	if postData, err = json.Marshal(list); err != nil {
		return
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/scripts/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/scripts/unspent", apiEndpoint, c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return
	}
	if len(resp) == 0 {
		return nil, ErrScriptNotFound
	}
	err = json.Unmarshal([]byte(resp), &response)
	return
}
