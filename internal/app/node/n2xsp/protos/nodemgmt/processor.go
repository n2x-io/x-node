package nodemgmt

import (
	"context"

	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/xlog"
)

func Processor(ctx context.Context, pdu *n2xsp.NodeMgmtPDU) {
	if pdu == nil {
		return
	}

	var err error

	switch pdu.Type {
	case n2xsp.NodeMgmtMsgType_NODE_CONFIG:
		err = n2xpNodeConfig(ctx, pdu)
	case n2xsp.NodeMgmtMsgType_NODE_HOST_METRICS_REQUEST:
		err = n2xpHostMetricsRequest(pdu)
	case n2xsp.NodeMgmtMsgType_NODE_NET_CONNTRACK_STATE_REQUEST:
		err = n2xpNetConntrackStateRequest(pdu)
	case n2xsp.NodeMgmtMsgType_NODE_NET_CONNTRACK_LOG_REQUEST:
		err = n2xpNetConntrackLogRequest(pdu)
	case n2xsp.NodeMgmtMsgType_NODE_NET_TRAFFIC_METRICS_REQUEST:
		err = n2xpNetTrafficMetricsRequest(pdu)
	case n2xsp.NodeMgmtMsgType_NODE_HOST_SECURITY_REQUEST:
		err = n2xpHostSecurityReportRequest(pdu)
	}

	if err != nil {
		xlog.Errorf("[n2xp] Unable to process n2xp nodeMgmtPDU (%s): %v",
			pdu.Type.String(), err)
	}
}
