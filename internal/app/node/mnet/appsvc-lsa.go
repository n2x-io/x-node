package mnet

import (
	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-api-go/grpc/network/routing"
	"n2x.dev/x-lib/pkg/n2xp/queuing"
)

func (ln *localNode) SendAppSvcLSAs(n2xID string) {
	if ln == nil {
		return
	}

	if ln.Node().Cfg.DisableNetworking || ln.Router() == nil {
		return
	}

	if !ln.initialized {
		return
	}

	n := ln.Node()
	if n == nil {
		return
	}

	for _, as := range ln.Router().RIB().GetNodeAppSvcs() {
		lsa := &routing.LSA{
			Type: routing.LSAType_APPSVC_LSA,
			AppSvcLSA: &routing.AppSvcLSA{
				AppSvc:      as,
				P2PHostID:   n.Agent.P2PHostID,
				Priority:    n.Cfg.Priority,
				IPv6:        ln.Router().GlobalIPv6(),
				Connections: ln.Router().GetConnections(),
			},
		}

		queuing.TxControlQueue <- &n2xsp.Payload{
			SrcID: n2xID,
			Type:  n2xsp.PDUType_ROUTING,
			RoutingPDU: &n2xsp.RoutingPDU{
				Type: n2xsp.RoutingMsgType_ROUTING_LSA,
				LSA:  lsa,
			},
		}
	}
}
