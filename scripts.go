package whatsonchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// scriptListEnvelope is the response envelope returned by the newer paginated
// single-script endpoints, e.g. /script/{hash}/unspent/all and the
// confirmed/unconfirmed history and unspent endpoints:
//
//	{"script":"...","result":[...],"error":""}
//
// The deprecated single-script endpoints (/script/{hash}/unspent and
// /script/{hash}/history) instead return a bare JSON array. unmarshalScriptList
// accepts either shape.
type scriptListEnvelope struct {
	Script string     `json:"script"`
	Result ScriptList `json:"result"`
	Error  string     `json:"error"`
}

// requestScriptList performs a GET request and decodes a ScriptList, tolerating
// both the bare-array form (deprecated endpoints) and the {script,result,error}
// envelope (newer paginated endpoints).
func (c *Client) requestScriptList(ctx context.Context, url string) (ScriptList, error) {
	resp, statusCode, err := c.request(ctx, url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	if err = checkStatusCode(statusCode, resp); err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, ErrScriptNotFound
	}
	return unmarshalScriptList(resp)
}

// unmarshalScriptList decodes a ScriptList from either a bare JSON array or a
// {script,result,error} envelope. When the envelope carries a non-empty error
// field, it is surfaced as an error.
func unmarshalScriptList(data []byte) (ScriptList, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && trimmed[0] == '[' {
		var list ScriptList
		if err := json.Unmarshal(trimmed, &list); err != nil {
			return nil, err
		}
		return list, nil
	}

	var envelope scriptListEnvelope
	if err := json.Unmarshal(trimmed, &envelope); err != nil {
		return nil, err
	}
	if envelope.Error != "" {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, envelope.Error)
	}
	return envelope.Result, nil
}

// GetScriptHistory retrieves the transaction history for a script.
//
// Deprecated: GetScriptHistory uses a combined history endpoint that is no longer in the API.
// Use GetScriptConfirmedHistory and GetScriptUnconfirmedHistory instead.
//
// For more information: https://docs.whatsonchain.com/#get-script-history
func (c *Client) GetScriptHistory(ctx context.Context, scriptHash string) (ScriptList, error) {
	url := c.buildURL("/script/%s/history", scriptHash)
	return c.requestScriptList(ctx, url)
}

// GetScriptUnspentTransactions this endpoint retrieves ordered list of UTXOs
//
// For more information: https://docs.whatsonchain.com/#get-script-unspent-transactions
func (c *Client) GetScriptUnspentTransactions(ctx context.Context, scriptHash string) (ScriptList, error) {
	url := c.buildURL("/script/%s/unspent/all", scriptHash)
	return c.requestScriptList(ctx, url)
}

// BulkScriptUnspentTransactions retrieves unspent transactions for multiple scripts.
// Max of 20 scripts at a time.
//
// Deprecated: BulkScriptUnspentTransactions uses a combined unspent endpoint that is no longer in the API.
// Use BulkScriptConfirmedUTXOs and BulkScriptUnconfirmedUTXOs instead.
//
// For more information: https://docs.whatsonchain.com/#bulk-script-unspent-transactions
func (c *Client) BulkScriptUnspentTransactions(ctx context.Context, list *ScriptsList) (BulkScriptUnspentResponse, error) {
	if list == nil {
		return nil, ErrMissingRequest
	}
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
	return c.requestScriptList(ctx, url)
}

// BulkScriptUnconfirmedUTXOs retrieves unconfirmed UTXOs for multiple scripts
// Max of 20 scripts at a time
//
// For more information: https://docs.whatsonchain.com/#bulk-unconfirmed-script-utxos
func (c *Client) BulkScriptUnconfirmedUTXOs(ctx context.Context, list *ScriptsList) (BulkScriptUnspentResponse, error) {
	if list == nil {
		return nil, ErrMissingRequest
	}
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
	return c.requestScriptList(ctx, url)
}

// BulkScriptConfirmedUTXOs retrieves confirmed UTXOs for multiple scripts
// Max of 20 scripts at a time
//
// For more information: https://docs.whatsonchain.com/#bulk-confirmed-script-utxos
func (c *Client) BulkScriptConfirmedUTXOs(ctx context.Context, list *ScriptsList) (BulkScriptUnspentResponse, error) {
	if list == nil {
		return nil, ErrMissingRequest
	}
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
	return c.requestScriptList(ctx, url)
}

// BulkScriptUnconfirmedHistory will fetch unconfirmed history for multiple scripts in a single request
// Max of 20 scripts at a time
//
// For more information: https://docs.whatsonchain.com/api/script#bulk-unconfirmed-script-history
func (c *Client) BulkScriptUnconfirmedHistory(ctx context.Context, list *ScriptsList) (BulkScriptHistoryResponse, error) {
	if list == nil {
		return nil, ErrMissingRequest
	}
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
	return c.requestScriptList(ctx, url)
}

// BulkScriptConfirmedHistory will fetch confirmed history for multiple scripts in a single request
// Max of 20 scripts at a time
//
// For more information: https://docs.whatsonchain.com/api/script#bulk-confirmed-script-history
func (c *Client) BulkScriptConfirmedHistory(ctx context.Context, list *ScriptsList) (BulkScriptHistoryResponse, error) {
	if list == nil {
		return nil, ErrMissingRequest
	}
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
