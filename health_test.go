package whatsonchain

import (
	"fmt"
	"testing"
)

// TestClient_GetHealth tests the GetHealth()
func TestClient_GetHealth(t *testing.T) {
	// Skip this test in short mode (not needed)
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// Create a new client object to handle your queries (supply an API Key)
	client, err := NewClient(NetworkMain, nil)
	if err != nil {
		t.Fatal(err)
	}

	var resp string
	if resp, err = client.GetHealth(); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}
	if resp != "Whats On Chain" {
		t.Fatal("expected value was wrong", resp)
	}

}

// ExampleClient_GetHealth example using GetHealth()
func ExampleClient_GetHealth() {
	client, _ := NewClient(NetworkMain, nil)
	resp, _ := client.GetHealth()
	fmt.Println(resp)
	// Output:Whats On Chain
}
