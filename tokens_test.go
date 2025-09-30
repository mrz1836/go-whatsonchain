package whatsonchain

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPTokensValid for mocking valid token requests
type mockHTTPTokensValid struct{}

// Do is a mock http request
func (m *mockHTTPTokensValid) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, ErrMissingRequest
	}

	// 1Sat Ordinals endpoints
	if strings.Contains(req.URL.String(), "/bsv/") && strings.Contains(req.URL.String(), "/token/1satordinals/") {
		// Check for "invalid" to trigger 404
		if strings.Contains(req.URL.String(), "/invalid/") {
			resp.StatusCode = http.StatusNotFound
			resp.Body = io.NopCloser(bytes.NewBufferString(`{"error":"Not found"}`))
			return resp, nil
		}

		resp.StatusCode = http.StatusOK

		if strings.Contains(req.URL.String(), "/origin") {
			// Mock GetOneSatOrdinalByOrigin
			mockToken := OneSatOrdinalToken{
				Outpoint: "test_outpoint",
				Origin:   "827748:753:0",
				Height:   123456,
				Idx:      1,
				Lock:     "test_lock",
				Spend:    "test_spend",
				Data:     "test_data",
				Listing:  false,
				Bsv20:    true,
			}
			jsonBytes, err := json.Marshal(mockToken)
			if err != nil {
				jsonBytes = []byte("{}")
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		} else if strings.Contains(req.URL.String(), "/content") {
			// Mock GetOneSatOrdinalContent
			mockContent := OneSatOrdinalContent{
				Content: []byte("test content"),
				Type:    "text/plain",
			}
			jsonBytes, err := json.Marshal(mockContent)
			if err != nil {
				jsonBytes = []byte("{}")
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		} else if strings.Contains(req.URL.String(), "/latest") {
			// Mock GetOneSatOrdinalLatest
			mockLatest := OneSatOrdinalLatest{
				TxID:   "test_txid",
				Vout:   0,
				Height: 123456,
				Idx:    1,
			}
			jsonBytes, err := json.Marshal(mockLatest)
			if err != nil {
				jsonBytes = []byte("{}")
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		} else if strings.Contains(req.URL.String(), "/history") {
			// Mock GetOneSatOrdinalHistory
			mockHistory := []*OneSatOrdinalHistory{
				{
					TxID:   "test_txid1",
					Vout:   0,
					Height: 123456,
					Idx:    1,
				},
				{
					TxID:   "test_txid2",
					Vout:   1,
					Height: 123457,
					Idx:    2,
				},
			}
			jsonBytes, err := json.Marshal(mockHistory)
			if err != nil {
				jsonBytes = []byte("[]")
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		} else if strings.Contains(req.URL.String(), "/tx/") {
			// Mock GetOneSatOrdinalsByTxID
			mockTokens := []*OneSatOrdinalToken{
				{
					Outpoint: "test_outpoint1",
					Origin:   "827748:753:0",
					Height:   123456,
					Idx:      1,
				},
				{
					Outpoint: "test_outpoint2",
					Origin:   "827748:753:1",
					Height:   123456,
					Idx:      2,
				},
			}
			jsonBytes, err := json.Marshal(mockTokens)
			if err != nil {
				jsonBytes = []byte("[]")
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		} else {
			// Mock GetOneSatOrdinalByOutpoint
			mockToken := OneSatOrdinalToken{
				Outpoint: "test_outpoint",
				Origin:   "827748:753:0",
				Height:   123456,
				Idx:      1,
				Lock:     "test_lock",
				Spend:    "test_spend",
				Data:     "test_data",
				Listing:  false,
				Bsv20:    true,
			}
			jsonBytes, err := json.Marshal(mockToken)
			if err != nil {
				jsonBytes = []byte("{}")
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		}
	}

	// 1Sat Ordinals stats
	if strings.Contains(req.URL.String(), "/bsv/") && strings.Contains(req.URL.String(), "/tokens/1satordinals") {
		resp.StatusCode = http.StatusOK
		mockStats := OneSatOrdinalStats{
			Pending:   10,
			Confirmed: 1000,
		}
		jsonBytes, err := json.Marshal(mockStats)
		if err != nil {
			jsonBytes = []byte("{}")
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	}

	// STAS token endpoints
	if strings.Contains(req.URL.String(), "/bsv/") && strings.Contains(req.URL.String(), "/tokens") && !strings.Contains(req.URL.String(), "1satordinals") && !strings.Contains(req.URL.String(), "stas") {
		resp.StatusCode = http.StatusOK
		// Mock GetAllSTASTokens
		mockTokens := []*STASToken{
			{
				ContractID:        "test_contract_id",
				Symbol:            "TEST",
				IssuerPK:          "test_issuer_pk",
				IsZeroSupplyToken: false,
				ProtocolID:        "STAS",
				TotalSupply:       1000000,
				CirculatingSupply: 500000,
				DecimalPrecision:  8,
				Name:              "Test Token",
				Description:       "A test token",
				TokenType:         "utility",
				SatsPerToken:      1,
				LifeCycleComplete: false,
			},
		}
		jsonBytes, err := json.Marshal(mockTokens)
		if err != nil {
			jsonBytes = []byte("[]")
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	}

	// STAS token by ID
	if strings.Contains(req.URL.String(), "/bsv/") && strings.Contains(req.URL.String(), "/token/") && !strings.Contains(req.URL.String(), "1satordinals") && !strings.Contains(req.URL.String(), "/tx") {
		resp.StatusCode = http.StatusOK
		mockToken := STASToken{
			ContractID:        "test_contract_id",
			Symbol:            "TEST",
			IssuerPK:          "test_issuer_pk",
			IsZeroSupplyToken: false,
			ProtocolID:        "STAS",
			TotalSupply:       1000000,
			CirculatingSupply: 500000,
			DecimalPrecision:  8,
			Name:              "Test Token",
			Description:       "A test token",
			TokenType:         "utility",
			SatsPerToken:      1,
			LifeCycleComplete: false,
		}
		jsonBytes, err := json.Marshal(mockToken)
		if err != nil {
			jsonBytes = []byte("{}")
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	}

	// Token UTXOs for address
	if strings.Contains(req.URL.String(), "/bsv/") && strings.Contains(req.URL.String(), "/address/") && strings.Contains(req.URL.String(), "/tokens/unspent") {
		resp.StatusCode = http.StatusOK
		mockUTXOs := []*STASTokenUTXO{
			{
				TxID:       "test_txid",
				Vout:       0,
				Amount:     1000,
				Script:     "test_script",
				ContractID: "test_contract_id",
				Symbol:     "TEST",
				Value:      1000,
				Height:     123456,
			},
		}
		jsonBytes, err := json.Marshal(mockUTXOs)
		if err != nil {
			jsonBytes = []byte("[]")
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	}

	// Address token balance
	if strings.Contains(req.URL.String(), "/bsv/") && strings.Contains(req.URL.String(), "/address/") && strings.Contains(req.URL.String(), "/tokens") && !strings.Contains(req.URL.String(), "/unspent") {
		resp.StatusCode = http.StatusOK
		mockBalance := STASTokenBalance{
			Address: "test_address",
			Tokens: []STASTokenBalanceInfo{
				{
					ContractID: "test_contract_id",
					Symbol:     "TEST",
					Balance:    1000,
					Decimal:    8,
				},
			},
		}
		jsonBytes, err := json.Marshal(mockBalance)
		if err != nil {
			jsonBytes = []byte("{}")
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	}

	// Token transactions (exclude 1satordinals which has its own endpoint)
	if strings.Contains(req.URL.String(), "/bsv/") && strings.Contains(req.URL.String(), "/token/") && strings.Contains(req.URL.String(), "/tx") && !strings.Contains(req.URL.String(), "1satordinals") {
		resp.StatusCode = http.StatusOK
		mockTxList := TxList{
			{
				TxID:        "test_txid1",
				BlockHeight: 123456,
			},
			{
				TxID:        "test_txid2",
				BlockHeight: 123457,
			},
		}
		jsonBytes, err := json.Marshal(mockTxList)
		if err != nil {
			jsonBytes = []byte("[]")
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	}

	// STAS stats
	if strings.Contains(req.URL.String(), "/bsv/") && strings.Contains(req.URL.String(), "/tokens/stas") {
		resp.StatusCode = http.StatusOK
		mockStats := STASStats{
			Tokens:      100,
			Issuers:     50,
			TotalSupply: 1000000000,
		}
		jsonBytes, err := json.Marshal(mockStats)
		if err != nil {
			jsonBytes = []byte("{}")
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	}

	return resp, nil
}

// TestClient_GetOneSatOrdinalByOrigin tests the GetOneSatOrdinalByOrigin method
func TestClient_GetOneSatOrdinalByOrigin(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetOneSatOrdinalByOrigin(context.Background(), "827748:753:0")
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	token, err := client.GetOneSatOrdinalByOrigin(context.Background(), "827748:753:0")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if token.Origin != "827748:753:0" {
		t.Fatalf("expected origin: 827748:753:0, got: %s", token.Origin)
	}

	// Test 404 not found
	client = newMockClientBSV(&mockHTTPTokensValid{})
	_, err = client.GetOneSatOrdinalByOrigin(context.Background(), "invalid")
	if !errors.Is(err, ErrTokenNotFound) {
		t.Fatalf("expected ErrTokenNotFound, got: %v", err)
	}
}

// TestClient_GetOneSatOrdinalByOutpoint tests the GetOneSatOrdinalByOutpoint method
func TestClient_GetOneSatOrdinalByOutpoint(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetOneSatOrdinalByOutpoint(context.Background(), "test_outpoint")
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	token, err := client.GetOneSatOrdinalByOutpoint(context.Background(), "test_outpoint")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if token.Outpoint != "test_outpoint" {
		t.Fatalf("expected outpoint: test_outpoint, got: %s", token.Outpoint)
	}
}

// TestClient_GetOneSatOrdinalContent tests the GetOneSatOrdinalContent method
func TestClient_GetOneSatOrdinalContent(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetOneSatOrdinalContent(context.Background(), "test_outpoint")
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	content, err := client.GetOneSatOrdinalContent(context.Background(), "test_outpoint")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if content.Type != "text/plain" {
		t.Fatalf("expected type: text/plain, got: %s", content.Type)
	}
}

// TestClient_GetOneSatOrdinalLatest tests the GetOneSatOrdinalLatest method
func TestClient_GetOneSatOrdinalLatest(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetOneSatOrdinalLatest(context.Background(), "test_outpoint")
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	latest, err := client.GetOneSatOrdinalLatest(context.Background(), "test_outpoint")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if latest.TxID != "test_txid" {
		t.Fatalf("expected txid: test_txid, got: %s", latest.TxID)
	}
}

// TestClient_GetOneSatOrdinalHistory tests the GetOneSatOrdinalHistory method
func TestClient_GetOneSatOrdinalHistory(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetOneSatOrdinalHistory(context.Background(), "test_outpoint")
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	history, err := client.GetOneSatOrdinalHistory(context.Background(), "test_outpoint")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if len(history) != 2 {
		t.Fatalf("expected 2 history records, got: %d", len(history))
	}

	if history[0].TxID != "test_txid1" {
		t.Fatalf("expected first txid: test_txid1, got: %s", history[0].TxID)
	}
}

// TestClient_GetOneSatOrdinalsByTxID tests the GetOneSatOrdinalsByTxID method
func TestClient_GetOneSatOrdinalsByTxID(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetOneSatOrdinalsByTxID(context.Background(), "test_txid")
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	tokens, err := client.GetOneSatOrdinalsByTxID(context.Background(), "test_txid")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got: %d", len(tokens))
	}

	if tokens[0].Outpoint != "test_outpoint1" {
		t.Fatalf("expected first outpoint: test_outpoint1, got: %s", tokens[0].Outpoint)
	}
}

// TestClient_GetOneSatOrdinalsStats tests the GetOneSatOrdinalsStats method
func TestClient_GetOneSatOrdinalsStats(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetOneSatOrdinalsStats(context.Background())
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	stats, err := client.GetOneSatOrdinalsStats(context.Background())
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if stats.Pending != 10 {
		t.Fatalf("expected pending: 10, got: %d", stats.Pending)
	}

	if stats.Confirmed != 1000 {
		t.Fatalf("expected confirmed: 1000, got: %d", stats.Confirmed)
	}
}

// TestClient_GetAllSTASTokens tests the GetAllSTASTokens method
func TestClient_GetAllSTASTokens(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetAllSTASTokens(context.Background())
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	tokens, err := client.GetAllSTASTokens(context.Background())
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got: %d", len(tokens))
	}

	if tokens[0].Symbol != "TEST" {
		t.Fatalf("expected symbol: TEST, got: %s", tokens[0].Symbol)
	}
}

// TestClient_GetSTASTokenByID tests the GetSTASTokenByID method
func TestClient_GetSTASTokenByID(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetSTASTokenByID(context.Background(), "test_contract", "TEST")
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	token, err := client.GetSTASTokenByID(context.Background(), "test_contract", "TEST")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if token.Symbol != "TEST" {
		t.Fatalf("expected symbol: TEST, got: %s", token.Symbol)
	}
}

// TestClient_GetTokenUTXOsForAddress tests the GetTokenUTXOsForAddress method
func TestClient_GetTokenUTXOsForAddress(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetTokenUTXOsForAddress(context.Background(), "test_address")
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	utxos, err := client.GetTokenUTXOsForAddress(context.Background(), "test_address")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if len(utxos) != 1 {
		t.Fatalf("expected 1 UTXO, got: %d", len(utxos))
	}

	if utxos[0].Symbol != "TEST" {
		t.Fatalf("expected symbol: TEST, got: %s", utxos[0].Symbol)
	}
}

// TestClient_GetAddressTokenBalance tests the GetAddressTokenBalance method
func TestClient_GetAddressTokenBalance(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetAddressTokenBalance(context.Background(), "test_address")
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	balance, err := client.GetAddressTokenBalance(context.Background(), "test_address")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if balance.Address != "test_address" {
		t.Fatalf("expected address: test_address, got: %s", balance.Address)
	}

	if len(balance.Tokens) != 1 {
		t.Fatalf("expected 1 token balance, got: %d", len(balance.Tokens))
	}
}

// TestClient_GetTokenTransactions tests the GetTokenTransactions method
func TestClient_GetTokenTransactions(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetTokenTransactions(context.Background(), "test_contract", "TEST")
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	transactions, err := client.GetTokenTransactions(context.Background(), "test_contract", "TEST")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if len(transactions) != 2 {
		t.Fatalf("expected 2 transactions, got: %d", len(transactions))
	}

	if transactions[0].TxID != "test_txid1" {
		t.Fatalf("expected first txid: test_txid1, got: %s", transactions[0].TxID)
	}
}

// TestClient_GetSTASStats tests the GetSTASStats method
func TestClient_GetSTASStats(t *testing.T) {
	// Test BSV chain requirement
	btcClient, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain))
	if err != nil {
		t.Fatal(err)
	}
	_, err = btcClient.GetSTASStats(context.Background())
	if !errors.Is(err, ErrBSVChainRequired) {
		t.Fatalf("expected BSV chain required error, got: %v", err)
	}

	// Test valid request
	client := newMockClientBSV(&mockHTTPTokensValid{})
	stats, err := client.GetSTASStats(context.Background())
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	if stats.Tokens != 100 {
		t.Fatalf("expected tokens: 100, got: %d", stats.Tokens)
	}

	if stats.Issuers != 50 {
		t.Fatalf("expected issuers: 50, got: %d", stats.Issuers)
	}

	if stats.TotalSupply != 1000000000 {
		t.Fatalf("expected total supply: 1000000000, got: %d", stats.TotalSupply)
	}
}
