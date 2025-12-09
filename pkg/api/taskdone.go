package api

import (
	"FinalProject/pkg/db"
	"net/http"
	"time"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJSON(w, map[string]string{"error": "method not allowed"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "id not found"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err == db.ErrTaskNotFound {
			w.WriteHeader(http.StatusNotFound)
			writeJSON(w, map[string]string{"error": "task not found"})
			return
		}

		if err == db.ErrWrongTaskID {
			w.WriteHeader(http.StatusNotFound)
			writeJSON(w, map[string]string{"error": "task not found"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": "database error"})
		return
	}

	now := time.Now()

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			if err == db.ErrDeleteNoRows {
				w.WriteHeader(http.StatusNotFound)
				writeJSON(w, map[string]string{"error": "task not found"})
				return
			}

			if err == db.ErrWrongTaskID {
				w.WriteHeader(http.StatusNotFound)
				writeJSON(w, map[string]string{"error": "task not found"})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, map[string]string{"error": "database error"})
			return
		}
	} else {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}

		err = db.UpdateDate(nextDate, id)
		if err != nil {
			if err == db.ErrUpdateDateNoRows {
				w.WriteHeader(http.StatusNotFound)
				writeJSON(w, map[string]string{"error": "task not found"})
				return
			}

			if err == db.ErrWrongTaskID {
				w.WriteHeader(http.StatusNotFound)
				writeJSON(w, map[string]string{"error": "task not found"})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, map[string]string{"error": "database error"})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]interface{}{})
}
