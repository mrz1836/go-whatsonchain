package whatsonchain

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test error variables
var (
	errMissingRequest  = errors.New("missing request")
	errBadRequest      = errors.New("bad request")
	errNoValidResponse = errors.New("no valid response found")
)

// mockHTTPAddresses for mocking requests
type mockHTTPAddresses struct{}

// Do is a mock http request
func (m *mockHTTPAddresses) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, errMissingRequest
	}

	//
	// Address Info
	//

	// Valid (info)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/info") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`{"isvalid": true,"address": "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA","scriptPubKey": "76a9143d0e5368bdadddca108a0fe44739919274c726c788ac","ismine": false,"iswatchonly": false,"isscript": false}`))
	}

	// Invalid (info) return an error
	if strings.Contains(req.URL.String(), "/error/info") {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, errMissingRequest
	}

	// Valid (but invalid bsv address)
	if strings.Contains(req.URL.String(), "/16ZqP5invalid/info") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`{"isvalid": false,"address": "","scriptPubKey": "","ismine": false,"iswatchonly": false,"isscript": false}`))
	}

	// Not found
	if strings.Contains(req.URL.String(), "/notFound/info") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, nil
	}

	//
	// Address Balance
	//

	// Valid (balance)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/balance") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`{"confirmed": 10102050381,"unconfirmed": 123}`))
	}

	// Invalid (balance) return an error
	if strings.Contains(req.URL.String(), "/16ZqP5invalid/balance") {
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, errBadRequest
	}

	// Not found
	if strings.Contains(req.URL.String(), "/16ZqP5notFound/balance") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, nil
	}

	//
	// Address History
	//

	// Valid (history)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"tx_hash": "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1","height": 563052},{"tx_hash": "1c312435789754392f92ffcb64e1248e17da47bed179abfd27e6003c775e0e04","height": 565076}]`))
	}

	// Valid (history) (no results)
	if strings.Contains(req.URL.String(), "/1NfHy82RqJVGEau9u5DwFRyGc6QKwDuQeT/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[]`))
	}

	// Invalid (history) return an error
	if strings.Contains(req.URL.String(), "/16ZqP5invalid/history") {
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, errBadRequest
	}

	// Not found
	if strings.Contains(req.URL.String(), "/16ZqP5notFound/history") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, nil
	}

	//
	// Address unspent
	//

	// Valid (unspent/all)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/unspent/all") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"height": 639302,"tx_pos": 3,"tx_hash": "33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3","value": 2451680},{"height": 639601,"tx_pos": 3,"tx_hash": "4805041897a2ae59ffca85f0deb46e89d73d1ba4478bbd9c0fcd76ba0985ded2","value": 2744764},{"height": 640276,"tx_pos": 3,"tx_hash": "2493ff4cbca16b892ac641b7f2cb6d4388e75cb3f8963c291183f2bf0b27f415","value": 2568774}]`))
	}

	// Valid (unspent/all) (no results)
	if strings.Contains(req.URL.String(), "/1NfHy82RqJVGEau9u5DwFRyGc6QKwDuQeT/unspent/all") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[]`))
	}

	// Invalid (unspent/all) return an error
	if strings.Contains(req.URL.String(), "/16ZqP5invalid/unspent/all") {
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, errBadRequest
	}

	// Not found
	if strings.Contains(req.URL.String(), "/16ZqP5notFound/unspent/all") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, nil
	}

	//
	// Address bulk balance
	//

	// Valid (unspent)
	if strings.Contains(req.URL.String(), "/addresses/balance") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"address":"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP","balance":{"confirmed":0,"unconfirmed":0},"error":""},{"address":"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob","balance":{"confirmed":301995631,"unconfirmed":0},"error":""}]`))
	}

	//
	// Address bulk utxo
	//

	// Valid (unspent/all)
	if strings.Contains(req.URL.String(), "/addresses/unspent/all") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"address":"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP","unspent":[],"error":""},{"address":"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob","unspent":[{"height":658677,"tx_pos":1,"tx_hash":"be97e63bf79a961c69bc09d73cbef18232c7962fdced58244ed4014ba7e342b9","value":39799008},{"height":658726,"tx_pos":1,"tx_hash":"e5a7bc2338287fc0f3e38dff696b9ba41b3950c12bd8b7b1d92f3b0c056b4255","value":19110368},{"height":661623,"tx_pos":1,"tx_hash":"adea0a707c16f712f7a8faacfc8759d0b9d148693545c83511be1e2ed7fab4aa","value":19599008},{"height":661746,"tx_pos":1,"tx_hash":"afd619887ba5de9eb0e6076b7a37e96625227791520092fe142366c5c631c79e","value":44764416},{"height":661989,"tx_pos":1,"tx_hash":"15dcf82c9c461f3cb430e5ada855483c9e6c01bf4bc6fe667f3b798bd9f44acb","value":16658528},{"height":662494,"tx_pos":1,"tx_hash":"db8872fb1315e7f62013657d68db1871859624991a3ed77265aa85b8fdc768e5","value":32237986},{"height":662783,"tx_pos":4,"tx_hash":"2fa2686c61b6df1796717ca6d5f1934f0c39a5f8d2e42a6f213e76cb2ae66b54","value":10000000},{"height":662789,"tx_pos":4,"tx_hash":"c868a0616836bb017f956ce846ca6f3c56a985955742bd0fea22840a9d0168df","value":10000000},{"height":662791,"tx_pos":2,"tx_hash":"95c69649798b1d66e37318f2d65374095f6e0cd1675d1214402bd8e6002bf424","value":10000000},{"height":662794,"tx_pos":4,"tx_hash":"69720b7e41ca113d5fa988f5f4fd635d398459ba8e6ebb5d0a3a8f42097f1dcd","value":10000000},{"height":662853,"tx_pos":1,"tx_hash":"c42a4752f2551c195b27d016f8e522e724ad99b81d7e2459630b53e5a06178f3","value":5045746},{"height":662897,"tx_pos":1,"tx_hash":"2ca9c744b857a46266e4c0ac827db254eef54fa19530d3f21e460fe8d9445844","value":11405228},{"height":662992,"tx_pos":1,"tx_hash":"01c33159e9e00a7cb07248926e1ff8ed2d6a2450565fc0c27c30600457b2e572","value":47748859},{"height":663033,"tx_pos":1,"tx_hash":"09c5c72f807e572a5ac96e809d1c10b5bf27d63099cb4a6d871b74d459778bde","value":13629728},{"height":663034,"tx_pos":1,"tx_hash":"c8d3137f13ce2a4b8bfd919210c233a14a565c87e7c1ef4a693e6576adcc0419","value":7393008},{"height":663095,"tx_pos":1,"tx_hash":"58f416f323ae5b4d104b6246fca84ec4b1a6bb5a26174a732801e48008d02bbc","value":4603748}],"error":""}]`))
	}

	//
	// Address download statement
	//

	// Valid (download statement)
	if strings.Contains(req.URL.String(), "/statement/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`%PDF-1.4
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
5 0 obj`))
	}

	// Valid (download statement) (invalid address)
	if strings.Contains(req.URL.String(), "/statement/invalid") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`%PDF-1.4
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
5 0 obj`))
	}

	//
	// Bulk Tx Data
	//

	var data TxHashes
	if strings.Contains(req.URL.String(), "/test/txs") {

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&data)
		if err != nil {
			return resp, err
		}

		// Valid - for AddressDetails
		if strings.Contains(data.TxIDs[0], "33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3") {
			resp.StatusCode = http.StatusOK
			resp.Body = io.NopCloser(strings.NewReader(`[{"hex":"","txid":"33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3","hash":"33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3","version":1,"size":440,"locktime":0,"vin":[{"coinbase":"","txid":"fabe0b5d0979e068dce986692d1c5620f37383657a2fe7969f1cfe4a81b7f517","vout":3,"scriptSig":{"asm":"30450221008f74bb75c331cb7902a4e7539ee60fafe2c9a73d325aba6fc3ff9105ed91e219022064e65a5662c0593086ab05a0131e5abac5ef249f5f33c74351c2bed653da269f[ALL|FORKID] 026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103","hex":"4830450221008f74bb75c331cb7902a4e7539ee60fafe2c9a73d325aba6fc3ff9105ed91e219022064e65a5662c0593086ab05a0131e5abac5ef249f5f33c74351c2bed653da269f4121026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103"},"sequence":4294967295}],"vout":[{"value":0,"n":0,"scriptPubKey":{"asm":"0 OP_RETURN 3150755161374b36324d694b43747373534c4b79316b683536575755374d74555235 5522771 7368801 746f6e6963706f77 1701869940 6f666665725f636c69636b 6f666665725f636f6e6669675f6964 56 6f666665725f73657373696f6e5f6964 33316566313830633732363465303032373836333261306131613830313835313336363236336537306361383233353138373664636436386563666163623365","hex":"006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f69640138106f666665725f73657373696f6e5f69644033316566313830633732363465303032373836333261306131613830313835313336363236336537306361383233353138373664636436386563666163623365","type":"nulldata","opReturn":{"type":"bitcom","action":"","text":"","parts":null},"isTruncated":false}},{"value":0.00000549,"n":1,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 09cc4559bdcb84cb35c107743f0dbb10d66679cc OP_EQUALVERIFY OP_CHECKSIG","hex":"76a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1tonicZQwN2BNKhVwPXqh8ez3q56y1EYw"],"opReturn":null,"isTruncated":false}},{"value":0.00005489,"n":2,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 b5195bf7db0652f536a7dddbe36a99a091125468 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914b5195bf7db0652f536a7dddbe36a99a09112546888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1HWZgiMKQKPSkLzT7hipS22AvkQZJsyxmT"],"opReturn":null,"isTruncated":false}},{"value":0.0245168,"n":3,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 bf49c6a5406675e174f4f6a83b3d94dd9d845398 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914bf49c6a5406675e174f4f6a83b3d94dd9d84539888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1JSSSgcyufLgbXFw6WAXyXgBrmgFpnqXWh"],"opReturn":null,"isTruncated":false}}],"blockhash":"0000000000000000026b9da3860e4c8ee351a7af46da6042eaa5d110113b9fad","confirmations":1348,"time":1592122768,"blocktime":1592122768},{"hex":"","txid":"4805041897a2ae59ffca85f0deb46e89d73d1ba4478bbd9c0fcd76ba0985ded2","hash":"4805041897a2ae59ffca85f0deb46e89d73d1ba4478bbd9c0fcd76ba0985ded2","version":1,"size":439,"locktime":0,"vin":[{"coinbase":"","txid":"5a45b8415e5c1740353cfb011d29e04ec104865be6560dff5bd6cb31db75d559","vout":3,"scriptSig":{"asm":"3044022008e2417d072cfbb95d4e04c7e6e6ab70e415a379fb912cb2e0503e3df0ae0d2002201f9fcbf6c65ba6624fe0669d08155ed7c0d19c28be72daf3e00de2613656f955[ALL|FORKID] 026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103","hex":"473044022008e2417d072cfbb95d4e04c7e6e6ab70e415a379fb912cb2e0503e3df0ae0d2002201f9fcbf6c65ba6624fe0669d08155ed7c0d19c28be72daf3e00de2613656f9554121026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103"},"sequence":4294967295}],"vout":[{"value":0,"n":0,"scriptPubKey":{"asm":"0 OP_RETURN 3150755161374b36324d694b43747373534c4b79316b683536575755374d74555235 5522771 7368801 746f6e6963706f77 1701869940 6f666665725f636c69636b 6f666665725f636f6e6669675f6964 56 6f666665725f73657373696f6e5f6964 38313865386165656339353733646431333439373334366135363464633461623035353062333039383830373563393733316631643063653731336536353335","hex":"006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f69640138106f666665725f73657373696f6e5f69644038313865386165656339353733646431333439373334366135363464633461623035353062333039383830373563393733316631643063653731336536353335","type":"nulldata","opReturn":{"type":"bitcom","action":"","text":"","parts":null},"isTruncated":false}},{"value":0.00000573,"n":1,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 09cc4559bdcb84cb35c107743f0dbb10d66679cc OP_EQUALVERIFY OP_CHECKSIG","hex":"76a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1tonicZQwN2BNKhVwPXqh8ez3q56y1EYw"],"opReturn":null,"isTruncated":false}},{"value":0.00005726,"n":2,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 819ae5a5cbb078e96379b8eb25c29d6f7b28c412 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914819ae5a5cbb078e96379b8eb25c29d6f7b28c41288ac","reqSigs":1,"type":"pubkeyhash","addresses":["1CpHjBbHoWzbrqQsPeZ39GLUXejZce9mBs"],"opReturn":null,"isTruncated":false}},{"value":0.02744764,"n":3,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 bf49c6a5406675e174f4f6a83b3d94dd9d845398 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914bf49c6a5406675e174f4f6a83b3d94dd9d84539888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1JSSSgcyufLgbXFw6WAXyXgBrmgFpnqXWh"],"opReturn":null,"isTruncated":false}}],"blockhash":"000000000000000003d684082ab45014f89a7f8e5e35ec94fcb4aa8b5f00c01e","confirmations":1049,"time":1592307236,"blocktime":1592307236},{"hex":"","txid":"2493ff4cbca16b892ac641b7f2cb6d4388e75cb3f8963c291183f2bf0b27f415","hash":"2493ff4cbca16b892ac641b7f2cb6d4388e75cb3f8963c291183f2bf0b27f415","version":1,"size":439,"locktime":0,"vin":[{"coinbase":"","txid":"2ebc8f094fdc012f7d9a0ed39999dcf0318665830f7d5f113af0d1c79fba2f8e","vout":3,"scriptSig":{"asm":"30440220010a62c1d79afcc274b8db821cba1f093c316d67d505a3900c231ae6dfb2dd51022031fe80787c531e1c890754d2cafdc624f3446e4d1bdca18ade83cabd3a2317ac[ALL|FORKID] 026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103","hex":"4730440220010a62c1d79afcc274b8db821cba1f093c316d67d505a3900c231ae6dfb2dd51022031fe80787c531e1c890754d2cafdc624f3446e4d1bdca18ade83cabd3a2317ac4121026d6fc8f05b630e637507084b1678ec753c75b9e050312919e1d973224c5c3103"},"sequence":4294967295}],"vout":[{"value":0,"n":0,"scriptPubKey":{"asm":"0 OP_RETURN 3150755161374b36324d694b43747373534c4b79316b683536575755374d74555235 5522771 7368801 746f6e6963706f77 1701869940 6f666665725f636c69636b 6f666665725f636f6e6669675f6964 56 6f666665725f73657373696f6e5f6964 35656237306231653930306535616437626335663961333663653861643435623664336435636337666466393437343762623364326461663732636631356533","hex":"006a223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540361707008746f6e6963706f7704747970650b6f666665725f636c69636b0f6f666665725f636f6e6669675f69640138106f666665725f73657373696f6e5f69644035656237306231653930306535616437626335663961333663653861643435623664336435636337666466393437343762623364326461663732636631356533","type":"nulldata","opReturn":{"type":"bitcom","action":"","text":"","parts":null},"isTruncated":false}},{"value":0.00000572,"n":1,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 09cc4559bdcb84cb35c107743f0dbb10d66679cc OP_EQUALVERIFY OP_CHECKSIG","hex":"76a91409cc4559bdcb84cb35c107743f0dbb10d66679cc88ac","reqSigs":1,"type":"pubkeyhash","addresses":["1tonicZQwN2BNKhVwPXqh8ez3q56y1EYw"],"opReturn":null,"isTruncated":false}},{"value":0.00005716,"n":2,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 0405a52b27214920873fa222071a8ec9610317a4 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a9140405a52b27214920873fa222071a8ec9610317a488ac","reqSigs":1,"type":"pubkeyhash","addresses":["1NGU17f9HTyv3LffW4zxukSEwsxwf4d53"],"opReturn":null,"isTruncated":false}},{"value":0.02568774,"n":3,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 bf49c6a5406675e174f4f6a83b3d94dd9d845398 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914bf49c6a5406675e174f4f6a83b3d94dd9d84539888ac","reqSigs":1,"type":"pubkeyhash","addresses":["1JSSSgcyufLgbXFw6WAXyXgBrmgFpnqXWh"],"opReturn":null,"isTruncated":false}}],"blockhash":"00000000000000000087222006199927280a010d0db21c6d253409f3e960c7bf","confirmations":374,"time":1592698834,"blocktime":1592698834}]`))
		}
	}

	//
	// Address Used
	//

	// Valid (used)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/used") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`{"used": true}`))
	}

	// Not found
	if strings.Contains(req.URL.String(), "/16ZqP5notFound/used") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, nil
	}

	//
	// Address Scripts
	//

	// Valid (scripts)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/scripts") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`{"scripts": ["76a9143d0e5368bdadddca108a0fe44739919274c726c788ac"]}`))
	}

	// Not found
	if strings.Contains(req.URL.String(), "/16ZqP5notFound/scripts") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, nil
	}

	//
	// Address Unconfirmed Balance
	//

	// Valid (unconfirmed/balance)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/unconfirmed/balance") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`{"balance": 123}`))
	}

	// Not found
	if strings.Contains(req.URL.String(), "/16ZqP5notFound/unconfirmed/balance") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, nil
	}

	//
	// Address Confirmed Balance
	//

	// Valid (confirmed/balance)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/confirmed/balance") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`{"balance": 10102050381}`))
	}

	// Not found
	if strings.Contains(req.URL.String(), "/16ZqP5notFound/confirmed/balance") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, nil
	}

	//
	// Address Unconfirmed History
	//

	// Valid (unconfirmed/history)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/unconfirmed/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"tx_hash": "unconfirmed123","height": 0}]`))
	}

	// Not found
	if strings.Contains(req.URL.String(), "/16ZqP5notFound/unconfirmed/history") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, nil
	}

	//
	// Address Confirmed History
	//

	// Valid (confirmed/history)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/confirmed/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"tx_hash": "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1","height": 563052}]`))
	}

	// Not found
	if strings.Contains(req.URL.String(), "/16ZqP5notFound/confirmed/history") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, nil
	}

	//
	// Bulk Addresses Unconfirmed Balance
	//

	// Valid (addresses/unconfirmed/balance)
	if strings.Contains(req.URL.String(), "/addresses/unconfirmed/balance") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"address":"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP","balance":{"confirmed":12812324,"unconfirmed":7340}},{"address":"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob","balance":{"confirmed":140041,"unconfirmed":0}}]`))
	}

	//
	// Bulk Addresses Confirmed Balance
	//

	// Valid (addresses/confirmed/balance)
	if strings.Contains(req.URL.String(), "/addresses/confirmed/balance") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"address":"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP","balance":{"confirmed":12812324,"unconfirmed":0}},{"address":"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob","balance":{"confirmed":140041,"unconfirmed":0}}]`))
	}

	//
	// Bulk Addresses Unconfirmed History
	//

	// Valid (addresses/unconfirmed/history)
	if strings.Contains(req.URL.String(), "/addresses/unconfirmed/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"address":"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP","history":[{"tx_hash":"unconfirmed123","height":0}]},{"address":"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob","history":[]}]`))
	}

	//
	// Bulk Addresses Confirmed History
	//

	// Valid (addresses/confirmed/history)
	if strings.Contains(req.URL.String(), "/addresses/confirmed/history") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"address":"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP","history":[{"tx_hash":"6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1","height":563052}]},{"address":"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob","history":[]}]`))
	}

	//
	// Bulk Addresses History All
	//

	// Valid (addresses/history/all)
	if strings.Contains(req.URL.String(), "/addresses/history/all") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"address":"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP","history":[{"tx_hash":"6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1","height":563052},{"tx_hash":"unconfirmed123","height":0}]},{"address":"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob","history":[]}]`))
	}

	//
	// Address Unconfirmed UTXOs
	//

	// Valid (unconfirmed/unspent)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/unconfirmed/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"tx_hash":"unconfirmed_utxo_123","height":0,"tx_pos":1,"value":50000}]`))
	}

	// Not found
	if strings.Contains(req.URL.String(), "/16ZqP5notFound/unconfirmed/unspent") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, nil
	}

	//
	// Address Confirmed UTXOs
	//

	// Valid (confirmed/unspent)
	if strings.Contains(req.URL.String(), "/16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA/confirmed/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"tx_hash":"6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1","height":563052,"tx_pos":0,"value":100000}]`))
	}

	// Not found
	if strings.Contains(req.URL.String(), "/16ZqP5notFound/confirmed/unspent") {
		resp.StatusCode = http.StatusNotFound
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, nil
	}

	//
	// Bulk Address Unconfirmed UTXOs
	//

	// Valid (addresses/unconfirmed/unspent)
	if strings.Contains(req.URL.String(), "/addresses/unconfirmed/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"address":"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP","unspent":[{"tx_hash":"unconfirmed_bulk_123","height":0,"tx_pos":1,"value":25000}]},{"address":"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob","unspent":[]}]`))
	}

	//
	// Bulk Address Confirmed UTXOs
	//

	// Valid (addresses/confirmed/unspent)
	if strings.Contains(req.URL.String(), "/addresses/confirmed/unspent") {
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(strings.NewReader(`[{"address":"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP","unspent":[{"tx_hash":"confirmed_bulk_456","height":563052,"tx_pos":0,"value":75000}]},{"address":"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob","unspent":[]}]`))
	}

	// Default is valid
	return resp, nil
}

