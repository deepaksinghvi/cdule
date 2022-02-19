package tests

import (
	"github.com/deepaksinghvi/cdule/pkg/cdule"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var testJobData map[string]string

/*
TestJob jobData map holds the data in the format of
	jobData := make(map[string]string)
	jobData["one"] = "1"
	jobData["two"] = "2"
	jobData["three"] = "3"
jobData gets stored for every execution and gets updated as the next counter value on Execute() method call.
*/

type TestJob struct {
	Job cdule.Job
}

func (m TestJob) Execute(jobData map[string]string) {
	log.Info("In TestJob")
	for k, v := range jobData {
		valNum, err := strconv.Atoi(v)
		if nil == err {
			jobData[k] = strconv.Itoa(valNum + 1)
		} else {
			log.Error(err)
		}

	}
	testJobData = jobData
}

func (m TestJob) JobName() string {
	return "job.TestJob"
}

func (m TestJob) GetJobData() map[string]string {
	return testJobData
}
