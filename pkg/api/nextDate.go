package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return dateOnly.After(nowOnly)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule cannot be empty")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("invalid date format")
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", errors.New("invalid repeat format")
	}

	switch parts[0] {
	case "d":
		if len(parts) < 2 {
			return "", errors.New("days count required for 'd' rule")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("days must be a number")
		}
		if days <= 0 || days > 400 {
			return "", errors.New("days must be between 1 and 400")
		}

		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	default:
		return "", errors.New("unsupported repeat format")
	}
}
