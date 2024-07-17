package hsec

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-api-go/grpc/resources/nstore/hsecdb"
	"n2x.dev/x-lib/pkg/n2xp/queuing"
	"n2x.dev/x-lib/pkg/runtime"
	"n2x.dev/x-lib/pkg/xlog"
)

var RequestQueue = make(chan *hsecdb.HostSecurityReportRequest, 128)

func Scanner(w *runtime.Wrkr) {
	xlog.Infof("Started worker %s", w.Name)
	w.Running = true

	endCh := make(chan struct{}, 2)

	go func() {
		securityScannerCtl := make(chan struct{}, 1)
		go func() {
			time.Sleep(30 * time.Second)
			securityScannerCtl <- struct{}{}
		}()

		n2xID := viper.GetString("n2x.id")

		ticker := time.NewTicker(24 * 3600 * time.Second) // 24 hours
		defer ticker.Stop()

		for {
			select {
			case <-securityScannerCtl:
				if err := scan(); err != nil {
					xlog.Warnf("[host-security] Unable to complete security scan: %v", errors.Cause(err))
					continue
				}

			case <-ticker.C:
				securityScannerCtl <- struct{}{}

			case r := <-RequestQueue:
				hsr, err := readReportFile()
				if err != nil {
					xlog.Warnf("[host-security] Unable to get host security report: %v", errors.Cause(err))
				}

				hsrr := query(r, hsr) // hsr can be nil

				queuing.TxControlQueue <- &n2xsp.Payload{
					SrcID:           n2xID,
					DstControllerID: r.Request.ControllerID,
					Type:            n2xsp.PDUType_NODEMGMT,
					NodeMgmtPDU: &n2xsp.NodeMgmtPDU{
						Type:               n2xsp.NodeMgmtMsgType_NODE_HOST_SECURITY_RESPONSE,
						HsecReportResponse: hsrr,
					},
				}

			case <-endCh:
				// xlog.Warn("[host-security] Closing security scanner")
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