// mockHTTPAddressesErrors for mocking requests
type mockHTTPAddressesErrors struct{}

// Do is a mock http request
func (m *mockHTTPAddressesErrors) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusBadRequest

	// No req found
	if req == nil {
		return resp, errMissingRequest
	}

	// Invalid (info) return an error
	if strings.Contains(req.URL.String(), "/addresses/balance") {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, errMissingRequest
	}

	// Invalid (info) return an error
	if strings.Contains(req.URL.String(), "/addresses/unspent/all") {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = io.NopCloser(strings.NewReader(``))
		return resp, errMissingRequest
	}

	return nil, errNoValidResponse
}

// mockHTTPAddressesNotFound for mocking requests
type mockHTTPAddressesNotFound struct{}

// Do is a mock http request
func (m *mockHTTPAddressesNotFound) Do(req *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusNotFound

	// No req found
	if req == nil {
		return resp, errMissingRequest
	}

	// Always return empty body for not found
	resp.Body = io.NopCloser(strings.NewReader(``))
	return resp, nil
}

// TestClient_AddressInfo tests the AddressInfo()
func TestClient_AddressInfo(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
		input         string
		expected      string
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", false, http.StatusOK},
		{"16ZqP5invalid", "", false, http.StatusOK},
		{"error", "", true, http.StatusInternalServerError},
		{"notFound", "", true, http.StatusNotFound},
	}

	// Test all
	for _, test := range tests {
		if output, err := client.AddressInfo(ctx, test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted and [%s] expected", t.Name(), test.input, test.expected)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%v] error [%s]", t.Name(), test.input, test.expected, output, err.Error())
		} else if output != nil && output.Address != test.expected && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output.Address)
		} else if client.LastRequest().StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest().StatusCode, test.input)
		}
	}
}

