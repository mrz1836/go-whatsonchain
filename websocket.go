package whatsonchain

import (
	"fmt"

	"github.com/centrifugal/centrifuge-go"
)

// socketHandler describe the interface
type socketHandler interface {
	OnConnect(_ *centrifuge.Client, _ centrifuge.ConnectEvent)
	OnError(_ *centrifuge.Client, e centrifuge.ErrorEvent)
	OnDisconnect(_ *centrifuge.Client, e centrifuge.DisconnectEvent)
	OnMessage(_ *centrifuge.Client, e centrifuge.MessageEvent)
	OnServerPublish(_ *centrifuge.Client, e centrifuge.ServerPublishEvent)
}

// NewMempoolWebsocket instantiates a new websocket client to stream mempool transactions
func (c *Client) NewMempoolWebsocket(handler socketHandler) *centrifuge.Client {
	return newWebsocketClient(fmt.Sprintf("%s%s", socketEndpoint, "mempool"), handler)
}

func newWebsocketClient(url string, handler socketHandler) *centrifuge.Client {
	if url == "" || handler == nil {
		return nil
	}
	c := centrifuge.NewJsonClient(url, centrifuge.DefaultConfig())
	c.OnDisconnect(handler)
	c.OnConnect(handler)
	c.OnServerPublish(handler)
	c.OnError(handler)
	return c
}
