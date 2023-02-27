module github.com/mrz1836/go-whatsonchain

go 1.17

require (
	github.com/centrifugal/centrifuge-go v0.9.4
	github.com/gojektech/heimdall/v6 v6.1.0
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/centrifugal/protocol v0.9.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gojektech/valkyrie v0.0.0-20190210220504-8f62c1e7ba45 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/segmentio/encoding v0.3.6 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Breaking changes - needs a full refactor in WOC and BUX
replace github.com/centrifugal/centrifuge-go => github.com/centrifugal/centrifuge-go v0.8.3