// TestClient_AddressBalance tests the AddressBalance()
// Deprecated: This tests a deprecated method. Use AddressConfirmedBalance and AddressUnconfirmedBalance.
func TestClient_AddressBalance(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
		input         string
		confirmed     int64
		unconfirmed   int64
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", 10102050381, 123, false, http.StatusOK},
		{"16ZqP5invalid", 0, 0, true, http.StatusBadRequest},
		{"16ZqP5notFound", 0, 0, true, http.StatusNotFound},
	}

	// Test all
	for _, test := range tests {
		output, err := client.AddressBalance(ctx, test.input)

		if err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
			continue
		}

		if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
			continue
		}

		if client.LastRequest().StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest().StatusCode, test.input)
			continue
		}

		if !test.expectedError && output != nil {
			if output.Confirmed != test.confirmed {
				t.Errorf("%s Failed: [%s] inputted and [%d] confirm expected, received: [%d]", t.Name(), test.input, test.confirmed, output.Confirmed)
			}
			if output.Unconfirmed != test.unconfirmed {
				t.Errorf("%s Failed: [%s] inputted and [%d] unconfirmed expected, received: [%d]", t.Name(), test.input, test.unconfirmed, output.Unconfirmed)
			}
		}
	}
}

// TestClient_AddressHistory tests the AddressHistory()
// Deprecated: This tests a deprecated method. Use AddressConfirmedHistory and AddressUnconfirmedHistory.
func TestClient_AddressHistory(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
		input         string
		txHash        string
		height        int64
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1", 563052, false, http.StatusOK},
		{"1NfHy82RqJVGEau9u5DwFRyGc6QKwDuQeT", "", 0, false, http.StatusOK},
		{"16ZqP5invalid", "", 0, true, http.StatusBadRequest},
		{"16ZqP5notFound", "", 0, true, http.StatusNotFound},
	}

	// Test all
	for _, test := range tests {
		output, err := client.AddressHistory(ctx, test.input)

		if err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
			continue
		}

		if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
			continue
		}

		if client.LastRequest().StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest().StatusCode, test.input)
			continue
		}

		if !test.expectedError && len(output) > 0 {
			if output[0].TxHash != test.txHash {
				t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.txHash, output[0].TxHash)
			}
			if output[0].Height != test.height {
				t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
			}
		}
	}
}

