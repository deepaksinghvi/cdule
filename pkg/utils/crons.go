package utils

const (
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
)
