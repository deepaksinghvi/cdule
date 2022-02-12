package watcher

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/deepaksinghvi/cdule/pkg/job"
	"github.com/deepaksinghvi/cdule/pkg/model"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var jobRegistry = make(map[string]reflect.Type)

var ScheduleParser cron.Parser

func registerType(job job.Job) {
	t := reflect.TypeOf(job).Elem()
	jobRegistry[job.JobName()] = t
}

type AbstractJob struct {
	Job     job.Job
	JobData map[string]string
}

func (j *AbstractJob) Register(job job.Job) {
	registerType(job)

}

func NewJob(job job.Job, jobData map[string]string) *AbstractJob {
	aj := &AbstractJob{
		Job:     job,
		JobData: jobData,
	}
	return aj
}
func (j *AbstractJob) Build(cronExpression string) (*model.Job, error) {
	// register job, this is used later to get the type of a job
	registerType(j.Job)

	newJob := &model.Job{
		Model:          model.Model{},
		JobName:        j.Job.JobName(),
		GroupName:      "",
		CronExpression: cronExpression,
		Expired:        false,
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
	schedule1 := &model.Schedule{
		ExecutionID: nextRunTime,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		DeletedAt:   gorm.DeletedAt{},
		WorkerID:    WorkerID,
		JobID:       job.ID,
	}
	_, err = model.CduleRepos.CduleRepository.CreateSchedule(schedule1)
	if err != nil {
		log.Error(err.Error())
		return job, err
	}
	return job, err
}
