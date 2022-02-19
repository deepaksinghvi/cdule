package model

import (
	"gorm.io/gorm"

	"github.com/deepaksinghvi/cdule/pkg"
)

type cduleRepository struct {
	DB *gorm.DB
}

func NewCduleRepository(db *gorm.DB) CduleRepository {
	return cduleRepository{
		DB: db,
	}
}

type CduleRepository interface {
	CreateWorker(worker *Worker) (*Worker, error)
	UpdateWorker(worker *Worker) (*Worker, error)
	GetWorker(workerID string) (*Worker, error)
	DeleteWorker(workerID string) (*Worker, error)

	CreateJob(job *Job) (*Job, error)
	UpdateJob(job *Job) (*Job, error)
	GetJob(jobID int64) (*Job, error)
	GetJobByName(name string) (*Job, error)
	GetJobs(workerID string) ([]Job, error)
	DeleteJob(jobID int64) (*Job, error)

	CreateJobHistory(jobHistory *JobHistory) (*JobHistory, error)
	UpdateJobHistory(jobHistory *JobHistory) (*JobHistory, error)
	GetJobHistory(jobID int64) ([]JobHistory, error)
	GetJobHistoryForSchedule(scheduleID int64) (*JobHistory, error)
	DeleteJobHistory(jobID int64) ([]JobHistory, error)

	CreateSchedule(schedule *Schedule) (*Schedule, error)
	UpdateSchedule(schedule *Schedule) (*Schedule, error)
	GetSchedule(scheduleID int64) (*Schedule, error)
	GetScheduleBetween(scheduleStart, scheduleEnd int64) ([]Schedule, error)
	GetSchedulesForJob(jobID int64) ([]Schedule, error)
	GetSchedulesForWorker(workerID string) ([]Schedule, error)
	DeleteScheduleForJob(jobID int64) ([]Schedule, error)
	DeleteScheduleForWorker(workerID string) ([]Schedule, error)
}

func (c cduleRepository) CreateWorker(worker *Worker) (*Worker, error) {
	if err := c.DB.Create(worker).Error; err != nil {
		return nil, err
	}
	return worker, nil
}

func (c cduleRepository) UpdateWorker(worker *Worker) (*Worker, error) {
	if err := c.DB.Updates(worker).Error; err != nil {
		return nil, err
	}
	return worker, nil
}

func (c cduleRepository) GetWorker(workerID string) (*Worker, error) {
	var worker Worker
	if err := c.DB.Where("worker_id = ?", workerID).Find(&worker).Error; err != nil {
		return nil, err
	}
	if worker.WorkerID == pkg.EMPTYSTRING {
		return nil, nil
	}
	return &worker, nil
}

func (c cduleRepository) DeleteWorker(workerID string) (*Worker, error) {
	var worker Worker
	if err := c.DB.Where("worker_id = ?", workerID).First(&worker).Error; err != nil {
		return nil, err
	}
	if err := c.DB.Delete(&worker).Error; err != nil {
		return nil, err
	}
	return &worker, nil
}

