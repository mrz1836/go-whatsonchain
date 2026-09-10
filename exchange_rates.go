package whatsonchain

import (
	"context"
	"net/http"
)

// GetExchangeRate this endpoint provides the exchange rate for BSV
//
// For more information: https://docs.whatsonchain.com/#get-exchange-rate
func (c *Client) GetExchangeRate(ctx context.Context) (*ExchangeRate, error) {
	url := c.buildURL("/exchangerate")
	return requestAndUnmarshal[ExchangeRate](ctx, c, url, http.MethodGet, nil, ErrExchangeRateNotFound)
}

// GetHistoricalExchangeRate this endpoint provides historical exchange rates for BSV
// within a specified time range
//
// For more information: https://docs.whatsonchain.com/#get-historical-exchange-rate
func (c *Client) GetHistoricalExchangeRate(ctx context.Context, from, to int64) ([]*HistoricalExchangeRate, error) {
	url := c.buildURL("/exchangerate/historical?from=%d&to=%d", from, to)
	return requestAndUnmarshalSlice[*HistoricalExchangeRate](ctx, c, url, http.MethodGet, nil, ErrExchangeRateNotFound)
}
