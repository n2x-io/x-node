package routing

import (
	"context"

	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/xlog"
)

func Processor(ctx context.Context, pdu *n2xsp.RoutingPDU) {
	if pdu == nil {
		return
	}

	var err error

	switch pdu.Type {
	case n2xsp.RoutingMsgType_ROUTING_STATUS:
		err = n2xpRoutingStatus(ctx, pdu)
	case n2xsp.RoutingMsgType_ROUTING_APPSVC:
		err = n2xpRoutingAppSvcConfig(ctx, pdu)
	}

	if err != nil {
		xlog.Errorf("[n2xp] Unable to process n2xp routingPDU (%s): %v", pdu.Type.String(), err)
	}
}
