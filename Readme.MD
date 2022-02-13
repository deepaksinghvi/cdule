# cdule (pronounce as Schedule)

Golang based scheduler library with database support.
Users could use any database which is supported by [gorm.io](https://gorm.io/).

### To Download the cdule library
```
go get github.com/deepaksinghvi/cdule
```

### Usage Instruction
Step 1 : Configure you database or create the same as mentioned in the below configuration

```
cduletype: DATABASE
dburl: postgres://cduleuser:cdulepassword@localhost:5432/cdule?sslmode=disable
cduleconsistency: AT_MOST_ONCE
```

Step 2 : User would have to implement Job interface

```
package job

type Job interface {
	Execute()
	JobName() string
}

```
Users have to implement Execute() and JobName() method.
Execute is where the main logic for the Job resides and JobName() is implementing to define a JobName.

JobName must be unique and should not be registered for any other job and stored in database.

Example MyJob

```
package job

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

var myJobData map[string]string

type MyJob struct {
	Job Job
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
```
Step 3: Build the job to submit
```
myJob := job.MyJob{}
jobData := make(map[string]string)
jobData["one"] = "1"
jobData["two"] = "2"
jobData["three"] = "3"
jobModel, err := watcher.NewJob(&myJob, jobData).Build(utils.EveryMinute)
```

### Demo Project
This demo describes how cdule library can be used.

https://github.com/deepaksinghvi/cduledemo

### Database Schema
For the sample app, postgresql is used but users can use any db which is supported by gorm.io.
#### DB Tables
* jobs : To store unique jobs.
* job_histories : To store job history with status as result.
* schedules : To store schedule for every next run.
* workers : To store the worker nodes and their health check.


![dbschema.png](resources/dbschema.png)


### Sample Cron
Users can use the pre-defined crons or use their own which are the [standard cron](https://en.wikipedia.org/wiki/Cron) expression
```
	EveryMinute              = "0 * * ? * *"
	EveryEvenMinute          = "0 */2 * ? * *"
	EveryUnEvenMinute        = "0 1/2 * ? * *"
	EveryTwoMinutes          = "0 */2 * ? * *"
	EveryHourAtMin153045     = "0 15,30,45 * ? * *"
	EveryHour                = "0 0 * ? * *"
	EveryEvenHour            = "0 0 0/2 ? * *"
	EveryUnEvenHour          = "0 0 1/2 ? * *"
	EveryThreeHours          = "0 0 */3 ? * *"
	EveryTwelveHours         = "0 0 */12 ? * *"
	EveryDayAtMidNight       = "0 0 0 * * ?"
	EveryDayAtOneAM          = "0 0 1 * * ?"
	EveryDayAtSixAM          = "0 0 6 * * ?"
	EverySundayAtNoon        = "0 0 12 ? * "
	EveryMondayAtNoon        = "0 0 12 ? *"
	EveryWeekDayAtNoon       = "0 0 12 ? * MON-FRI"
	EveryWeekEndAtNoon       = "0 0 12 ? * SUN,SAT"
	EveryMonthOnFirstAtNoon  = "0 0 12 1 * ?"
	EveryMonthOnSecondAtNoon = "0 0 12 2 * ?"
```


### This library is built using 

* Cron parser using the library [robfig/cron](https://github.com/robfig/cron)  
* ORM usign the library [gorm.io](https://gorm.io/)