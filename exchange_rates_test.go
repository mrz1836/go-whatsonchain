package whatsonchain

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPExchangeValid for mocking requests
type mockHTTPExchangeValid struct{}

// Do is a mock http request
func (m *mockHTTPExchangeValid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	// Valid (historical exchange rate)
	if strings.Contains(req.URL.String(), "/exchangerate/historical") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`[{"rate":38.542,"time":1660139745,"currency":"USD"},{"rate":39.123,"time":1660312545,"currency":"USD"}]`))
	}
	// Valid (exchange rate)
	if strings.Contains(req.URL.String(), "/exchangerate") && !strings.Contains(req.URL.String(), "/historical") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`{"rate":38.542,"time":1668439893,"currency":"USD"}`))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPExchangeInvalid for mocking requests
type mockHTTPExchangeInvalid struct{}

// Do is a mock http request
func (m *mockHTTPExchangeInvalid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	// Invalid (exchange rate)
	if strings.Contains(req.URL.String(), "/exchangerate") {
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, ErrBadRequest
	}

	// Default is valid
	return resp, nil
}

// mockHTTPExchangeNotFound for mocking requests
type mockHTTPExchangeNotFound struct{}

// Do is a mock http request
func (m *mockHTTPExchangeNotFound) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusNotFound

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	// Invalid (exchange rate)
	if strings.Contains(req.URL.String(), "/exchangerate") {
		resp.Body = io.NopCloser(bytes.NewBufferString(""))
		return resp, nil
	}

	// Default is valid
	return resp, nil
}

// TestClient_GetExchangeRate tests the GetExchangeRate()
func TestClient_GetExchangeRate(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPExchangeValid{})
	ctx := context.Background()

	// Test the valid response
	info, err := client.GetExchangeRate(ctx)
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if info == nil {
		t.Errorf("%s Failed: info was nil", t.Name())
	} else if info.Currency != "USD" {
		t.Errorf("%s Failed: currency was [%s] expected [%s]", t.Name(), info.Currency, "USD")
	} else if info.Rate != 38.542 {
		t.Errorf("%s Failed: currency was [%v] expected [%v]", t.Name(), info.Rate, 38.542)
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPExchangeInvalid{})

	// Test invalid response
	_, err = client.GetExchangeRate(ctx)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}

	// New not found mock client
	client = newMockClient(&mockHTTPExchangeNotFound{})

	// Test invalid response
	_, err = client.GetExchangeRate(ctx)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}

// TestClient_GetHistoricalExchangeRate tests the GetHistoricalExchangeRate()
func TestClient_GetHistoricalExchangeRate(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPExchangeValid{})
	ctx := context.Background()
	from := int64(1660139745)
	to := int64(1660312545)

	// Test the valid response
	rates, err := client.GetHistoricalExchangeRate(ctx, from, to)
	if err != nil {
		t.Errorf("%s Failed: error [%s]", t.Name(), err.Error())
	} else if rates == nil {
		t.Errorf("%s Failed: rates was nil", t.Name())
	} else if len(rates) != 2 {
		t.Errorf("%s Failed: expected 2 rates, got [%d]", t.Name(), len(rates))
	} else if rates[0].Currency != "USD" {
		t.Errorf("%s Failed: first rate currency was [%s] expected [%s]", t.Name(), rates[0].Currency, "USD")
	} else if rates[0].Rate != 38.542 {
		t.Errorf("%s Failed: first rate was [%v] expected [%v]", t.Name(), rates[0].Rate, 38.542)
	} else if rates[0].Time != 1660139745 {
		t.Errorf("%s Failed: first rate time was [%d] expected [%d]", t.Name(), rates[0].Time, 1660139745)
	} else if rates[1].Rate != 39.123 {
		t.Errorf("%s Failed: second rate was [%v] expected [%v]", t.Name(), rates[1].Rate, 39.123)
	}

	// New invalid mock client
	client = newMockClient(&mockHTTPExchangeInvalid{})

	// Test invalid response
	_, err = client.GetHistoricalExchangeRate(ctx, from, to)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}

	// New not found mock client
	client = newMockClient(&mockHTTPExchangeNotFound{})

	// Test invalid response
	_, err = client.GetHistoricalExchangeRate(ctx, from, to)
	if err == nil {
		t.Errorf("%s Failed: error should have occurred", t.Name())
	}
}
