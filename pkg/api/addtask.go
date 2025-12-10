package api

import (
	"FinalProject/pkg/db"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("writeJSON error: %v", err)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "JSON decoding error"})
		return
	}

	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "no task title"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": "database error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	idStr := strconv.FormatInt(id, 10)
	writeJSON(w, map[string]string{"id": idStr})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return errors.New("invalid data format")
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

		if afterNow(now, t) {
			task.Date = next
		}
	} else {
		if afterNow(now, t) {
			task.Date = now.Format(DateFormat)
		}
	}

	return nil
}
