package whatsonchain

import "testing"

// TestClient_GetBlockByHash tests the GetBlockByHash()
func TestClient_GetBlockByHash(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *BlockInfo
	hash := "0000000000000000025b8506c83450afe84f0318775a52c7b91ee64aad0d5a23"
	if resp, err = client.GetBlockByHash(hash); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.Hash != "0000000000000000025b8506c83450afe84f0318775a52c7b91ee64aad0d5a23" {
		t.Fatal("failed to get the block hash", resp.Hash)
	}

}

// TestClient_GetBlockByHeight tests the GetBlockByHeight()
func TestClient_GetBlockByHeight(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *BlockInfo
	height := int64(604648)
	if resp, err = client.GetBlockByHeight(height); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if resp.Hash != "0000000000000000025b8506c83450afe84f0318775a52c7b91ee64aad0d5a23" {
		t.Fatal("failed to get the block hash", resp.Hash)
	}

}

// TestClient_GetBlockPages tests the GetBlockPages()
func TestClient_GetBlockPages(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *BlockPagesInfo
	hash := "000000000000000000885a4d8e9912f085b42288adc58b3ee5830a7da9f4fef4"
	if resp, err = client.GetBlockPages(hash, 1); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(*resp) == 0 {
		t.Fatal("no transactions found", resp)
	}

}
