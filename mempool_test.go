package whatsonchain

import (
	"strings"
	"testing"
)

// TestClient_GetMempoolInfo tests the GetMempoolInfo()
func TestClient_GetMempoolInfo(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *MempoolInfo
	if resp, err = client.GetMempoolInfo(); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.MaxMempool == 0 {
		t.Fatal("possibly failed getting mempool info?")
	}

	// Cannot test other values as they change frequently
}

// TestClient_GetMempoolTransactions tests the GetMempoolTransactions()
func TestClient_GetMempoolTransactions(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var transactions []string
	if transactions, err = client.GetMempoolTransactions(); err != nil {

		// Skip if it takes more than 10 seconds (the request works, but too many Txs)
		if strings.Contains(err.Error(), "deadline exceeded") {
			return
		}
		t.Fatal("error occurred: " + err.Error())
	}

	if len(transactions) == 0 {
		t.Fatal("no transactions found in mempool?")
	}

	// Cannot test other values as they change frequently
}
