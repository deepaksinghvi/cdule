package main

import (
	"github.com/deepaksinghvi/cdule/pkg/utils"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"

	"github.com/deepaksinghvi/cdule/pkg/cdule"
)

/*
TODO This TestJob schedule in main program is for the development debugging.
*/
func main() {
	c := cdule.Cdule{}
	c.NewCdule()

	myJob := TestJob{}
	jobData := make(map[string]string)
	jobData["one"] = "1"
	jobData["two"] = "2"
	jobData["three"] = "3"
	jobRecordExpected, err := cdule.NewJob(&myJob, jobData).Build(utils.EveryMinute)
	if nil != err {
		log.Error(err)
	}
	log.Infof("%v", jobRecordExpected)
	time.Sleep(5 * time.Minute)
	c.StopWatcher()
}

var testJobData map[string]string

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
