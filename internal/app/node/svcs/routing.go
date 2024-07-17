package svcs

import (
	"time"

	"github.com/spf13/viper"
	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/n2xp/queuing"
	"n2x.dev/x-lib/pkg/runtime"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/mnet"
	"n2x.dev/x-node/internal/app/node/n2xsp/protos/routing"
)

// RoutingAgent runs routing engine
func RoutingAgent(w *runtime.Wrkr) {
	xlog.Infof("Started worker %s", w.Name)
	w.Running = true

	endCh := make(chan struct{}, 2)

	go func() {
		n2xID := viper.GetString("n2x.id")

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if mnet.LocalNode().Node().Cfg.DisableNetworking ||
					mnet.LocalNode().Router() == nil {
					continue
				}

				if !routing.ServiceEnabled {
					continue
				}

				xlog.Debug("Sending routing LSAs")

				lsa := mnet.LocalNode().GetNodeLSA()
				if lsa == nil {
					continue
				}

				queuing.TxControlQueue <- &n2xsp.Payload{
					SrcID: n2xID,
					Type:  n2xsp.PDUType_ROUTING,
					RoutingPDU: &n2xsp.RoutingPDU{
						Type: n2xsp.RoutingMsgType_ROUTING_LSA,
						LSA:  lsa,
					},
				}

				mnet.LocalNode().SendAppSvcLSAs(n2xID)
			case <-endCh:
				// xlog.Warn("Closing rtRequest send stream")
				return
			}
		}
	}()

	<-w.QuitChan

	endCh <- struct{}{}

	w.WG.Done()
	w.Running = false
	xlog.Infof("Stopped worker %s", w.Name)
}
