package cdule

import (
	"time"

	"github.com/deepaksinghvi/cdule/pkg/model"
	"github.com/deepaksinghvi/cdule/pkg/watcher"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Cdule struct {
	*watcher.WorkerWatcher
	*watcher.ScheduleWatcher
}

func (cdule Cdule) NewCdule() {
	model.ConnectDataBase()
	worker, err := model.CduleRepos.CduleRepository.GetWorker(watcher.WorkerID)
	if nil != err {
		log.Errorf("Error getting workder %s ", err.Error())
	}
	if nil != worker {
		worker.UpdatedAt = time.Now()
		model.CduleRepos.CduleRepository.UpdateWorker(worker)
	} else {
		// First time cdule started on a worker node
		worker := model.Worker{
			WorkerID:  watcher.WorkerID,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		}
		model.CduleRepos.CduleRepository.CreateWorker(&worker)
	}
	/*myJob := job.MyJob{}
	jobModel, err := watcher.NewJob(&myJob, nil).Build(utils.EveryMinute)
	log.Info(jobModel)

	panicJob := job.PanicJob{}
	panicJobModel, err := watcher.NewJob(&panicJob, nil).Build(utils.EveryMinute)
	log.Info(panicJobModel)*/

	cdule.createWatcherAndWaitForSignal()
}

func (cdule Cdule) createWatcherAndWaitForSignal() {
	/*c := make(chan os.Signal)
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
func createWorkerWatcher() *watcher.WorkerWatcher {
	workerWatcher := &watcher.WorkerWatcher{
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

func createSchedulerWatcher() *watcher.ScheduleWatcher {
	scheduleWatcher := &watcher.ScheduleWatcher{
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
