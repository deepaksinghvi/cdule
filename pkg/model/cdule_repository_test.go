package model

import (
	"encoding/json"
	l "log"
	"os"
	"testing"
	"time"

	"github.com/deepaksinghvi/cdule/pkg"
	"github.com/deepaksinghvi/cdule/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	approxTime = cmp.Options{cmpopts.EquateApproxTime(time.Second)}
)

func TestRepository_Job(t *testing.T) {
	err := DBConn()
	require.NoError(t, err)
	testJob, err := createTestJob()
	require.NoError(t, err)

	expectedResult, err := CduleRepos.CduleRepository.CreateJob(testJob)

	actualResult, err := CduleRepos.CduleRepository.GetJob(expectedResult.ID)

	if diff := cmp.Diff(expectedResult, actualResult, approxTime); diff != "" {
		t.Fatalf("mismatch (-expectedResult, +actRes):\n%s", diff)
	}
	expectedResult.Expired = true
	_, err = CduleRepos.CduleRepository.UpdateJob(expectedResult)

	actualResult, err = CduleRepos.CduleRepository.GetJobByName("job.RepoTestJob")

	require.Equal(t, expectedResult.Expired, actualResult.Expired)

	actualResult, err = CduleRepos.CduleRepository.DeleteJob(expectedResult.ID)

	require.Equal(t, expectedResult.JobName, actualResult.JobName)
}

func TestRepository_JobHistory(t *testing.T) {
	err := DBConn()
	require.NoError(t, err)
	testJobHistory, err := createTestJobHistory()
	require.NoError(t, err)

	expectedResult, err := CduleRepos.CduleRepository.CreateJobHistory(testJobHistory)

	actualResultJobHistoryArray, err := CduleRepos.CduleRepository.GetJobHistory(expectedResult.JobID)

	require.Equal(t, expectedResult.Status, actualResultJobHistoryArray[0].Status)
	require.Equal(t, expectedResult.JobID, actualResultJobHistoryArray[0].JobID)
	require.Equal(t, expectedResult.ExecutionID, actualResultJobHistoryArray[0].ExecutionID)

	expectedResult.Status = JobStatusInProgress
	_, err = CduleRepos.CduleRepository.UpdateJobHistory(expectedResult)

	actualResult, err := CduleRepos.CduleRepository.GetJobHistoryForSchedule(testJobHistory.ExecutionID)

	require.Equal(t, expectedResult.Status, actualResult.Status)

	actualResultJobHistoryArray, err = CduleRepos.CduleRepository.DeleteJobHistory(expectedResult.JobID)

	require.Equal(t, expectedResult.ExecutionID, actualResultJobHistoryArray[0].ExecutionID)
}

func TestRepository_Schedule(t *testing.T) {
	err := DBConn()
	require.NoError(t, err)
	schedule, err := createTestSchedule()
	require.NoError(t, err)

	expectedResult, err := CduleRepos.CduleRepository.CreateSchedule(schedule)
	actualResult, err := CduleRepos.CduleRepository.GetSchedule(expectedResult.ExecutionID)
	if diff := cmp.Diff(expectedResult, actualResult, approxTime); diff != "" {
		t.Fatalf("mismatch (-expectedResult, +actRes):\n%s", diff)
	}

	data := make(map[string]string)
	data["a"] = "xyz"
	jobDataMapStr, err := mapToString(data)
	expectedResult.JobData = jobDataMapStr

	_, err = CduleRepos.CduleRepository.UpdateSchedule(expectedResult)
	actualResultScheduleArray, err := CduleRepos.CduleRepository.GetSchedulesForJob(schedule.JobID)
	require.Equal(t, expectedResult.JobData, actualResultScheduleArray[0].JobData)

	actualResultScheduleArray, err = CduleRepos.CduleRepository.DeleteScheduleForJob(schedule.JobID)
	require.Equal(t, expectedResult.ExecutionID, actualResultScheduleArray[0].ExecutionID)
	schedule.JobID = 3
	expectedResult, err = CduleRepos.CduleRepository.CreateSchedule(schedule)
	actualResultScheduleArray, err = CduleRepos.CduleRepository.DeleteScheduleForWorker("dsinghvi-host")
	require.Equal(t, expectedResult.ExecutionID, actualResultScheduleArray[0].ExecutionID)
}

