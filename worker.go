package linda

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type worker struct {
	process
}

type workerFunc func(...interface{}) error

func newWorker(id string) (*worker, error) {
	process, err := newProcess("worker" + id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &worker{
		process: *process,
	}, nil
}

func (w *worker) work(jobs <-chan *Job, monitor *sync.WaitGroup) {
	monitor.Add(1)
	go func() {
		logrus.Debugf("worker {%s} start...", w)
		defer func() {
			logrus.Debugf("worker {%s} stop...", w)
			defer monitor.Done()
		}()
		for job := range jobs {
			if workerFunc, ok := workers[job.Payload.Class]; ok {
				if err := w.run(job, workerFunc); err != nil {
					logrus.Error(err)
				}
				if err := w.nextRun(job); err != nil {
					logrus.Error(err)
				}
			} else {
				logrus.Errorf("no worker for job {%s}", job)
			}
		}
	}()
}

func (w *worker) run(job *Job, workerFunc workerFunc) error {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error(r)
		}
	}()
	logrus.Infof("run job {%s}", job)
	return workerFunc(job.Payload.Args...)
}

func (w *worker) nextRun(job *Job) error {
	var err error
	if job.Period == 0 {
		err = brokerConn.Delete(job.Queue, job)
	} else {
		err = brokerConn.Release(job.Queue, job, job.Period)
	}
	return err
}
