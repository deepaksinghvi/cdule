package utils

const (
	// EveryMinute cron expression
	EveryMinute = "0 * * ? * *"
	// EveryEvenMinute cron expression
	EveryEvenMinute = "0 */2 * ? * *"
	// EveryUnEvenMinute cron expression
	EveryUnEvenMinute = "0 1/2 * ? * *"
	// EveryTwoMinutes cron expression
	EveryTwoMinutes = "0 */2 * ? * *"
	// EveryHourAtMin153045 cron expression
	EveryHourAtMin153045 = "0 15,30,45 * ? * *"
	// EveryHour cron expression
	EveryHour = "0 0 * ? * *"
	// EveryEvenHour cron expression
	EveryEvenHour = "0 0 0/2 ? * *"
	// EveryUnEvenHour cron expression
	EveryUnEvenHour = "0 0 1/2 ? * *"
	// EveryThreeHours cron expression
	EveryThreeHours = "0 0 */3 ? * *"
	// EveryTwelveHours cron expression
	EveryTwelveHours = "0 0 */12 ? * *"
	// EveryDayAtMidNight cron expression
	EveryDayAtMidNight = "0 0 0 * * ?"
	// EveryDayAtOneAM cron expression
	EveryDayAtOneAM = "0 0 1 * * ?"
	// EveryDayAtSixAM cron expression
	EveryDayAtSixAM = "0 0 6 * * ?"
	// EverySundayAtNoon cron expression
	EverySundayAtNoon = "0 0 12 ? * "
	// EveryMondayAtNoon cron expression
	EveryMondayAtNoon = "0 0 12 ? *"
	// EveryWeekDayAtNoon cron expression
	EveryWeekDayAtNoon = "0 0 12 ? * MON-FRI"
	// EveryWeekEndAtNoon cron expression
	EveryWeekEndAtNoon = "0 0 12 ? * SUN,SAT"
	// EveryMonthOnFirstAtNoon cron expression
	EveryMonthOnFirstAtNoon = "0 0 12 1 * ?"
	// EveryMonthOnSecondAtNoon cron expression
	EveryMonthOnSecondAtNoon = "0 0 12 2 * ?"
)
