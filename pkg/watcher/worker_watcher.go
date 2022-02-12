package watcher

import (
	"fmt"
	"github.com/deepaksinghvi/cdule/pkg/model"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

var WorkerID string

type WorkerWatcher struct {
	Closed chan struct{}
	WG     sync.WaitGroup
	Ticker *time.Ticker
}

func init() {
	WorkerID = getWorkerID()
}

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

func getWorkerID() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return hostname
}
