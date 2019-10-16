package whatsonchain

import "testing"

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
