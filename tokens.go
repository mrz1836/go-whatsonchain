package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetOneSatOrdinalByOrigin gets a 1Sat Ordinal token by origin (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-by-origin
func (c *Client) GetOneSatOrdinalByOrigin(ctx context.Context, origin string) (token *OneSatOrdinalToken, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/token/1satordinals/<origin>/origin
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/token/1satordinals/%s/origin", apiEndpointBase, c.Chain(), c.Network(), origin),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	if c.LastRequest().StatusCode == http.StatusNotFound {
		return nil, ErrTokenNotFound
	}

	var response OneSatOrdinalToken
	err = json.Unmarshal([]byte(resp), &response)
	return &response, err
}

// GetOneSatOrdinalByOutpoint gets a 1Sat Ordinal token by outpoint (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-by-outpoint
func (c *Client) GetOneSatOrdinalByOutpoint(ctx context.Context, outpoint string) (token *OneSatOrdinalToken, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/token/1satordinals/<outpoint>
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/token/1satordinals/%s", apiEndpointBase, c.Chain(), c.Network(), outpoint),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var response OneSatOrdinalToken
	err = json.Unmarshal([]byte(resp), &response)
	return &response, err
}

// GetOneSatOrdinalContent gets content data for a 1Sat Ordinal token (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-content
func (c *Client) GetOneSatOrdinalContent(ctx context.Context, outpoint string) (content *OneSatOrdinalContent, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/token/1satordinals/<outpoint>/content
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/token/1satordinals/%s/content", apiEndpointBase, c.Chain(), c.Network(), outpoint),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var response OneSatOrdinalContent
	err = json.Unmarshal([]byte(resp), &response)
	return &response, err
}

// GetOneSatOrdinalLatest gets the latest transfer of a 1Sat Ordinal token (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-latest-transfer
func (c *Client) GetOneSatOrdinalLatest(ctx context.Context, outpoint string) (latest *OneSatOrdinalLatest, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/token/1satordinals/<outpoint>/latest
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/token/1satordinals/%s/latest", apiEndpointBase, c.Chain(), c.Network(), outpoint),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var response OneSatOrdinalLatest
	err = json.Unmarshal([]byte(resp), &response)
	return &response, err
}

// GetOneSatOrdinalHistory gets transfer history of a 1Sat Ordinal token (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-transfers-history
func (c *Client) GetOneSatOrdinalHistory(ctx context.Context, outpoint string) (history []*OneSatOrdinalHistory, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/token/1satordinals/<outpoint>/history
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/token/1satordinals/%s/history", apiEndpointBase, c.Chain(), c.Network(), outpoint),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var response []*OneSatOrdinalHistory
	err = json.Unmarshal([]byte(resp), &response)
	return response, err
}

// GetOneSatOrdinalsByTxID gets all 1Sat Ordinal tokens in a transaction (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-tokens-by-txid
func (c *Client) GetOneSatOrdinalsByTxID(ctx context.Context, txid string) (tokens []*OneSatOrdinalToken, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/token/1satordinals/tx/<txid>
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/token/1satordinals/tx/%s", apiEndpointBase, c.Chain(), c.Network(), txid),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var response []*OneSatOrdinalToken
	err = json.Unmarshal([]byte(resp), &response)
	return response, err
}

// GetOneSatOrdinalsStats gets statistics for 1Sat Ordinals (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-stats
func (c *Client) GetOneSatOrdinalsStats(ctx context.Context) (stats *OneSatOrdinalStats, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tokens/1satordinals
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/tokens/1satordinals", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var response OneSatOrdinalStats
	err = json.Unmarshal([]byte(resp), &response)
	return &response, err
}

// GetAllSTASTokens gets all STAS tokens (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/stas#get-all-tokens
func (c *Client) GetAllSTASTokens(ctx context.Context) (tokens []*STASToken, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tokens
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/tokens", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var response []*STASToken
	err = json.Unmarshal([]byte(resp), &response)
	return response, err
}

// GetSTASTokenByID gets a STAS token by contract ID and symbol (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/stas#get-token-by-id
func (c *Client) GetSTASTokenByID(ctx context.Context, contractID, symbol string) (token *STASToken, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/token/<contractId>/<symbol>
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/token/%s/%s", apiEndpointBase, c.Chain(), c.Network(), contractID, symbol),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var response STASToken
	err = json.Unmarshal([]byte(resp), &response)
	return &response, err
}

// GetTokenUTXOsForAddress gets token UTXOs for an address (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/stas#get-token-utxos-for-address
func (c *Client) GetTokenUTXOsForAddress(ctx context.Context, address string) (utxos []*STASTokenUTXO, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/tokens/unspent
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/tokens/unspent", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var response []*STASTokenUTXO
	err = json.Unmarshal([]byte(resp), &response)
	return response, err
}

// GetAddressTokenBalance gets token balance for an address (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/stas#get-address-token-balance
func (c *Client) GetAddressTokenBalance(ctx context.Context, address string) (balance *STASTokenBalance, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/address/<address>/tokens
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/address/%s/tokens", apiEndpointBase, c.Chain(), c.Network(), address),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var response STASTokenBalance
	err = json.Unmarshal([]byte(resp), &response)
	return &response, err
}

// GetTokenTransactions gets transactions for a STAS token (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/stas#get-token-transactions
func (c *Client) GetTokenTransactions(ctx context.Context, contractID, symbol string) (transactions TxList, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/token/<contractId>/<symbol>/tx
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/token/%s/%s/tx", apiEndpointBase, c.Chain(), c.Network(), contractID, symbol),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var response TxList
	err = json.Unmarshal([]byte(resp), &response)
	return response, err
}

// GetSTASStats gets statistics for STAS tokens (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/stas#get-stats
func (c *Client) GetSTASStats(ctx context.Context) (stats *STASStats, err error) {
	// Only available for BSV
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/tokens/stas
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/tokens/stas", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodGet, nil,
	); err != nil {
		return nil, err
	}

	var response STASStats
	err = json.Unmarshal([]byte(resp), &response)
	return &response, err
}
