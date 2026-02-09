package whatsonchain

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// TestClient_GetBlockStats tests the method GetBlockStats()
func TestClient_GetBlockStats(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		mockClient := &mockHTTPValidChain{}
		client, err := NewClient(context.Background(), WithChain(ChainBSV), WithNetwork(NetworkMain), WithHTTPClient(mockClient))
		if err != nil {
			t.Fatal(err)
		}
		mockClient.SetResponse(func(_ *http.Request) (*http.Response, error) {
			resp := newHTTPResponse(`{"height":698730,"hash":"000000000000000002799ba826646d0060b06779e7bde9622145b410f114c1fb","version":536870912,"size":1234567,"weight":1234567,"merkleroot":"abcd1234","timestamp":1609459200,"mediantime":1609459100,"nonce":123456789,"bits":"207fffff","difficulty":1.0,"chainwork":"00000000000000000000000000000000000000000000000000000000deadbeef","tx_count":100,"total_size":1234567,"total_fees":50000,"subsidy_total":625000000,"subsidy_address":0,"subsidy_miner":625000000,"miner_name":"Test Miner","miner_address":"1TestMinerAddress","fee_rate_avg":10.5,"fee_rate_min":1.0,"fee_rate_max":50.0,"fee_rate_median":8.0,"fee_rate_stddev":5.2,"input_count":200,"output_count":150,"utxo_increase":50,"utxo_size_inc":2500}`)
			resp.StatusCode = http.StatusOK
			return resp, nil
		})

		blockStats, err := client.GetBlockStats(context.Background(), 698730)
		if err != nil {
			t.Errorf("%s Failed: error [%s] inputted", t.Name(), err.Error())
		} else if blockStats == nil {
			t.Errorf("%s Failed: blockStats was nil", t.Name())
		} else if blockStats.Height != 698730 {
			t.Errorf("%s Failed: expected height [%d] got [%d]", t.Name(), 698730, blockStats.Height)
		}
	})

	t.Run("http error", func(t *testing.T) {
		client, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain), WithHTTPClient(&mockHTTPError{}))
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.GetBlockStats(context.Background(), 698730)
		if err == nil {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), "698730")
		}
	})
}

// TestClient_GetBlockStatsByHash tests the method GetBlockStatsByHash()
func TestClient_GetBlockStatsByHash(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		mockClient := &mockHTTPValidChain{}
		client, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain), WithHTTPClient(mockClient))
		if err != nil {
			t.Fatal(err)
		}
		mockClient.SetResponse(func(_ *http.Request) (*http.Response, error) {
			resp := newHTTPResponse(`{"height":698730,"hash":"000000000000000002799ba826646d0060b06779e7bde9622145b410f114c1fb","version":536870912,"size":1234567,"weight":1234567,"merkleroot":"abcd1234","timestamp":1609459200,"mediantime":1609459100,"nonce":123456789,"bits":"207fffff","difficulty":1.0,"chainwork":"00000000000000000000000000000000000000000000000000000000deadbeef","tx_count":100,"total_size":1234567,"total_fees":50000,"subsidy_total":625000000,"subsidy_address":0,"subsidy_miner":625000000,"miner_name":"Test Miner","miner_address":"1TestMinerAddress","fee_rate_avg":10.5,"fee_rate_min":1.0,"fee_rate_max":50.0,"fee_rate_median":8.0,"fee_rate_stddev":5.2,"input_count":200,"output_count":150,"utxo_increase":50,"utxo_size_inc":2500}`)
			resp.StatusCode = http.StatusOK
			return resp, nil
		})

		blockStats, err := client.GetBlockStatsByHash(context.Background(), "000000000000000002799ba826646d0060b06779e7bde9622145b410f114c1fb")
		if err != nil {
			t.Errorf("%s Failed: error [%s] inputted", t.Name(), err.Error())
		} else if blockStats == nil {
			t.Errorf("%s Failed: blockStats was nil", t.Name())
		} else if blockStats.Hash != "000000000000000002799ba826646d0060b06779e7bde9622145b410f114c1fb" {
			t.Errorf("%s Failed: expected hash [%s] got [%s]", t.Name(), "000000000000000002799ba826646d0060b06779e7bde9622145b410f114c1fb", blockStats.Hash)
		}
	})

	t.Run("http error", func(t *testing.T) {
		client, err := NewClient(context.Background(), WithChain(ChainBSV), WithNetwork(NetworkMain), WithHTTPClient(&mockHTTPError{}))
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.GetBlockStatsByHash(context.Background(), "invalid-hash")
		if err == nil {
			t.Errorf("%s Failed: expected to throw an error, no error inputted", t.Name())
		}
	})
}

