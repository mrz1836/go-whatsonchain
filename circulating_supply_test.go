package whatsonchain

import (
	"fmt"
	"log"
	"testing"
)

// TestClient_GetCirculatingSupply tests the GetCirculatingSupply()
func TestClient_GetCirculatingSupply(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var supply float64
	if supply, err = client.GetCirculatingSupply(); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if supply <= 0 {
		t.Fatal("failed to get circulating supply")
	}

}

// ExampleClient_GetCirculatingSupply example using GetCirculatingSupply()
func ExampleClient_GetCirculatingSupply() {
	client, _ := NewClient(NetworkMain, nil)
	supply, _ := client.GetCirculatingSupply()
	log.Printf("%f", supply)
	fmt.Println("18225787.5")
	// Output:18225787.5
}
