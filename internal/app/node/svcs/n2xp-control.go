package svcs

import (
	"context"
	"io"
	"time"

	"github.com/spf13/viper"
	n2xsp_pb "n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/n2xp/queuing"
	"n2x.dev/x-lib/pkg/runtime"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/mnet"
	"n2x.dev/x-node/internal/app/node/n2xsp"
)

// Control method implementation of NetworkAPI gRPC Service
func NetworkControl(w *runtime.Wrkr) {
	xlog.Infof("Started worker %s", w.Name)
	w.Running = true

	endCh := make(chan struct{}, 2)

	go func() {
		stream, err := w.NxNC.Control(context.Background())
		if err != nil {
			xlog.Errorf("Unable to get n2xp stream from controller: %v", err)
			mnet.LocalNode().Connection().Watcher() <- struct{}{}
			return
		}

		go func() {
			for {
				payload, err := stream.Recv()
				if err == io.EOF {
					// xlog.Warnf("Ended (io.EOF) n2xp stream: %v", err)
					mnet.LocalNode().Connection().Watcher() <- struct{}{}
					break
				}
				if err != nil {
					// xlog.Warnf("Unable to receive n2xp payload: %v", err)
					mnet.LocalNode().Connection().Watcher() <- struct{}{}
					break
				}

				// if !serviceEnabled {
				// 	continue
				// }

				n2xsp.Preprocessor(context.TODO(), payload)
			}
			if err := stream.CloseSend(); err != nil {
				xlog.Errorf("Unable to close n2xp stream: %v", err)
			}
			endCh <- struct{}{}
			// xlog.Warn("Closing n2xp recv stream")
		}()

		go func() {
			for {
				select {
				case payload := <-n2xsp.RxQueue:
					xlog.Debug("[n2xp] Received n2xp payload on queue")
					go n2xsp.Processor(context.TODO(), payload)

				case payload := <-queuing.TxControlQueue:
					if err := stream.Send(payload); err != nil {
						// xlog.Warnf("[n2xp] Unable to send n2xp payload: %v", err)

						mnet.LocalNode().Connection().Watcher() <- struct{}{}

						if err := stream.CloseSend(); err != nil {
							xlog.Errorf("Unable to close n2xp stream: %v", err)
						}
						return
					}
				case <-endCh:
					// xlog.Warn("Closing n2xp send stream")
					return
				}
			}
		}()

		queuing.TxControlQueue <- &n2xsp_pb.Payload{
			SrcID: viper.GetString("n2x.id"),
			Type:  n2xsp_pb.PDUType_NODEMGMT,
			NodeMgmtPDU: &n2xsp_pb.NodeMgmtPDU{
				Type:    n2xsp_pb.NodeMgmtMsgType_NODE_INIT,
				NodeReq: mnet.LocalNode().NodeReq(),
			},
		}
	}()

	go n2xpCtl()

	<-w.QuitChan

	endCh <- struct{}{}

	w.WG.Done()
	w.Running = false
	xlog.Infof("Stopped worker %s", w.Name)
}

var n2xpCtlRun bool

func n2xpCtl() {
	n2xID := viper.GetString("n2x.id")

	if !mnet.LocalNode().Node().Cfg.DisableNetworking || mnet.LocalNode().Router() != nil {
		return
	}

	if !n2xpCtlRun {
		n2xpCtlRun = true
		for {
			queuing.TxControlQueue <- &n2xsp_pb.Payload{
				SrcID: n2xID,
				Type:  n2xsp_pb.PDUType_SESSION,
				SessionPDU: &n2xsp_pb.SessionPDU{
					Type:      n2xsp_pb.SessionMsgType_SESSION_KEEPALIVE,
					SessionID: n2xID,
				},
			}
			time.Sleep(30 * time.Second)
		}
	}
}
