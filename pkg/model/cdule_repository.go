package model

import (
	"gorm.io/gorm"
	"time"

	"github.com/deepaksinghvi/cdule/pkg"
)

type cduleRepository struct {
	DB    *gorm.DB
	Heart time.Duration
}

// NewCduleRepository cdule repository
func NewCduleRepository(db *gorm.DB, heart time.Duration) CduleRepository {
	return cduleRepository{
		DB:    db,
		Heart: heart,
	}
}

// CduleRepository cdule repository interface
type CduleRepository interface {
	CreateWorker(worker *Worker) (*Worker, error)
	UpdateWorker(worker *Worker) (*Worker, error)
	GetWorker(workerID string) (*Worker, error)
	GetWorkers() ([]Worker, error)
	DeleteWorker(workerID string) (*Worker, error)

	CreateJob(job *Job) (*Job, error)
	UpdateJob(job *Job) (*Job, error)
	GetJob(jobID int64) (*Job, error)
	GetJobByName(name string) (*Job, error)
	DeleteJob(jobID int64) (*Job, error)

	CreateJobHistory(jobHistory *JobHistory) (*JobHistory, error)
	UpdateJobHistory(jobHistory *JobHistory) (*JobHistory, error)
	GetJobHistory(jobID int64) ([]JobHistory, error)
	GetJobHistoryWithLimit(jobID int64, limit int) ([]JobHistory, error)
	GetJobHistoryForSchedule(scheduleID int64) (*JobHistory, error)
	DeleteJobHistory(jobID int64) ([]JobHistory, error)

	CreateSchedule(schedule *Schedule) (*Schedule, error)
	UpdateSchedule(schedule *Schedule) (*Schedule, error)
	GetSchedule(scheduleID int64) (*Schedule, error)
	GetScheduleBetween(scheduleStart, scheduleEnd int64, workerID string) ([]Schedule, error)
	GetSchedulesForJob(jobID int64) ([]Schedule, error)
	GetSchedulesForWorker(workerID string) ([]Schedule, error)
	DeleteScheduleForJob(jobID int64) ([]Schedule, error)
	DeleteScheduleForWorker(workerID string) ([]Schedule, error)
}

// CreateWorker to create a worker
func (c cduleRepository) CreateWorker(worker *Worker) (*Worker, error) {
	if err := c.DB.Create(worker).Error; err != nil {
		return nil, err
	}
	return worker, nil
}

// UpdateWorker to update a worker
func (c cduleRepository) UpdateWorker(worker *Worker) (*Worker, error) {
	if err := c.DB.Updates(worker).Error; err != nil {
		return nil, err
	}
	return worker, nil
}

// GetWorker to get a worker
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

// GetWorkers to get a list of workers
func (c cduleRepository) GetWorkers() ([]Worker, error) {
	var workers []Worker
	if err := c.DB.Find(&workers).Error; err != nil {
		return workers, err
	}
	return workers, nil
}

// GetAliveWorkers to get a list of alive workers
func (c cduleRepository) GetAliveWorkers() ([]Worker, error) {
	var workers []Worker
	// updated_at gt 3 heart means alive
	available := time.Now().Add(-3 * c.Heart)
	if err := c.DB.Where("updated_at > ?", available).Find(&workers).Error; err != nil {
		return workers, err
	}
	return workers, nil
}

// DeleteWorker to delete a worker
func (c cduleRepository) DeleteWorker(workerID string) (*Worker, error) {
	worker, err := c.GetWorker(workerID)
	if err = c.DB.Delete(&worker).Error; err != nil {
		return nil, err
	}
	return worker, nil
}

// CreateJob to create a job
func (c cduleRepository) CreateJob(job *Job) (*Job, error) {
	if err := c.DB.Create(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

// UpdateJob to update a job
func (c cduleRepository) UpdateJob(job *Job) (*Job, error) {
	if err := c.DB.Updates(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

// GetJob to get a job based on ID
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

// GetJobByName to get a job based on Name
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

// DeleteJob to get a job based on ID
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

// CreateJobHistory to create a JobHistory
func (c cduleRepository) CreateJobHistory(jobHistory *JobHistory) (*JobHistory, error) {
	if err := c.DB.Create(jobHistory).Error; err != nil {
		return nil, err
	}
	return jobHistory, nil
}

// UpdateJobHistory to update a JobHistory
func (c cduleRepository) UpdateJobHistory(jobHistory *JobHistory) (*JobHistory, error) {
	if err := c.DB.Updates(jobHistory).Error; err != nil {
		return nil, err
	}
	return jobHistory, nil
}

// GetJobHistory to get a JobHistory by JobID
func (c cduleRepository) GetJobHistory(jobID int64) ([]JobHistory, error) {
	var jobHistories []JobHistory
	if err := c.DB.Where("job_id = ?", jobID).First(&jobHistories).Error; err != nil {
		return nil, err
	}
	return jobHistories, nil
}

// GetJobHistoryWithLimit to get a JobHistory by JobID and limit
func (c cduleRepository) GetJobHistoryWithLimit(jobID int64, limit int) ([]JobHistory, error) {
	var jobHistories []JobHistory
	if err := c.DB.Where("job_id = ?", jobID).Limit(limit).Find(&jobHistories).Error; err != nil {
		return nil, err
	}
	return jobHistories, nil
}

// GetJobHistoryForSchedule to get a JobHistory by scheduleID
func (c cduleRepository) GetJobHistoryForSchedule(scheduleID int64) (*JobHistory, error) {
	var jobHistory JobHistory
	if err := c.DB.Where("execution_id = ?", scheduleID).First(&jobHistory).Error; err != nil {
		return nil, err
	}
	return &jobHistory, nil
}

// DeleteJobHistory to delete a JobHistory by jobID
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

// CreateSchedule to create a schedule
func (c cduleRepository) CreateSchedule(schedule *Schedule) (*Schedule, error) {
	if err := c.DB.Create(schedule).Error; err != nil {
		return nil, err
	}
	return schedule, nil
}

// UpdateSchedule to update a schedule
func (c cduleRepository) UpdateSchedule(schedule *Schedule) (*Schedule, error) {
	if err := c.DB.Updates(schedule).Error; err != nil {
		return nil, err
	}
	return schedule, nil
}

// GetSchedule to get a schedule by executionID
func (c cduleRepository) GetSchedule(executionID int64) (*Schedule, error) {
	var schedule Schedule
	if err := c.DB.Where("execution_id = ?", executionID).Find(&schedule).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

// GetScheduleBetween to get a schedule between scheduleStart and scheduleEnd and by workerID
func (c cduleRepository) GetScheduleBetween(scheduleStart, scheduleEnd int64, workerID string) ([]Schedule, error) {
	var schedules []Schedule
	if err := c.DB.Where("execution_id >= ? and execution_id <= ? and worker_id = ?", scheduleStart, scheduleEnd, workerID).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// GetSchedulesForJob to get a schedules by jobID
func (c cduleRepository) GetSchedulesForJob(jobID int64) ([]Schedule, error) {
	var schedules []Schedule
	if err := c.DB.Where("job_id = ?", jobID).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// GetSchedulesForWorker to get a schedules by workerID
func (c cduleRepository) GetSchedulesForWorker(workerID string) ([]Schedule, error) {
	var schedules []Schedule
	if err := c.DB.Where("worker_id = ?", workerID).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// DeleteScheduleForJob to delete a schedules by jobID
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

// DeleteScheduleForWorker to delete a schedules by workerID
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
