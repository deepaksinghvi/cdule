package cdule

import (
	"encoding/json"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/deepaksinghvi/cdule/pkg"
	"github.com/deepaksinghvi/cdule/pkg/model"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ScheduleWatcher struct {
	Closed chan struct{}
	WG     sync.WaitGroup
	Ticker *time.Ticker
}

var lastScheduleExecutionTime int64
var nextScheduleExecutionTime int64

func (t *ScheduleWatcher) Run() {
	for {
		select {
		case <-t.Closed:
			return
		case <-t.Ticker.C:
			now := time.Now()
			lastScheduleExecutionTime = now.Add(-1 * time.Minute).UnixNano()
			nextScheduleExecutionTime = now.UnixNano()

			log.Infof("lastScheduleExecutionTime %d", lastScheduleExecutionTime)
			log.Infof("nextScheduleExecutionTime %d", nextScheduleExecutionTime)
			runNextScheduleJobs(lastScheduleExecutionTime, nextScheduleExecutionTime)
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
	workers, err := model.CduleRepos.CduleRepository.GetWorkers()
	if nil != err {
		log.Error(err)
		return
	}
	for _, schedule := range schedules {
		log.Infof("Schedule ID Exeuction Time %d for Job ID: %d", schedule.ExecutionID, schedule.JobID)
		scheduledJob, err := model.CduleRepos.CduleRepository.GetJob(schedule.JobID)
		if nil != err {
			log.Errorf("Error while running Schedule for %d : %s", schedule.JobID, err.Error())
			continue
		}
		jobDataStr := schedule.JobData
		var jobDataMap map[string]string
		if pkg.EMPTYSTRING != jobDataStr {
			err = json.Unmarshal([]byte(jobDataStr), &jobDataMap)
			if nil != err {
				log.Error(err)
				continue
			}
		}
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
					jobInstance.(Job).Execute(jobDataMap)
					jobDataMap = jobInstance.(Job).GetJobData()
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

				jobDataMap = executeJob(jobInstance, jobHistory, &jobDataMap)
			}

			// Calculate the next schedule for the current job
			storedJob, err := model.CduleRepos.CduleRepository.GetJobByName(jobInstance.(Job).JobName())
			if err != nil {
				log.Error(err.Error())
				return
			}
			jobDataBytes, err := json.Marshal(jobDataMap)
			if nil != err {
				log.Errorf("Error %s for JobName %s and Schedule ID %d ", err.Error(), storedJob.JobName, schedule.ExecutionID)
			}
			if string(jobDataBytes) != pkg.EMPTYSTRING {
				jobDataStr = string(jobDataBytes)
			}
			SchedulerParser, err := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow).Parse(storedJob.CronExpression)
			if err != nil {
				log.Error(err.Error())
				return
			}
			nextRunTime := SchedulerParser.Next(time.Now()).UnixNano()
			workerIDForNextRun, _ := findNextAvailableWorker(workers, schedule)
			newSchedule := model.Schedule{
				ExecutionID: nextRunTime,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				DeletedAt:   gorm.DeletedAt{},
				WorkerID:    workerIDForNextRun,
				JobID:       schedule.JobID,
				JobData:     jobDataStr,
			}
			model.CduleRepos.CduleRepository.CreateSchedule(&newSchedule)

		}
	}
	log.Infof("Scheduler Completed For StartTime %d To EndTime %d", scheduleStart, scheduleEnd)
}

type WorkerJobCount struct {
	WorkerID string `json:"worker_id"`
	Count    int64  `json:"count"`
}

type workerJobCountList []WorkerJobCount

func (w workerJobCountList) Len() int {
	return len(w)
}

func (w workerJobCountList) Less(i, j int) bool {
	return w[i].Count > w[j].Count
}

func (w workerJobCountList) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}
func findNextAvailableWorker(workers []model.Worker, schedule model.Schedule) (string, error) {
	workerName := schedule.WorkerID
	var result []WorkerJobCount
	model.DB.Raw("SELECT worker_id, count(1) FROM job_histories WHERE job_id = ? group by worker_id", schedule.JobID).Scan(&result)
	//sort.Sort(workerJobCountList(result))
	if len(result) <= 0 {
		return workerName, nil
	}
	sort.Slice(result[:], func(i, j int) bool {
		return result[i].Count < result[j].Count
	})
	log.Infof("Next Job scheduled for JobID %d with worker %s", schedule.JobID, result[0].WorkerID)
	return result[0].WorkerID, nil
}

/*
For go 1.17 following method can be used.
func executeJob(jobInstance interface{}, jobHistory *model.JobHistory, jobDataMap map[string]string) {
	defer panicRecovery(jobHistory)
	jobInstance.(job.Job).Execute(jobDataMap)
}
*/

/*
cdule library has been built and developed using go 1.18 (go1.18beta2), if you need to use it for 1.17
then build from source by uncommenting the above method and comment the following
*/
func executeJob(jobInstance any, jobHistory *model.JobHistory, jobDataMap *map[string]string) map[string]string {
	defer panicRecovery(jobHistory)
	jobInstance.(Job).Execute(*jobDataMap)
	return jobInstance.(Job).GetJobData()
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
