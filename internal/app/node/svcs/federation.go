package svcs

import (
	"time"

	"n2x.dev/x-lib/pkg/errors"
	"n2x.dev/x-lib/pkg/runtime"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/mnet"
	"n2x.dev/x-node/internal/app/node/mnet/connection"
)

var federationMonitorCh = make(chan struct{}, 1)

func FederationMonitor(w *runtime.Wrkr) {
	xlog.Infof("Started worker %s", w.Name)
	w.Running = true

	endCh := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case <-federationMonitorCh:
				xlog.Debug("Federation Monitor: Updating controller list...")
				if err := connection.FederationUpdate(w.NxNC); err != nil {
					xlog.Errorf("Unable to update controller list: %v", errors.Cause(err))
					mnet.LocalNode().Connection().Watcher() <- struct{}{}
					return
				}

			case <-endCh:
				// xlog.Debug("Closing federation monitor")
				return
			}
		}
	}()

	go federationMonitorCtl()

	<-w.QuitChan

	endCh <- struct{}{}

	w.WG.Done()
	w.Running = false
	xlog.Infof("Stopped worker %s", w.Name)
}

var federationMonitorRun bool

func federationMonitorCtl() {
	if !federationMonitorRun {
		federationMonitorRun = true
		go func() {
			for {
				federationMonitorCh <- struct{}{}
				time.Sleep(3600 * time.Second)
			}
		}()
	}
}
