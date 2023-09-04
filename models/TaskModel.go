package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus int

const (
	TaskStatusUndefined TaskStatus = iota
	TaskStatusPending
	TaskStatusCompleted
)

type Task struct {
	gorm.Model
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Status      TaskStatus `json:"status"`
	UserID      uint       `json:"-"`
}
