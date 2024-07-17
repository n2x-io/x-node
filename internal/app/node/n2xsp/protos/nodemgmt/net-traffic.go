package nodemgmt

import (
	"fmt"

	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/kvstore/db/netflowdb"
)

func n2xpNetTrafficMetricsRequest(pdu *n2xsp.NodeMgmtPDU) error {
	if pdu.NetTrafficMetricsRequest == nil {
		return fmt.Errorf("null netTrafficMetricsRequest")
	}
	req := pdu.NetTrafficMetricsRequest

	xlog.Debugf("[n2xp] Received new traffic metrics request..")

	netflowdb.RequestQueue <- req

	return nil
}
