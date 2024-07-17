package connection

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/grpc"
	"n2x.dev/x-api-go/grpc/resources/iam/auth"
	"n2x.dev/x-api-go/grpc/resources/topology"
	"n2x.dev/x-api-go/grpc/rpc"
	"n2x.dev/x-lib/pkg/errors"
	"n2x.dev/x-lib/pkg/xlog"
)

type Interface interface {
	Node() *topology.Node
	NetworkClient() rpc.NetworkAPIClient
	NewRoutingSession(locationID string) error
	RIBDataMsgRxQueue() <-chan []byte
	Watcher() chan struct{}
	GetExternalIPv4() string
	Close()
}

type session struct {
	connection *connection
	locationID string
	watcherCh  chan struct{}
	endCh      chan struct{}
}

type connection struct {
	defaultControllerEndpoint string

	authKey    *auth.AuthKey
	authSecret string
	node       *topology.Node

	grpcClientConn *grpc.ClientConn
	nxnc           rpc.NetworkAPIClient
	externalIPv4   string

	mqttClient mqtt.Client

	initialized bool
}

func New() Interface {
	s := &session{
		connection: &connection{},
		watcherCh:  make(chan struct{}, 64),
		endCh:      make(chan struct{}, 1),
	}

	s.connection.new()

	go s.connWatcher()

	return s
}

func (s *session) Node() *topology.Node {
	if s.connection == nil {
		return nil
	}

	return s.connection.node
}

func (s *session) NetworkClient() rpc.NetworkAPIClient {
	if s.connection == nil {
		return nil
	}

	return s.connection.nxnc
}

func (s *session) NewRoutingSession(locationID string) error {
	s.locationID = locationID

	accountID := s.connection.node.AccountID
	tenantID := s.connection.node.TenantID
	netID := s.connection.node.Cfg.NetID

	netRoutingTopic := fmt.Sprintf("%s/%s/%s", accountID, tenantID, netID)

	if err := s.connection.newRoutingSubscription(netRoutingTopic); err != nil {
		return errors.Wrapf(err, "[%v] function s.connection.newRoutingSubscription()", errors.Trace())
	}

	locationsTopic := fmt.Sprintf("locations/%s", locationID)

	if err := s.connection.newRoutingSubscription(locationsTopic); err != nil {
		return errors.Wrapf(err, "[%v] function s.connection.newRoutingSubscription()", errors.Trace())
	}

	return nil
}

func (s *session) RIBDataMsgRxQueue() <-chan []byte {
	return ribDataMsgRxQueue
}

func (s *session) Watcher() chan struct{} {
	return s.watcherCh
}

func (s *session) GetExternalIPv4() string {
	return s.connection.externalIPv4
}

func (s *session) Close() {
	// ends connection watcher
	s.endCh <- struct{}{}

	// close connection
	if s.connection == nil {
		return
	}

	if err := s.connection.grpcClientConn.Close(); err != nil {
		xlog.Errorf("Unable to close gRPC network connection: %v", err)
	}
	s.connection.nxnc = nil

	if s.connection.mqttClient != nil {
		s.connection.mqttClient.Disconnect(250)
		s.connection.mqttClient = nil
	}
}
