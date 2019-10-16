package whatsonchain

import (
	"fmt"
	"testing"
)

// TestNewClient test new client
func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	if client.Parameters.Network != NetworkMain {
		t.Fatal("expected value to be main")
	}
}

// ExampleNewClient example using NewClient()
func ExampleNewClient() {
	client, _ := NewClient()
	fmt.Println(client.Parameters.Network)
	// Output:main
}

// BenchmarkNewClient benchmarks the NewClient method
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewClient()
	}
}

// TestClient_GetHealth tests the GetHealth()
func TestClient_GetHealth(t *testing.T) {
	// Skip tis test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	var resp string
	resp, err = client.GetHealth()
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}
	if resp != "Whats On Chain" {
		t.Fatal("expected value was wrong", resp)
	}

}

// ExampleClient_GetHealth example using GetHealth()
func ExampleClient_GetHealth() {
	client, _ := NewClient()
	resp, _ := client.GetHealth()
	fmt.Println(resp)
	// Output:Whats On Chain
}

// TestClient_GetChainInfo tests the GetChainInfo()
func TestClient_GetChainInfo(t *testing.T) {
	// Skip tis test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	var resp *ChainInfo
	resp, err = client.GetChainInfo()
	if err != nil {
		t.Fatal("error occurred: " + err.Error())
	}

	if len(resp.BestBlockHash) == 0 {
		t.Fatal("failed to get best block hash")
	}

}
