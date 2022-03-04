package cdule

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/deepaksinghvi/cdule/pkg"
	"github.com/deepaksinghvi/cdule/pkg/model"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// JobRegistry job registry
var JobRegistry = make(map[string]reflect.Type)

// ScheduleParser cron parser
var ScheduleParser cron.Parser

func registerType(job Job) {
	t := reflect.TypeOf(job).Elem()
	JobRegistry[job.JobName()] = t
}

// AbstractJob for holding job and jobdata
type AbstractJob struct {
	Job     Job
	JobData map[string]string
}

// NewJob to create new abstract job
func NewJob(job Job, jobData map[string]string) *AbstractJob {
	aj := &AbstractJob{
		Job:     job,
		JobData: jobData,
	}
	return aj
}

// Build to build job and store in the database
func (j *AbstractJob) Build(cronExpression string) (*model.Job, error) {
	// register job, this is used later to get the type of a job
	registerType(j.Job)
	newJobModel, err := model.CduleRepos.CduleRepository.GetJobByName(j.Job.JobName())
	if nil != newJobModel || nil != err {
		return nil, fmt.Errorf("job with Name: %s already exists", newJobModel.JobName)
	}
	jobDataBytes, err := json.Marshal(j.JobData)
	/*if nil != err {
		log.Errorf("Error %s for JobName %s", err.Error(), j.Job.JobName())
		return nil, fmt.Errorf("invalid Job Data %v", j.JobData)
	}*/
	var jobDataStr = ""
	if string(jobDataBytes) != pkg.EMPTYSTRING {
		jobDataStr = string(jobDataBytes)
	}
	newJob := &model.Job{
		Model:          model.Model{},
		JobName:        j.Job.JobName(),
		GroupName:      "",
		CronExpression: cronExpression,
		Expired:        false,
		JobData:        jobDataStr,
	}
	SchedulerParser, err := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow).Parse(newJob.CronExpression)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	nextRunTime := SchedulerParser.Next(time.Now()).UnixNano()
	job, err := model.CduleRepos.CduleRepository.CreateJob(newJob)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	firstSchedule := &model.Schedule{
		ExecutionID: nextRunTime,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		DeletedAt:   gorm.DeletedAt{},
		WorkerID:    WorkerID,
		JobID:       job.ID,
		JobData:     job.JobData,
	}
	_, err = model.CduleRepos.CduleRepository.CreateSchedule(firstSchedule)
	if err != nil {
		log.Error(err.Error())
		return job, err
	}
	log.Infof("*** Job Scheduled Info ***\n JobName: %s,\n Schedule Cron: %s,\n Job Scheduled Time: %d,\n Worker: %s ",
		newJob.JobName, newJob.CronExpression, firstSchedule.ExecutionID, firstSchedule.WorkerID)
	return job, err
}
