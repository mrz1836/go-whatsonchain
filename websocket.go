package whatsonchain

import (
	"github.com/centrifugal/centrifuge-go"
)

const (
	socketEndpointMempool = "mempool"
)

// SocketHandler describe the interface
type SocketHandler interface {
	OnConnect(*centrifuge.Client, centrifuge.ConnectEvent)
	OnDisconnect(*centrifuge.Client, centrifuge.DisconnectEvent)
	OnError(*centrifuge.Client, centrifuge.ErrorEvent)
	OnMessage(*centrifuge.Client, centrifuge.MessageEvent)
	OnServerJoin(*centrifuge.Client, centrifuge.ServerJoinEvent)
	OnServerLeave(*centrifuge.Client, centrifuge.ServerLeaveEvent)
	OnServerPublish(*centrifuge.Client, centrifuge.ServerPublishEvent)
	OnServerSubscribe(*centrifuge.Client, centrifuge.ServerSubscribeEvent)
	OnServerUnsubscribe(*centrifuge.Client, centrifuge.ServerUnsubscribeEvent)
	OnPublish(*centrifuge.Subscription, centrifuge.PublishEvent)
	OnJoin(*centrifuge.Subscription, centrifuge.JoinEvent)
	OnLeave(*centrifuge.Subscription, centrifuge.LeaveEvent)
	OnSubscribeSuccess(*centrifuge.Subscription, centrifuge.SubscribeSuccessEvent)
	OnSubscribeError(*centrifuge.Subscription, centrifuge.SubscribeErrorEvent)
	OnUnsubscribe(*centrifuge.Subscription, centrifuge.UnsubscribeEvent)
}

// NewMempoolWebsocket instantiates a new websocket client to stream mempool transactions
func (c *Client) NewMempoolWebsocket(handler SocketHandler) *centrifuge.Client {
	return newWebsocketClient(socketEndpoint+socketEndpointMempool, handler)
}

// newWebsocketClient will create a new websocket client
func newWebsocketClient(url string, handler SocketHandler) (client *centrifuge.Client) {
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
