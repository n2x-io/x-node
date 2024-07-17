package ops

import (
	"github.com/spf13/viper"
	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-api-go/grpc/resources/ops"
)

func newWorkflowResponse(pdu *n2xsp.WorkflowPDU) *n2xsp.Payload {
	n2xID := viper.GetString("n2x.id")

	return &n2xsp.Payload{
		SrcID: n2xID,
		Type:  n2xsp.PDUType_WORKFLOW,
		WorkflowPDU: &n2xsp.WorkflowPDU{
			Type: n2xsp.WorkflowMsgType_WORKFLOW_RESPONSE,
			Workflow: &ops.Workflow{
				AccountID:   pdu.Workflow.AccountID,
				TenantID:    pdu.Workflow.TenantID,
				ProjectID:   pdu.Workflow.ProjectID,
				WorkflowID:  pdu.Workflow.WorkflowID,
				Name:        pdu.Workflow.Name,
				Description: pdu.Workflow.Description,
				Notify:      pdu.Workflow.Notify,
				TaskLogs:    pdu.Workflow.TaskLogs,
			},
		},
	}

	// return &n2xsp.Payload{
	// 	SrcID:       p.DstID,
	// 	DstID:       p.SrcID,
	// 	RequesterID: p.RequesterID,
	// 	Interactive: p.Interactive,
	// 	PayloadType: n2xsp.PayloadType_WORKFLOW_RESPONSE,
	// 	Workflow: &ops.Workflow{
	// 		AccountID: p.Workflow.AccountID,
	// 		TenantID:  p.Workflow.TenantID,
	// 		ProjectID:  p.Workflow.ProjectID,
	// 		WorkflowID: p.Workflow.WorkflowID,
	//      Name:       p.Workflow.Name,
	//      Description: p.Workflow.Description,
	// 		OwnerUserID:      p.Workflow.OwnerUserID,
	// 		Notify:     p.Workflow.Notify,
	// 		TaskLogs: p.Workflow.TaskLogs,
	// 	},
	// }
}
