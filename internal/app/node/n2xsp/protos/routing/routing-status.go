package routing

import (
	"context"
	"fmt"
	"os"

	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/xlog"
)

var ServiceEnabled bool = true
var disabledRetries int

func n2xpRoutingStatus(ctx context.Context, pdu *n2xsp.RoutingPDU) error {
	if pdu.Status == nil {
		return fmt.Errorf("null status")
	}
	s := pdu.Status

	if s.Disabled {
		xlog.Alert("Service is DISABLED.")
		xlog.Alert("Please contact n2x customer service urgently.")
	}
	// } else if s.OverLimit {
	// 	xlog.Alert("Account over tier limits. Service is DISABLED.")
	// 	xlog.Alert("If you are on the Free Plan, make sure you")
	// 	xlog.Alert("are not exceeding its limits. If not, please")
	// 	xlog.Alert("contact n2x customer service urgently.")
	// }

	// if s.Disabled || s.OverLimit {
	if s.Disabled {
		ServiceEnabled = false

		disabledRetries++
		if disabledRetries > 10 {
			os.Exit(1)
		}
		return nil
	} else {
		ServiceEnabled = true
		disabledRetries = 0
	}

	return nil
}