// TestClient_GetMinerBlocksStats tests the method GetMinerBlocksStats()
func TestClient_GetMinerBlocksStats(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		mockClient := &mockHTTPValidChain{}
		client, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain), WithHTTPClient(mockClient))
		if err != nil {
			t.Fatal(err)
		}
		mockClient.SetResponse(func(_ *http.Request) (*http.Response, error) {
			resp := newHTTPResponse(`[{"name":"Test Miner 1","address":"1TestMiner1Address","block_count":10,"percentage":50.0},{"name":"Test Miner 2","address":"1TestMiner2Address","block_count":8,"percentage":40.0}]`)
			resp.StatusCode = http.StatusOK
			return resp, nil
		})

		minerStats, err := client.GetMinerBlocksStats(context.Background(), 1)
		if err != nil {
			t.Errorf("%s Failed: error [%s] inputted", t.Name(), err.Error())
		} else if minerStats == nil {
			t.Errorf("%s Failed: minerStats was nil", t.Name())
		} else if len(minerStats) != 2 {
			t.Errorf("%s Failed: expected 2 miners, got [%d]", t.Name(), len(minerStats))
		} else if minerStats[0].Name != "Test Miner 1" {
			t.Errorf("%s Failed: expected name [%s] got [%s]", t.Name(), "Test Miner 1", minerStats[0].Name)
		}
	})

	t.Run("http error", func(t *testing.T) {
		client, err := NewClient(context.Background(), WithChain(ChainBSV), WithNetwork(NetworkMain), WithHTTPClient(&mockHTTPError{}))
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.GetMinerBlocksStats(context.Background(), 1)
		if err == nil {
			t.Errorf("%s Failed: expected to throw an error, no error inputted", t.Name())
		}
	})
}

// TestClient_GetMinerFeesStats tests the method GetMinerFeesStats()
func TestClient_GetMinerFeesStats(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		mockClient := &mockHTTPValidChain{}
		client, err := NewClient(context.Background(), WithChain(ChainBSV), WithNetwork(NetworkMain), WithHTTPClient(mockClient))
		if err != nil {
			t.Fatal(err)
		}
		mockClient.SetResponse(func(_ *http.Request) (*http.Response, error) {
			resp := newHTTPResponse(`[{"miner":"Test Miner","min_fee_rate":10.5}]`)
			resp.StatusCode = http.StatusOK
			return resp, nil
		})

		feeStats, err := client.GetMinerFeesStats(context.Background(), 1714608000, 1714653060)
		if err != nil {
			t.Errorf("%s Failed: error [%s] inputted", t.Name(), err.Error())
		} else if feeStats == nil {
			t.Errorf("%s Failed: feeStats was nil", t.Name())
		} else if len(feeStats) != 1 {
			t.Errorf("%s Failed: expected 1 fee stat, got [%d]", t.Name(), len(feeStats))
		} else if feeStats[0].Miner != "Test Miner" {
			t.Errorf("%s Failed: expected miner [%s] got [%s]", t.Name(), "Test Miner", feeStats[0].Miner)
		} else if feeStats[0].MinFeeRate != 10.5 {
			t.Errorf("%s Failed: expected min fee rate [%f] got [%f]", t.Name(), 10.5, feeStats[0].MinFeeRate)
		}
	})

	t.Run("http error", func(t *testing.T) {
		client, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain), WithHTTPClient(&mockHTTPError{}))
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.GetMinerFeesStats(context.Background(), 1714608000, 1714653060)
		if err == nil {
			t.Errorf("%s Failed: expected to throw an error, no error inputted", t.Name())
		}
	})
}

