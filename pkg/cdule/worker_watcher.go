package cdule

import (
	"sync"
	"time"

	"github.com/deepaksinghvi/cdule/pkg/model"

	log "github.com/sirupsen/logrus"
)

// WorkerWatcher struct
type WorkerWatcher struct {
	Closed chan struct{}
	WG     sync.WaitGroup
	Ticker *time.Ticker
}

// Run to run watcher in a continuous loop
func (t *WorkerWatcher) Run() {
	for {
		select {
		case <-t.Closed:
			return
		case <-t.Ticker.C:
			healthCheckUpdate()
		}
	}
}

// Stop to stop worker watcher
func (t *WorkerWatcher) Stop() {
	close(t.Closed)
	t.WG.Wait()
}

func healthCheckUpdate() {
	worker, err := model.CduleRepos.CduleRepository.GetWorker(WorkerID)
	if nil != err {
		log.Errorf("Error getting workder %s ", err.Error())
	}
	if nil != worker {
		worker.UpdatedAt = time.Now()
		model.CduleRepos.CduleRepository.UpdateWorker(worker)
		log.Debugf("Health check updated for worker_id %s updated", WorkerID)
		return
	}
	log.Warningf("Health check update failed for worker_id %s", WorkerID)
}
