package mnet

import (
	"context"
	"time"

	"n2x.dev/x-api-go/grpc/network/nac"
	"n2x.dev/x-api-go/grpc/resources/topology"
	"n2x.dev/x-lib/pkg/errors"
)

func (ln *localNode) NodeReq() *topology.NodeReq {
	return &topology.NodeReq{
		AccountID: ln.node.AccountID,
		TenantID:  ln.node.TenantID,
		NodeID:    ln.node.NodeID,
	}
}

func (ln *localNode) Node() *topology.Node {
	ln.node.Endpoints = ln.endpoints.endpt

	return ln.node
}

func (ln *localNode) registerNode() error {
	n := ln.Node()

	nrReq := &nac.NodeRegRequest{
		Node: n,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := ln.Connection().NetworkClient().RegisterNode(ctx, nrReq)
	if err != nil {
		return errors.Wrapf(err, "[%v] function ln.Connection().NetworkClient().RegisterNode()", errors.Trace())
	}

	if ln.initialized {
		return nil
	}

	if n.Cfg.DisableNetworking || ln.Router() == nil {
		ln.initialized = true

		return nil
	}

	// initialize rib data structure
	mq := ln.Connection().RIBDataMsgRxQueue()
	ln.Router().RIB().Initialize(mq, n, r.RoutingDomain, r.NetworkPolicy)

	// subscribe to mqtt routing topics
	if err := ln.Connection().NewRoutingSession(r.RoutingDomain.LocationID); err != nil {
		return errors.Wrapf(err, "[%v] function ln.Connection().NewRoutingSession()", errors.Trace())
	}

	for _, as := range r.AppSvcs {
		ln.Router().RIB().AddNodeAppSvc(as)
	}

	ln.initialized = true

	return nil
}
