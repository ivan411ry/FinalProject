package api

import (
	"FinalProject/pkg/db"
	"net/http"
)

const TasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJSON(w, map[string]string{"error": "method not allowed"})
		return
	}

	tasks, err := db.Tasks(TasksLimit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": "database error"})
		return
	}
	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
