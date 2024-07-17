package mnet

import (
	"fmt"

	"github.com/spf13/viper"
	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-api-go/grpc/resources/topology"
	"n2x.dev/x-lib/pkg/errors"
	"n2x.dev/x-lib/pkg/n2x"
	"n2x.dev/x-lib/pkg/runtime"
	"n2x.dev/x-lib/pkg/version"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/hstat"
	"n2x.dev/x-node/internal/app/node/mnet/connection"
	"n2x.dev/x-node/internal/app/node/mnet/maddr"
	"n2x.dev/x-node/internal/app/node/mnet/router"
)

func Init() error {
	conn := connection.New()

	if err := cfgInit(conn); err != nil {
		return errors.Wrapf(err, "[%v] function cfgInit()", errors.Trace())
	}

	return nil
}

func NewCfg(nodeConfig *n2xsp.NodeMgmtConfig) error {
	switch nodeConfig.Type {
	case n2xsp.NodeMgmtConfigActionType_CFG_METADATA:
		LocalNode().Node().Cfg.NodeName = nodeConfig.Cfg.NodeName
		LocalNode().Node().Cfg.Description = nodeConfig.Cfg.Description

		viper.Set("nodeName", nodeConfig.Cfg.NodeName)
	case n2xsp.NodeMgmtConfigActionType_CFG_NETWORKING:
		xlog.Info("Reconfiguring networking...")

		LocalNode().Close()

		conn := connection.New()

		if err := cfgInit(conn); err != nil {
			return errors.Wrapf(err, "[%v] function cfgInit()", errors.Trace())
		}

		runtime.NetworkWrkrReconnect(conn.NetworkClient())
	case n2xsp.NodeMgmtConfigActionType_CFG_MANAGEMENT:
		if nodeConfig.Cfg.Management == nil {
			return fmt.Errorf("invalid management config")
		}

		xlog.Info("Applying new configuration...")

		viper.Set("management.disableOps", nodeConfig.Cfg.Management.DisableOps)
		viper.Set("management.disableExec", nodeConfig.Cfg.Management.DisableExec)
		viper.Set("management.disableTransfer", nodeConfig.Cfg.Management.DisableTransfer)
		viper.Set("management.disablePortForwarding", nodeConfig.Cfg.Management.DisablePortForwarding)
	}

	return nil
}

func cfgInit(conn connection.Interface) error {
	hostID := viper.GetString("host.id")
	port := viper.GetInt("port")
	dnsPort := viper.GetInt("dnsPort")
	rtExported := viper.GetStringSlice("routes.export")
	rtImported := viper.GetStringSlice("routes.import")

	n := conn.Node()

	n2xID, err := n2x.GetID(&topology.NodeReq{
		AccountID: n.AccountID,
		TenantID:  n.TenantID,
		NodeID:    n.NodeID,
	})
	if err != nil {
		return errors.Wrapf(err, "[%v] function n2x.GetID()", errors.Trace())
	}

	viper.Set("n2x.id", n2xID.String())

	if n.Agent.DevMode {
		viper.Set("version.branch", "dev")
	}

	setCfgVars(n.Cfg)

	nodeName := n.Cfg.NodeName

	n.Agent.ExternalIPv4 = conn.GetExternalIPv4()
	n.Agent.CanRelay = len(n.Agent.ExternalIPv4) > 0

	if n.Type == topology.NodeType_K8S_GATEWAY {
		n.Cfg.DisableNetworking = false
	}

	var rtr router.Interface

	if n.Cfg.DisableNetworking {
		xlog.Info("Networking disabled")
		rtr = nil
		n.State = topology.NodeDeploymentState_STUBBY_MODE
		n.NodeDeploymentState = topology.NodeDeploymentState_STUBBY_MODE.String()
	} else {
		localForwarding := true

		rtr = router.New(n.Agent.ExternalIPv4, n.Cfg.SubnetID, port, localForwarding, rtImported, rtExported)
		if err := rtr.Init(); err != nil {
			return errors.Wrapf(err, "[%v] function r.Init()", errors.Trace())
		}

		maddrs := maddr.GetGlobalUnicastAddrStrings(rtr.P2PHost().Addrs()...)

		if len(maddrs) > 0 {
			xlog.Info("Node multi-addresses:")
			for _, ma := range maddrs {
				xlog.Infof(" => %s", ma)
			}
		}
		xlog.Debugf("p2pHostID: %s", rtr.P2PHost().ID().String())

		n.Agent.P2PHostID = rtr.P2PHost().ID().String()
		n.Agent.MAddrs = maddrs
		n.Agent.Port = int32(port)
		n.Agent.Routes = &topology.Routes{
			Export: rtExported,
			Import: rtImported,
		}
		n.State = topology.NodeDeploymentState_CONNECTED
		n.NodeDeploymentState = topology.NodeDeploymentState_CONNECTED.String()
	}

	n.Agent.Hostname = hostID
	n.Agent.DNSPort = int32(dnsPort)
	n.Agent.Metrics = &topology.AgentMetrics{}
	n.Agent.Version = version.GetVersion()
	// n.Agent.DevMode =
	n.Endpoints = make(map[string]*topology.Endpoint)

	n.Class = topology.NodeClass_COMPUTE_NODE
	n.NodeClass = topology.NodeClass_COMPUTE_NODE.String()

	localnode = &localNode{
		node: n,
		endpoints: &endpointsMap{
			endpt: make(map[string]*topology.Endpoint),
		},
		connection: conn,
		router:     rtr,
		stats: hstat.Init(&topology.NodeReq{
			AccountID: n.AccountID,
			TenantID:  n.TenantID,
			NodeID:    n.NodeID,
		}),
	}

	if n.Cfg.DisableNetworking {
		if err := localnode.registerNode(); err != nil {
			return errors.Wrapf(err, "[%v] function localnode.registerNode()", errors.Trace())
		}
	} else {
		dnsName := nodeName
		if n.KubernetesAttrs != nil {
			if len(n.KubernetesAttrs.Namespace) > 0 {
				dnsName = fmt.Sprintf("%s.%s", nodeName, n.KubernetesAttrs.Namespace)
			}
		}
		endpointID := dnsName
		if _, err := localnode.AddNetworkEndpoint(endpointID, dnsName); err != nil {
			return errors.Wrapf(err, "[%v] function localnode.AddNetworkEndpoint()", errors.Trace())
		}
	}

	xlog.Infof("Node %s initialized", nodeName)

	return nil
}

func setCfgVars(cfg *topology.NodeCfg) {
	viper.Set("nodeName", cfg.NodeName)

	viper.Set("management.disableOps", cfg.Management.DisableOps)
	viper.Set("management.disableExec", cfg.Management.DisableExec)
	viper.Set("management.disableTransfer", cfg.Management.DisableTransfer)
	viper.Set("management.disablePortForwarding", cfg.Management.DisablePortForwarding)

	if viper.GetBool("maintenance.disableAutoUpdate") {
		viper.Set("maintenance.autoUpdate", false)
	} else {
		viper.Set("maintenance.autoUpdate", cfg.Maintenance.AutoUpdate)
	}
	viper.Set("maintenance.schedule.hour", int(cfg.Maintenance.Schedule.Hour))
	viper.Set("maintenance.schedule.minute", int(cfg.Maintenance.Schedule.Minute))
}
