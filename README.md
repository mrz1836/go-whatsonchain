# go-whatsonchain
**go-whatsonchain** is the unofficial golang implementation for the whatsonchain.com API

[![Build Status](https://travis-ci.org/mrz1836/go-whatsonchain.svg?branch=master)](https://travis-ci.org/mrz1836/go-whatsonchain)
[![Report](https://goreportcard.com/badge/github.com/mrz1836/go-whatsonchain?style=flat)](https://goreportcard.com/report/github.com/mrz1836/go-whatsonchain)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/01708ca3079e4933bafb3b39fe2aaa9d)](https://www.codacy.com/app/mrz1818/go-whatsonchain?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mrz1836/go-whatsonchain&amp;utm_campaign=Badge_Grade)
[![Release](https://img.shields.io/github/release-pre/mrz1836/go-whatsonchain.svg?style=flat)](https://github.com/mrz1836/go-whatsonchain/releases)
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

**go-whatsonchain** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy) and [dep](https://github.com/golang/dep).
```bash
$ go get -u github.com/mrz1836/go-whatsonchain
```

Updating dependencies in **go-whatsonchain**:
```bash
$ cd ../go-whatsonchain
$ dep ensure -update -v
```

## Documentation
You can view the generated [documentation here](https://godoc.org/github.com/mrz1836/go-whatsonchain).

### Features
- Complete coverage for the [whatsonchain.com](https://developers.whatsonchain.com/) API
- Client is completely configurable
- Customize API Key and User Agent per request
- Using [heimdall http client](https://github.com/gojek/heimdall) with exponential backoff & more

## Examples & Tests
All unit tests and [examples](whatsonchain_test.go) run via [Travis CI](https://travis-ci.org/mrz1836/go-whatsonchain) and uses [Go version 1.13.x](https://golang.org/doc/go1.13). View the [deployment configuration file](.travis.yml).

- [helper examples & tests](helper_test.go)
- [whatsonchain examples &  tests](whatsonchain_test.go)
- [response tests](response_test.go)

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
- View the [whatsonchain examples & benchmarks](whatsonchain_test.go)
- View the [helper examples & benchmarks](helper_test.go)
- View the [response tests](response_test.go)

Basic implementation:
```golang
package main

import (
	"fmt"
	"log"
	"os"
	"github.com/mrz1836/go-whatsonchain"
)

func main() {

    // Create a client with your api key
    client, _ := whatsonchain.NewClient("your-api-key")


    // Use the whatsonchain response
    fmt.Println(response.Person.Names[0].Display)
    // Output: Jeff Preston Bezos
}
```

## Maintainers

[@MrZ1836](https://github.com/mrz1836)

## Contributing

If you're looking for a python version, checkout [this package](https://github.com/AustEcon/whatsonchain)

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

Support the development of this project 🙏

[![Donate](https://img.shields.io/badge/donate-bitcoin-brightgreen.svg)](https://mrz1818.com/?tab=tips&af=go-whatsonchain)

## License

![License](https://img.shields.io/github/license/mrz1836/go-whatsonchain.svg?style=flat)
