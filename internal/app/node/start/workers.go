package start

import (
	"n2x.dev/x-api-go/grpc/rpc"
	// "n2x.dev/x-lib/pkg/n2xp"
	"n2x.dev/x-lib/pkg/runtime"
	"n2x.dev/x-lib/pkg/update"
	"n2x.dev/x-node/internal/app/node/hsec"
	"n2x.dev/x-node/internal/app/node/ops"
	"n2x.dev/x-node/internal/app/node/svcs"
)

const (
	errorEventsHandler = iota
	// networkErrorEventsHandler
	// n2xDispatcher
	// n2xProcessor
	dnsAgent
	metricsAgent
	n2xpController
	routingAgent
	cronAgent
	atdAgent
	k8sConnector
	// proxy64gc
	federationMonitor
	securityScanner
	updateAgent
	// bgpAgent
)

func initWrkrs(nxnc rpc.NetworkAPIClient) {
	runtime.RegisterWrkr(
		errorEventsHandler,
		runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xErrorEventsHandler"),
		runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, runtime.ErrorEventsHandler),
	)
	// runtime.RegisterWrkr(
	// 	networkErrorEventsHandler,
	// 	runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xNetworkErrorEventsHandler"),
	// 	runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, svcs.NetworkErrorEventsHandler),
	// )
	// runtime.RegisterWrkr(
	// 	n2xDispatcher,
	// 	runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xDispatcher"),
	// 	runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, n2xp.Dispatcher),
	// )
	// runtime.RegisterWrkr(
	// 	n2xProcessor,
	// 	runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xProcessor"),
	// 	runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, svcs.MMPProcessor),
	// )
	runtime.RegisterWrkr(
		dnsAgent,
		runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xDNSAgent"),
		runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, svcs.DNSAgent),
		runtime.SetWrkrOpt(runtime.WrkrOptNxNetworkClient, nxnc),
	)
	runtime.RegisterWrkr(
		metricsAgent,
		runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xMetricsAgent"),
		runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, svcs.MetricsAgent),
		runtime.SetWrkrOpt(runtime.WrkrOptNxNetworkClient, nxnc),
	)
	runtime.RegisterWrkr(
		n2xpController,
		runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xpController"),
		runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, svcs.NetworkControl),
		runtime.SetWrkrOpt(runtime.WrkrOptNxNetworkClient, nxnc),
	)
	runtime.RegisterWrkr(
		routingAgent,
		runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xRoutingAgent"),
		runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, svcs.RoutingAgent),
		runtime.SetWrkrOpt(runtime.WrkrOptNxNetworkClient, nxnc),
	)
	runtime.RegisterWrkr(
		cronAgent,
		runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xCron"),
		runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, ops.Cron),
	)
	runtime.RegisterWrkr(
		atdAgent,
		runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xAtd"),
		runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, ops.Atd),
	)
	runtime.RegisterWrkr(
		k8sConnector,
		runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xKubernetesGateway"),
		runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, svcs.KubernetesConnector),
		runtime.SetWrkrOpt(runtime.WrkrOptNxNetworkClient, nxnc),
	)
	// runtime.RegisterWrkr(
	// 	proxy64gc,
	// 	runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xProxy64GC"),
	// 	runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, svcs.Proxy64GC),
	// )
	runtime.RegisterWrkr(
		federationMonitor,
		runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xFederationMonitor"),
		runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, svcs.FederationMonitor),
		runtime.SetWrkrOpt(runtime.WrkrOptNxNetworkClient, nxnc),
	)
	runtime.RegisterWrkr(
		securityScanner,
		runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xSecurityScanner"),
		runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, hsec.Scanner),
	)
	runtime.RegisterWrkr(
		updateAgent,
		runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xUpdateAgent"),
		runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, update.UpdateAgent),
	)
	// runtime.RegisterWrkr(
	// 	bgpAgent,
	// 	runtime.SetWrkrOpt(runtime.WrkrOptName, "n2xBGPAgent"),
	// 	runtime.SetWrkrOpt(runtime.WrkrOptStartFunc, svcs.BGPAgent),
	// )
}
