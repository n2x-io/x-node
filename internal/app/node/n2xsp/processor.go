package n2xsp

import (
	"context"

	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-node/internal/app/node/n2xsp/protos/nodemgmt"
	"n2x.dev/x-node/internal/app/node/n2xsp/protos/routing"
	"n2x.dev/x-node/internal/app/node/n2xsp/protos/workflow"
)

var RxQueue = make(chan *n2xsp.Payload, 128)

func Processor(ctx context.Context, p *n2xsp.Payload) {
	switch p.Type {
	case n2xsp.PDUType_ROUTING:
		routing.Processor(ctx, p.RoutingPDU)
	case n2xsp.PDUType_NODEMGMT:
		nodemgmt.Processor(ctx, p.NodeMgmtPDU)
	case n2xsp.PDUType_WORKFLOW:
		workflow.Processor(ctx, p.WorkflowPDU)
	}
}
