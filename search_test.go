package whatsonchain

import "testing"

// TestClient_GetExplorerLinks tests the GetExplorerLinks()
func TestClient_GetExplorerLinks(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Test searching for block hash
	query := "000000000000000002080d0ad78d08691d956d08fb8556339b6dd84fbbfdf1bc"
	var results SearchResults
	if results, err = client.GetExplorerLinks(query); err != nil {
		t.Fatal("error occurred", err.Error())
	}
	if results.Results == nil {
		t.Fatal("expected to get results")
	}
	if results.Results[0].Type != "block" {
		t.Fatal("expected to get back block result", results.Results[0].Type)
	}
	if len(results.Results[0].URL) == 0 {
		t.Fatal("missing url", results.Results[0].URL)
	}

	// Test searching for address
	query = "1GJ3x5bcEnKMnzNFPPELDfXUCwKEaLHM5H"
	if results, err = client.GetExplorerLinks(query); err != nil {
		t.Fatal("error occurred", err.Error())
	}
	if results.Results == nil {
		t.Fatal("expected to get results")
	}
	if results.Results[0].Type != "address" {
		t.Fatal("expected to get address result", results.Results[0].Type)
	}
	if len(results.Results[0].URL) == 0 {
		t.Fatal("missing url", results.Results[0].URL)
	}

	// Test searching for transaction
	query = "6a7c821fd13c5cec773f7e221479651804197866469e92a4d6d47e1fd34d090d"
	if results, err = client.GetExplorerLinks(query); err != nil {
		t.Fatal("error occurred", err.Error())
	}
	if results.Results == nil {
		t.Fatal("expected to get results")
	}
	if results.Results[0].Type != "tx" {
		t.Fatal("expected to get tx result", results.Results[0].Type)
	}
	if len(results.Results[0].URL) == 0 {
		t.Fatal("missing url", results.Results[0].URL)
	}
}
