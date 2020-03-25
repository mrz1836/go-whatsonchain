package whatsonchain

import "testing"

// TestClient_GetExchangeRate tests the GetExchangeRate()
func TestClient_GetExchangeRate(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *ExchangeRate
	if resp, err = client.GetExchangeRate(); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(resp.Currency) == 0 {
		t.Fatal("missing currency", resp.Currency)
	}

	if resp.Currency != "USD" {
		t.Fatal("expected USD - might change in the future?", resp.Currency)
	}

	if len(resp.Rate) == 0 {
		t.Fatal("missing the rate", resp.Rate)
	}
}
