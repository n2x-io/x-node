package nodemgmt

import (
	"context"
	"fmt"

	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/errors"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/mnet"
)

func n2xpNodeConfig(ctx context.Context, pdu *n2xsp.NodeMgmtPDU) error {
	if pdu.NodeConfig == nil {
		return fmt.Errorf("null nodeConfig")
	}

	xlog.Infof("[n2xp] Received new configuration..")

	if err := mnet.NewCfg(pdu.NodeConfig); err != nil {
		return errors.Wrapf(err, "[%v] function mnet.NewCfg()", errors.Trace())
	}

	return nil
}
