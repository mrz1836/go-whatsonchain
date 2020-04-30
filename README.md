# go-whatsonchain
> The unofficial golang implementation for the [whatsonchain.com API](https://developers.whatsonchain.com/)

[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-whatsonchain)](https://golang.org/)
[![Build Status](https://travis-ci.com/mrz1836/go-whatsonchain.svg?branch=master&v=2)](https://travis-ci.com/mrz1836/go-whatsonchain)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-whatsonchain?style=flat&v=2)](https://goreportcard.com/report/github.com/mrz1836/go-whatsonchain)
[![codecov](https://codecov.io/gh/mrz1836/go-whatsonchain/branch/master/graph/badge.svg)](https://codecov.io/gh/mrz1836/go-whatsonchain)
[![Release](https://img.shields.io/github/release-pre/mrz1836/go-whatsonchain.svg?style=flat&v=1)](https://github.com/mrz1836/go-whatsonchain/releases)
[![GoDoc](https://godoc.org/github.com/mrz1836/go-whatsonchain?status.svg&style=flat)](https://pkg.go.dev/github.com/mrz1836/go-whatsonchain)

## Table of Contents
- [Installation](#installation)
- [Documentation](#documentation)
- [Examples & Tests](#examples--tests)
- [Benchmarks](#benchmarks)
- [Code Standards](#code-standards)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Installation

**go-whatsonchain** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/mrz1836/go-whatsonchain
```

## Documentation
You can view the generated [documentation here](https://pkg.go.dev/github.com/mrz1836/go-whatsonchain).

### Features
- [Client](client.go) is completely configurable
- Customize the network per request (`main`, `test` or `stn`)
- Using [heimdall http client](https://github.com/gojek/heimdall) with exponential backoff & more
- Current coverage for the [whatsonchain.com API](https://developers.whatsonchain.com/)
    - [x] Health
        - [x] Get API Status
    - [x] Chain Info
        - [x] Get Blockchain Info
        - [x] Get Circulating Supply
    - [x] Block
        - [x] Get by Hash
        - [x] Get by Height
        - [x] Get Block Pages
    - [x] Transaction
        - [x] Get by TX Hash
        - [x] Broadcast Transaction
        - [x] Bulk Broadcast
        - [x] Bulk Transaction Details
        - [x] Decode Transaction
        - [x] Download Receipt
        - [x] Get Raw Transaction Data
        - [x] Get Raw Transaction Output
        - [x] Get Merkle Proof
    - [x] Mempool
        - [x] Get Mempool Info
        - [x] Get Mempool Transactions
    - [x] Address
        - [x] Get Address Info
        - [x] Get Balance
        - [x] Get History
        - [x] Get Unspent Transactions
        - [x] Get Unspent Transaction Details (Custom)
        - [x] Download Statement
    - [x] Script
        - [x] Get Script History
        - [x] Get Script Unspent Transactions
    - [x] Exchange Rate
        - [x] Get Exchange Rate
    - [x] Search
        - [x] Get Explorer Links

<details>
<summary><strong><code>Library Deployment</code></strong></summary>

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to Github and can be installed via: `brew install goreleaser`.

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production.
</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
bench                          Run all benchmarks in the Go application
clean                          Remove previous builds and any test cache data
clean-mods                     Remove all the Go mod cache
coverage                       Shows the test coverage
godocs                         Sync the latest tag with GoDocs
help                           Show all make commands available
lint                           Run the Go lint application
release                        Full production release (creates release in Github)
release-test                   Full production test release (everything except deploy)
release-snap                   Test the full release (build binaries)
run-examples                   Runs all the examples
tag                            Generate a new tag and push (IE: make tag version=0.0.0)
tag-remove                     Remove a tag if found (IE: make tag-remove version=0.0.0)
tag-update                     Update an existing tag to current commit (IE: make tag-update version=0.0.0)
test                           Runs vet, lint and ALL tests
test-short                     Runs vet, lint and tests (excludes integration tests)
update                         Update all project dependencies
update-releaser                Update the goreleaser application
vet                            Run the Go vet application
```
</details>

## Examples & Tests
All unit tests and [examples](whatsonchain_test.go) run via [Travis CI](https://travis-ci.org/mrz1836/go-whatsonchain) and uses [Go version 1.14.x](https://golang.org/doc/go1.14). View the [deployment configuration file](.travis.yml).

Run all tests (including integration tests)
```shell script
make test
```

Run tests (excluding integration tests)
```shell script
make test-short
```

## Benchmarks
Run the Go [benchmarks](whatsonchain_test.go):
```shell script
make bench
```

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

## Usage
- View the [whatsonchain examples](whatsonchain_test.go)

Basic implementation:
```go
package main

import (

    "fmt"
    
    "github.com/mrz1836/go-whatsonchain"
)

func main() {

    // Create a client
    client, _ := whatsonchain.NewClient(whatsonchain.NetworkMain, nil)

    // Get a balance for an address
    balance, _ := client.AddressBalance("16ZqP5Tb22KJuvSAbjNkoiZs13mmRmexZA")
    fmt.Println("confirmed balance", balance.Confirmed)
}
```

## Maintainers

| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:---:|
| [MrZ](https://github.com/mrz1836) |

## Contributing

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

Support the development of this project 🙏 and support [WhatsOnChain](https://tncpw.co/65733e42)

[![Donate](https://img.shields.io/badge/donate-bitcoin-brightgreen.svg)](https://mrz1818.com/?tab=tips&af=go-whatsonchain) [![Support](https://img.shields.io/badge/support-WhatsOnChain-blue.svg)](https://whatsonchain.com/support)

#### Credits

[WhatsOnChain](https://tncpw.co/65733e42) for their hard work on the Whatsonchain API

[AustEcon's Python Version](https://github.com/AustEcon/whatsonchain)

## License

![License](https://img.shields.io/github/license/mrz1836/go-whatsonchain.svg?style=flat)