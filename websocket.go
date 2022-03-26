package whatsonchain

import (
	"github.com/centrifugal/centrifuge-go"
)

const (
	socketEndpointMempool = "mempool"
)

// socketHandler describe the interface
type socketHandler interface {
	OnConnect(*centrifuge.Client, centrifuge.ConnectEvent)
	OnError(*centrifuge.Client, centrifuge.ErrorEvent)
	OnDisconnect(*centrifuge.Client, centrifuge.DisconnectEvent)
	OnMessage(*centrifuge.Client, centrifuge.MessageEvent)
	OnServerPublish(*centrifuge.Client, centrifuge.ServerPublishEvent)
}

// NewMempoolWebsocket instantiates a new websocket client to stream mempool transactions
func (c *Client) NewMempoolWebsocket(handler socketHandler) *centrifuge.Client {
	return newWebsocketClient(socketEndpoint+socketEndpointMempool, handler)
}

// newWebsocketClient will create a new websocket client
func newWebsocketClient(url string, handler socketHandler) (client *centrifuge.Client) {
	if url == "" || handler == nil {
		return
	}
	if client = centrifuge.NewJsonClient(
		url, centrifuge.DefaultConfig(),
	); client == nil {
		return
	}
	client.OnDisconnect(handler)
	client.OnConnect(handler)
	client.OnServerPublish(handler)
	client.OnError(handler)
	return
}
