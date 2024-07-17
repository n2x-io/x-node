package workflow

import (
	"context"

	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/ops"
)

func Processor(ctx context.Context, pdu *n2xsp.WorkflowPDU) {
	if pdu == nil {
		return
	}

	var err error

	switch pdu.Type {
	case n2xsp.WorkflowMsgType_WORKFLOW_EXPEDITE:
		err = ops.WorkflowExpedite(ctx, pdu)
	case n2xsp.WorkflowMsgType_WORKFLOW_SCHEDULE:
		err = ops.WorkflowSchedule(ctx, pdu)
	}

	if err != nil {
		xlog.Errorf("[n2xp] Unable to process n2xp workflowPDU (%s): %v",
			pdu.Type.String(), err)
	}
}