func (c cduleRepository) CreateJob(job *Job) (*Job, error) {
	if err := c.DB.Create(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (c cduleRepository) UpdateJob(job *Job) (*Job, error) {
	if err := c.DB.Updates(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (c cduleRepository) GetJob(jobID int64) (*Job, error) {
	var job Job
	if err := c.DB.Where("id = ?", jobID).Find(&job).Error; err != nil {
		return nil, err
	}
	if job.ID == 0 {
		return nil, nil
	}
	return &job, nil
}

func (c cduleRepository) GetJobByName(jobName string) (*Job, error) {
	var job Job
	if err := c.DB.Where("job_name = ?", jobName).Find(&job).Error; err != nil {
		return nil, err
	}
	if job.ID == 0 {
		return nil, nil
	}
	return &job, nil
}
func (c cduleRepository) GetJobs(workerID string) ([]Job, error) {
	var jobs []Job
	if err := c.DB.Where("worker_id = ? and expired=false", workerID).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (c cduleRepository) DeleteJob(jobID int64) (*Job, error) {
	var job Job
	if err := c.DB.Where("id = ?", jobID).First(&job).Error; err != nil {
		return nil, err
	}
	if err := c.DB.Delete(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (c cduleRepository) CreateJobHistory(jobHistory *JobHistory) (*JobHistory, error) {
	if err := c.DB.Create(jobHistory).Error; err != nil {
		return nil, err
	}
	return jobHistory, nil
}

func (c cduleRepository) UpdateJobHistory(jobHistory *JobHistory) (*JobHistory, error) {
	if err := c.DB.Updates(jobHistory).Error; err != nil {
		return nil, err
	}
	return jobHistory, nil
}

func (c cduleRepository) GetJobHistory(jobID int64) ([]JobHistory, error) {
	var jobHistories []JobHistory
	if err := c.DB.Where("job_id = ?", jobID).First(&jobHistories).Error; err != nil {
		return nil, err
	}
	return jobHistories, nil
}

func (c cduleRepository) GetJobHistoryForSchedule(scheduleID int64) (*JobHistory, error) {
	var jobHistory JobHistory
	if err := c.DB.Where("execution_id = ?", scheduleID).First(&jobHistory).Error; err != nil {
		return nil, err
	}
	return &jobHistory, nil
}

func (c cduleRepository) DeleteJobHistory(jobID int64) ([]JobHistory, error) {
	jobHistories, err := c.GetJobHistory(jobID)
	if nil != err {
		return nil, err
	}
	if err := c.DB.Delete(&jobHistories).Error; err != nil {
		return nil, err
	}
	return jobHistories, nil
}

func (c cduleRepository) CreateSchedule(schedule *Schedule) (*Schedule, error) {
	if err := c.DB.Create(schedule).Error; err != nil {
		return nil, err
	}
	return schedule, nil
}

func (c cduleRepository) UpdateSchedule(schedule *Schedule) (*Schedule, error) {
	if err := c.DB.Updates(schedule).Error; err != nil {
		return nil, err
	}
	return schedule, nil
}

func (c cduleRepository) GetSchedule(executionID int64) (*Schedule, error) {
	var schedule Schedule
	if err := c.DB.Where("execution_id = ?", executionID).Find(&schedule).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (c cduleRepository) GetScheduleBetween(scheduleStart, scheduleEnd int64) ([]Schedule, error) {
	var schedules []Schedule
	if err := c.DB.Where("execution_id >= ? and execution_id <= ?", scheduleStart, scheduleEnd).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (c cduleRepository) GetSchedulesForJob(jobID int64) ([]Schedule, error) {
	var schedules []Schedule
	if err := c.DB.Where("job_id = ?", jobID).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (c cduleRepository) GetSchedulesForWorker(workerID string) ([]Schedule, error) {
	var schedules []Schedule
	if err := c.DB.Where("worker_id = ?", workerID).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (c cduleRepository) DeleteScheduleForJob(jobID int64) ([]Schedule, error) {
	schedules, err := c.GetSchedulesForJob(jobID)
	if nil != err {
		return nil, err
	}
	for _, schedule := range schedules {
		if err := c.DB.Where("job_id = ? and execution_id = ?",
			schedule.JobID, schedule.ExecutionID).Delete(&Schedule{}).Error; err != nil {
			return nil, err
		}
	}
	return schedules, nil
}

func (c cduleRepository) DeleteScheduleForWorker(workerID string) ([]Schedule, error) {
	schedules, err := c.GetSchedulesForWorker(workerID)
	if nil != err {
		return nil, err
	}
	for _, schedule := range schedules {
		if err := c.DB.Where("job_id = ? and execution_id = ?",
			schedule.JobID, schedule.ExecutionID).Delete(&Schedule{}).Error; err != nil {
			return nil, err
		}
	}
	return schedules, nil
}
