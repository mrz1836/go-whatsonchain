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
		fmt.Sprintf("%s%s/%s/script/%s/history", apiEndpointBase, c.Chain(), c.Network(), scriptHash),
		http.MethodGet, nil,
	); err != nil {
		return history, err
	}
	if len(resp) == 0 {
		return nil, ErrScriptNotFound
	}
	err = json.Unmarshal([]byte(resp), &history)
	return history, err
}

// GetScriptUnspentTransactions this endpoint retrieves ordered list of UTXOs
//
// For more information: https://developers.whatsonchain.com/#get-script-unspent-transactions
func (c *Client) GetScriptUnspentTransactions(ctx context.Context,
	scriptHash string,
) (scriptList ScriptList, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/script/<scriptHash>/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/script/%s/unspent", apiEndpointBase, c.Chain(), c.Network(), scriptHash),
		http.MethodGet, nil,
	); err != nil {
		return scriptList, err
	}
	if len(resp) == 0 {
		if c.LastRequest().StatusCode == http.StatusNotFound {
			return nil, ErrScriptNotFound
		}
		return scriptList, nil
	}
	err = json.Unmarshal([]byte(resp), &scriptList)

	return scriptList, err
}

// BulkScriptUnspentTransactions will fetch UTXOs for multiple scripts in a single request
// Max of 20 scripts at a time
//
// For more information: https://developers.whatsonchain.com/#bulk-script-unspent-transactions
func (c *Client) BulkScriptUnspentTransactions(ctx context.Context,
	list *ScriptsList,
) (response BulkScriptUnspentResponse, err error) {
	// The max limit by WOC
	if len(list.Scripts) > MaxScriptsForLookup {
		return nil, fmt.Errorf("%w: %d scripts requested, max is %d", ErrMaxScriptsExceeded, len(list.Scripts), MaxScriptsForLookup)
	}

	// Get the JSON
	var postData []byte
	if postData, err = json.Marshal(list); err != nil {
		return response, err
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/scripts/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/scripts/unspent", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodPost, postData,
	); err != nil {
		return response, err
	}
	if len(resp) == 0 {
		return nil, ErrScriptNotFound
	}
	err = json.Unmarshal([]byte(resp), &response)
	return response, err
}