// TestClient_AddressUnspentTransactions tests the AddressUnspentTransactions()
func TestClient_AddressUnspentTransactions(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
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
		{"16ZqP5notFound", "", 0, 0, true, http.StatusNotFound},
	}

	// Test all
	for _, test := range tests {
		output, err := client.AddressUnspentTransactions(ctx, test.input)

		if err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
			continue
		}

		if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
			continue
		}

		if client.LastRequest().StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest().StatusCode, test.input)
			continue
		}

		if !test.expectedError && len(output) > 0 {
			if output[0].TxHash != test.txHash {
				t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.txHash, output[0].TxHash)
			}
			if output[0].Height != test.height {
				t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
			}
			if output[0].Value != test.value {
				t.Errorf("%s Failed: [%s] inputted and [%d] value expected, received: [%d]", t.Name(), test.input, test.value, output[0].Value)
			}
		}
	}
}

// TestClient_AddressUnspentTransactions tests the AddressUnspentTransactions()
func TestClient_AddressUnspentTransactionDetails(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
		input         string
		txHash        string ``
		height        int64
		expectedError bool
		statusCode    int
	}{
		{"16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA", "33b9432a0ea203bbb6ec00592622cf6e90223849e4c9a76447a19a3ed43907d3", 639302, false, http.StatusOK},
		{"16ZqP5notFound", "", 0, true, http.StatusNotFound},
		{"16ZqP5invalid", "", 0, true, http.StatusBadRequest},
	}

	// Test all
	for _, test := range tests {
		output, err := client.AddressUnspentTransactionDetails(ctx, test.input, 5)

		if err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
			continue
		}

		if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
			continue
		}

		if client.LastRequest().StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest().StatusCode, test.input)
			continue
		}

		if !test.expectedError && len(output) > 0 {
			if output[0].TxHash != test.txHash {
				t.Errorf("%s Failed: [%s] inputted and [%s] hash expected, received: [%s]", t.Name(), test.input, test.txHash, output[0].TxHash)
			}
			if output[0].Height != test.height {
				t.Errorf("%s Failed: [%s] inputted and [%d] height expected, received: [%d]", t.Name(), test.input, test.height, output[0].Height)
			}
		}
	}
}

