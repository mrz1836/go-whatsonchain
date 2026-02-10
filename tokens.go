package whatsonchain

import (
	"context"
	"net/http"
)

// GetOneSatOrdinalByOrigin gets a 1Sat Ordinal token by origin (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-by-origin
func (c *Client) GetOneSatOrdinalByOrigin(ctx context.Context, origin string) (*OneSatOrdinalToken, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/token/1satordinals/%s/origin", origin)
	return requestAndUnmarshal[OneSatOrdinalToken](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}

// GetOneSatOrdinalByOutpoint gets a 1Sat Ordinal token by outpoint (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-by-outpoint
func (c *Client) GetOneSatOrdinalByOutpoint(ctx context.Context, outpoint string) (*OneSatOrdinalToken, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/token/1satordinals/%s", outpoint)
	return requestAndUnmarshal[OneSatOrdinalToken](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}

// GetOneSatOrdinalContent gets content data for a 1Sat Ordinal token (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-content
func (c *Client) GetOneSatOrdinalContent(ctx context.Context, outpoint string) (*OneSatOrdinalContent, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/token/1satordinals/%s/content", outpoint)
	return requestAndUnmarshal[OneSatOrdinalContent](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}

// GetOneSatOrdinalLatest gets the latest transfer of a 1Sat Ordinal token (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-latest-transfer
func (c *Client) GetOneSatOrdinalLatest(ctx context.Context, outpoint string) (*OneSatOrdinalLatest, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/token/1satordinals/%s/latest", outpoint)
	return requestAndUnmarshal[OneSatOrdinalLatest](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}

// GetOneSatOrdinalHistory gets transfer history of a 1Sat Ordinal token (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-transfers-history
func (c *Client) GetOneSatOrdinalHistory(ctx context.Context, outpoint string) ([]*OneSatOrdinalHistory, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/token/1satordinals/%s/history", outpoint)
	return requestAndUnmarshalSlice[*OneSatOrdinalHistory](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}

// GetOneSatOrdinalsByTxID gets all 1Sat Ordinal tokens in a transaction (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-tokens-by-txid
func (c *Client) GetOneSatOrdinalsByTxID(ctx context.Context, txid string) ([]*OneSatOrdinalToken, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/token/1satordinals/tx/%s", txid)
	return requestAndUnmarshalSlice[*OneSatOrdinalToken](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}

// GetOneSatOrdinalsStats gets statistics for 1Sat Ordinals (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-stats
func (c *Client) GetOneSatOrdinalsStats(ctx context.Context) (*OneSatOrdinalStats, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/tokens/1satordinals")
	return requestAndUnmarshal[OneSatOrdinalStats](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}

// GetAllSTASTokens gets all STAS tokens (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/stas#get-all-tokens
func (c *Client) GetAllSTASTokens(ctx context.Context) ([]*STASToken, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/tokens")
	return requestAndUnmarshalSlice[*STASToken](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}

// GetSTASTokenByID gets a STAS token by contract ID and symbol (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/stas#get-token-by-id
func (c *Client) GetSTASTokenByID(ctx context.Context, contractID, symbol string) (*STASToken, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/token/%s/%s", contractID, symbol)
	return requestAndUnmarshal[STASToken](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}

// GetTokenUTXOsForAddress gets token UTXOs for an address (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/stas#get-token-utxos-for-address
func (c *Client) GetTokenUTXOsForAddress(ctx context.Context, address string) ([]*STASTokenUTXO, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/address/%s/tokens/unspent", address)
	return requestAndUnmarshalSlice[*STASTokenUTXO](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}

// GetAddressTokenBalance gets token balance for an address (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/stas#get-address-token-balance
func (c *Client) GetAddressTokenBalance(ctx context.Context, address string) (*STASTokenBalance, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/address/%s/tokens", address)
	return requestAndUnmarshal[STASTokenBalance](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}

// GetTokenTransactions gets transactions for a STAS token (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/stas#get-token-transactions
func (c *Client) GetTokenTransactions(ctx context.Context, contractID, symbol string) (TxList, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/token/%s/%s/tx", contractID, symbol)
	return requestAndUnmarshalSlice[*TxInfo](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}

// GetSTASStats gets statistics for STAS tokens (BSV-only endpoint)
//
// For more information: https://docs.whatsonchain.com/api/tokens/stas#get-stats
func (c *Client) GetSTASStats(ctx context.Context) (*STASStats, error) {
	if c.Chain() != ChainBSV {
		return nil, ErrBSVChainRequired
	}

	url := c.buildURL("/tokens/stas")
	return requestAndUnmarshal[STASStats](ctx, c, url, http.MethodGet, nil, ErrTokenNotFound)
}
