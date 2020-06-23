package whatsonchain

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// mockHTTP for mocking requests
type mockHTTPAddresses struct{}

// Do is a mock http request
func (m *mockHTTPAddresses) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, fmt.Errorf("missing request")
	}

	//
	// Address Info
	//

	// Valid (info)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/info") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"isvalid": true,"address": "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA","scriptPubKey": "76a9143d0e5368bdadddca108a0fe44739919274c726c788ac","ismine": false,"iswatchonly": false,"isscript": false}`)))
	}

	// Invalid (info) return an error
	if strings.Contains(req.URL.String(), "/error/info") {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("missing request")
	}

	// Valid (but invalid bsv address)
	if strings.Contains(req.URL.String(), "/16ZqP5invalid/info") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"isvalid": false,"address": "","scriptPubKey": "","ismine": false,"iswatchonly": false,"isscript": false}`)))
	}

	//
	// Address Balance
	//

	// Valid (balance)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/balance") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"confirmed": 10102050381,"unconfirmed": 123}`)))
	}

	// Invalid (balance) return an error
	if strings.Contains(req.URL.String(), "/16ZqP5invalid/balance") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	//
	// Address History
	//

	// Valid (history)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"tx_hash": "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1","height": 563052},{"tx_hash": "1c312435789754392f92ffcb64e1248e17da47bed179abfd27e6003c775e0e04","height": 565076}]`)))
	}

	// Valid (history) (no results)
	if strings.Contains(req.URL.String(), "/1NfHy82RqJVGEau9u5DwFRyGc6QKwDuQeT/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[]`)))
	}

	// Invalid (history) return an error
	if strings.Contains(req.URL.String(), "/16ZqP5invalid/history") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	//
	// Address unspent
	//

	// Valid (unspent)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[{"height": 639302,"tx_pos": 3,"tx_hash": "33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3","value": 2451680},{"height": 639601,"tx_pos": 3,"tx_hash": "4805041897a2ae59ffca85f0deb46e89d73d1ba4478bbd9c0fcd76ba0985ded2","value": 2744764},{"height": 640276,"tx_pos": 3,"tx_hash": "2493ff4cbca16b892ac641b7f2cb6d4388e75cb3f8963c291183f2bf0b27f415","value": 2568774}]`)))
	}

	// Valid (unspent) (no results)
	if strings.Contains(req.URL.String(), "/1NfHy82RqJVGEau9u5DwFRyGc6QKwDuQeT/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`[]`)))
	}

	// Invalid (unspent) return an error
	if strings.Contains(req.URL.String(), "/16ZqP5invalid/unspent") {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(``)))
		return resp, fmt.Errorf("bad request")
	}

	//
	// Address download statement
	//

	// Valid (download statement)
	if strings.Contains(req.URL.String(), "/statement/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`%PDF-1.4
%Óëéá
1 0 obj
<</Creator (Chromium)
/Producer (Skia/PDF m73)
/CreationDate (D:20200622155222+00'00')
/ModDate (D:20200622155222+00'00')>>
endobj
3 0 obj
<</ca 1
/BM /Normal>>
endobj
5 0 obj`)))
	}

	// Valid (download statement) (invalid address)
	if strings.Contains(req.URL.String(), "/statement/invalid") {
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`%PDF-1.4
%Óëéá
1 0 obj
<</Creator (Chromium)
/Producer (Skia/PDF m73)
/CreationDate (D:20200622155222+00'00')
/ModDate (D:20200622155222+00'00')>>
endobj
3 0 obj
<</ca 1
/BM /Normal>>
endobj
invalid
5 0 obj`)))
	}

	// Default is valid
	return resp, nil
}

// TestClient_AddressInfo tests the AddressInfo()
func TestClient_AddressInfo(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})

	// Create the list of tests
	var tests = []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", false, http.StatusOK},
		{"16ZqP5invalid", "", false, http.StatusOK},
		{"error", "", true, http.StatusInternalServerError},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.AddressInfo(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted and [%s] expected", t.Name(), test.input, test.expected)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v] error [%s]", t.Name(), test.input, test.expected, output, err.Error())
		} else if output != nil && output.Address != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output.Address)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_AddressBalance tests the AddressBalance()
func TestClient_AddressBalance(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})

	// Create the list of tests
	var tests = []struct {
		input         string
		confirmed     int64
		unconfirmed   int64
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", 10102050381, 123, false, http.StatusOK},
		{"16ZqP5invalid", 0, 0, true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.AddressBalance(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if output != nil && output.Confirmed != test.confirmed && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%d] confirm expected, received: [%d]", t.Name(), test.input, test.confirmed, output.Confirmed)
		} else if output != nil && output.Unconfirmed != test.unconfirmed && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%d] unconfirmed expected, received: [%d]", t.Name(), test.input, test.unconfirmed, output.Unconfirmed)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_AddressHistory tests the AddressHistory()
func TestClient_AddressHistory(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})

	// Create the list of tests
	var tests = []struct {
		input         string
		txHash        string
		height        int64
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1", 563052, false, http.StatusOK},
		{"1NfHy82RqJVGEau9u5DwFRyGc6QKwDuQeT", "", 0, false, http.StatusOK},
		{"16ZqP5invalid", "", 0, true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.AddressHistory(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if output != nil && len(output) > 0 && output[0].TxHash != test.txHash && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.txHash, output[0].TxHash)
		} else if output != nil && len(output) > 0 && output[0].Height != test.height && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_AddressUnspentTransactions tests the AddressUnspentTransactions()
func TestClient_AddressUnspentTransactions(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})

	// Create the list of tests
	var tests = []struct {
		input         string
		txHash        string
		height        int64
		value         int64
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", "33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3", 639302, 2451680, false, http.StatusOK},
		{"1NfHy82RqJVGEau9u5DwFRyGc6QKwDuQeT", "", 0, 0, false, http.StatusOK},
		{"16ZqP5invalid", "", 0, 0, true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.AddressUnspentTransactions(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if output != nil && len(output) > 0 && !test.expectedError {
			if output[0].TxHash != test.txHash {
				t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.txHash, output[0].TxHash)
			} else if output[0].Height != test.height {
				t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
			} else if output[0].Value != test.value {
				t.Errorf("%s Failed: [%s] inputted and [%d] value expected, received: [%d]", t.Name(), test.input, test.value, output[0].Value)
			}
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}

// TestClient_AddressUnspentTransactions tests the AddressUnspentTransactions()
func TestClient_AddressUnspentTransactionDetails(t *testing.T) {
	t.Parallel()

	// New mock client
	// client := newMockClient(t, &mockHTTPAddresses{})

	// var resp AddressHistory
	/*address := "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA"
	// var history AddressHistory
	// if history, err = client.AddressUnspentTransactionDetails(address, 5); err != nil {
	if _, err = client.AddressUnspentTransactionDetails(address, 5); err != nil {
		t.Fatal("error occurred: " + err.Error())
	}
	*/
	/*if len(history) == 0 {
		t.Fatal("no utxos found")
	}*/
}

// TestClient_DownloadStatement tests the DownloadStatement()
func TestClient_DownloadStatement(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})

	// Create the list of tests
	var tests = []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", "PDF", false, http.StatusOK},
		{"invalid", "invalid", false, http.StatusOK},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.DownloadStatement(test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if !strings.Contains(output, test.expected) && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output)
		} else if client.LastRequest.StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest.StatusCode, test.input)
		}
	}
}
