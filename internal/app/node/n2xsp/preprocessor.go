package n2xsp

import (
	"context"

	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/n2xp/queuing"
)

func Preprocessor(ctx context.Context, p *n2xsp.Payload) {
	switch p.Type {
	case n2xsp.PDUType_ROUTING:
		if p.RoutingPDU == nil {
			return
		}

		switch p.RoutingPDU.Type {
		case n2xsp.RoutingMsgType_ROUTING_STATUS:
			RxQueue <- p
			return
		case n2xsp.RoutingMsgType_ROUTING_APPSVC:
			RxQueue <- p
			return
		}
	case n2xsp.PDUType_NODEMGMT:
		if p.NodeMgmtPDU == nil {
			return
		}

		switch p.NodeMgmtPDU.Type {
		case n2xsp.NodeMgmtMsgType_NODE_CONFIG:
			RxQueue <- p
			return
		case n2xsp.NodeMgmtMsgType_NODE_HOST_METRICS_REQUEST:
			RxQueue <- p
			return
		case n2xsp.NodeMgmtMsgType_NODE_NET_CONNTRACK_STATE_REQUEST:
			RxQueue <- p
			return
		case n2xsp.NodeMgmtMsgType_NODE_NET_CONNTRACK_LOG_REQUEST:
			RxQueue <- p
			return
		case n2xsp.NodeMgmtMsgType_NODE_NET_TRAFFIC_METRICS_REQUEST:
			RxQueue <- p
			return
		case n2xsp.NodeMgmtMsgType_NODE_HOST_SECURITY_REQUEST:
			RxQueue <- p
			return
		}
	case n2xsp.PDUType_WORKFLOW:
		if p.WorkflowPDU == nil {
			return
		}

		switch p.WorkflowPDU.Type {
		case n2xsp.WorkflowMsgType_WORKFLOW_EXPEDITE:
			RxQueue <- p
			return
		case n2xsp.WorkflowMsgType_WORKFLOW_SCHEDULE:
			RxQueue <- p
			return
		}
	case n2xsp.PDUType_EVENT:
		if p.EventPDU == nil {
			return
		}
	}

	queuing.RxControlQueue <- p
}
