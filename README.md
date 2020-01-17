# go-whatsonchain
**go-whatsonchain** is the unofficial golang implementation for the [whatsonchain.com API](https://developers.whatsonchain.com/)

[![Go](https://img.shields.io/github/go-mod/go-version/mrz1836/go-whatsonchain)](https://golang.org/)
[![Build Status](https://travis-ci.com/mrz1836/go-whatsonchain.svg?branch=master&v=2)](https://travis-ci.com/mrz1836/go-whatsonchain)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-whatsonchain?style=flat&v=2)](https://goreportcard.com/report/github.com/mrz1836/go-whatsonchain)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/f9815e59758743b9adca25c11558ab1c)](https://www.codacy.com/app/mrz1818/go-whatsonchain?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mrz1836/go-whatsonchain&amp;utm_campaign=Badge_Grade)
[![Release](https://img.shields.io/github/release-pre/mrz1836/go-whatsonchain.svg?style=flat&v=1)](https://github.com/mrz1836/go-whatsonchain/releases)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat)](https://github.com/RichardLitt/standard-readme)
[![GoDoc](https://godoc.org/github.com/mrz1836/go-whatsonchain?status.svg&style=flat)](https://godoc.org/github.com/mrz1836/go-whatsonchain)

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
```bash
$ go get -u github.com/mrz1836/go-whatsonchain
```

## Documentation
You can view the generated [documentation here](https://godoc.org/github.com/mrz1836/go-whatsonchain).

### Features
- [Client](client.go) is completely configurable
- Customize the network per request (`main`, `test` or `stn`)
- Using [heimdall http client](https://github.com/gojek/heimdall) with exponential backoff & more
- Current coverage for the [whatsonchain.com API](https://developers.whatsonchain.com/)
    - [x] Health
    - [x] Chain Info
    - [x] Blocks
    - [x] Transactions
    - [x] Addresses
    - [ ] Receipt
    - [ ] Statement
    - [ ] Mempool
    - [ ] Search

## Examples & Tests
All unit tests and [examples](whatsonchain_test.go) run via [Travis CI](https://travis-ci.org/mrz1836/go-whatsonchain) and uses [Go version 1.13.x](https://golang.org/doc/go1.13). View the [deployment configuration file](.travis.yml).

Run all tests (including integration tests)
```bash
$ cd ../go-whatsonchain
$ go test ./... -v
```

Run tests (excluding integration tests)
```bash
$ cd ../go-whatsonchain
$ go test ./... -v -test.short
```

## Benchmarks
Run the Go [benchmarks](whatsonchain_test.go):
```bash
$ cd ../go-whatsonchain
$ go test -bench . -benchmem
```

## Code Standards
Read more about this Go project's [code standards](CODE_STANDARDS.md).

## Usage
- View the [whatsonchain examples](whatsonchain_test.go)

Basic implementation:
```golang
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

Support the development of this project 🙏 & support [WhatsOnChain](https://tncpw.co/65733e42)

[![Donate](https://img.shields.io/badge/donate-bitcoin-brightgreen.svg)](https://mrz1818.com/?tab=tips&af=go-whatsonchain) [![Support](https://img.shields.io/badge/support-WhatsOnChain-blue.svg)](https://whatsonchain.com/support)

#### Credits

[WhatsOnChain](https://tncpw.co/65733e42) for their hard work on the Whatsonchain API

[AustEcon's Python Version](https://github.com/AustEcon/whatsonchain)

## License

![License](https://img.shields.io/github/license/mrz1836/go-whatsonchain.svg?style=flat)