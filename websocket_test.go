package whatsonchain

import (
	"testing"

	"github.com/centrifugal/centrifuge-go"
	"github.com/stretchr/testify/assert"
)

func TestClient_NewMempoolWebsocket(t *testing.T) {
	type args struct {
		handler SocketHandler
	}
	tests := []struct {
		name string
		args args
		want *centrifuge.Client
	}{
		{"nil handler should set nil client",
			args{
				handler: nil,
			},
			nil,
		},
	}
	client := newMockClient(&mockHTTPAddresses{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, client.NewMempoolWebsocket(tt.args.handler), "NewMempoolWebsocket(%v)", tt.args.handler)
		})
	}
}

func Test_newWebsocketClient(t *testing.T) {
	type args struct {
		url     string
		handler SocketHandler
	}
	tests := []struct {
		name string
		args args
		want *centrifuge.Client
	}{
		{
			name: "empty url should set nil client",
			args: args{
				url:     "",
				handler: nil,
			},
			want: nil,
		},
		{
			name: "nil handler should set nil client",
			args: args{
				url:     "wss://socket.whatsonchain.com/mempool",
				handler: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, newWebsocketClient(tt.args.url, tt.args.handler), "newWebsocketClient(%v, %v)", tt.args.url, tt.args.handler)
		})
	}
}
