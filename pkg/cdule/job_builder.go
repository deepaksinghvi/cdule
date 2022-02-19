package cdule

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/deepaksinghvi/cdule/pkg"
	"github.com/deepaksinghvi/cdule/pkg/model"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var jobRegistry = make(map[string]reflect.Type)

var ScheduleParser cron.Parser

func registerType(job Job) {
	t := reflect.TypeOf(job).Elem()
	jobRegistry[job.JobName()] = t
}

type AbstractJob struct {
	Job     Job
	JobData map[string]string
}

func NewJob(job Job, jobData map[string]string) *AbstractJob {
	aj := &AbstractJob{
		Job:     job,
		JobData: jobData,
	}
	return aj
}
func (j *AbstractJob) Build(cronExpression string) (*model.Job, error) {
	// register job, this is used later to get the type of a job
	registerType(j.Job)
	jobDataBytes, err := json.Marshal(j.JobData)
	if nil != err {
		log.Errorf("Error %s for JobName %s", err.Error(), j.Job.JobName())
		return nil, errors.New(fmt.Sprintf("Invalid Job Data %v", j.JobData))
	}
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
	newJobModel, err := model.CduleRepos.CduleRepository.GetJobByName(newJob.JobName)
	if nil != newJobModel {
		return nil, errors.New(fmt.Sprintf("Job with Name: %s already exists", newJob.JobName))
	}
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
		WorkerID:    getWorkerID(),
		JobID:       job.ID,
		JobData:     job.JobData,
	}
	_, err = model.CduleRepos.CduleRepository.CreateSchedule(firstSchedule)
	if err != nil {
		log.Error(err.Error())
		return job, err
	}
	return job, err
}