// TestClient_DownloadStatement tests the DownloadStatement()
func TestClient_DownloadStatement(t *testing.T) {
	t.Parallel()

	// New mock client
	client := newMockClient(&mockHTTPAddresses{})
	ctx := context.Background()

	// Create the list of tests
	tests := []struct {
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
		if output, err := client.DownloadStatement(ctx, test.input); err == nil && test.expectedError {
			t.Errorf("%s Failed: expected to throw an error, no error [%s] inputted", t.Name(), test.input)
		} else if err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted, received: [%v] error [%s]", t.Name(), test.input, output, err.Error())
		} else if !strings.Contains(output, test.expected) && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output)
		} else if client.LastRequest().StatusCode != test.statusCode {
			t.Errorf("%s Expected status code to be %d, got %d, [%s] inputted", t.Name(), test.statusCode, client.LastRequest().StatusCode, test.input)
		}
	}
}

// TestClient_BulkBalance tests the BulkBalance()
// Deprecated: This tests a deprecated method. Use BulkAddressConfirmedBalance and BulkAddressUnconfirmedBalance.
func TestClient_BulkBalance(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balances, err := client.BulkBalance(ctx, &AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}})
		require.NoError(t, err)
		assert.NotNil(t, balances)
		assert.Len(t, balances, 2)
	})

	t.Run("max addresses (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balances, err := client.BulkBalance(ctx, &AddressList{Addresses: []string{
			"1",
			"2",
			"3",
			"4",
			"5",
			"6",
			"7",
			"8",
			"9",
			"10",
			"11",
			"12",
			"13",
			"14",
			"15",
			"16",
			"17",
			"18",
			"19",
			"20",
			"21",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})

	t.Run("bad response (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddressesErrors{})
		ctx := context.Background()
		balances, err := client.BulkBalance(ctx, &AddressList{Addresses: []string{
			"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})

	t.Run("not found", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddressesNotFound{})
		ctx := context.Background()
		balances, err := client.BulkBalance(ctx, &AddressList{Addresses: []string{
			"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})
}

// TestClient_BulkUnspentTransactionsProcessor tests the BulkUnspentTransactionsProcessor()
// Deprecated: This tests a deprecated method. Use BulkAddressConfirmedUTXOs and BulkAddressUnconfirmedUTXOs.
func TestClient_BulkUnspentTransactionsProcessor(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balances, err := client.BulkUnspentTransactionsProcessor(ctx, &AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}})
		require.NoError(t, err)
		assert.NotNil(t, balances)
		assert.Len(t, balances, 2)
	})

	t.Run("over max addresses (no error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balances, err := client.BulkUnspentTransactionsProcessor(ctx, &AddressList{Addresses: []string{
			"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
			"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob",
			"1Bo1qfAm93cgqrd4vjTN1CeNXwjByjfrDC",
			"1AXQmPt6eyU1ZSt2bSvDiV1PJctpLbEZ3u",
			"1GhWikHYDvYRiN37KjDfc6ba6CkaTAZmHG",
			"1BzUYnHr6tY2uAkydt9M8ozctM4e8keW9G",
			"1AU4yMBFnnB8SWjy7nofZcPDRd8x8pJdY5",
			"18x1r2cL1CGjoMbKn5sq3BuDfYFdbjdK3U",
			"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
			"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob",
			"1Bo1qfAm93cgqrd4vjTN1CeNXwjByjfrDC",
			"1AXQmPt6eyU1ZSt2bSvDiV1PJctpLbEZ3u",
			"1GhWikHYDvYRiN37KjDfc6ba6CkaTAZmHG",
			"1BzUYnHr6tY2uAkydt9M8ozctM4e8keW9G",
			"1AU4yMBFnnB8SWjy7nofZcPDRd8x8pJdY5",
			"18x1r2cL1CGjoMbKn5sq3BuDfYFdbjdK3U",
			"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
			"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob",
			"1Bo1qfAm93cgqrd4vjTN1CeNXwjByjfrDC",
			"1AXQmPt6eyU1ZSt2bSvDiV1PJctpLbEZ3u",
			"1GhWikHYDvYRiN37KjDfc6ba6CkaTAZmHG",
			"1BzUYnHr6tY2uAkydt9M8ozctM4e8keW9G",
			"1AU4yMBFnnB8SWjy7nofZcPDRd8x8pJdY5",
			"18x1r2cL1CGjoMbKn5sq3BuDfYFdbjdK3U",
			"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
			"1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob",
			"1Bo1qfAm93cgqrd4vjTN1CeNXwjByjfrDC",
			"1AXQmPt6eyU1ZSt2bSvDiV1PJctpLbEZ3u",
			"1GhWikHYDvYRiN37KjDfc6ba6CkaTAZmHG",
			"1BzUYnHr6tY2uAkydt9M8ozctM4e8keW9G",
			"1AU4yMBFnnB8SWjy7nofZcPDRd8x8pJdY5",
			"18x1r2cL1CGjoMbKn5sq3BuDfYFdbjdK3U",
		}})
		require.NoError(t, err)
		assert.NotNil(t, balances)
	})

	t.Run("bad response (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddressesErrors{})
		ctx := context.Background()
		balances, err := client.BulkUnspentTransactionsProcessor(ctx, &AddressList{Addresses: []string{
			"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})

	t.Run("not found", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddressesNotFound{})
		ctx := context.Background()
		balances, err := client.BulkUnspentTransactionsProcessor(ctx, &AddressList{Addresses: []string{
			"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})
}

// TestClient_BulkUnspentTransactions tests the BulkUnspentTransactions()
// Deprecated: This tests a deprecated method. Use BulkAddressConfirmedUTXOs and BulkAddressUnconfirmedUTXOs.
func TestClient_BulkUnspentTransactions(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balances, err := client.BulkUnspentTransactions(ctx, &AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}})
		require.NoError(t, err)
		assert.NotNil(t, balances)
		assert.Len(t, balances, 2)
	})

	t.Run("max addresses (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balances, err := client.BulkUnspentTransactions(ctx, &AddressList{Addresses: []string{
			"1",
			"2",
			"3",
			"4",
			"5",
			"6",
			"7",
			"8",
			"9",
			"10",
			"11",
			"12",
			"13",
			"14",
			"15",
			"16",
			"17",
			"18",
			"19",
			"20",
			"21",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})

	t.Run("bad response (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddressesErrors{})
		ctx := context.Background()
		balances, err := client.BulkUnspentTransactions(ctx, &AddressList{Addresses: []string{
			"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})

	t.Run("not found", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddressesNotFound{})
		ctx := context.Background()
		balances, err := client.BulkUnspentTransactions(ctx, &AddressList{Addresses: []string{
			"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})
}

// TestClient_AddressUsed tests the AddressUsed()
func TestClient_AddressUsed(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		used, err := client.AddressUsed(ctx, "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
		require.NoError(t, err)
		assert.NotNil(t, used)
		assert.True(t, used.Used)
	})

	t.Run("not found", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		used, err := client.AddressUsed(ctx, "16ZqP5notFound")
		require.Error(t, err)
		assert.Nil(t, used)
	})
}

// TestClient_AddressScripts tests the AddressScripts()
func TestClient_AddressScripts(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		scripts, err := client.AddressScripts(ctx, "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
		require.NoError(t, err)
		assert.NotNil(t, scripts)
		assert.Len(t, scripts.Scripts, 1)
		assert.Equal(t, "76a9143d0e5368bdadddca108a0fe44739919274c726c788ac", scripts.Scripts[0])
	})

	t.Run("not found", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		scripts, err := client.AddressScripts(ctx, "16ZqP5notFound")
		require.Error(t, err)
		assert.Nil(t, scripts)
	})
}

// TestClient_AddressUnconfirmedBalance tests the AddressUnconfirmedBalance()
func TestClient_AddressUnconfirmedBalance(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balance, err := client.AddressUnconfirmedBalance(ctx, "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
		require.NoError(t, err)
		assert.NotNil(t, balance)
		assert.Equal(t, int64(123), balance.Balance)
	})

	t.Run("not found", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balance, err := client.AddressUnconfirmedBalance(ctx, "16ZqP5notFound")
		require.Error(t, err)
		assert.Nil(t, balance)
	})
}

// TestClient_AddressConfirmedBalance tests the AddressConfirmedBalance()
func TestClient_AddressConfirmedBalance(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balance, err := client.AddressConfirmedBalance(ctx, "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
		require.NoError(t, err)
		assert.NotNil(t, balance)
		assert.Equal(t, int64(10102050381), balance.Balance)
	})

	t.Run("not found", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balance, err := client.AddressConfirmedBalance(ctx, "16ZqP5notFound")
		require.Error(t, err)
		assert.Nil(t, balance)
	})
}

