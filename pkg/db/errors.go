package db

import "errors"

// Ошибки задач
var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrWrongTaskID      = errors.New("wrong task id")
	ErrNoRowsAffected   = errors.New("incorrect id for updating task")
	ErrDeleteNoRows     = errors.New("incorrect id for deleting task")
	ErrUpdateDateNoRows = errors.New("incorrect id for updating date")
)
