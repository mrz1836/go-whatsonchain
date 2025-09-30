package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetExchangeRate this endpoint provides exchange rate for BSV/BTC
//
// For more information: https://developers.whatsonchain.com/#get-exchange-rate
func (c *Client) GetExchangeRate(ctx context.Context) (rate *ExchangeRate, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/exchangerate
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/exchangerate", apiEndpointBase, c.Chain(), c.Network()),
		http.MethodGet, nil,
	); err != nil {
		return rate, err
	}
	if len(resp) == 0 {
		return nil, ErrExchangeRateNotFound
	}
	err = json.Unmarshal([]byte(resp), &rate)
	return rate, err
}

// GetHistoricalExchangeRate this endpoint provides historical exchange rates for BSV/BTC
// within a specified time range
//
// For more information: https://developers.whatsonchain.com/#get-historical-exchange-rate
func (c *Client) GetHistoricalExchangeRate(ctx context.Context, from, to int64) (rates []*HistoricalExchangeRate, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/exchangerate/historical?from=<from>&to=<to>
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/%s/exchangerate/historical?from=%d&to=%d", apiEndpointBase, c.Chain(), c.Network(), from, to),
		http.MethodGet, nil,
	); err != nil {
		return rates, err
	}
	if len(resp) == 0 {
		return nil, ErrExchangeRateNotFound
	}
	err = json.Unmarshal([]byte(resp), &rates)
	return rates, err
}
