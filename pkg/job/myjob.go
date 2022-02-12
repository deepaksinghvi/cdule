package job

import (
	log "github.com/sirupsen/logrus"
)

type MyJob struct {
	Job Job
}

func (m MyJob) Execute() {
	log.Info("In MyJob")
}

func (m MyJob) JobName() string {
	return "job.MyJob"
}
