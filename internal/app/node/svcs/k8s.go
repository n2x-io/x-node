package svcs

import (
	"n2x.dev/x-lib/pkg/runtime"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/k8s"
	"n2x.dev/x-node/internal/app/node/mnet"
)

func KubernetesConnector(w *runtime.Wrkr) {
	xlog.Infof("Started worker %s", w.Name)
	w.Running = true

	quitKCtrlCh := make(chan struct{})

	go func() {
		if !mnet.LocalNode().IsK8sGwEnabled() {
			xlog.Info("Kubernetes controller not enabled")
			return
		}

		if err := k8s.KubernetesController(quitKCtrlCh); err != nil {
			xlog.Tracef("Unable to start kubernetes controller: %v", err)
			xlog.Info("Kubernetes environment not detected")
		} else {
			xlog.Info("Kubernetes controller running")
		}
	}()

	<-w.QuitChan
	close(quitKCtrlCh)

	w.WG.Done()
	w.Running = false
	xlog.Infof("Stopped worker %s", w.Name)
}