// TestClient_GetMinerSummaryStats tests the method GetMinerSummaryStats()
func TestClient_GetMinerSummaryStats(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		mockClient := &mockHTTPValidChain{}
		client, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain), WithHTTPClient(mockClient))
		if err != nil {
			t.Fatal(err)
		}
		mockClient.SetResponse(func(_ *http.Request) (*http.Response, error) {
			resp := newHTTPResponse(`{"days":90,"miners":[{"name":"Test Miner 1","address":"1TestMiner1Address","block_count":100,"percentage":60.0},{"name":"Test Miner 2","address":"1TestMiner2Address","block_count":66,"percentage":40.0}]}`)
			resp.StatusCode = http.StatusOK
			return resp, nil
		})

		summaryStats, err := client.GetMinerSummaryStats(context.Background(), 90)
		if err != nil {
			t.Errorf("%s Failed: error [%s] inputted", t.Name(), err.Error())
		} else if summaryStats == nil {
			t.Errorf("%s Failed: summaryStats was nil", t.Name())
		} else if summaryStats.Days != 90 {
			t.Errorf("%s Failed: expected days [%d] got [%d]", t.Name(), 90, summaryStats.Days)
		} else if len(summaryStats.Miners) != 2 {
			t.Errorf("%s Failed: expected 2 miners, got [%d]", t.Name(), len(summaryStats.Miners))
		}
	})

	t.Run("http error", func(t *testing.T) {
		client, err := NewClient(context.Background(), WithChain(ChainBSV), WithNetwork(NetworkMain), WithHTTPClient(&mockHTTPError{}))
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.GetMinerSummaryStats(context.Background(), 90)
		if err == nil {
			t.Errorf("%s Failed: expected to throw an error, no error inputted", t.Name())
		}
	})
}

// TestClient_GetTagCountByHeight tests the method GetTagCountByHeight()
func TestClient_GetTagCountByHeight(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		mockClient := &mockHTTPValidChain{}
		client, err := NewClient(context.Background(), WithChain(ChainBSV), WithNetwork(NetworkMain), WithHTTPClient(mockClient))
		if err != nil {
			t.Fatal(err)
		}
		mockClient.SetResponse(func(_ *http.Request) (*http.Response, error) {
			resp := newHTTPResponse(`{"height":762291,"hash":"000000000000000000abcdef1234567890abcdef1234567890abcdef12345678","tag_counts":{"tag1":10,"tag2":5,"tag3":15}}`)
			resp.StatusCode = http.StatusOK
			return resp, nil
		})

		tagCount, err := client.GetTagCountByHeight(context.Background(), 762291)
		if err != nil {
			t.Errorf("%s Failed: error [%s] inputted", t.Name(), err.Error())
		} else if tagCount == nil {
			t.Errorf("%s Failed: tagCount was nil", t.Name())
		} else if tagCount.Height != 762291 {
			t.Errorf("%s Failed: expected height [%d] got [%d]", t.Name(), 762291, tagCount.Height)
		} else if len(tagCount.TagCounts) != 3 {
			t.Errorf("%s Failed: expected 3 tag counts, got [%d]", t.Name(), len(tagCount.TagCounts))
		} else if tagCount.TagCounts["tag1"] != 10 {
			t.Errorf("%s Failed: expected tag1 count [%d] got [%d]", t.Name(), 10, tagCount.TagCounts["tag1"])
		}
	})

	t.Run("http error", func(t *testing.T) {
		client, err := NewClient(context.Background(), WithChain(ChainBTC), WithNetwork(NetworkMain), WithHTTPClient(&mockHTTPError{}))
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.GetTagCountByHeight(context.Background(), 762291)
		if err == nil {
			t.Errorf("%s Failed: expected to throw an error, no error inputted", t.Name())
		}
	})
}

