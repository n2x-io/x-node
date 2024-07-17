package rib

import "n2x.dev/x-api-go/grpc/network/routing"

type RouteEventType int

const (
	RouteEventTypeADD RouteEventType = iota
	RouteEventTypeDELETE
)

type RouteEvent struct {
	Type      RouteEventType
	Addr      string
	RouteType routing.RouteType
}

var routeEventQueue = make(chan *RouteEvent, 32)
var relayConnQueue = make(chan *routing.NetHop, 16)
var routerConnQueue = make(chan *routing.NetHop, 16)
var proxyConnQueue = make(chan *routing.NetHop, 16)

func evtAddRoute(addr string, routeType routing.RouteType) {
	if routeType == routing.RouteType_DYNAMIC ||
		routeType == routing.RouteType_PROXY {
		// only STATIC or CONNECTED routes are processed
		return
	}

	routeEventQueue <- &RouteEvent{
		Type:      RouteEventTypeADD,
		Addr:      addr,
		RouteType: routeType,
	}
}

func evtDeleteRoute(addr string, routeType routing.RouteType) {
	if routeType == routing.RouteType_DYNAMIC ||
		routeType == routing.RouteType_PROXY {
		// only STATIC or CONNECTED routes are processed
		return
	}

	routeEventQueue <- &RouteEvent{
		Type:      RouteEventTypeDELETE,
		Addr:      addr,
		RouteType: routeType,
	}
}

func evtRelay(nh *routing.NetHop) {
	relayConnQueue <- nh
}

func evtRouter(nh *routing.NetHop) {
	routerConnQueue <- nh
}

func evtProxy(nh *routing.NetHop) {
	proxyConnQueue <- nh
}
