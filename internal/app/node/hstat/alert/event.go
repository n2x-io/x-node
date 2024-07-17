package alert

import (
	"github.com/spf13/viper"
	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-api-go/grpc/resources/events"
	"n2x.dev/x-lib/pkg/n2xp/queuing"
	"n2x.dev/x-lib/pkg/xlog"
)

func newAlertEvent(evt *events.Event) {
	n2xID := viper.GetString("n2x.id")

	xlog.Debugf("[event] New event from srcID %s", n2xID)

	queuing.TxControlQueue <- &n2xsp.Payload{
		SrcID: n2xID,
		Type:  n2xsp.PDUType_EVENT,
		EventPDU: &n2xsp.EventPDU{
			Event: evt,
		},
	}
}