// TestStatsEndpoints_ValidChains tests that stats endpoints work for both BSV and BTC
func TestStatsEndpoints_ValidChains(t *testing.T) {
	testCases := []struct {
		name  string
		chain ChainType
	}{
		{"BSV chain", ChainBSV},
		{"BTC chain", ChainBTC},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockHTTPValidChain{}
			client, err := NewClient(context.Background(), WithChain(tc.chain), WithNetwork(NetworkMain), WithHTTPClient(mockClient))
			if err != nil {
				t.Fatal(err)
			}
			mockClient.SetResponse(func(req *http.Request) (*http.Response, error) {
				// Check if the URL contains the correct chain
				expectedChain := string(tc.chain)
				if !strings.Contains(req.URL.String(), expectedChain) {
					t.Errorf("Expected URL to contain chain [%s], got URL: %s", expectedChain, req.URL.String())
				}

				resp := newHTTPResponse(`{"height":698730}`)
				resp.StatusCode = http.StatusOK
				return resp, nil
			})

			// Test that the stats endpoints work for both chains
			_, err = client.GetBlockStats(context.Background(), 698730)
			if err != nil {
				t.Errorf("GetBlockStats failed for %s: %s", tc.chain, err.Error())
			}

			_, err = client.GetBlockStatsByHash(context.Background(), "testHash")
			if err != nil {
				t.Errorf("GetBlockStatsByHash failed for %s: %s", tc.chain, err.Error())
			}

			mockClient.SetResponse(func(_ *http.Request) (*http.Response, error) {
				resp := newHTTPResponse(`[]`)
				resp.StatusCode = http.StatusOK
				return resp, nil
			})

			_, err = client.GetMinerBlocksStats(context.Background(), 1)
			if err != nil {
				t.Errorf("GetMinerBlocksStats failed for %s: %s", tc.chain, err.Error())
			}

			_, err = client.GetMinerFeesStats(context.Background(), 1714608000, 1714653060)
			if err != nil {
				t.Errorf("GetMinerFeesStats failed for %s: %s", tc.chain, err.Error())
			}

			mockClient.SetResponse(func(_ *http.Request) (*http.Response, error) {
				resp := newHTTPResponse(`{"days":90,"miners":[]}`)
				resp.StatusCode = http.StatusOK
				return resp, nil
			})

			_, err = client.GetMinerSummaryStats(context.Background(), 90)
			if err != nil {
				t.Errorf("GetMinerSummaryStats failed for %s: %s", tc.chain, err.Error())
			}

			mockClient.SetResponse(func(_ *http.Request) (*http.Response, error) {
				resp := newHTTPResponse(`{"height":762291,"hash":"test","tag_counts":{}}`)
				resp.StatusCode = http.StatusOK
				return resp, nil
			})

			_, err = client.GetTagCountByHeight(context.Background(), 762291)
			if err != nil {
				t.Errorf("GetTagCountByHeight failed for %s: %s", tc.chain, err.Error())
			}
		})
	}
}

// mockHTTPValidChain is for testing with a chain-specific mock
type mockHTTPValidChain struct {
	responseFunc func(*http.Request) (*http.Response, error)
}

func (m *mockHTTPValidChain) Do(req *http.Request) (*http.Response, error) {
	if m.responseFunc != nil {
		return m.responseFunc(req)
	}
	resp := newHTTPResponse("")
	resp.StatusCode = http.StatusOK
	return resp, nil
}

func (m *mockHTTPValidChain) SetResponse(f func(*http.Request) (*http.Response, error)) {
	m.responseFunc = f
}

// newHTTPResponse is a helper for tests
func newHTTPResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

// mockHTTPError is for testing error conditions
type mockHTTPError struct{}

var errHTTP = errors.New("http error")

func (m *mockHTTPError) Do(_ *http.Request) (*http.Response, error) {
	return nil, errHTTP
}