// TestClient_AddressUnconfirmedHistory tests the AddressUnconfirmedHistory()
func TestClient_AddressUnconfirmedHistory(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		history, err := client.AddressUnconfirmedHistory(ctx, "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
		require.NoError(t, err)
		assert.NotNil(t, history)
		assert.Len(t, history, 1)
		assert.Equal(t, "unconfirmed123", history[0].TxHash)
		assert.Equal(t, int64(0), history[0].Height)
	})

	t.Run("not found", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		history, err := client.AddressUnconfirmedHistory(ctx, "16ZqP5notFound")
		require.Error(t, err)
		assert.Nil(t, history)
	})
}

// TestClient_AddressConfirmedHistory tests the AddressConfirmedHistory()
func TestClient_AddressConfirmedHistory(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		history, err := client.AddressConfirmedHistory(ctx, "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
		require.NoError(t, err)
		assert.NotNil(t, history)
		assert.Len(t, history, 1)
		assert.Equal(t, "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1", history[0].TxHash)
		assert.Equal(t, int64(563052), history[0].Height)
	})

	t.Run("not found", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		history, err := client.AddressConfirmedHistory(ctx, "16ZqP5notFound")
		require.Error(t, err)
		assert.Nil(t, history)
	})
}

// TestClient_BulkAddressUnconfirmedBalance tests the BulkAddressUnconfirmedBalance()
func TestClient_BulkAddressUnconfirmedBalance(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balances, err := client.BulkAddressUnconfirmedBalance(ctx, &AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}})
		require.NoError(t, err)
		assert.NotNil(t, balances)
		assert.Len(t, balances, 2)
		assert.Equal(t, "16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", balances[0].Address)
	})

	t.Run("max addresses (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balances, err := client.BulkAddressUnconfirmedBalance(ctx, &AddressList{Addresses: []string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})
}

