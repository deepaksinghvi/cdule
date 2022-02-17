package tests

import (
	"github.com/deepaksinghvi/cdule/pkg/job"
	log "github.com/sirupsen/logrus"
)

type PanicJob struct {
	Job job.Job
}

func (j PanicJob) Execute(jobData map[string]string) {
	log.Info("In MyJPanicJob")
	log.Infof("JobData %v", jobData)
	a := 0
	i := 100 / a
	log.Infof("i: %d", i)
}

func (j PanicJob) JobName() string {
	return "job.PanicJob"
}

func (j PanicJob) GetJobData() map[string]string {
	return nil
}
