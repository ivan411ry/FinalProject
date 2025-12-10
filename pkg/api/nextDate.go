package api

import (
	"errors"
	"log"
	"net/http"
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
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if dstart == "" {
		http.Error(w, "date parameter is required", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		http.Error(w, "repeat parameter is required", http.StatusBadRequest)
		return
	}

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		parsedNow, err := time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid now date format", http.StatusBadRequest)
			return
		}
		now = parsedNow
	}
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	result, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(result)); err != nil {
		log.Printf("NextDateHandler: write error: %v", err)
	}
}