// TestClient_BulkAddressConfirmedBalance tests the BulkAddressConfirmedBalance()
func TestClient_BulkAddressConfirmedBalance(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balances, err := client.BulkAddressConfirmedBalance(ctx, &AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}})
		require.NoError(t, err)
		assert.NotNil(t, balances)
		assert.Len(t, balances, 2)
		assert.Equal(t, "16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", balances[0].Address)
	})

	t.Run("max addresses (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		balances, err := client.BulkAddressConfirmedBalance(ctx, &AddressList{Addresses: []string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
		}})
		require.Error(t, err)
		assert.Nil(t, balances)
	})
}

// TestClient_BulkAddressUnconfirmedHistory tests the BulkAddressUnconfirmedHistory()
func TestClient_BulkAddressUnconfirmedHistory(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		history, err := client.BulkAddressUnconfirmedHistory(ctx, &AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}})
		require.NoError(t, err)
		assert.NotNil(t, history)
		assert.Len(t, history, 2)
		assert.Equal(t, "16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", history[0].Address)
		assert.Len(t, history[0].History, 1)
		assert.Equal(t, "unconfirmed123", history[0].History[0].TxHash)
	})

	t.Run("max addresses (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		history, err := client.BulkAddressUnconfirmedHistory(ctx, &AddressList{Addresses: []string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
		}})
		require.Error(t, err)
		assert.Nil(t, history)
	})
}

// TestClient_BulkAddressConfirmedHistory tests the BulkAddressConfirmedHistory()
func TestClient_BulkAddressConfirmedHistory(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		history, err := client.BulkAddressConfirmedHistory(ctx, &AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}})
		require.NoError(t, err)
		assert.NotNil(t, history)
		assert.Len(t, history, 2)
		assert.Equal(t, "16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", history[0].Address)
		assert.Len(t, history[0].History, 1)
		assert.Equal(t, "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1", history[0].History[0].TxHash)
	})

	t.Run("max addresses (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		history, err := client.BulkAddressConfirmedHistory(ctx, &AddressList{Addresses: []string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
		}})
		require.Error(t, err)
		assert.Nil(t, history)
	})
}

// TestClient_BulkAddressHistory tests the BulkAddressHistory()
func TestClient_BulkAddressHistory(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		history, err := client.BulkAddressHistory(ctx, &AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}})
		require.NoError(t, err)
		assert.NotNil(t, history)
		assert.Len(t, history, 2)
		assert.Equal(t, "16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", history[0].Address)
		assert.Len(t, history[0].History, 2) // Should have both confirmed and unconfirmed
		assert.Equal(t, "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1", history[0].History[0].TxHash)
		assert.Equal(t, "unconfirmed123", history[0].History[1].TxHash)
	})

	t.Run("max addresses (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		history, err := client.BulkAddressHistory(ctx, &AddressList{Addresses: []string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
		}})
		require.Error(t, err)
		assert.Nil(t, history)
	})
}

// TestClient_AddressUnconfirmedUTXOs tests the AddressUnconfirmedUTXOs()
func TestClient_AddressUnconfirmedUTXOs(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		utxos, err := client.AddressUnconfirmedUTXOs(ctx, "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
		require.NoError(t, err)
		assert.NotNil(t, utxos)
		assert.Len(t, utxos, 1)
		assert.Equal(t, "unconfirmed_utxo_123", utxos[0].TxHash)
		assert.Equal(t, int64(0), utxos[0].Height)
		assert.Equal(t, int64(1), utxos[0].TxPos)
		assert.Equal(t, int64(50000), utxos[0].Value)
	})

	t.Run("not found", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		utxos, err := client.AddressUnconfirmedUTXOs(ctx, "16ZqP5notFound")
		require.Error(t, err)
		assert.Nil(t, utxos)
	})
}

// TestClient_AddressConfirmedUTXOs tests the AddressConfirmedUTXOs()
func TestClient_AddressConfirmedUTXOs(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		utxos, err := client.AddressConfirmedUTXOs(ctx, "16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
		require.NoError(t, err)
		assert.NotNil(t, utxos)
		assert.Len(t, utxos, 1)
		assert.Equal(t, "6b22c47e7956e5404e05c3dc87dc9f46e929acfd46c8dd7813a34e1218d2f9d1", utxos[0].TxHash)
		assert.Equal(t, int64(563052), utxos[0].Height)
		assert.Equal(t, int64(0), utxos[0].TxPos)
		assert.Equal(t, int64(100000), utxos[0].Value)
	})

	t.Run("not found", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		utxos, err := client.AddressConfirmedUTXOs(ctx, "16ZqP5notFound")
		require.Error(t, err)
		assert.Nil(t, utxos)
	})
}

// TestClient_BulkAddressUnconfirmedUTXOs tests the BulkAddressUnconfirmedUTXOs()
func TestClient_BulkAddressUnconfirmedUTXOs(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		response, err := client.BulkAddressUnconfirmedUTXOs(ctx, &AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}})
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response, 2)
		assert.Equal(t, "16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", response[0].Address)
		assert.Len(t, response[0].Utxos, 1)
		assert.Equal(t, "unconfirmed_bulk_123", response[0].Utxos[0].TxHash)
		assert.Equal(t, int64(0), response[0].Utxos[0].Height)
		assert.Equal(t, int64(25000), response[0].Utxos[0].Value)
	})

	t.Run("max addresses (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		response, err := client.BulkAddressUnconfirmedUTXOs(ctx, &AddressList{Addresses: []string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
		}})
		require.Error(t, err)
		assert.Nil(t, response)
	})
}

