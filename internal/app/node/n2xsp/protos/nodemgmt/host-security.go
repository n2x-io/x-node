package nodemgmt

import (
	"fmt"

	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/hsec"
)

func n2xpHostSecurityReportRequest(pdu *n2xsp.NodeMgmtPDU) error {
	if pdu.HsecReportRequest == nil {
		return fmt.Errorf("null hsecReportRequest")
	}
	req := pdu.HsecReportRequest

	xlog.Debugf("[n2xp] Received new host security report request..")

	hsec.RequestQueue <- req

	return nil
}
