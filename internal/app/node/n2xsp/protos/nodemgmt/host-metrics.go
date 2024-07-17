package nodemgmt

import (
	"fmt"

	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/kvstore/db/metricsdb"
)

func n2xpHostMetricsRequest(pdu *n2xsp.NodeMgmtPDU) error {
	if pdu.HostMetricsRequest == nil {
		return fmt.Errorf("null hostMetrisRequest")
	}
	req := pdu.HostMetricsRequest

	xlog.Debugf("[n2xp] Received new host metrics request..")

	metricsdb.RequestQueue <- req

	return nil
}