// TestClient_BulkAddressConfirmedUTXOs tests the BulkAddressConfirmedUTXOs()
func TestClient_BulkAddressConfirmedUTXOs(t *testing.T) {
	t.Parallel()

	t.Run("valid response", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		response, err := client.BulkAddressConfirmedUTXOs(ctx, &AddressList{Addresses: []string{"16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", "1KGHhLTQaPr4LErrvbAuGE62yPpDoRwrob"}})
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response, 2)
		assert.Equal(t, "16ZBEb7pp6mx5EAGrdeKivztd5eRJFuvYP", response[0].Address)
		assert.Len(t, response[0].Utxos, 1)
		assert.Equal(t, "confirmed_bulk_456", response[0].Utxos[0].TxHash)
		assert.Equal(t, int64(563052), response[0].Utxos[0].Height)
		assert.Equal(t, int64(75000), response[0].Utxos[0].Value)
	})

	t.Run("max addresses (error)", func(t *testing.T) {
		client := newMockClient(&mockHTTPAddresses{})
		ctx := context.Background()
		response, err := client.BulkAddressConfirmedUTXOs(ctx, &AddressList{Addresses: []string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21",
		}})
		require.Error(t, err)
		assert.Nil(t, response)
	})
}

// TestClient_AddressUnconfirmedUTXOs_EmptyResponse tests AddressUnconfirmedUTXOs with empty response
func TestClient_AddressUnconfirmedUTXOs_EmptyResponse(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPAddressesNotFound{})
	ctx := context.Background()
	history, err := client.AddressUnconfirmedUTXOs(ctx, "notFound")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAddressNotFound)
	assert.Nil(t, history)
}

// TestClient_AddressConfirmedUTXOs_EmptyResponse tests AddressConfirmedUTXOs with empty response
func TestClient_AddressConfirmedUTXOs_EmptyResponse(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPAddressesNotFound{})
	ctx := context.Background()
	history, err := client.AddressConfirmedUTXOs(ctx, "notFound")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAddressNotFound)
	assert.Nil(t, history)
}

// TestClient_AddressUsed_EmptyResponse tests AddressUsed with empty response
func TestClient_AddressUsed_EmptyResponse(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPAddressesNotFound{})
	ctx := context.Background()
	_, err := client.AddressUsed(ctx, "notFound")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAddressNotFound)
}

// TestClient_AddressScripts_EmptyResponse tests AddressScripts with empty response
func TestClient_AddressScripts_EmptyResponse(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPAddressesNotFound{})
	ctx := context.Background()
	scripts, err := client.AddressScripts(ctx, "notFound")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAddressNotFound)
	assert.Nil(t, scripts)
}

// TestClient_BulkAddressUnconfirmedBalance_EmptyResponse tests BulkAddressUnconfirmedBalance with empty response
func TestClient_BulkAddressUnconfirmedBalance_EmptyResponse(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPAddressesNotFound{})
	ctx := context.Background()
	balances, err := client.BulkAddressUnconfirmedBalance(ctx, &AddressList{Addresses: []string{"notFound"}})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAddressNotFound)
	assert.Nil(t, balances)
}

// TestClient_BulkAddressConfirmedBalance_EmptyResponse tests BulkAddressConfirmedBalance with empty response
func TestClient_BulkAddressConfirmedBalance_EmptyResponse(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPAddressesNotFound{})
	ctx := context.Background()
	balances, err := client.BulkAddressConfirmedBalance(ctx, &AddressList{Addresses: []string{"notFound"}})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAddressNotFound)
	assert.Nil(t, balances)
}

// TestClient_BulkAddressUnconfirmedHistory_EmptyResponse tests BulkAddressUnconfirmedHistory with empty response
func TestClient_BulkAddressUnconfirmedHistory_EmptyResponse(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPAddressesNotFound{})
	ctx := context.Background()
	history, err := client.BulkAddressUnconfirmedHistory(ctx, &AddressList{Addresses: []string{"notFound"}})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAddressNotFound)
	assert.Nil(t, history)
}

// TestClient_BulkAddressConfirmedHistory_EmptyResponse tests BulkAddressConfirmedHistory with empty response
func TestClient_BulkAddressConfirmedHistory_EmptyResponse(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPAddressesNotFound{})
	ctx := context.Background()
	history, err := client.BulkAddressConfirmedHistory(ctx, &AddressList{Addresses: []string{"notFound"}})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAddressNotFound)
	assert.Nil(t, history)
}

// TestClient_BulkAddressHistory_EmptyResponse tests BulkAddressHistory with empty response
func TestClient_BulkAddressHistory_EmptyResponse(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPAddressesNotFound{})
	ctx := context.Background()
	history, err := client.BulkAddressHistory(ctx, &AddressList{Addresses: []string{"notFound"}})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAddressNotFound)
	assert.Nil(t, history)
}

// TestClient_BulkAddressUnconfirmedUTXOs_EmptyResponse tests BulkAddressUnconfirmedUTXOs with empty response
func TestClient_BulkAddressUnconfirmedUTXOs_EmptyResponse(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPAddressesNotFound{})
	ctx := context.Background()
	response, err := client.BulkAddressUnconfirmedUTXOs(ctx, &AddressList{Addresses: []string{"notFound"}})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAddressNotFound)
	assert.Nil(t, response)
}

// TestClient_BulkAddressConfirmedUTXOs_EmptyResponse tests BulkAddressConfirmedUTXOs with empty response
func TestClient_BulkAddressConfirmedUTXOs_EmptyResponse(t *testing.T) {
	t.Parallel()

	client := newMockClient(&mockHTTPAddressesNotFound{})
	ctx := context.Background()
	response, err := client.BulkAddressConfirmedUTXOs(ctx, &AddressList{Addresses: []string{"notFound"}})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrAddressNotFound)
	assert.Nil(t, response)
}
