package api

import (
	"FinalProject/pkg/db"
	"encoding/json"
	"net/http"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTask(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		handleUpdateTask(w, r)
	case http.MethodDelete:
		handleDeleteTask(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJSON(w, map[string]string{"error": "method not allowed"})
	}
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "task id not found"})
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

	writeJSON(w, task)
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "JSON decoding error"})
		return
	}

	if task.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "task id not found"})
		return
	}

	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "task title not found"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		if err == db.ErrNoRowsAffected {
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

	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]interface{}{})
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "id not found"})
		return
	}

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

	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]interface{}{})
}
