package nodemgmt

import (
	"fmt"

	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/kvstore/db/ctlogdb"
	"n2x.dev/x-node/internal/app/node/mnet/router/conntrack"
)

func n2xpNetConntrackStateRequest(pdu *n2xsp.NodeMgmtPDU) error {
	if pdu.NetCtStateRequest == nil {
		return fmt.Errorf("null netCtStateRequest")
	}
	req := pdu.NetCtStateRequest

	xlog.Debugf("[n2xp] Received new conntrack state request..")

	conntrack.RequestQueue <- req

	return nil
}

func n2xpNetConntrackLogRequest(pdu *n2xsp.NodeMgmtPDU) error {
	if pdu.NetCtLogRequest == nil {
		return fmt.Errorf("null netCtLogRequest")
	}
	req := pdu.NetCtLogRequest

	xlog.Debugf("[n2xp] Received new conntrack log request..")

	ctlogdb.RequestQueue <- req

	return nil
}
