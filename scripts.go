package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetScriptHistory retrieves the transaction history for a script.
//
// Deprecated: GetScriptHistory uses a combined history endpoint that is no longer in the API.
// Use GetScriptConfirmedHistory and GetScriptUnconfirmedHistory instead.
//
// For more information: https://docs.whatsonchain.com/#get-script-history
func (c *Client) GetScriptHistory(ctx context.Context, scriptHash string) (ScriptList, error) {
	url := c.buildURL("/script/%s/history", scriptHash)
	return requestAndUnmarshalSlice[*ScriptRecord](ctx, c, url, http.MethodGet, nil, ErrScriptNotFound)
}

// GetScriptUnspentTransactions this endpoint retrieves ordered list of UTXOs
//
// For more information: https://docs.whatsonchain.com/#get-script-unspent-transactions
func (c *Client) GetScriptUnspentTransactions(ctx context.Context, scriptHash string) (ScriptList, error) {
	url := c.buildURL("/script/%s/unspent/all", scriptHash)
	return requestAndUnmarshalSlice[*ScriptRecord](ctx, c, url, http.MethodGet, nil, ErrScriptNotFound)
}

// BulkScriptUnspentTransactions retrieves unspent transactions for multiple scripts.
// Max of 20 scripts at a time.
//
// Deprecated: BulkScriptUnspentTransactions uses a combined unspent endpoint that is no longer in the API.
// Use BulkScriptConfirmedUTXOs and BulkScriptUnconfirmedUTXOs instead.
//
// For more information: https://docs.whatsonchain.com/#bulk-script-unspent-transactions
func (c *Client) BulkScriptUnspentTransactions(ctx context.Context, list *ScriptsList) (BulkScriptUnspentResponse, error) {
	if len(list.Scripts) > MaxScriptsForLookup {
		return nil, fmt.Errorf("%w: %d scripts requested, max is %d", ErrMaxScriptsExceeded, len(list.Scripts), MaxScriptsForLookup)
	}

	postData, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/scripts/unspent/all")
	return requestAndUnmarshalSlice[*BulkScriptResponseRecord](ctx, c, url, http.MethodPost, postData, ErrScriptNotFound)
}

// ScriptUnconfirmedUTXOs retrieves unconfirmed UTXOs for a script
//
// For more information: https://docs.whatsonchain.com/#get-unconfirmed-script-utxos
func (c *Client) ScriptUnconfirmedUTXOs(ctx context.Context, scriptHash string) (ScriptList, error) {
	url := c.buildURL("/script/%s/unconfirmed/unspent", scriptHash)
	return requestAndUnmarshalSlice[*ScriptRecord](ctx, c, url, http.MethodGet, nil, ErrScriptNotFound)
}

// BulkScriptUnconfirmedUTXOs retrieves unconfirmed UTXOs for multiple scripts
// Max of 20 scripts at a time
//
// For more information: https://docs.whatsonchain.com/#bulk-unconfirmed-script-utxos
func (c *Client) BulkScriptUnconfirmedUTXOs(ctx context.Context, list *ScriptsList) (BulkScriptUnspentResponse, error) {
	if len(list.Scripts) > MaxScriptsForLookup {
		return nil, fmt.Errorf("%w: %d scripts requested, max is %d", ErrMaxScriptsExceeded, len(list.Scripts), MaxScriptsForLookup)
	}

	postData, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/scripts/unconfirmed/unspent")
	return requestAndUnmarshalSlice[*BulkScriptResponseRecord](ctx, c, url, http.MethodPost, postData, ErrScriptNotFound)
}

// ScriptConfirmedUTXOs retrieves confirmed UTXOs for a script
//
// For more information: https://docs.whatsonchain.com/#get-confirmed-script-utxos
func (c *Client) ScriptConfirmedUTXOs(ctx context.Context, scriptHash string) (ScriptList, error) {
	url := c.buildURL("/script/%s/confirmed/unspent", scriptHash)
	return requestAndUnmarshalSlice[*ScriptRecord](ctx, c, url, http.MethodGet, nil, ErrScriptNotFound)
}

