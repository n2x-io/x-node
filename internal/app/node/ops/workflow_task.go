package ops

import (
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"n2x.dev/x-api-go/grpc/common/status"
	"n2x.dev/x-api-go/grpc/resources/ops"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/mnet"
)

func runWorkflowTask(wf *ops.Workflow, t *ops.Task) *ops.TaskLog {
	var statusMsg string
	var statusCode status.StatusCode
	var resultStatus ops.CommandResultStatus

	if t.Command == nil {
		return nil
	}

	xlog.Infof("[ops] Executing workflow %s, task %s", wf.WorkflowID, t.Name)

	c := t.Command

	// execute the command
	cmd := exec.Command(c.Cmd, c.Args...)
	cmd.Stdin = nil

	t1 := time.Now()
	out, err := cmd.CombinedOutput()
	if err != nil {
		statusCode = status.StatusCode_FAILED
		statusMsg = err.Error()
		resultStatus = ops.CommandResultStatus_FAILED

		xlog.Errorf("Unable to run command %s: %v", c.Cmd, err)
	} else {
		statusCode = status.StatusCode_OK
		statusMsg = "OK"
		resultStatus = ops.CommandResultStatus_EXECUTED
	}

	n2xID := viper.GetString("n2x.id")
	n := mnet.LocalNode().Node()

	return &ops.TaskLog{
		AccountID:       wf.AccountID,
		TenantID:        wf.TenantID,
		ProjectID:       wf.ProjectID,
		WorkflowID:      wf.WorkflowID,
		TaskLogID:       uuid.New().String(),
		TaskName:        t.Name,
		TaskDescription: t.Description,
		// NetID:           n.NetID,
		// SubnetID:        n.SubnetID,
		NodeID:   n.NodeID,
		NodeName: n.Cfg.NodeName,
		Status: &status.StatusResponse{
			SourceID:  n2xID,
			Code:      statusCode,
			Message:   statusMsg,
			Timestamp: time.Now().UnixMilli(),
		},
		Result: &ops.CommandResult{
			Status:   resultStatus,
			Duration: int64(time.Since(t1).Seconds()),
		},
		StdoutStderr: out,
	}
}
