package peer

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/net/swarm"
	"n2x.dev/x-lib/pkg/errors"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/mnet/p2p"
	"n2x.dev/x-node/internal/app/node/mnet/p2p/conn"
)

func NewStream(p2pHost host.Host, hop *NetHop) (network.Stream, error) {
	peerInfo, err := connectPeer(p2pHost, hop)
	if err != nil {
		return nil, errors.Wrapf(err, "[%v] function connectPeer()", errors.Trace())
	}

	conns := p2pHost.Network().ConnsToPeer(peerInfo.ID)

	xlog.Infof("Peer %s CONNECTED (%d conns)", peerInfo.ID.ShortString(), len(conns))

	streams := make([]network.Stream, 0)

	limitedConnection := false

	for _, c := range conns {
		streams = append(streams, c.GetStreams()...)

		if c.Stat().Limited {
			limitedConnection = true
		}

		conn.Log(c)
	}

	if len(streams) > 0 {
		return streams[0], nil
	}

	ctx := context.TODO() // context for direct connection
	if limitedConnection {
		ctx = network.WithAllowLimitedConn(ctx, "n2x") // context for relayed connection
	}

	s, err := p2pHost.NewStream(ctx, peerInfo.ID, p2p.ProtocolID)
	if err != nil {
		return nil, errors.Wrapf(err, "[%v] function p2pHost.NewStream()", errors.Trace())
	}

	return s, nil
}

func connectPeer(p2pHost host.Host, hop *NetHop) (*peer.AddrInfo, error) {
	// try direct/relayed connection

	pm := newPeerAddrInfoMapFromNetHop(hop)

	// fmt.Println("----- pm - start -----")
	// pm.show()
	// fmt.Println("----- pm - end -----")

	peerInfo := connectPeerGroup(p2pHost, pm.peer)
	if peerInfo != nil {
		return peerInfo, nil
	}

	xlog.Warn("Unable to connect to peer internally, trying via default routers...")

	// try connection via default routers

	pm = newRtrPeerAddrInfoMapFromNetHop(hop)

	peerInfo = connectPeerGroup(p2pHost, pm.peer)
	if peerInfo != nil {
		return peerInfo, nil
	}

	return nil, fmt.Errorf("unable to connect to peer")
}

func connectPeerGroup(p2pHost host.Host, peers map[peer.ID]*peer.AddrInfo) *peer.AddrInfo {
	for _, peerInfo := range peers {
		if peerInfo.ID == p2pHost.ID() {
			continue
		}

		if err := connect(p2pHost, peerInfo); err != nil {
			continue
		}

		return peerInfo
	}

	return nil
}

func connect(p2pHost host.Host, peerInfo *peer.AddrInfo) error {
	p2pHost.Network().(*swarm.Swarm).Backoff().Clear(peerInfo.ID)
	if err := p2pHost.Connect(context.TODO(), *peerInfo); err != nil {
		xlog.Tracef("Unable to connect to peer %s: %v", peerInfo.ID.ShortString(), err)
		return errors.Wrapf(err, "[%v] function p2pHost.Connect()", errors.Trace())
	}
	p2pHost.Network().(*swarm.Swarm).Backoff().Clear(peerInfo.ID)

	return nil
}
