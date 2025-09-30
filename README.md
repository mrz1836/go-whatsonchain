# 📨 go-whatsonchain
> The Unofficial Golang SDK for the [whatsonchain.com API](https://docs.whatsonchain.com/) supporting both **[BSV](https://bsvblockchain.org/)** and **[BTC](https://thatsbtcnotbitcoin.com/)** blockchains

<table>
  <thead>
    <tr>
      <th>CI&nbsp;/&nbsp;CD</th>
      <th>Quality&nbsp;&amp;&nbsp;Security</th>
      <th>Docs&nbsp;&amp;&nbsp;Meta</th>
      <th>Community</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td valign="top" align="left">
        <a href="https://github.com/mrz1836/go-whatsonchain/releases">
          <img src="https://img.shields.io/github/release-pre/mrz1836/go-whatsonchain?logo=github&style=flat" alt="Latest Release">
        </a><br/>
        <a href="https://github.com/mrz1836/go-whatsonchain/actions">
          <img src="https://img.shields.io/github/actions/workflow/status/mrz1836/go-whatsonchain/fortress.yml?branch=master&logo=github&style=flat" alt="Build Status">
        </a><br/>
		<a href="https://github.com/mrz1836/go-whatsonchain/actions">
          <img src="https://github.com/mrz1836/go-whatsonchain/actions/workflows/codeql-analysis.yml/badge.svg?style=flat" alt="CodeQL">
        </a><br/>
        <a href="https://github.com/mrz1836/go-whatsonchain/commits/master">
		  <img src="https://img.shields.io/github/last-commit/mrz1836/go-whatsonchain?style=flat&logo=clockify&logoColor=white" alt="Last commit">
		</a>
      </td>
      <td valign="top" align="left">
        <a href="https://goreportcard.com/report/github.com/mrz1836/go-whatsonchain">
          <img src="https://goreportcard.com/badge/github.com/mrz1836/go-whatsonchain?style=flat" alt="Go Report Card">
        </a><br/>
		<a href="https://codecov.io/gh/mrz1836/go-whatsonchain">
          <img src="https://codecov.io/gh/mrz1836/go-whatsonchain/branch/master/graph/badge.svg?style=flat" alt="Code Coverage">
        </a><br/>
		<a href="https://scorecard.dev/viewer/?uri=github.com/mrz1836/go-whatsonchain">
          <img src="https://api.scorecard.dev/projects/github.com/mrz1836/go-whatsonchain/badge?logo=springsecurity&logoColor=white" alt="OpenSSF Scorecard">
        </a><br/>
		<a href=".github/SECURITY.md">
          <img src="https://img.shields.io/badge/security-policy-blue?style=flat&logo=springsecurity&logoColor=white" alt="Security policy">
        </a>
      </td>
      <td valign="top" align="left">
        <a href="https://golang.org/">
          <img src="https://img.shields.io/github/go-mod/go-version/mrz1836/go-whatsonchain?style=flat" alt="Go version">
        </a><br/>
        <a href="https://pkg.go.dev/github.com/mrz1836/go-whatsonchain?tab=doc">
          <img src="https://pkg.go.dev/badge/github.com/mrz1836/go-whatsonchain.svg?style=flat" alt="Go docs">
        </a><br/>
        <a href=".github/AGENTS.md">
          <img src="https://img.shields.io/badge/AGENTS.md-found-40b814?style=flat&logo=openai" alt="AGENTS.md rules">
        </a><br/>
        <a href="https://github.com/mrz1836/mage-x">
          <img src="https://img.shields.io/badge/Mage-supported-brightgreen?style=flat&logo=go&logoColor=white" alt="MAGE-X Supported">
        </a><br/>
		<a href=".github/dependabot.yml">
          <img src="https://img.shields.io/badge/dependencies-automatic-blue?logo=dependabot&style=flat" alt="Dependabot">
        </a>
      </td>
      <td valign="top" align="left">
        <a href="https://docs.whatsonchain.com/">
          <img src="https://img.shields.io/badge/API-docs-FFDD00?style=flat&logo=postman&logoColor=white" alt="go-whatsonchain API docs">
        </a><br/>
        <a href="https://github.com/mrz1836/go-whatsonchain/graphs/contributors">
          <img src="https://img.shields.io/github/contributors/mrz1836/go-whatsonchain?style=flat&logo=contentful&logoColor=white" alt="Contributors">
        </a><br/>
        <a href="https://github.com/sponsors/mrz1836">
          <img src="https://img.shields.io/badge/sponsor-MrZ-181717.svg?logo=github&style=flat" alt="Sponsor">
        </a><br/>
        <a href="https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=go-whatsonchain&utm_term=go-whatsonchain&utm_content=go-whatsonchain">
          <img src="https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat" alt="Donate Bitcoin">
        </a>
      </td>
    </tr>
  </tbody>
</table>

<br/>

## 🗂️ Table of Contents
* [Installation](#-installation)
* [Usage](#-usage)
* [Documentation](#-documentation)
* [Examples & Tests](#-examples--tests)
* [Benchmarks](#-benchmarks)
* [Code Standards](#-code-standards)
* [AI Compliance](#-ai-compliance)
* [Maintainers](#-maintainers)
* [Contributing](#-contributing)
* [License](#-license)

<br/>

## 📦 Installation

**go-whatsonchain** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get github.com/mrz1836/go-whatsonchain
```

<br/>

## 💡 Usage

### Quick Start

```go
package main

import (
	"context"
	"log"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create a client with default options (BSV mainnet)
	client, err := whatsonchain.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("client loaded", client.UserAgent())
	log.Println("Chain:", client.Chain(), "Network:", client.Network())
}
```

### Configuration Options

The library uses functional options for clean and flexible configuration:

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create a client with custom options
	client, err := whatsonchain.NewClient(
		context.Background(),
		whatsonchain.WithChain(whatsonchain.ChainBSV),
		whatsonchain.WithNetwork(whatsonchain.NetworkMain),
		whatsonchain.WithAPIKey("your-secret-key"),
		whatsonchain.WithUserAgent("my-app/1.0"),
		whatsonchain.WithRateLimit(10),
		whatsonchain.WithRequestTimeout(60*time.Second),
		whatsonchain.WithRequestRetryCount(3),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("client loaded with custom options")
}
```

### Available Options

- `WithChain(chain)` - Set blockchain (ChainBSV or ChainBTC)
- `WithNetwork(network)` - Set network (NetworkMain, NetworkTest, NetworkStn)
- `WithAPIKey(key)` - Set API key for authenticated requests
- `WithUserAgent(agent)` - Set custom user agent
- `WithRateLimit(limit)` - Set rate limit per second
- `WithHTTPClient(client)` - Use custom HTTP client
- `WithRequestTimeout(timeout)` - Set request timeout
- `WithRequestRetryCount(count)` - Set retry count for failed requests
- `WithBackoff(initial, max, factor, jitter)` - Configure exponential backoff
- `WithDialer(keepAlive, timeout)` - Configure dialer settings
- `WithTransport(idle, tls, expect, maxIdle)` - Configure transport settings

### Multi-Chain Support

#### BSV Client

```go
package main

import (
	"context"
	"log"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create BSV client
	client, err := whatsonchain.NewClient(
		context.Background(),
		whatsonchain.WithChain(whatsonchain.ChainBSV),
		whatsonchain.WithNetwork(whatsonchain.NetworkMain),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Use BSV-specific methods
	opReturnData, err := client.GetOpReturnData(context.Background(), "your-tx-hash")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OP_RETURN data:", opReturnData)

	// Use shared methods (work for both BSV and BTC)
	chainInfo, err := client.GetChainInfo(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("BSV Chain Info: %+v", chainInfo)
}
```

#### BTC Client

```go
package main

import (
	"context"
	"log"

	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	// Create BTC client
	client, err := whatsonchain.NewClient(
		context.Background(),
		whatsonchain.WithChain(whatsonchain.ChainBTC),
		whatsonchain.WithNetwork(whatsonchain.NetworkMain),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Use BTC-specific methods
	blockStats, err := client.GetBlockStats(context.Background(), 700000)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block Stats: %+v", blockStats)

	// Get miner statistics
	minerStats, err := client.GetMinerBlocksStats(context.Background(), 7) // last 7 days
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Miner Stats: %+v", minerStats)

	// Use shared methods (work for both BSV and BTC)
	chainInfo, err := client.GetChainInfo(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("BTC Chain Info: %+v", chainInfo)
}
```

<br/>

## 📚 Documentation

### Features
- **Multi-blockchain support** - Seamless switching between [BSV](https://bsvblockchain.org/) and [BTC](https://thatsbtcnotbitcoin.com/) blockchains with a single client
- **Production-ready HTTP client** - Built-in exponential backoff with configurable retry logic and crypto-secure jitter to handle transient failures gracefully
- **Intelligent rate limiting** - Per-second request throttling with automatic sleep intervals to stay within API quotas
- **Zero external dependencies** - Pure Go implementation with no production dependencies (testify only for testing)
- **Comprehensive API coverage** - 135+ endpoints (71 BSV, 64 BTC) fully implemented and tested
- **Flexible configuration** - Functional options pattern for clean, type-safe client initialization
- **Enterprise-grade transport** - Fine-grained control over timeouts, keep-alives, connection pooling, and TLS handshake settings
- **Network flexibility** - Switch between mainnet, testnet, and STN per client or per request

View the generated [documentation](https://pkg.go.dev/github.com/mrz1836/go-whatsonchain?tab=doc)

> **Heads up!** `go-whatsonchain` is intentionally light on dependencies. The only
external package it uses is the excellent `testify` suite—and that's just for
our tests. You can drop this library into your projects without dragging along
extra baggage.

<br/>

<details>
<summary><strong><code>Supported API Coverage</code></strong></summary>
<br/>

**Coverage Summary:** 135 total endpoints (71 BSV + 64 BTC) from the [whatsonchain.com API](https://docs.whatsonchain.com/)

**Quick Navigation:** [BSV API](#bsv-api-71-endpoints) • [BTC API](#btc-api-64-endpoints) • [WebSockets](#websockets)

---

## BSV API (71 endpoints)

### ✅ Health (1 endpoint)
- [x] [Get Health Status](https://docs.whatsonchain.com/api/health) - `/woc`

### ✅ Chain Info (4 endpoints)
- [x] [Get Blockchain Info](https://docs.whatsonchain.com/api/chain-info#get-blockchain-info) - `/chain/info`
- [x] [Get Circulating Supply](https://docs.whatsonchain.com/api/chain-info#get-circulating-supply) - `/circulatingsupply`
- [x] [Get Chain Tips](https://docs.whatsonchain.com/api/chain-info#get-chain-tips) - `/chain/tips`
- [x] [Get Peer Info](https://docs.whatsonchain.com/api/chain-info#get-peer-info) - `/peer/info`

### ✅ Block (7 endpoints)
- [x] [Get Block by Hash](https://docs.whatsonchain.com/api/block#get-by-hash) - `/block/hash/{hash}`
- [x] [Get Block by Height](https://docs.whatsonchain.com/api/block#get-by-height) - `/block/height/{height}`
- [x] [Get Block Pages](https://docs.whatsonchain.com/api/block#get-block-pages) - `/block/hash/{hash}/page/{page}`
- [x] [Get Block Headers](https://docs.whatsonchain.com/api/block#get-headers) - `/block/headers`
- [x] [Get Header by Hash or Height](https://docs.whatsonchain.com/api/block#get-header-by-hash) - `/block/{hash}/header`
- [x] [Get Header Bytes File Links](https://docs.whatsonchain.com/api/block#get-header-bytes) - `/block/headers/resources`
- [x] [Get Latest Header Bytes](https://docs.whatsonchain.com/api/block#get-latest-headers) - `/block/headers/latest`

### ✅ Transaction (13 endpoints)
- [x] [Get Transaction by Hash](https://docs.whatsonchain.com/api/transaction#get-by-tx-hash) - `/tx/hash/{hash}`
- [x] [Get Transaction Propagation Status](https://docs.whatsonchain.com/api/transaction#get-tx-propagation) - `/tx/hash/{hash}/propagation`
- [x] [Broadcast Transaction](https://docs.whatsonchain.com/api/transaction#broadcast-transaction) - `/tx/raw` (POST)
- [x] [Bulk Transaction Details](https://docs.whatsonchain.com/api/transaction#bulk-transaction-details) - `/txs` (POST)
- [x] [Bulk Transaction Status](https://docs.whatsonchain.com/api/transaction#bulk-transaction-status) - `/txs/status` (POST)
- [x] [Decode Transaction](https://docs.whatsonchain.com/api/transaction#decode-transaction) - `/tx/decode` (POST)
- [x] [Download Receipt (PDF)](https://docs.whatsonchain.com/api/transaction#download-receipt) - Receipt download
- [x] [Get Transaction as Binary](https://docs.whatsonchain.com/api/transaction#get-tx-binary) - `/tx/{hash}/bin`
- [x] [Get Raw Transaction Data (Hex)](https://docs.whatsonchain.com/api/transaction#get-raw-tx-data) - `/tx/{hash}/hex`
- [x] [Bulk Raw Transaction Data](https://docs.whatsonchain.com/api/transaction#bulk-raw-tx-data) - `/txs/hex` (POST)
- [x] [Get Raw Transaction Output](https://docs.whatsonchain.com/api/transaction#get-raw-tx-output) - `/tx/{hash}/out/{index}/hex`
- [x] [Bulk Raw Transaction Output Data](https://docs.whatsonchain.com/api/transaction#bulk-raw-tx-output) - `/txs/vouts/hex` (POST)
- [x] [Get Merkle Proof (TSC)](https://docs.whatsonchain.com/api/transaction#get-merkle-proof) - `/tx/{hash}/proof/tsc`

### ✅ Mempool (2 endpoints)
- [x] [Get Mempool Info](https://docs.whatsonchain.com/api/mempool#get-mempool-info) - `/mempool/info`
- [x] [Get Mempool Transactions](https://docs.whatsonchain.com/api/mempool#get-mempool-transactions) - `/mempool/raw`

### ✅ (Un)Spent Transaction Outputs (14 endpoints)
- [x] [Get Unspent UTXOs by Address](https://docs.whatsonchain.com/api/address#get-unspent-transactions) - `/address/{address}/unspent/all`
- [x] [Get Unconfirmed UTXOs by Address](https://docs.whatsonchain.com/api/address#get-unconfirmed-utxos) - `/address/{address}/unconfirmed/unspent`
- [x] [Bulk Unconfirmed UTXOs by Address](https://docs.whatsonchain.com/api/address#bulk-unconfirmed-utxos) - `/addresses/unconfirmed/unspent` (POST)
- [x] [Get Confirmed UTXOs by Address](https://docs.whatsonchain.com/api/address#get-confirmed-utxos) - `/address/{address}/confirmed/unspent`
- [x] [Bulk Confirmed UTXOs by Address](https://docs.whatsonchain.com/api/address#bulk-confirmed-utxos) - `/addresses/confirmed/unspent` (POST)
- [x] [Get Unspent UTXOs by Script](https://docs.whatsonchain.com/api/script#get-script-unspent-transactions) - `/script/{script}/unspent/all`
- [x] [Get Unconfirmed UTXOs by Script](https://docs.whatsonchain.com/api/script#get-unconfirmed-script-utxos) - `/script/{script}/unconfirmed/unspent`
- [x] [Bulk Unconfirmed UTXOs by Script](https://docs.whatsonchain.com/api/script#bulk-unconfirmed-script-utxos) - `/scripts/unconfirmed/unspent` (POST)
- [x] [Get Confirmed UTXOs by Script](https://docs.whatsonchain.com/api/script#get-confirmed-script-utxos) - `/script/{script}/confirmed/unspent`
- [x] [Bulk Confirmed UTXOs by Script](https://docs.whatsonchain.com/api/script#bulk-confirmed-script-utxos) - `/scripts/confirmed/unspent` (POST)
- [x] [Get Unconfirmed Spent Output](https://docs.whatsonchain.com/api/utxo#get-unconfirmed-spent) - `/tx/{hash}/{index}/unconfirmed/spent`
- [x] [Get Confirmed Spent Output](https://docs.whatsonchain.com/api/utxo#get-confirmed-spent) - `/tx/{hash}/{index}/confirmed/spent`
- [x] [Get Spent Transaction Output](https://docs.whatsonchain.com/api/utxo#get-spent-output) - `/tx/{hash}/{index}/spent`
- [x] [Bulk Spent Transaction Outputs](https://docs.whatsonchain.com/api/utxo#bulk-spent-outputs) - `/utxos/spent` (POST)

### ✅ Address (13 endpoints)
- [x] [Get Address Info](https://docs.whatsonchain.com/api/address#get-address-info) - `/address/{address}/info`
- [x] [Get Address Usage Status](https://docs.whatsonchain.com/api/address#get-address-usage) - `/address/{address}/used`
- [x] [Get Associated Scripthashes](https://docs.whatsonchain.com/api/address#get-associated-scripthashes) - `/address/{address}/scripts`
- [x] [Download Statement (PDF)](https://docs.whatsonchain.com/api/address#download-statement) - Statement download
- [x] [Get Unconfirmed Balance](https://docs.whatsonchain.com/api/address#get-unconfirmed-balance) - `/address/{address}/unconfirmed/balance`
- [x] [Bulk Unconfirmed Balance](https://docs.whatsonchain.com/api/address#bulk-unconfirmed-balance) - `/addresses/unconfirmed/balance` (POST)
- [x] [Get Confirmed Balance](https://docs.whatsonchain.com/api/address#get-confirmed-balance) - `/address/{address}/confirmed/balance`
- [x] [Bulk Confirmed Balance](https://docs.whatsonchain.com/api/address#bulk-confirmed-balance) - `/addresses/confirmed/balance` (POST)
- [x] [Get Unconfirmed History](https://docs.whatsonchain.com/api/address#get-unconfirmed-history) - `/address/{address}/unconfirmed/history`
- [x] [Bulk Unconfirmed History](https://docs.whatsonchain.com/api/address#bulk-unconfirmed-history) - `/addresses/unconfirmed/history` (POST)
- [x] [Get Confirmed History](https://docs.whatsonchain.com/api/address#get-confirmed-history) - `/address/{address}/confirmed/history`
- [x] [Bulk Confirmed History](https://docs.whatsonchain.com/api/address#bulk-confirmed-history) - `/addresses/confirmed/history` (POST)
- [x] [Bulk History (All)](https://docs.whatsonchain.com/api/address#bulk-history) - `/addresses/history/all` (POST)

### ✅ Script (6 endpoints)
- [x] [Get Script Usage Status](https://docs.whatsonchain.com/api/script#get-script-usage) - `/script/{script}/used`
- [x] [Get Unconfirmed Script History](https://docs.whatsonchain.com/api/script#get-unconfirmed-script-history) - `/script/{script}/unconfirmed/history`
- [x] [Bulk Unconfirmed Script History](https://docs.whatsonchain.com/api/script#bulk-unconfirmed-script-history) - `/scripts/unconfirmed/history` (POST)
- [x] [Get Confirmed Script History](https://docs.whatsonchain.com/api/script#get-confirmed-script-history) - `/script/{script}/confirmed/history`
- [x] [Bulk Confirmed Script History](https://docs.whatsonchain.com/api/script#bulk-confirmed-script-history) - `/scripts/confirmed/history` (POST)

### ✅ Exchange Rate (2 endpoints)
- [x] [Get Current Exchange Rate](https://docs.whatsonchain.com/api/exchange-rate#get-exchange-rate) - `/exchangerate`
- [x] [Get Historical Exchange Rate](https://docs.whatsonchain.com/api/exchange-rate#get-historical-exchange-rate) - `/exchangerate/historical`

### ✅ Search (1 endpoint)
- [x] [Get Explorer Links](https://docs.whatsonchain.com/api/search#get-explorer-links) - `/search/links` (POST)

### ✅ On-Chain Data (1 endpoint)
- [x] [Get OP_RETURN Data](https://docs.whatsonchain.com/api/onchain-data#get-opreturn-data) - `/tx/{hash}/opreturn`

### ✅ Stats (6 endpoints)
- [x] [Get Block Stats by Height](https://docs.whatsonchain.com/api/stats#get-block-stats-by-height) - `/block/height/{height}/stats`
- [x] [Get Block Stats by Hash](https://docs.whatsonchain.com/api/stats#get-block-stats-by-hash) - `/block/hash/{hash}/stats`
- [x] [Get Miner Block Stats](https://docs.whatsonchain.com/api/stats#get-miner-block-stats) - `/miner/blocks/stats`
- [x] [Get Miner Minimum Fee Rate Stats](https://docs.whatsonchain.com/api/stats#get-miner-fee-stats) - `/miner/fees`
- [x] [Get Miner Summary Stats](https://docs.whatsonchain.com/api/stats#get-miner-summary-stats) - `/miner/summary/stats`
- [x] [Get Tag Count by Height](https://docs.whatsonchain.com/api/stats#get-tag-count-by-height) - `/block/tagcount/height/{height}/stats`

### ✅ Tokens (13 endpoints)

#### 1Sat Ordinals (7 endpoints)
- [x] [Get Token by Origin](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-by-origin) - `/token/1satordinals/{origin}/origin`
- [x] [Get Token by Outpoint](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-by-outpoint) - `/token/1satordinals/{outpoint}`
- [x] [Get Token Content](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-content) - `/token/1satordinals/{outpoint}/content`
- [x] [Get Token Latest Transfer](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-latest-transfer) - `/token/1satordinals/{outpoint}/latest`
- [x] [Get Token Transfer History](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-token-transfers-history) - `/token/1satordinals/{outpoint}/history`
- [x] [Get Tokens by TxID](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-tokens-by-txid) - `/token/1satordinals/tx/{txid}`
- [x] [Get 1Sat Ordinals Stats](https://docs.whatsonchain.com/api/tokens/1sat-ordinals#get-stats) - `/tokens/1satordinals`

#### STAS v0 (6 endpoints)
- [x] [Get All STAS Tokens](https://docs.whatsonchain.com/api/tokens/stas#get-all-tokens) - `/tokens`
- [x] [Get STAS Token by ID](https://docs.whatsonchain.com/api/tokens/stas#get-token-by-id) - `/token/{contractId}/{symbol}`
- [x] [Get Token UTXOs for Address](https://docs.whatsonchain.com/api/tokens/stas#get-token-utxos-for-address) - `/address/{address}/tokens/unspent`
- [x] [Get Address Token Balance](https://docs.whatsonchain.com/api/tokens/stas#get-address-token-balance) - `/address/{address}/tokens`
- [x] [Get Token Transactions](https://docs.whatsonchain.com/api/tokens/stas#get-token-transactions) - `/token/{contractId}/{symbol}/tx`
- [x] [Get STAS Stats](https://docs.whatsonchain.com/api/tokens/stas#get-stats) - `/tokens/stas`

---

## BTC API (64 endpoints)

### ✅ Health (1 endpoint)
- [x] [Get Health Status](https://docs.whatsonchain.com/api/btc/health) - `/woc`

### ✅ Chain Info (4 endpoints)
- [x] [Get Blockchain Info](https://docs.whatsonchain.com/api/btc/chain-info#get-blockchain-info) - `/chain/info`
- [x] [Get Circulating Supply](https://docs.whatsonchain.com/api/btc/chain-info#get-circulating-supply) - `/circulatingsupply`
- [x] [Get Chain Tips](https://docs.whatsonchain.com/api/btc/chain-info#get-chain-tips) - `/chain/tips`
- [x] [Get Peer Info](https://docs.whatsonchain.com/api/btc/chain-info#get-peer-info) - `/peer/info`

### ✅ Block (7 endpoints)
- [x] [Get Block by Hash](https://docs.whatsonchain.com/api/btc/block#get-by-hash) - `/block/hash/{hash}`
- [x] [Get Block by Height](https://docs.whatsonchain.com/api/btc/block#get-by-height) - `/block/height/{height}`
- [x] [Get Block Pages](https://docs.whatsonchain.com/api/btc/block#get-block-pages) - `/block/hash/{hash}/page/{page}`
- [x] [Get Block Headers](https://docs.whatsonchain.com/api/btc/block#get-headers) - `/block/headers`
- [x] [Get Header by Hash or Height](https://docs.whatsonchain.com/api/btc/block#get-header-by-hash) - `/block/{hash}/header`
- [x] [Get Header Bytes File Links](https://docs.whatsonchain.com/api/btc/block#get-header-bytes) - `/block/headers/resources`
- [x] [Get Latest Header Bytes](https://docs.whatsonchain.com/api/btc/block#get-latest-headers) - `/block/headers/latest`

### ✅ Transaction (9 endpoints)
- [x] [Get Transaction by Hash](https://docs.whatsonchain.com/api/btc/transaction#get-by-tx-hash) - `/tx/hash/{hash}`
- [x] [Bulk Transaction Details](https://docs.whatsonchain.com/api/btc/transaction#bulk-transaction-details) - `/txs` (POST)
- [x] [Bulk Transaction Status](https://docs.whatsonchain.com/api/btc/transaction#bulk-transaction-status) - `/txs/status` (POST)
- [x] [Get Transaction as Binary](https://docs.whatsonchain.com/api/btc/transaction#get-tx-binary) - `/tx/{hash}/bin`
- [x] [Get Raw Transaction Data](https://docs.whatsonchain.com/api/btc/transaction#get-raw-tx-data) - `/tx/{hash}/hex`
- [x] [Bulk Raw Transaction Data](https://docs.whatsonchain.com/api/btc/transaction#bulk-raw-tx-data) - `/txs/hex` (POST)
- [x] [Get Raw Transaction Output](https://docs.whatsonchain.com/api/btc/transaction#get-raw-tx-output) - `/tx/{hash}/out/{index}/hex`
- [x] [Bulk Raw Transaction Output Data](https://docs.whatsonchain.com/api/btc/transaction#bulk-raw-tx-output) - `/txs/vouts/hex` (POST)
- [x] [Decode Transaction](https://docs.whatsonchain.com/api/btc/transaction#decode-transaction) - `/tx/decode` (POST)

### ✅ Mempool (2 endpoints)
- [x] [Get Mempool Info](https://docs.whatsonchain.com/api/btc/mempool#get-mempool-info) - `/mempool/info`
- [x] [Get Mempool Transactions](https://docs.whatsonchain.com/api/btc/mempool#get-mempool-transactions) - `/mempool/raw`

### ✅ (Un)Spent Transaction Outputs (14 endpoints)
- [x] [Get Unspent UTXOs by Address](https://docs.whatsonchain.com/api/btc/address#get-unspent-transactions) - `/address/{address}/unspent/all`
- [x] [Get Unconfirmed UTXOs by Address](https://docs.whatsonchain.com/api/btc/address#get-unconfirmed-utxos) - `/address/{address}/unconfirmed/unspent`
- [x] [Bulk Unconfirmed UTXOs by Address](https://docs.whatsonchain.com/api/btc/address#bulk-unconfirmed-utxos) - `/addresses/unconfirmed/unspent` (POST)
- [x] [Get Confirmed UTXOs by Address](https://docs.whatsonchain.com/api/btc/address#get-confirmed-utxos) - `/address/{address}/confirmed/unspent`
- [x] [Bulk Confirmed UTXOs by Address](https://docs.whatsonchain.com/api/btc/address#bulk-confirmed-utxos) - `/addresses/confirmed/unspent` (POST)
- [x] [Get Unspent UTXOs by Script](https://docs.whatsonchain.com/api/btc/script#get-script-unspent-transactions) - `/script/{script}/unspent/all`
- [x] [Get Unconfirmed UTXOs by Script](https://docs.whatsonchain.com/api/btc/script#get-unconfirmed-script-utxos) - `/script/{script}/unconfirmed/unspent`
- [x] [Bulk Unconfirmed UTXOs by Script](https://docs.whatsonchain.com/api/btc/script#bulk-unconfirmed-script-utxos) - `/scripts/unconfirmed/unspent` (POST)
- [x] [Get Confirmed UTXOs by Script](https://docs.whatsonchain.com/api/btc/script#get-confirmed-script-utxos) - `/script/{script}/confirmed/unspent`
- [x] [Bulk Confirmed UTXOs by Script](https://docs.whatsonchain.com/api/btc/script#bulk-confirmed-script-utxos) - `/scripts/confirmed/unspent` (POST)
- [x] [Get Unconfirmed Spent Output](https://docs.whatsonchain.com/api/btc/utxo#get-unconfirmed-spent) - `/tx/{hash}/{index}/unconfirmed/spent`
- [x] [Get Confirmed Spent Output](https://docs.whatsonchain.com/api/btc/utxo#get-confirmed-spent) - `/tx/{hash}/{index}/confirmed/spent`
- [x] [Get Spent Transaction Output](https://docs.whatsonchain.com/api/btc/utxo#get-spent-output) - `/tx/{hash}/{index}/spent`
- [x] [Bulk Spent Transaction Outputs](https://docs.whatsonchain.com/api/btc/utxo#bulk-spent-outputs) - `/utxos/spent` (POST)

### ✅ Address (12 endpoints)
- [x] [Get Address Info](https://docs.whatsonchain.com/api/btc/address#get-address-info) - `/address/{address}/info`
- [x] [Get Address Usage Status](https://docs.whatsonchain.com/api/btc/address#get-address-usage) - `/address/{address}/used`
- [x] [Get Associated Scripthashes](https://docs.whatsonchain.com/api/btc/address#get-associated-scripthashes) - `/address/{address}/scripts`
- [x] [Get Unconfirmed Balance](https://docs.whatsonchain.com/api/btc/address#get-unconfirmed-balance) - `/address/{address}/unconfirmed/balance`
- [x] [Bulk Unconfirmed Balance](https://docs.whatsonchain.com/api/btc/address#bulk-unconfirmed-balance) - `/addresses/unconfirmed/balance` (POST)
- [x] [Get Confirmed Balance](https://docs.whatsonchain.com/api/btc/address#get-confirmed-balance) - `/address/{address}/confirmed/balance`
- [x] [Bulk Confirmed Balance](https://docs.whatsonchain.com/api/btc/address#bulk-confirmed-balance) - `/addresses/confirmed/balance` (POST)
- [x] [Get Unconfirmed History](https://docs.whatsonchain.com/api/btc/address#get-unconfirmed-history) - `/address/{address}/unconfirmed/history`
- [x] [Bulk Unconfirmed History](https://docs.whatsonchain.com/api/btc/address#bulk-unconfirmed-history) - `/addresses/unconfirmed/history` (POST)
- [x] [Get Confirmed History](https://docs.whatsonchain.com/api/btc/address#get-confirmed-history) - `/address/{address}/confirmed/history`
- [x] [Bulk Confirmed History](https://docs.whatsonchain.com/api/btc/address#bulk-confirmed-history) - `/addresses/confirmed/history` (POST)
- [x] [Bulk History (All)](https://docs.whatsonchain.com/api/btc/address#bulk-history) - `/addresses/history/all` (POST)

### ✅ Script (6 endpoints)
- [x] [Get Script Usage Status](https://docs.whatsonchain.com/api/btc/script#get-script-usage) - `/script/{script}/used`
- [x] [Get Unconfirmed Script History](https://docs.whatsonchain.com/api/btc/script#get-unconfirmed-script-history) - `/script/{script}/unconfirmed/history`
- [x] [Bulk Unconfirmed Script History](https://docs.whatsonchain.com/api/btc/script#bulk-unconfirmed-script-history) - `/scripts/unconfirmed/history` (POST)
- [x] [Get Confirmed Script History](https://docs.whatsonchain.com/api/btc/script#get-confirmed-script-history) - `/script/{script}/confirmed/history`
- [x] [Bulk Confirmed Script History](https://docs.whatsonchain.com/api/btc/script#bulk-confirmed-script-history) - `/scripts/confirmed/history` (POST)

### ✅ Exchange Rate (2 endpoints)
- [x] [Get Current Exchange Rate](https://docs.whatsonchain.com/api/btc/exchange-rate#get-exchange-rate) - `/exchangerate`
- [x] [Get Historical Exchange Rate](https://docs.whatsonchain.com/api/btc/exchange-rate#get-historical-exchange-rate) - `/exchangerate/historical`

### ✅ Search (1 endpoint)
- [x] [Get Explorer Links](https://docs.whatsonchain.com/api/btc/search#get-explorer-links) - `/search/links` (POST)

### ✅ Stats (6 endpoints)
- [x] [Get Block Stats by Height](https://docs.whatsonchain.com/api/btc/stats#get-block-stats-by-height) - `/block/height/{height}/stats`
- [x] [Get Block Stats by Hash](https://docs.whatsonchain.com/api/btc/stats#get-block-stats-by-hash) - `/block/hash/{hash}/stats`
- [x] [Get Miner Block Stats](https://docs.whatsonchain.com/api/btc/stats#get-miner-block-stats) - `/miner/blocks/stats`
- [x] [Get Miner Minimum Fee Rate Stats](https://docs.whatsonchain.com/api/btc/stats#get-miner-fee-stats) - `/miner/fees`
- [x] [Get Miner Summary Stats](https://docs.whatsonchain.com/api/btc/stats#get-miner-summary-stats) - `/miner/summary/stats`
- [x] [Get Tag Count by Height](https://docs.whatsonchain.com/api/btc/stats#get-tag-count-by-height) - `/block/tagcount/height/{height}/stats`

---

## WebSockets
- [ ] [New block header event](https://docs.whatsonchain.com/api/websockets#new-block-header)
- [ ] [Block headers history](https://docs.whatsonchain.com/api/websockets#block-headers-history)
- [ ] [Block transactions](https://docs.whatsonchain.com/api/websockets#block-transactions)
- [ ] [Mempool transactions](https://docs.whatsonchain.com/api/websockets#mempool-transactions)
- [ ] [Confirmed transactions](https://docs.whatsonchain.com/api/websockets#confirmed-transactions)
- [ ] [Chain Stats](https://docs.whatsonchain.com/api/websockets#chain-stats)
- [ ] [Customized events](https://docs.whatsonchain.com/api/websockets#customized-events)
</details>

<details>
<summary><strong><code>Development Setup (Getting Started)</code></strong></summary>
<br/>

Install [MAGE-X](https://github.com/mrz1836/mage-x) build tool for development:

```bash
# Install MAGE-X for development and building
go install github.com/mrz1836/mage-x/cmd/magex@latest
magex update:install
```
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

This project uses [goreleaser](https://github.com/goreleaser/goreleaser) for streamlined binary and library deployment to GitHub. To get started, install it via:

```bash
brew install goreleaser
```

The release process is defined in the [.goreleaser.yml](.goreleaser.yml) configuration file.

Then create and push a new Git tag using:

```bash
magex version:bump bump=patch push
```

This process ensures consistent, repeatable releases with properly versioned artifacts and citation metadata.

</details>

<details>
<summary><strong><code>Build Commands</code></strong></summary>
<br/>

View all build commands

```bash script
magex help
```

</details>

<details>
<summary><strong><code>GitHub Workflows</code></strong></summary>
<br/>


### 🎛️ The Workflow Control Center

All GitHub Actions workflows in this repository are powered by configuration files: [**.env.base**](.github/.env.base) (default configuration) and optionally **.env.custom** (project-specific overrides) – your one-stop shop for tweaking CI/CD behavior without touching a single YAML file! 🎯

**Configuration Files:**
- **[.env.base](.github/.env.base)** – Default configuration that works for most Go projects
- **[.env.custom](.github/.env.custom)** – Optional project-specific overrides

This magical file controls everything from:
- **🚀 Go version matrix** (test on multiple versions or just one)
- **🏃 Runner selection** (Ubuntu or macOS, your wallet decides)
- **🔬 Feature toggles** (coverage, fuzzing, linting, race detection, benchmarks)
- **🛡️ Security tool versions** (gitleaks, nancy, govulncheck)
- **🤖 Auto-merge behaviors** (how aggressive should the bots be?)
- **🏷️ PR management rules** (size labels, auto-assignment, welcome messages)

> **Pro tip:** Want to disable code coverage? Just add `ENABLE_CODE_COVERAGE=false` to your .env.custom to override the default in .env.base and push. No YAML archaeology required!

<br/>

| Workflow Name                                                                      | Description                                                                                                            |
|------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------|
| [auto-merge-on-approval.yml](.github/workflows/auto-merge-on-approval.yml)         | Automatically merges PRs after approval and all required checks, following strict rules.                               |
| [codeql-analysis.yml](.github/workflows/codeql-analysis.yml)                       | Analyzes code for security vulnerabilities using [GitHub CodeQL](https://codeql.github.com/).                          |
| [dependabot-auto-merge.yml](.github/workflows/dependabot-auto-merge.yml)           | Automatically merges [Dependabot](https://github.com/dependabot) PRs that meet all requirements.                       |
| [fortress.yml](.github/workflows/fortress.yml)                                     | Runs the GoFortress security and testing workflow, including linting, testing, releasing, and vulnerability checks.    |
| [pull-request-management.yml](.github/workflows/pull-request-management.yml)       | Labels PRs by branch prefix, assigns a default user if none is assigned, and welcomes new contributors with a comment. |
| [scorecard.yml](.github/workflows/scorecard.yml)                                   | Runs [OpenSSF](https://openssf.org/) Scorecard to assess supply chain security.                                        |
| [stale.yml](.github/workflows/stale-check.yml)                                     | Warns about (and optionally closes) inactive issues and PRs on a schedule or manual trigger.                           |
| [sync-labels.yml](.github/workflows/sync-labels.yml)                               | Keeps GitHub labels in sync with the declarative manifest at [`.github/labels.yml`](./.github/labels.yml).             |

</details>

<details>
<summary><strong><code>Updating Dependencies</code></strong></summary>
<br/>

To update all dependencies (Go modules, linters, and related tools), run:

```bash
magex deps:update
```

This command ensures all dependencies are brought up to date in a single step, including Go modules and any managed tools. It is the recommended way to keep your development environment and CI in sync with the latest versions.

</details>

<br/>

## Examples & Tests
## 🧪 Examples & Tests

All unit tests and fuzz tests run via [GitHub Actions](https://github.com/mrz1836/go-whatsonchain/actions) and use [Go version 1.24.x](https://go.dev/doc/go1.24). View the [configuration file](.github/workflows/fortress.yml).

Run all tests (fast):

```bash script
magex test
```

Run all tests with race detector (slower):
```bash script
magex test:race
```

<br/>

## ⚡ Benchmarks

Run the Go benchmarks:

```bash script
magex bench
```

<br/>

## 🛠️ Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## 🤖 AI Compliance
This project documents expectations for AI assistants using a few dedicated files:

- [AGENTS.md](.github/AGENTS.md) — canonical rules for coding style, workflows, and pull requests used by [Codex](https://chatgpt.com/codex).
- [CLAUDE.md](.github/CLAUDE.md) — quick checklist for the [Claude](https://www.anthropic.com/product) agent.
- [.cursorrules](.cursorrules) — machine-readable subset of the policies for [Cursor](https://www.cursor.so/) and similar tools.
- [sweep.yaml](.github/sweep.yaml) — rules for [Sweep](https://github.com/sweepai/sweep), a tool for code review and pull request management.

Edit `AGENTS.md` first when adjusting these policies, and keep the other files in sync within the same pull request.

<br/>

## 👥 Maintainers
| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:------------------------------------------------------------------------------------------------:|
|                                [MrZ](https://github.com/mrz1836)                                 |

<br/>

## 🤝 Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and please follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/mrz1836) :clap:
or by making a [**bitcoin donation**](https://mrz1818.com/?tab=tips&utm_source=github&utm_medium=sponsor-link&utm_campaign=go-whatsonchain&utm_term=go-whatsonchain&utm_content=go-whatsonchain) to ensure this journey continues indefinitely! :rocket:


[![Stars](https://img.shields.io/github/stars/mrz1836/go-whatsonchain?label=Please%20like%20us&style=social)](https://github.com/mrz1836/go-whatsonchain/stargazers)

<br/>

## 📝 License

[![License](https://img.shields.io/github/license/mrz1836/go-whatsonchain.svg?style=flat)](LICENSE)

