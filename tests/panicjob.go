package tests

import (
	"github.com/deepaksinghvi/cdule/pkg/cdule"
	log "github.com/sirupsen/logrus"
)

/*
PanicJob would be creating a panic because of divide by zero error.
This job is to check that execution of a job happens and Execute(...) does not abort the program because of error raised
in Execute() method call.
*/
type PanicJob struct {
	Job cdule.Job
}

// Execute to execute a job
func (j PanicJob) Execute(jobData map[string]string) {
	log.Info("In MyJPanicJob")
	log.Infof("JobData %v", jobData)
	a := 0
	i := 100 / a
	log.Infof("i: %d", i)
}

// JobName name of the job
func (j PanicJob) JobName() string {
	return "job.PanicJob"
}

// GetJobData job data of the job
func (j PanicJob) GetJobData() map[string]string {
	return nil
}
