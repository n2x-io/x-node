package router

import (
	"net/netip"

	"n2x.dev/x-api-go/grpc/resources/topology"
	"n2x.dev/x-lib/pkg/ipnet"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/mnet/router/conntrack"
)

func (r *router) packetFilter(conn *conntrack.Connection, pktlen int) bool {
	// check ipDst rib entry
	if err := r.RIB().CheckIPDst(&conn.DstAddr); err != nil {
		xlog.Warnf("[pf] Unable to check IP dstAddr: %v", err)
		return true // drop the pkt
	}

	// check iap traffic
	if conn.AF == ipnet.AddressFamilyIPv6 && ipnet.IsIAPIPv6Addr(conn.SrcIP.String()) {
		return false // accept the pkt
	}

	// check configured network filters
	if r.policyFilter(conn) {
		// packet dropped by policy

		// check conntrack table
		if !conn.IsActive(pktlen) {
			// packet dropped by conntrack
			xlog.Warnf("[pf] Dropping %s packet from %s:%d to %s:%d",
				conn.Proto.String(),
				conn.SrcIP.String(),
				conn.SrcPort,
				conn.DstAddr.String(),
				conn.DstPort,
			)
			return true // drop the pkt
		}
	}

	return false // accept the pkt
}

func (r *router) policyFilter(conn *conntrack.Connection) bool {
	// filter specific proto packet types
	if conn.ProtoFilter() {
		return true // invalid proto packet type, drop the pkt
	}

	// get network policy
	p := r.RIB().GetPolicy(r.subnetID)
	if p == nil {
		return true // no policy, drop the pkt
	}

	// check configured network filters
	for _, f := range p.NetworkFilters {
		// proto
		if !(f.Proto == topology.Protocol_ANY || ipnet.IPProtocol(f.Proto.String()) == conn.Proto) {
			continue
		}
		// fmt.Printf("*** MATCHED Proto (%s): %s\n", f.Proto.String(), conn.proto.String())

		// dstPort
		if !(f.DstPort == 0 || f.DstPort == uint32(conn.DstPort)) {
			continue
		}
		// fmt.Printf("*** MATCHED DstPort (%d): %d\n", f.DstPort, conn.DstPort)

		// srcIP
		srcIPNet, err := netip.ParsePrefix(f.SrcIPNet)
		if err != nil {
			xlog.Warnf("[pf] Unable to parse filter srcIPNet prefix: %v", err)
			continue
		}

		// dstIP
		dstIPNet, err := netip.ParsePrefix(f.DstIPNet)
		if err != nil {
			xlog.Warnf("[pf] Unable to parse filter dstIPNet prefix: %v", err)
			continue
		}

		if srcIPNet.Contains(conn.SrcIP) && dstIPNet.Contains(conn.DstIP) {
			// fmt.Printf("*** MATCHED SrcIP (%s): %s\n", srcIPNet.String(), conn.SrcIP.String())
			// fmt.Printf("*** MATCHED DstIP (%s): %s\n", dstIPNet.String(), conn.DstIP.String())
			switch f.Policy {
			case topology.SecurityPolicy_ACCEPT:
				return false // accept the pkt
			case topology.SecurityPolicy_DROP:
				return true // drop the pkt
			}
		}
	}

	// fmt.Printf("+++ DEFAULT Policy: %s\n", p.DefaultPolicy.String())
	switch p.DefaultPolicy {
	case topology.SecurityPolicy_ACCEPT:
		return false // accept the pkt
	case topology.SecurityPolicy_DROP:
		return true // drop the pkt
	}

	return true // drop the pkt
}
