package api

import (
	"net/http"
	"time"
)

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
	w.Write([]byte(result))
}
