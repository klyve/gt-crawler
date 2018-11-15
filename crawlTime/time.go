package crawlTime

import (
	"strconv"
	"strings"
	"time"
)

// GetDate returns a Time object. Use daysToAdd to either add or subtract days from Time.
// For subtraction, use negative numbers
func GetDate(daysToAdd int) (date time.Time) {
	today := time.Now()
	return today.Add(time.Hour * 24 * time.Duration(daysToAdd))
}

func FindPreviousDate(date string) (previousDate string){
	parts := strings.Split(date, "-")

	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	day, _ := strconv.Atoi(parts[2])


	dateParsed := time.Date(year, time.Month(month), day, 0,0, 0, 0, time.Local)
	dateParsed = dateParsed.Add(-(time.Hour * 24) * 1)

	previousDate = strconv.Itoa(dateParsed.Year()) + "-" + strconv.Itoa(int(dateParsed.Month())) + "-" + strconv.Itoa(dateParsed.Day())
	return
}

func GetDateString(addedDays int) (date string) {
	tObj := GetDate(addedDays)
	date = strconv.Itoa(tObj.Year()) + "-" + strconv.Itoa(int(tObj.Month())) + "-" + strconv.Itoa(tObj.Day())
	return
}