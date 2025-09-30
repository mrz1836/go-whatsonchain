# CLAUDE.md

> **Quick Reference Guide for Claude Code**
> This document helps Claude Code understand, navigate, and contribute to **go-whatsonchain** effectively.

## 🎯 What This Project Does

**go-whatsonchain** is the unofficial Go SDK for the [whatsonchain.com API](https://docs.whatsonchain.com/), providing programmatic access to blockchain data for both **BSV** (Bitcoin SV) and **BTC** (Bitcoin Core) blockchains.

### Key Capabilities
- **135 API endpoints**: 71 BSV + 64 BTC endpoints covering blocks, transactions, addresses, UTXOs, mempool, tokens, and stats
- **Multi-chain support**: Switch between BSV and BTC with a single configuration option
- **Zero dependencies**: No external dependencies except `testify` for tests
- **Production-ready**: Retry logic, exponential backoff, rate limiting, context support throughout
- **Type-safe**: Comprehensive type definitions for all API responses

### Use Cases
Developers use this library to:
- Query blockchain data (blocks, transactions, addresses)
- Broadcast transactions to the network
- Track UTXOs and balances
- Monitor mempool and network stats
- Work with BSV-specific features (OP_RETURN, 1Sat Ordinals, STAS tokens)

---

## 🏗️ Architecture Overview

### Core Design Patterns

#### 1. **Interface-Based Service Organization**
The codebase uses interface segregation to organize functionality:

```go
type ClientInterface interface {
    AddressService      // Address-related operations
    BlockService        // Block queries
    TransactionService  // Transaction operations
    ChainService        // Chain info & network data
    MempoolService      // Mempool operations
    ScriptService       // Script-based queries
    StatsService        // Statistics endpoints
    TokenService        // Token operations (BSV)
    BSVService         // BSV-specific methods
    BTCService         // BTC-specific methods
    // + Getters/Setters
}
```

**Why this matters**: When adding new functionality, implement the appropriate service interface. All services are implemented on the same `*Client` struct.

#### 2. **Functional Options Pattern**
Client configuration uses functional options for clean, extensible API:

```go
client, err := whatsonchain.NewClient(
    context.Background(),
    whatsonchain.WithChain(whatsonchain.ChainBSV),
    whatsonchain.WithNetwork(whatsonchain.NetworkMain),
    whatsonchain.WithAPIKey("your-key"),
    whatsonchain.WithRateLimit(10),
)
```

**Implementation**: All options are defined in `client.go` as `ClientOption` functions that modify `clientOptions`.

#### 3. **Generic Request Helpers**
The codebase uses Go 1.18+ generics to avoid repetition:

```go
// request_helpers.go
func requestAndUnmarshal[T any](ctx context.Context, c *Client, url, method string, payload []byte, emptyErr error) (*T, error)
func requestAndUnmarshalSlice[T any](ctx context.Context, c *Client, url, method string, payload []byte, emptyErr error) ([]T, error)
```

**When to use**: Most API methods follow this pattern. See `transactions.go`, `blocks.go`, `addresses.go` for examples.

#### 4. **URL Construction**
Centralized URL building with chain/network awareness:

```go
// url_builder.go
func (c *Client) buildURL(path string, args ...interface{}) string {
    baseURL := fmt.Sprintf("%s%s/%s", apiEndpointBase, c.Chain(), c.Network())
    if len(args) > 0 {
        path = fmt.Sprintf(path, args...)
    }
    return baseURL + path
}
```

**Usage**: `url := c.buildURL("/tx/hash/%s", txHash)` → `https://api.whatsonchain.com/v1/bsv/main/tx/hash/abc123`

#### 5. **Context-First Design**
Every API method accepts `context.Context` as the first parameter for cancellation, timeouts, and tracing.

---

## 📁 File Organization

```
.
├── client.go              # Client struct, options, configuration
├── interface.go           # All service interfaces
├── definitions.go         # Type definitions for API responses
├── whatsonchain.go        # NewClient constructor, core request logic
├── url_builder.go         # Centralized URL construction
├── request_helpers.go     # Generic request/unmarshal helpers
├── errors.go              # Custom error definitions
├── http_client.go         # Retry logic & exponential backoff
├── addresses.go           # AddressService implementation
├── blocks.go              # BlockService implementation
├── transactions.go        # TransactionService implementation
├── chain_info.go          # ChainService implementation
├── mempool.go             # MempoolService implementation
├── scripts.go             # ScriptService implementation
├── stats.go               # StatsService implementation
├── tokens.go              # TokenService implementation (BSV)
├── bsv_specific.go        # BSV-only methods
├── btc_specific.go        # BTC-only methods (currently minimal)
├── exchange_rates.go      # Exchange rate endpoints
├── search.go              # Search/explorer endpoints
├── health.go              # Health check endpoint
└── examples/              # Usage examples
```

---

## 🔧 Working with This Codebase

### Adding a New Endpoint

1. **Determine the service category**: Address, Block, Transaction, etc.
2. **Add method to interface** in `interface.go`:
   ```go
   type TransactionService interface {
       // ... existing methods
       GetNewFeature(ctx context.Context, id string) (*FeatureResponse, error)
   }
   ```

3. **Define response type** in `definitions.go`:
   ```go
   type FeatureResponse struct {
       ID     string `json:"id"`
       Status string `json:"status"`
   }
   ```

4. **Implement method** in the appropriate file (e.g., `transactions.go`):
   ```go
   func (c *Client) GetNewFeature(ctx context.Context, id string) (*FeatureResponse, error) {
       url := c.buildURL("/feature/%s", id)
       return requestAndUnmarshal[FeatureResponse](ctx, c, url, http.MethodGet, nil, ErrFeatureNotFound)
   }
   ```

5. **Add error constant** to `errors.go` if needed:
   ```go
   var ErrFeatureNotFound = errors.New("feature not found")
   ```

6. **Write tests** in `*_test.go` with table-driven tests
7. **Add example** in `examples/` directory if useful

### Chain-Specific Features

For BSV-only or BTC-only endpoints:

```go
func (c *Client) GetBSVSpecificFeature(ctx context.Context) error {
    if c.Chain() != ChainBSV {
        return ErrBSVChainRequired
    }
    // implementation
}
```

See `bsv_specific.go:18` for `GetOpReturnData()` as an example.

### Bulk Operations with Rate Limiting

When implementing bulk processors (methods that handle batching automatically):

```go
func (c *Client) BulkOperationProcessor(ctx context.Context, items []string) (results Results, err error) {
    // 1. Chunk items into batches (typically MaxTransactionsUTXO = 20)
    batches := chunkItems(items, MaxTransactionsUTXO)

    var currentRateLimit int
    for _, batch := range batches {
        // 2. Process batch
        result, err := c.BulkOperation(ctx, batch)
        if err != nil {
            return nil, err
        }
        results = append(results, result...)

        // 3. Respect rate limiting
        currentRateLimit++
        if currentRateLimit >= c.RateLimit() {
            time.Sleep(1 * time.Second)
            currentRateLimit = 0
        }
    }
    return results, nil
}
```

See `transactions.go:38` for `BulkTransactionDetailsProcessor()` as a reference.

---

## 🧪 Testing

### Test Structure
- **Table-driven tests**: Use `testCase` structs with `name`, `args`, `want`, `wantErr`
- **Mock HTTP client**: Tests use a mock client to avoid real API calls
- **Interface compliance**: `interface_test.go` ensures `Client` implements all service interfaces

### Running Tests
```bash
magex test           # Fast tests
magex test:race      # With race detector
magex bench          # Benchmarks
magex coverage       # Coverage report
```

### Example Test Pattern
```go
func TestClient_GetTxByHash(t *testing.T) {
    tests := []struct {
        name    string
        hash    string
        wantErr bool
    }{
        {"valid transaction", "abc123", false},
        {"not found", "missing", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

---

## 🎨 Code Style & Conventions

### General Principles
- **Context-first**: All API methods accept `context.Context` as first parameter
- **Error handling**: Return `error` as last return value; use custom errors from `errors.go`
- **Pointer receivers**: All methods use pointer receivers on `*Client`
- **No panics**: Never panic; return errors instead
- **Godoc comments**: Every exported function, type, and constant must have a godoc comment

### Naming Conventions
- **Methods**: Use descriptive names like `GetTxByHash`, `BulkTransactionDetails`
- **Types**: Use meaningful names like `TxInfo`, `BlockInfo`, `AddressBalance`
- **Constants**: Use `MaxTransactionsUTXO`, `ChainBSV`, `NetworkMain`
- **Errors**: Prefix with `Err`: `ErrTransactionNotFound`, `ErrMaxUTXOsExceeded`

### Documentation
Every exported function should have:
```go
// GetTxByHash retrieves transaction details by hash
//
// For more information: https://docs.whatsonchain.com/#get-by-tx-hash
func (c *Client) GetTxByHash(ctx context.Context, hash string) (*TxInfo, error)
```

---

## 🚨 Common Gotchas

1. **Chain-specific methods**: Always check if operation requires BSV/BTC and return appropriate error
2. **Max limits**: Respect API limits (MaxTransactionsUTXO=20, MaxAddressesForLookup=20)
3. **Empty responses**: Use the `emptyErr` parameter in request helpers for proper error handling
4. **Rate limiting**: Bulk processors must implement rate limiting to avoid API throttling
5. **URL building**: Always use `c.buildURL()` to ensure chain/network are included
6. **Context usage**: Pass context through to all HTTP requests for proper cancellation

---

## 📚 Reference Documentation

### For Detailed Standards
All technical conventions for this repository are documented in [**AGENTS.md**](AGENTS.md), which covers:
- Go coding essentials (interfaces, goroutines, error handling)
- Testing standards and coverage requirements
- Commit and branch conventions
- Pull request guidelines
- Release workflow and versioning
- CI/CD and validation requirements
- Dependency management
- Security practices

**When in doubt, consult AGENTS.md first.**

### Key Files to Reference
- **`client.go`**: Client configuration, options pattern, getters/setters
- **`interface.go`**: Complete service interface definitions
- **`definitions.go`**: All type definitions for API responses
- **`request_helpers.go`**: Generic request patterns
- **`transactions.go`**: Example of bulk operations with rate limiting
- **`examples/`**: Real-world usage patterns

---

## 🤝 Contributing Workflow

1. **Create feature branch**: `git checkout -b feature/new-endpoint`
2. **Make changes**: Follow patterns in existing code
3. **Write tests**: Maintain 100% coverage for new code
4. **Run validation**: `magex test` and `magex lint`
5. **Commit**: Follow [commit conventions](tech-conventions/commit-branch-conventions.md)
6. **Create PR**: Reference the [PR guidelines](tech-conventions/pull-request-guidelines.md)

### Before Committing
```bash
magex test           # Ensure all tests pass
magex lint           # Run linters (golangci-lint)
magex fmt            # Format code
```

---

## 💡 Pro Tips for Claude Code

1. **When adding endpoints**: Use existing implementations as templates. `transactions.go` is comprehensive.
2. **For bulk operations**: Always implement both single and "Processor" versions (see `BulkTransactionDetails` vs `BulkTransactionDetailsProcessor`)
3. **Error handling**: Use custom errors from `errors.go`; add new ones if needed
4. **Testing**: Mirror the existing test patterns; use table-driven tests
5. **Chain support**: If an endpoint works for both chains, implement it generally. If chain-specific, guard with chain checks.
6. **Documentation**: Include link to WhatsonChain API docs in godoc comments
7. **Zero dependencies**: Don't add external dependencies without strong justification

---

## 📖 Quick Examples

### Basic Usage
```go
// Create client for BSV mainnet
client, _ := whatsonchain.NewClient(
    context.Background(),
    whatsonchain.WithChain(whatsonchain.ChainBSV),
    whatsonchain.WithNetwork(whatsonchain.NetworkMain),
)

// Get transaction
tx, _ := client.GetTxByHash(ctx, "abc123...")

// Get address balance
balance, _ := client.AddressBalance(ctx, "1A1zP1...")

// Broadcast transaction
txID, _ := client.BroadcastTx(ctx, "01000000...")
```

### Advanced Configuration
```go
client, _ := whatsonchain.NewClient(
    context.Background(),
    whatsonchain.WithChain(whatsonchain.ChainBSV),
    whatsonchain.WithAPIKey("your-key"),
    whatsonchain.WithRateLimit(10),
    whatsonchain.WithRequestTimeout(60*time.Second),
    whatsonchain.WithRequestRetryCount(3),
    whatsonchain.WithBackoff(2*time.Millisecond, 10*time.Millisecond, 2.0, 2*time.Millisecond),
)
```

### Bulk Operations
```go
// Process 100+ transactions automatically in batches of 20
hashes := &whatsonchain.TxHashes{
    TxIDs: []string{"hash1", "hash2", ..., "hash150"},
}
txList, _ := client.BulkTransactionDetailsProcessor(ctx, hashes)
```

---

## 🔗 External Resources

- **API Documentation**: https://docs.whatsonchain.com/
- **GoDoc**: https://pkg.go.dev/github.com/mrz1836/go-whatsonchain
- **Repository**: https://github.com/mrz1836/go-whatsonchain
- **Technical Conventions**: [.github/AGENTS.md](AGENTS.md)

---

**Last Updated**: 2025-09-30
**For Questions**: See [CONTRIBUTING.md](CONTRIBUTING.md) or open an issue
