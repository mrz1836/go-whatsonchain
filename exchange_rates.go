package whatsonchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetExchangeRate this endpoint provides exchange rate for BSV
//
// For more information: https://developers.whatsonchain.com/#get-exchange-rate
func (c *Client) GetExchangeRate(ctx context.Context) (rate *ExchangeRate, err error) {
	var resp string
	// https://api.whatsonchain.com/v1/bsv/<network>/exchangerate
	if resp, err = c.request(
		ctx,
		fmt.Sprintf("%s%s/exchangerate", apiEndpoint, c.Network()),
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
