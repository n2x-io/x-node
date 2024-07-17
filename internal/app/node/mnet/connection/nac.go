package connection

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"n2x.dev/x-api-go/grpc/network/nac"
	"n2x.dev/x-api-go/grpc/resources/iam/auth"
	"n2x.dev/x-api-go/grpc/resources/topology"
	"n2x.dev/x-lib/pkg/errors"
	"n2x.dev/x-lib/pkg/grpc/client"
	"n2x.dev/x-lib/pkg/xlog"
)

func (c *connection) networkAdmissionRequest() error {
	token := viper.GetString("token")
	if len(token) == 0 {
		return fmt.Errorf("authorization token not found")
	}

	nodeToken, err := decodeNodeToken(token)
	if err != nil {
		return errors.Wrapf(err, "[%v] function decodeNodeToken()", errors.Trace())
	}

	controllerEndpoint := nodeToken.ControllerEndpoint
	authKey := &auth.AuthKey{
		Key: "no-auth",
	}
	authSecret := ""

	nc, grpcConn, err := client.NewNetworkAPIClient(controllerEndpoint, authKey, authSecret)
	if err != nil {
		return errors.Wrapf(err, "[%v] function client.NewNetworkAPIClient()", errors.Trace())
	}
	defer grpcConn.Close()

	naReq := &nac.NetworkAdmissionRequest{
		NodeToken:  token,
		NodeHostID: viper.GetString("host.id"),
	}

	r, err := nc.NetworkAdmissionControl(context.TODO(), naReq)
	if err != nil {
		return errors.Wrapf(err, "[%v] function nc.NetworkAdmissionControl()", errors.Trace())
	}

	switch r.Result {
	case nac.NetworkAdmissionResult_AUTHORIZED:
		c.defaultControllerEndpoint = controllerEndpoint
		c.authKey = r.AuthKey
		c.authSecret = ""
		c.node = r.Node
	case nac.NetworkAdmissionResult_UNAUTHORIZED:
		xlog.Alert("Network access NOT AUTHORIZED")
		grpcConn.Close()
		os.Exit(1)
	case nac.NetworkAdmissionResult_SERVICE_DISABLED:
		xlog.Alert("Service is DISABLED.")
		xlog.Alert("Please contact n2x customer service urgently.")
		time.Sleep(15 * time.Minute)
		return fmt.Errorf("SERVICE_DISABLED")
	}

	return nil
}

func decodeNodeToken(token string) (*topology.NodeTokenPayload, error) {
	if len(token) == 0 {
		return nil, fmt.Errorf("unable to decode nodeToken: invalid data")
	}

	// nodeToken

	jsonDataNodeToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, errors.Wrapf(err, "[%v] function base64.StdEncoding.DecodeString()", errors.Trace())
	}

	var nodeToken topology.NodeTokenPayload

	if err := json.Unmarshal(jsonDataNodeToken, &nodeToken); err != nil {
		return nil, errors.Wrapf(err, "[%v] function json.Unmarshal()", errors.Trace())
	}

	return &nodeToken, nil
}
