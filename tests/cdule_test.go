package tests

import (
	"os"
	"testing"
	"time"

	"github.com/deepaksinghvi/cdule/pkg/cdule"
	"github.com/deepaksinghvi/cdule/pkg/model"
	"github.com/deepaksinghvi/cdule/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	approxTime = cmp.Options{cmpopts.EquateApproxTime(time.Second)}
)

func Test_BuildNewJob(t *testing.T) {
	err, cdule := createScheduler()
	jobRecordExpected, err := createTestJob()
	require.NoError(t, err)
	jobRecordActual, err := model.CduleRepos.CduleRepository.GetJobByName(jobRecordExpected.JobName)
	require.NoError(t, err)
	if diff := cmp.Diff(jobRecordExpected, jobRecordActual, approxTime); diff != "" {
		t.Fatalf("mismatch (-expectedResult, +actRes):\n%s", diff)
	}

	jobRecordActual.Expired = true
	jobRecordActual, err = model.CduleRepos.CduleRepository.UpdateJob(jobRecordActual)
	require.NoError(t, err)
	cdule.StopWatcher()
}

func Test_BuildNewJobExecution(t *testing.T) {
	err, cdule := createScheduler()

	jobRecordExpected, err := createTestJob()
	require.NoError(t, err)
	time.Sleep(2 * time.Minute)
	schedules, err := model.CduleRepos.CduleRepository.GetSchedulesForJob(jobRecordExpected.ID)
	require.NoError(t, err)
	// Job is expected to run every minute so it should have atleast 2 schedules.
	require.NotEqual(t, 1, len(schedules))
	cdule.StopWatcher()
}

func createScheduler() (error, cdule.Cdule) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	os.Remove(dirname + "/sqlite.db")

	cdule := cdule.Cdule{}
	cdule.NewCdule("./resources", "config_in_memory")
	return err, cdule
}

func createTestJob() (*model.Job, error) {
	myJob := TestJob{}
	jobData := make(map[string]string)
	jobData["one"] = "1"
	jobData["two"] = "2"
	jobData["three"] = "3"
	jobRecordExpected, err := cdule.NewJob(&myJob, jobData).Build(utils.EveryMinute)
	return jobRecordExpected, err
}
