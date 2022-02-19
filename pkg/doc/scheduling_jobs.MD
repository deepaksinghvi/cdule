In order to schedule jobs with cdule, user needs to
1. Configure persistence
2. Implement cdule.Job Interface &
3. Schedule job with required cron expression.

Job will be persisted in the jobs table.\
Next execution would be persisted in schedules tables.\
Job history would be persisted and maintained in job_histories table.


## Configuration

User needs to create a resources/config.yml in their project home directory with the followling keys

* cduletype
* dburl
* cduleconsistency

cduletype is used to specify whether it is an In-Memory or Database based configuration. Possible values are DATABASE and MEMORY.
dburl is the database connection url.
cduleconsistency is for reserved for future usage.

### config.yml for postgressql based configuration
```
cduletype: DATABASE
dburl: postgres://cduleuser:cdulepassword@localhost:5432/cdule?sslmode=disable
cduleconsistency: AT_MOST_ONCE
```

### config.yml for sqlite based in-memory configuration
```
cduletype: MEMORY
dburl: /Users/dsinghvi/sqlite.db
cduleconsistency: AT_MOST_ONCE
```


## Job Interface Implementation

```
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
```

## Schedule a Job
It is expected that testJob will be Executed five times, once for every minute and program will exit. TestJob jobData map holds the data in the format of map[string]string where gets stored for every execution and gets updated as the next counter value on Execute() method call.

```
cdule := cdule.Cdule{}
cdule.NewCdule()
testJob := TestJob{}
jobData := make(map[string]string)
jobData["one"] = "1"
jobData["two"] = "2"
jobData["three"] = "3"
cdule.NewJob(&testJob, jobData).Build(utils.EveryMinute)

time.Sleep(5 * time.Minute)
cdule.StopWatcher()
```