func TestRepository_Worker(t *testing.T) {
	err := DBConn()
	require.NoError(t, err)
	testWorker, err := createTestWorker()
	require.NoError(t, err)

	expectedResult, err := CduleRepos.CduleRepository.CreateWorker(testWorker)

	actualResult, err := CduleRepos.CduleRepository.GetWorker(expectedResult.WorkerID)

	if diff := cmp.Diff(expectedResult, actualResult, approxTime); diff != "" {
		t.Fatalf("mismatch (-expectedResult, +actRes):\n%s", diff)
	}
	expectedResult.UpdatedAt = time.Now()
	_, err = CduleRepos.CduleRepository.UpdateWorker(expectedResult)

	actualResult, err = CduleRepos.CduleRepository.GetWorker(testWorker.WorkerID)

	require.Equal(t, true, expectedResult.UpdatedAt.Equal(actualResult.UpdatedAt))

	actualResult, err = CduleRepos.CduleRepository.DeleteWorker(expectedResult.WorkerID)

	require.Equal(t, expectedResult.WorkerID, actualResult.WorkerID)
}

func createTestWorker() (*Worker, error) {
	return &Worker{
		WorkerID:  "dsinghvi-host",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}, nil
}

func createTestJob() (*Job, error) {
	jobDataStr, err := getJobDatMapAsString()
	if err != nil {
		return nil, err
	}
	return &Job{
		Model:          Model{},
		JobName:        "job.RepoTestJob",
		GroupName:      "",
		CronExpression: utils.EveryWeekDayAtNoon,
		Expired:        false,
		JobData:        jobDataStr,
	}, nil
}

func getJobDatMapAsString() (string, error) {
	data := make(map[string]string)
	data["a"] = "abc"
	return mapToString(data)
}

func mapToString(data map[string]string) (string, error) {
	var jobDataStr = ""
	jobDataBytes, err := json.Marshal(data)
	if nil != err {
		return jobDataStr, err
	}
	if string(jobDataBytes) != pkg.EMPTYSTRING {
		jobDataStr = string(jobDataBytes)
	}
	return jobDataStr, nil
}

func createTestJobHistory() (*JobHistory, error) {
	return &JobHistory{
		Model:       Model{},
		JobID:       2,
		ExecutionID: 34534543534,
		DeletedAt:   gorm.DeletedAt{},
		Status:      "NEW",
		WorkerID:    "dsinghvi-host",
		RetryCount:  0,
	}, nil
}

func createTestSchedule() (*Schedule, error) {
	jobDataMapStr, err := getJobDatMapAsString()
	if err != nil {
		return nil, err
	}
	schedule := &Schedule{
		ExecutionID: 34534543534,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   gorm.DeletedAt{},
		WorkerID:    "dsinghvi-host",
		JobID:       2,
		JobData:     jobDataMapStr,
	}
	return schedule, nil
}
func DBConn() error {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	sqlLogger := logger.New(
		l.New(os.Stdout, "\r\n", l.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		},
	)

	db.Logger = sqlLogger
	MigrateTestTables(db)
	CduleRepos = &Repositories{
		CduleRepository: NewCduleRepository(db),
		DB:              db,
	}
	return err
}

func MigrateTestTables(db *gorm.DB) {
	db.AutoMigrate(&Job{})
	db.AutoMigrate(&JobHistory{})
	db.AutoMigrate(&Schedule{})
	db.AutoMigrate(&Worker{})
}
