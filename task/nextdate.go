package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/smir7/scheduler/constans"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("repeat wrong format")
	}

	startDate, err := time.Parse(constans.DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("wrong format date: %v", err)
	}

	split := strings.Split(repeat, " ")

	switch split[0] {
	case "d":
		if len(split) < 2 {
			return "", fmt.Errorf("repeat wrong format")
		}
		nextDays, err := strconv.Atoi(split[1])

		if err != nil {
			return "", err
		}

		if nextDays > 400 {
			return "", fmt.Errorf("repeat wrong format")
		}

		newDate := startDate.AddDate(0, 0, nextDays)

		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, nextDays)
		}
		return newDate.Format(constans.DateFormat), nil

	case "y":
		newDate := startDate.AddDate(1, 0, 0)

		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
		return newDate.Format(constans.DateFormat), nil

	default:
		return "", fmt.Errorf("repeat wrong format")
	}
}
