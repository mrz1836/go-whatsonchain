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

	// Valid (exchange rate)
	if strings.Contains(req.URL.String(), "/exchangerate") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBuffer([]byte(`{"rate":38.542,"time":1668439893,"currency":"USD"}`)))
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
