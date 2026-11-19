package whatsonchain

import (
	"context"
	"net/http"
)

// GetExchangeRate this endpoint provides the exchange rate for BSV.
//
// It returns ErrExchangeRateNotFound when the API has no rate: a 404, an empty
// body, or a body without a positive rate (such as {} or null).
//
// For more information: https://docs.whatsonchain.com/#get-exchange-rate
func (c *Client) GetExchangeRate(ctx context.Context) (*ExchangeRate, error) {
	url := c.buildURL("/exchangerate")
	rate, err := requestAndUnmarshal[ExchangeRate](ctx, c, url, http.MethodGet, nil, ErrExchangeRateNotFound)
	if err != nil {
		return nil, err
	}

	// A body without a usable rate decodes to the zero value; never report it as a price
	if rate.Rate <= 0 {
		return nil, ErrExchangeRateNotFound
	}

	return rate, nil
}

// GetHistoricalExchangeRate this endpoint provides historical exchange rates for BSV
// within a specified time range
//
// For more information: https://docs.whatsonchain.com/#get-historical-exchange-rate
func (c *Client) GetHistoricalExchangeRate(ctx context.Context, from, to int64) ([]*HistoricalExchangeRate, error) {
	url := c.buildURL("/exchangerate/historical?from=%d&to=%d", from, to)
	return requestAndUnmarshalSlice[*HistoricalExchangeRate](ctx, c, url, http.MethodGet, nil, ErrExchangeRateNotFound)
}
