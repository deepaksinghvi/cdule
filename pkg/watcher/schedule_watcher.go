package watcher

import (
	"github.com/deepaksinghvi/cdule/pkg/job"
	"github.com/deepaksinghvi/cdule/pkg/model"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"reflect"
	"sync"
	"time"
)

type ScheduleWatcher struct {
	Closed chan struct{}
	WG     sync.WaitGroup
	Ticker *time.Ticker
}

var lastScheduleExeucitonTime int64
var nextScheduleExeuctionTime int64

func (t *ScheduleWatcher) Run() {
	for {
		select {
		case <-t.Closed:
			return
		case <-t.Ticker.C:
			now := time.Now()
			lastScheduleExeucitonTime = now.Add(-1 * time.Minute).UnixNano()
			nextScheduleExeuctionTime = now.UnixNano()

			log.Infof("lastScheduleExeucitonTime %d", lastScheduleExeucitonTime)
			log.Infof("nextScheduleExeuctionTime %d", nextScheduleExeuctionTime)
			runNextScheduleJobs(lastScheduleExeucitonTime, nextScheduleExeuctionTime)
		}
	}
}

func (t *ScheduleWatcher) Stop() {
	close(t.Closed)
	t.WG.Wait()
}

func runNextScheduleJobs(scheduleStart, scheduleEnd int64) {
	defer panicRecoveryForSchedule()
	schedules, err := model.CduleRepos.CduleRepository.GetScheduleBetween(scheduleStart, scheduleEnd)
	if nil != err {
		log.Error(err)
		return
	}
	for _, schedule := range schedules {
		log.Infof("Schedule ID Exeuction Time %d for Job ID: %d", schedule.ExecutionID, schedule.JobID)
	}

	for _, schedule := range schedules {
		scheduledJob, err := model.CduleRepos.CduleRepository.GetJob(schedule.JobID)
		var jobHistory *model.JobHistory
		if err == nil {
			jobHistory, err = model.CduleRepos.CduleRepository.GetJobHistoryForSchedule(schedule.ExecutionID)
			j := jobRegistry[scheduledJob.JobName]
			jobInstance := reflect.New(j).Elem().Interface()
			if err != nil && err.Error() == "record not found" && jobHistory != nil {
				// if job history was present but not executed
				if jobHistory.Status == model.JobStatusNew {
					jobHistory.Status = model.JobStatusInProgress
					model.CduleRepos.CduleRepository.UpdateJobHistory(jobHistory)
					jobInstance.(job.Job).Execute()
				}
			} else {
				// if job history is not there for this schedule, so this should be executed.
				jobHistory = &model.JobHistory{
					JobID:       schedule.JobID,
					ExecutionID: schedule.ExecutionID,
					DeletedAt:   gorm.DeletedAt{},
					Status:      model.JobStatusNew,
					WorkerID:    schedule.WorkerID,
					RetryCount:  0,
				}
				model.CduleRepos.CduleRepository.CreateJobHistory(jobHistory)
				jobHistory.Status = model.JobStatusInProgress
				model.CduleRepos.CduleRepository.UpdateJobHistory(jobHistory)

				executeJob(jobInstance, jobHistory)
			}

			// Calculate the next schedule for the current job
			storedJob, err := model.CduleRepos.CduleRepository.GetJobByName(jobInstance.(job.Job).JobName())
			if err != nil {
				log.Error(err.Error())
				return
			}
			SchedulerParser, err := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow).Parse(storedJob.CronExpression)
			if err != nil {
				log.Error(err.Error())
				return
			}
			nextRunTime := SchedulerParser.Next(time.Now()).UnixNano()
			newSchedule := model.Schedule{
				ExecutionID: nextRunTime,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				DeletedAt:   gorm.DeletedAt{},
				// TODO to be updated based on the number of jobs are getting executed by any worker. To load balance the job execution
				WorkerID: schedule.WorkerID,
				JobID:    schedule.JobID,
			}
			model.CduleRepos.CduleRepository.CreateSchedule(&newSchedule)

		}
	}
	log.Infof("Scheduler Completed For StartTime %d To EndTime %d", scheduleStart, scheduleEnd)
}

/*
For go 1.17 following method can be used.
func executeJob(jobInstance interface{}, jobHistory *model.JobHistory) {
	defer panicRecovery(jobHistory)
	jobInstance.(job.Job).Execute()
}
*/

/*
cdule library has been built and developed using go 1.18 (go1.18beta2), if you need to use it for 1.17
then build from source by uncommenting the above method and comment the follwoing
*/
func executeJob(jobInstance any, jobHistory *model.JobHistory) {
	defer panicRecovery(jobHistory)
	jobInstance.(job.Job).Execute()
}

// If there is any panic from Job Execution, set the JobStatus as FAILED
func panicRecovery(jobHistory *model.JobHistory) {
	// TODO should be handled for any panic and set the status as FAILED for job history with error message
	jobHistory.Status = model.JobStatusCompleted
	if r := recover(); r != nil {
		log.Warning("Recovered in panicRecovery for job execution ", r)
		jobHistory.Status = model.JobStatusFailed
	}
	model.CduleRepos.CduleRepository.UpdateJobHistory(jobHistory)
}

func panicRecoveryForSchedule() {
	if r := recover(); r != nil {
		log.Warning("Recovered in runNextScheduleJobs ", r)
	}
}
