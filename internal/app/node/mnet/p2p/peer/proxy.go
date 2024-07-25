package peer

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
)

func ProxyConnect(p2pHost host.Host, hop *NetHop) error {
	pm := newPeerAddrInfoMapFromNetHop(hop)

	peerInfo := connectPeerGroup(p2pHost, pm.peer)
	if peerInfo != nil {
		return nil
	}

	return fmt.Errorf("unable to connect to proxy")
}
