package job

import log "github.com/sirupsen/logrus"

type PanicJob struct {
	Job Job
}

func (j PanicJob) Execute() {
	log.Info("In MyJPanicJobob")
	a := 0
	i := 100 / a
	log.Infof("i: %d", i)
}

func (j PanicJob) JobName() string {
	return "job.PanicJob"
}
