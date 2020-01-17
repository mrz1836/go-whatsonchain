package whatsonchain

import (
	"fmt"
	"log"
	"testing"
)

// TestClient_GetChainInfo tests the GetChainInfo()
func TestClient_GetChainInfo(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp *ChainInfo
	if resp, err = client.GetChainInfo(); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(resp.BestBlockHash) == 0 {
		t.Fatal("failed to get best block hash")
	}

}

// ExampleClient_GetChainInfo example using GetChainInfo()
func ExampleClient_GetChainInfo() {
	client, _ := NewClient(NetworkMain, nil)
	resp, _ := client.GetChainInfo()
	log.Println(resp.BestBlockHash)
	fmt.Println("0000000000000000057d09c9d9928c53aaff1f6b019ead3ceed52aca8abbc1c9")
	// Output:0000000000000000057d09c9d9928c53aaff1f6b019ead3ceed52aca8abbc1c9
}
