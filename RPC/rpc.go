package RPCClient

import (
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
)

func StartWebsocket() error {
	return Connectivity.Socket.StartWebsocket()
}
