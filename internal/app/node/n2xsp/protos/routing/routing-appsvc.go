package routing

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-node/internal/app/node/mnet"
)

func n2xpRoutingAppSvcConfig(ctx context.Context, pdu *n2xsp.RoutingPDU) error {
	if pdu.AppSvcConfig == nil {
		return fmt.Errorf("null appSvcConfig")
	}

	ascfg := pdu.AppSvcConfig

	switch ascfg.Operation {
	case n2xsp.AppSvcConfigOperation_APPSVC_SET:
		mnet.LocalNode().Router().RIB().AddNodeAppSvc(ascfg.AppSvc)
	case n2xsp.AppSvcConfigOperation_APPSVC_UNSET:
		mnet.LocalNode().Router().RIB().RemoveNodeAppSvc(ascfg.AppSvc.AppSvcID)
	}

	if mnet.LocalNode().Node().Cfg.DisableNetworking || mnet.LocalNode().Router() == nil {
		return nil
	}

	if !ServiceEnabled {
		return nil
	}

	n2xID := viper.GetString("n2x.id")

	mnet.LocalNode().SendAppSvcLSAs(n2xID)

	return nil
}