// BulkScriptConfirmedUTXOs retrieves confirmed UTXOs for multiple scripts
// Max of 20 scripts at a time
//
// For more information: https://docs.whatsonchain.com/#bulk-confirmed-script-utxos
func (c *Client) BulkScriptConfirmedUTXOs(ctx context.Context, list *ScriptsList) (BulkScriptUnspentResponse, error) {
	if len(list.Scripts) > MaxScriptsForLookup {
		return nil, fmt.Errorf("%w: %d scripts requested, max is %d", ErrMaxScriptsExceeded, len(list.Scripts), MaxScriptsForLookup)
	}

	postData, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/scripts/confirmed/unspent")
	return requestAndUnmarshalSlice[*BulkScriptResponseRecord](ctx, c, url, http.MethodPost, postData, ErrScriptNotFound)
}

// GetScriptUsed this endpoint determines if a script has been used in any transaction
//
// For more information: https://docs.whatsonchain.com/api/script#get-script-usage
func (c *Client) GetScriptUsed(ctx context.Context, scriptHash string) (bool, error) {
	url := c.buildURL("/script/%s/used", scriptHash)
	resp, err := requestString(ctx, c, url)
	if err != nil {
		return false, err
	}
	if len(resp) == 0 {
		return false, ErrScriptNotFound
	}
	// The response is a simple boolean string "true" or "false"
	return resp == "true", nil
}

// GetScriptUnconfirmedHistory this endpoint retrieves unconfirmed script transactions
//
// For more information: https://docs.whatsonchain.com/api/script#get-unconfirmed-script-history
func (c *Client) GetScriptUnconfirmedHistory(ctx context.Context, scriptHash string) (ScriptList, error) {
	url := c.buildURL("/script/%s/unconfirmed/history", scriptHash)
	return requestAndUnmarshalSlice[*ScriptRecord](ctx, c, url, http.MethodGet, nil, ErrScriptNotFound)
}

// BulkScriptUnconfirmedHistory will fetch unconfirmed history for multiple scripts in a single request
// Max of 20 scripts at a time
//
// For more information: https://docs.whatsonchain.com/api/script#bulk-unconfirmed-script-history
func (c *Client) BulkScriptUnconfirmedHistory(ctx context.Context, list *ScriptsList) (BulkScriptHistoryResponse, error) {
	if len(list.Scripts) > MaxScriptsForLookup {
		return nil, fmt.Errorf("%w: %d scripts requested, max is %d", ErrMaxScriptsExceeded, len(list.Scripts), MaxScriptsForLookup)
	}

	postData, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/scripts/unconfirmed/history")
	return requestAndUnmarshalSlice[*BulkScriptHistoryRecord](ctx, c, url, http.MethodPost, postData, ErrScriptNotFound)
}

// GetScriptConfirmedHistory this endpoint retrieves confirmed script transactions
//
// For more information: https://docs.whatsonchain.com/api/script#get-confirmed-script-history
func (c *Client) GetScriptConfirmedHistory(ctx context.Context, scriptHash string) (ScriptList, error) {
	url := c.buildURL("/script/%s/confirmed/history", scriptHash)
	return requestAndUnmarshalSlice[*ScriptRecord](ctx, c, url, http.MethodGet, nil, ErrScriptNotFound)
}

// BulkScriptConfirmedHistory will fetch confirmed history for multiple scripts in a single request
// Max of 20 scripts at a time
//
// For more information: https://docs.whatsonchain.com/api/script#bulk-confirmed-script-history
func (c *Client) BulkScriptConfirmedHistory(ctx context.Context, list *ScriptsList) (BulkScriptHistoryResponse, error) {
	if len(list.Scripts) > MaxScriptsForLookup {
		return nil, fmt.Errorf("%w: %d scripts requested, max is %d", ErrMaxScriptsExceeded, len(list.Scripts), MaxScriptsForLookup)
	}

	postData, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}

	url := c.buildURL("/scripts/confirmed/history")
	return requestAndUnmarshalSlice[*BulkScriptHistoryRecord](ctx, c, url, http.MethodPost, postData, ErrScriptNotFound)
}
