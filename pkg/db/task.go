package db

import (
	"database/sql"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err

}

func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	rows, err := db.Query(query, limit)
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var id int

		task := &Task{}
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []*Task{}, err
		}
		task.ID = strconv.Itoa(id)
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return []*Task{}, err
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return nil, ErrWrongTaskID
	}
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := db.QueryRow(query, taskID)
	task := &Task{}
	var dbID int
	err = row.Scan(&dbID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	task.ID = strconv.Itoa(dbID)
	return task, nil
}

func UpdateTask(task *Task) error {
	taskID, err := strconv.Atoi(task.ID)
	if err != nil {
		return ErrWrongTaskID
	}
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, taskID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrUpdateDateNoRows
	}
	return nil
}
func DeleteTask(id string) error {
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return ErrWrongTaskID
	}
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, taskID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrDeleteNoRows
	}
	return nil
}
func UpdateDate(next string, id string) error {
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return ErrWrongTaskID
	}
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, next, taskID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNoRowsAffected
	}
	return nil
}
