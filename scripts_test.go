package whatsonchain

import "testing"

// TestClient_GetScriptHistory tests the GetScriptHistory()
func TestClient_GetScriptHistory(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp ScriptList
	scriptHash := "995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3"
	if resp, err = client.GetScriptHistory(scriptHash); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(resp) == 0 {
		t.Fatal("failed to get history values", resp)
	}

	// todo: this might change! not the best test :P
	if resp[0].TxHash != "52dfceb815ad129a0fd946e3d665f44fa61f068135b9f38b05d3c697e11bad48" {
		t.Fatal("failed to get correct hash", resp, resp[0].TxHash)
	}

}

// TestClient_GetScriptUnspentTransactions tests the GetScriptUnspentTransactions()
func TestClient_GetScriptUnspentTransactions(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	// var resp ScriptList
	scriptHash := "995ea8d0f752f41cdd99bb9d54cb004709e04c7dc4088bcbbbb9ea5c390a43c3"
	if _, err = client.GetScriptUnspentTransactions(scriptHash); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	// todo: this can only be tested if there is a static UTXO for a script
	/*if len(resp) == 0 {
		t.Fatal("failed to get unspent values", resp)
	}*/

	/*if resp[0].TxHash != "52dfceb815ad129a0fd946e3d665f44fa61f068135b9f38b05d3c697e11bad48" {
		t.Fatal("failed to get correct hash", resp, resp[0].TxHash)
	}*/
}
