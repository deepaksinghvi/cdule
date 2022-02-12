package utils

import (
	"fmt"
	"time"
)

func GetFormattedTimeFromEpochMillis(z int64, zone string, style int) string {
	var x string
	secondsSinceEpoch := z / 1000
	unixTime := time.Unix(secondsSinceEpoch, 0)
	timeZoneLocation, err := time.LoadLocation(zone)
	if err != nil {
		fmt.Println("Error loading timezone:", err)
	}

	timeInZone := unixTime.In(timeZoneLocation)

	switch style {
	case 1:
		timeInZoneStyleOne := timeInZone.Format("Mon Jan 2 15:04:05")
		//Mon Aug 14 13:36:02
		return timeInZoneStyleOne
	case 2:
		timeInZoneStyleTwo := timeInZone.Format("02-01-2006 15:04:05")
		//14-08-2017 13:36:02
		return timeInZoneStyleTwo
	case 3:
		timeInZoneStyleThree := timeInZone.Format("2006-02-01 15:04:05")
		//2017-14-08 13:36:02
		return timeInZoneStyleThree
	}
	return x
}
