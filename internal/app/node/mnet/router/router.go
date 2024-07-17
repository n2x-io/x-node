package router

import (
	"bufio"
	"net/netip"
	"os"
	"sync"

	libp2pHost "github.com/libp2p/go-libp2p/core/host"
	"n2x.dev/x-lib/pkg/errors"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/mnet/p2p"
	"n2x.dev/x-node/internal/app/node/mnet/p2p/host"
	"n2x.dev/x-node/internal/app/node/mnet/router/conntrack"
	"n2x.dev/x-node/internal/app/node/mnet/router/rib"
)

type Interface interface {
	Init() error
	P2PHost() libp2pHost.Host
	// NetworkInterfaceName() string
	IPv4() string
	SetIPv4(ipv4 string)
	IPv6() string
	GlobalIPv6() string
	SetIPv6(ipv6 string)
	SetGlobalIPv6(ipv6 string)
	IP4AddrAdd(ipv4 string) error
	IP4AddrDel(ipv4 string) error
	IP6AddrAdd(ipv6 string) error
	IP6AddrDel(ipv6 string) error
	GetConnections() int32
	RIB() rib.Interface
	Disconnect()
}

type router struct {
	p2pHost libp2pHost.Host

	subnetID string

	port         int
	externalIPv4 string

	ipv4       string
	ipv6       string
	globalIPv6 string

	rib     rib.Interface
	routes  *routeMap
	dialing *dialMap
	streams *streamMap
	proxy64 *proxy64Map

	networkInterface *networkInterface
	localForwarding  bool

	connections int32

	evtProcessorCloseCh chan struct{}
}

type cidrIPDst string

type routeMap struct {
	imported []string
	exported []string
	// sync.RWMutex
}

type dialMap struct {
	addr map[netip.Addr]struct{}
	sync.RWMutex
}

type streamMap struct {
	tunnel map[netip.Addr]*bufio.ReadWriter
	sync.RWMutex
}

func (d cidrIPDst) string() string {
	return string(d)
}

func New(externalIPv4, subnetID string, port int, localForwarding bool, rtImported, rtExported []string) Interface {
	return &router{
		subnetID:     subnetID,
		port:         port,
		externalIPv4: externalIPv4,
		rib:          rib.New(),
		routes: &routeMap{
			imported: rtImported,
			exported: rtExported,
		},
		dialing: &dialMap{
			addr: make(map[netip.Addr]struct{}),
		},
		streams: &streamMap{
			tunnel: make(map[netip.Addr]*bufio.ReadWriter),
		},
		proxy64: &proxy64Map{
			vs:      make(map[proxy64VSID]*proxy64VS),
			closeCh: make(chan struct{}, 1),
		},
		localForwarding:     localForwarding,
		connections:         0,
		evtProcessorCloseCh: make(chan struct{}, 1),
	}
}

func (r *router) Init() error {
	var p2pHost libp2pHost.Host
	var err error

	// if r.canRelay && len(r.externalIPv4) > 0 {
	// 	xlog.Info("Initializing tier-1 node...")
	// 	p2pHost, err = host.New(host.P2PHostTypeRelayHost, r.port)
	// } else if len(r.externalIPv4) > 0 {
	// 	xlog.Info("Initializing basic node...")
	// 	p2pHost, err = host.New(host.P2PHostTypeBasicHost, r.port)
	// } else {
	// 	xlog.Info("Initializing hidden node...")
	// 	p2pHost, err = host.New(host.P2PHostTypeHiddenHost, r.port)
	// }

	if len(r.externalIPv4) > 0 {
		xlog.Info("Initializing tier-1 node...")
		p2pHost, err = host.New(host.P2PHostTypeRelayHost, r.port)
	} else {
		xlog.Info("Initializing basic node...")
		p2pHost, err = host.New(host.P2PHostTypeBasicHost, r.port)
	}
	if err != nil {
		return errors.Wrapf(err, "[%v] function host.New()", errors.Trace())
	}

	r.p2pHost = p2pHost

	// Set a function as stream handler.
	// This function is called when a peer connects, and starts a stream with this protocol.
	// Only applies on the receiving side.
	r.p2pHost.SetStreamHandler(p2p.ProtocolID, r.handleStream)

	// set network interface
	if r.localForwarding {
		if err := r.ifDown(); err != nil {
			xlog.Alertf("Unable to reset interface: %v", err)
			os.Exit(1)
		}

		if err := r.ifUp(); err != nil {
			return errors.Wrapf(err, "[%v] function r.ifUp()", errors.Trace())
		}

		go r.proxy64GC()
	}

	go r.eventProcessor(r.evtProcessorCloseCh)

	return nil
}

func (r *router) P2PHost() libp2pHost.Host {
	return r.p2pHost
}

/*
func (r *router) NetworkInterfaceName() string {
	return r.networkInterface.devName()
}
*/

func (r *router) IPv4() string {
	return r.ipv4
}

func (r *router) SetIPv4(ipv4 string) {
	r.ipv4 = ipv4
}

func (r *router) IPv6() string {
	return r.ipv6
}

func (r *router) GlobalIPv6() string {
	return r.globalIPv6
}

func (r *router) SetIPv6(ipv6 string) {
	r.ipv6 = ipv6
}

func (r *router) SetGlobalIPv6(ipv6 string) {
	r.globalIPv6 = ipv6
}

func (r *router) IP4AddrAdd(ipv4 string) error {
	if r.networkInterface == nil || !r.localForwarding {
		return nil
	}

	return r.networkInterface.ip4AddrAdd(ipv4)
}

func (r *router) IP4AddrDel(ipv4 string) error {
	if r.networkInterface == nil || !r.localForwarding {
		return nil
	}

	return r.networkInterface.ip4AddrDel(ipv4)
}

func (r *router) IP6AddrAdd(ipv6 string) error {
	if r.networkInterface == nil || !r.localForwarding {
		return nil
	}

	return r.networkInterface.ip6AddrAdd(ipv6)
}

func (r *router) IP6AddrDel(ipv6 string) error {
	if r.networkInterface == nil || !r.localForwarding {
		return nil
	}

	return r.networkInterface.ip6AddrDel(ipv6)
}

func (r *router) RIB() rib.Interface {
	return r.rib
}

func (r *router) GetConnections() int32 {
	return r.connections
}

func (r *router) Disconnect() {
	xlog.Info("Disconnecting...")

	if r.localForwarding && r.networkInterface != nil {
		r.networkInterface.closeCh <- struct{}{}

		if err := r.networkInterface.close(); err != nil {
			xlog.Warnf("Unable to close interface %s: %v", r.networkInterface.devName(), err)
		}

		r.proxy64.closeCh <- struct{}{}

		conntrack.Ctrl().Close()
	}

	if err := r.p2pHost.Close(); err != nil {
		xlog.Warnf("Unable to close p2pHost handler: %v", err)
		os.Exit(1)
	}

	r.rib.Close()

	r.evtProcessorCloseCh <- struct{}{}
}
