package tests

import (
	"github.com/deepaksinghvi/cdule/pkg/job"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var myJobData map[string]string

type MyJob struct {
	Job job.Job
}

func (m MyJob) Execute(jobData map[string]string) {
	log.Info("In MyJob")
	for k, v := range jobData {
		valNum, err := strconv.Atoi(v)
		if nil == err {
			jobData[k] = strconv.Itoa(valNum + 1)
		} else {
			log.Error(err)
		}

	}
	myJobData = jobData
}

func (m MyJob) JobName() string {
	return "job.MyJob"
}

func (m MyJob) GetJobData() map[string]string {
	return myJobData
}
