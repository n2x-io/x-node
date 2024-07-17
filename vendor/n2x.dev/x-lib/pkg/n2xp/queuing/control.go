package queuing

import "n2x.dev/x-api-go/grpc/network/n2xsp"

var RxControlQueue = make(chan *n2xsp.Payload, 128)
var TxControlQueue = make(chan *n2xsp.Payload, 128)
