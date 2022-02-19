package cdule

type Job interface {
	Execute(map[string]string)
	JobName() string
	GetJobData() map[string]string
}
