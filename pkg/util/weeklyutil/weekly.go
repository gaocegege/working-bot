package weeklyutil

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/go-github/github"
)

func GetWeeklyNumber(issue github.Issue) (int, error) {
	title := *issue.Title
	num, err := strconv.Atoi(title[7:len(title)])
	return num, err
}

func GenerateTimeFromNumber(number int) time.Time {
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	// Weekly-1 is 2019.04.08
	dateFor1 := time.Date(2019, time.April, 8, 12, 0, 0, 0, beijing)

	weeksToAdd := number
	dateFor1.AddDate(0, 0, 7*weeksToAdd)
	return dateFor1.AddDate(0, 0, 7*weeksToAdd)
}

func GetFileNameFromTime(date time.Time) string {
	return fmt.Sprintf("%d/%d-%.2d-%.2d-weekly.md", date.Year(), date.Year(), int(date.Month()), date.Day())
}

func GetFileNameFromIssue(issue github.Issue) (string, error) {
	num, err := GetWeeklyNumber(issue)
	if err != nil {
		return "", err
	}
	return GetFileNameFromTime(GenerateTimeFromNumber(num)), nil
}
