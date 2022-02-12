package job

type Job interface {
	Execute()
	JobName() string
}
