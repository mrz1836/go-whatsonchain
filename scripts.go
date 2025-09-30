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
	// https://api.whatsonchain.com/v1/bsv/<network>/script/<scriptHash>/unspent/all
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/script/%s/unspent/all", apiEndpointBase, c.Chain(), c.Network(), scriptHash),
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
	// https://api.whatsonchain.com/v1/bsv/<network>/scripts/unspent/all
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/scripts/unspent/all", apiEndpointBase, c.Chain(), c.Network()),
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

// ScriptUnconfirmedUTXOs retrieves unconfirmed UTXOs for a script
//
// For more information: https://developers.whatsonchain.com/#get-unconfirmed-script-utxos
func (c *Client) ScriptUnconfirmedUTXOs(ctx context.Context, scriptHash string) (scriptList ScriptList, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/script/<scriptHash>/unconfirmed/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/script/%s/unconfirmed/unspent", apiEndpointBase, c.Chain(), c.Network(), scriptHash),
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

// BulkScriptUnconfirmedUTXOs retrieves unconfirmed UTXOs for multiple scripts
// Max of 20 scripts at a time
//
// For more information: https://developers.whatsonchain.com/#bulk-unconfirmed-script-utxos
func (c *Client) BulkScriptUnconfirmedUTXOs(ctx context.Context, list *ScriptsList) (response BulkScriptUnspentResponse, err error) {
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
	// https://api.whatsonchain.com/v1/<chain>/<network>/scripts/unconfirmed/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/scripts/unconfirmed/unspent", apiEndpointBase, c.Chain(), c.Network()),
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

// ScriptConfirmedUTXOs retrieves confirmed UTXOs for a script
//
// For more information: https://developers.whatsonchain.com/#get-confirmed-script-utxos
func (c *Client) ScriptConfirmedUTXOs(ctx context.Context, scriptHash string) (scriptList ScriptList, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/script/<scriptHash>/confirmed/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/script/%s/confirmed/unspent", apiEndpointBase, c.Chain(), c.Network(), scriptHash),
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

// BulkScriptConfirmedUTXOs retrieves confirmed UTXOs for multiple scripts
// Max of 20 scripts at a time
//
// For more information: https://developers.whatsonchain.com/#bulk-confirmed-script-utxos
func (c *Client) BulkScriptConfirmedUTXOs(ctx context.Context, list *ScriptsList) (response BulkScriptUnspentResponse, err error) {
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
	// https://api.whatsonchain.com/v1/<chain>/<network>/scripts/confirmed/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/scripts/confirmed/unspent", apiEndpointBase, c.Chain(), c.Network()),
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

// GetScriptUsed this endpoint determines if a script has been used in any transaction
//
// For more information: https://docs.whatsonchain.com/api/script#get-script-usage
func (c *Client) GetScriptUsed(ctx context.Context, scriptHash string) (used bool, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/script/<scriptHash>/used
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/script/%s/used", apiEndpointBase, c.Chain(), c.Network(), scriptHash),
		http.MethodGet, nil,
	); err != nil {
		return used, err
	}
	if len(resp) == 0 {
		if c.LastRequest().StatusCode == http.StatusNotFound {
			return false, ErrScriptNotFound
		}
		return false, nil
	}
	// The response is a simple boolean string "true" or "false"
	if resp == "true" {
		return true, nil
	}
	return false, nil
}

// GetScriptUnconfirmedHistory this endpoint retrieves unconfirmed script transactions
//
// For more information: https://docs.whatsonchain.com/api/script#get-unconfirmed-script-history
func (c *Client) GetScriptUnconfirmedHistory(ctx context.Context, scriptHash string) (history ScriptList, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/script/<scriptHash>/unconfirmed/history
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/script/%s/unconfirmed/history", apiEndpointBase, c.Chain(), c.Network(), scriptHash),
		http.MethodGet, nil,
	); err != nil {
		return history, err
	}
	if len(resp) == 0 {
		if c.LastRequest().StatusCode == http.StatusNotFound {
			return nil, ErrScriptNotFound
		}
		return history, nil
	}
	err = json.Unmarshal([]byte(resp), &history)
	return history, err
}

// BulkScriptUnconfirmedHistory will fetch unconfirmed history for multiple scripts in a single request
// Max of 20 scripts at a time
//
// For more information: https://docs.whatsonchain.com/api/script#bulk-unconfirmed-script-history
func (c *Client) BulkScriptUnconfirmedHistory(ctx context.Context, list *ScriptsList) (response BulkScriptHistoryResponse, err error) {
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
	// https://api.whatsonchain.com/v1/<chain>/<network>/scripts/unconfirmed/history
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/scripts/unconfirmed/history", apiEndpointBase, c.Chain(), c.Network()),
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

// GetScriptConfirmedHistory this endpoint retrieves confirmed script transactions
//
// For more information: https://docs.whatsonchain.com/api/script#get-confirmed-script-history
func (c *Client) GetScriptConfirmedHistory(ctx context.Context, scriptHash string) (history ScriptList, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/<chain>/<network>/script/<scriptHash>/confirmed/history
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/script/%s/confirmed/history", apiEndpointBase, c.Chain(), c.Network(), scriptHash),
		http.MethodGet, nil,
	); err != nil {
		return history, err
	}
	if len(resp) == 0 {
		if c.LastRequest().StatusCode == http.StatusNotFound {
			return nil, ErrScriptNotFound
		}
		return history, nil
	}
	err = json.Unmarshal([]byte(resp), &history)
	return history, err
}

// BulkScriptConfirmedHistory will fetch confirmed history for multiple scripts in a single request
// Max of 20 scripts at a time
//
// For more information: https://docs.whatsonchain.com/api/script#bulk-confirmed-script-history
func (c *Client) BulkScriptConfirmedHistory(ctx context.Context, list *ScriptsList) (response BulkScriptHistoryResponse, err error) {
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
	// https://api.whatsonchain.com/v1/<chain>/<network>/scripts/confirmed/history
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/scripts/confirmed/history", apiEndpointBase, c.Chain(), c.Network()),
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
