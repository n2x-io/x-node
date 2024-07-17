package ops

import (
	"context"
	"sync"
	"time"

	"n2x.dev/x-api-go/grpc/common/datetime"
	"n2x.dev/x-api-go/grpc/network/n2xsp"
	"n2x.dev/x-lib/pkg/runtime"
	"n2x.dev/x-lib/pkg/utils"
	"n2x.dev/x-lib/pkg/xlog"
)

type atTime int64

type jobsAtTime struct {
	jobs    map[atTime][]*n2xsp.WorkflowPDU
	running map[atTime]bool
	sync.RWMutex
}

var atdCommandQueue = make(chan *n2xsp.WorkflowPDU, 128)

func Atd(w *runtime.Wrkr) {
	xlog.Infof("Started worker %s", w.Name)
	w.Running = true

	jobsAtd := newJobsAtTime()

	var wwg sync.WaitGroup
	waitc := make(chan struct{}, 1)

	wwg.Add(1)
	go atdRunner(jobsAtd, waitc, &wwg)

	for {
		select {
		case pdu := <-atdCommandQueue:
			xlog.Info("Received workflow on atdCommandQueue")

			wf := pdu.Workflow

			t, err := utils.GetDateTime(wf.Triggers.Schedule.DateTime)
			if err != nil {
				xlog.Errorf("Unable to schedule workflow %s: %v", wf.WorkflowID, err)
				continue
			}

			xlog.Infof("Updating schedule for workflow %s", wf.WorkflowID)

			jobsAtd.deleteAtJobs(t)

			if wf.Enabled {
				xlog.Infof("Enabling schedule for workflow %s", wf.WorkflowID)
				jobsAtd.setAtJobs(t, pdu)
			} else {
				xlog.Infof("Schedule for workflow %s has been disabled", wf.WorkflowID)
			}

		case <-w.QuitChan:
			xlog.Alert("Stopping workflow scheduler.. Be careful, there might be jobs running!")
			waitc <- struct{}{}
			wwg.Wait()
			w.WG.Done()
			w.Running = false
			xlog.Infof("Stopped worker %s", w.Name)
			return
		}
	}
}

func atdRunner(jobsAtd *jobsAtTime, waitc chan struct{}, wwg *sync.WaitGroup) {
	go func() {
		for {
			tm := time.Now()
			dateTime := &datetime.DateTime{
				Year:   int32(tm.Year()),
				Month:  int32(tm.Month()),
				Day:    int32(tm.Day()),
				Hour:   int32(tm.Hour()),
				Minute: int32(tm.Minute()),
				Second: 0,
			}
			t, err := utils.GetDateTime(dateTime)
			if err != nil {
				xlog.Errorf("Unable to get atTime: %v", err)
				continue
			}

			if pl := jobsAtd.runAtJobs(t); pl != nil {
				go func() {
					for _, p := range pl {
						if err := WorkflowExpedite(context.TODO(), p); err != nil {
							xlog.Errorf("Workflow %s finished abnormally: %v", p.Workflow.WorkflowID, err)
						}
					}
					jobsAtd.deleteAtJobs(t)
				}()
			}

			time.Sleep(time.Minute)
		}
	}()

	<-waitc
	xlog.Info("Stopped runner for scheduled workflows")
	wwg.Done()
}

func newJobsAtTime() *jobsAtTime {
	return &jobsAtTime{
		jobs:    make(map[atTime][]*n2xsp.WorkflowPDU),
		running: make(map[atTime]bool),
	}
}

func (j *jobsAtTime) setAtJobs(t int64, pdu *n2xsp.WorkflowPDU) {
	j.Lock()
	j.jobs[atTime(t)] = append(j.jobs[atTime(t)], pdu)
	j.running[atTime(t)] = false
	j.Unlock()
}

func (j *jobsAtTime) deleteAtJobs(t int64) {
	j.Lock()
	if _, ok := j.jobs[atTime(t)]; ok {
		delete(j.jobs, atTime(t))
	}
	if _, ok := j.running[atTime(t)]; ok {
		delete(j.running, atTime(t))
	}
	j.Unlock()
}

func (j *jobsAtTime) runAtJobs(t int64) []*n2xsp.WorkflowPDU {
	j.Lock()
	defer j.Unlock()

	pl, ok := j.jobs[atTime(t)]
	if ok {
		j.running[atTime(t)] = true

		return pl
	}

	return nil
}

func (j *jobsAtTime) atJobsRunning(t int64) bool {
	j.Lock()
	defer j.Unlock()

	if running, ok := j.running[atTime(t)]; ok {
		return running
	}

	return false
}
