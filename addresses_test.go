package whatsonchain

import "testing"

// TestClient_AddressInfo tests the AddressInfo()
func TestClient_AddressInfo(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *AddressInfo
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"
	if resp, err = client.AddressInfo(address); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.Address != address {
		t.Fatal("failed to get the address", address, resp.Address)
	}

}

// TestClient_AddressBalance tests the AddressBalance()
func TestClient_AddressBalance(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	//var resp *AddressBalance
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"
	if _, err = client.AddressBalance(address); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	// todo: not so easy to test for unless you leave a wallet with unspent transactions just for this test ;-)
	/*if resp.Unconfirmed != 0 {
		t.Fatal("failed to get the unconfirmed value", resp.Unconfirmed)
	}

	if resp.Confirmed != 0 {
		t.Fatal("failed to get the confirmed value", resp.Unconfirmed)
	}*/

}

// TestClient_AddressHistory tests the AddressHistory()
func TestClient_AddressHistory(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp AddressHistory
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"
	if resp, err = client.AddressHistory(address); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(resp) == 0 {
		t.Fatal("failed to get history values", resp)
	}

	// todo: this might change! not the best test :P
	if resp[0].TxHash != "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1" {
		t.Fatal("failed to get correct hash", resp, resp[0].TxHash)
	}

}

// TestClient_AddressUnspentTransactions tests the AddressUnspentTransactions()
func TestClient_AddressUnspentTransactions(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	//var resp AddressHistory
	address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"
	if _, err = client.AddressUnspentTransactions(address); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	// todo: not so easy to test for unless you leave a wallet with unspent transactions just for this test ;-)
	/*if len(resp) == 0 {
		t.Fatal("failed to get history values", resp)
	}*/

}
