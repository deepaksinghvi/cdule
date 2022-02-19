package cdule

import (
	"fmt"
	"os"
	"time"

	"github.com/deepaksinghvi/cdule/pkg/model"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var WorkerID string

type Cdule struct {
	*WorkerWatcher
	*ScheduleWatcher
}

func init() {
	WorkerID = getWorkerID()
}
func (cdule *Cdule) NewCduleWithWorker(workerName string, param ...string) {
	WorkerID = workerName
	cdule.NewCdule(param...)
}
func (cdule *Cdule) NewCdule(param ...string) {
	if nil == param {
		param = []string{"./resources", "config"} // default path for resources
	}
	model.ConnectDataBase(param)
	worker, err := model.CduleRepos.CduleRepository.GetWorker(WorkerID)
	if nil != err {
		log.Errorf("Error getting workder %s ", err.Error())
	}
	if nil != worker {
		worker.UpdatedAt = time.Now()
		model.CduleRepos.CduleRepository.UpdateWorker(worker)
	} else {
		// First time cdule started on a worker node
		worker := model.Worker{
			WorkerID:  WorkerID,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		}
		model.CduleRepos.CduleRepository.CreateWorker(&worker)
	}

	cdule.createWatcherAndWaitForSignal()
}

func (cdule *Cdule) createWatcherAndWaitForSignal() {
	/*
		schedule watcher stop logic to abort program with signal like ctrl + c
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)*/

	workerWatcher := createWorkerWatcher()
	schedulerWatcher := createSchedulerWatcher()
	cdule.WorkerWatcher = workerWatcher
	cdule.ScheduleWatcher = schedulerWatcher
	/*select {
	case sig := <-c:
		fmt.Printf("Received %s signal. Aborting...\n", sig)
		workerWatcher.Stop()
		schedulerWatcher.Stop()
	}*/
}

func (cdule Cdule) StopWatcher() {
	cdule.WorkerWatcher.Stop()
	cdule.ScheduleWatcher.Stop()
}
func createWorkerWatcher() *WorkerWatcher {
	workerWatcher := &WorkerWatcher{
		Closed: make(chan struct{}),
		Ticker: time.NewTicker(time.Second * 30), // used for worker health check update in db.
	}

	workerWatcher.WG.Add(1)
	go func() {
		defer workerWatcher.WG.Done()
		workerWatcher.Run()
	}()
	return workerWatcher
}

func createSchedulerWatcher() *ScheduleWatcher {
	scheduleWatcher := &ScheduleWatcher{
		Closed: make(chan struct{}),
		Ticker: time.NewTicker(time.Minute * 1), // used for worker health check update in db.
	}

	scheduleWatcher.WG.Add(1)
	go func() {
		defer scheduleWatcher.WG.Done()
		scheduleWatcher.Run()
	}()
	return scheduleWatcher
}

func getWorkerID() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return hostname
}
