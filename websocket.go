package whatsonchain

import (
	"fmt"

	"github.com/centrifugal/centrifuge-go"
)

type SocketHandler interface {
	OnConnect(_ *centrifuge.Client, _ centrifuge.ConnectEvent)
	OnError(_ *centrifuge.Client, e centrifuge.ErrorEvent)
	OnDisconnect(_ *centrifuge.Client, e centrifuge.DisconnectEvent)
	OnMessage(_ *centrifuge.Client, e centrifuge.MessageEvent)
	OnServerPublish(_ *centrifuge.Client, e centrifuge.ServerPublishEvent)
}

func (c *Client) NewMempoolWebsocket(handler SocketHandler) *centrifuge.Client {
	return newWebsocketClient(fmt.Sprintf("%s%s", socketEndpoint, "mempool"), handler)
}

func newWebsocketClient(url string, handler SocketHandler) *centrifuge.Client {
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
