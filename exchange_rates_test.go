package whatsonchain

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// TestClient_GetExchangeRate_NoRate checks that GetExchangeRate never reports a
// missing or non-positive rate as a price (issue #223). Before the fix, a 404 with
// a JSON body, or a 200 with {} or null, returned &ExchangeRate{Rate: 0} and no error.
func TestClient_GetExchangeRate_NoRate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		status   int
		body     string
		wantRate float64 // 0: ErrExchangeRateNotFound is expected
		wantBare bool    // the error must be ErrExchangeRateNotFound itself
	}{
		{"positive rate", http.StatusOK, `{"rate":38.542,"time":1668439893,"currency":"USD"}`, 38.542, false},
		{"tiny positive rate", http.StatusOK, `{"rate":0.0001,"time":1668439893,"currency":"USD"}`, 0.0001, false},
		{"zero rate", http.StatusOK, `{"rate":0,"time":1668439893,"currency":"USD"}`, 0, true},
		{"negative rate", http.StatusOK, `{"rate":-1.5,"time":1668439893,"currency":"USD"}`, 0, true},
		{"empty object", http.StatusOK, `{}`, 0, true},
		{"null", http.StatusOK, `null`, 0, true},
		{"rate missing", http.StatusOK, `{"time":1668439893,"currency":"USD"}`, 0, true},
		{"empty 200 body", http.StatusOK, ``, 0, true},
		{"empty 404 body", http.StatusNotFound, ``, 0, true},
		{"404 with JSON body", http.StatusNotFound, `{"error":"not found"}`, 0, false},
		{"404 with zero rate in body", http.StatusNotFound, `{"rate":0}`, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := newMockClient(&mockHTTPStatusCode{statusCode: tt.status, body: tt.body})
			rate, err := client.GetExchangeRate(context.Background())

			if tt.wantRate > 0 {
				require.NoError(t, err)
				require.NotNil(t, rate)
				assert.InDelta(t, tt.wantRate, rate.Rate, 1e-12)
				assert.Equal(t, "USD", rate.Currency)
				return
			}

			require.ErrorIs(t, err, ErrExchangeRateNotFound)
			require.NotErrorIs(t, err, ErrRequestFailed)
			assert.Nil(t, rate)
			if tt.wantBare {
				assert.Same(t, ErrExchangeRateNotFound, err)
			}
		})
	}
}

// TestClient_GetExchangeRate_Failures checks that real failures are not mistaken
// for a missing rate.
func TestClient_GetExchangeRate_Failures(t *testing.T) {
	t.Parallel()

	t.Run("server error", func(t *testing.T) {
		t.Parallel()

		client := newMockClient(&mockHTTPStatusCode{statusCode: http.StatusInternalServerError, body: "boom"})
		rate, err := client.GetExchangeRate(context.Background())
		require.ErrorIs(t, err, ErrRequestFailed)
		require.NotErrorIs(t, err, ErrExchangeRateNotFound)
		assert.Nil(t, rate)
	})

	t.Run("malformed JSON", func(t *testing.T) {
		t.Parallel()

		client := newMockClient(&mockHTTPStatusCode{statusCode: http.StatusOK, body: "not json"})
		rate, err := client.GetExchangeRate(context.Background())
		require.Error(t, err)
		require.NotErrorIs(t, err, ErrExchangeRateNotFound)
		assert.Nil(t, rate)
	})
}
