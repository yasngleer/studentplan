package api

import (
	"time"

	"github.com/yasngleer/studentplan/models"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" `
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserID uint `json:"user_id"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UpdateUserRequest struct {
	Password string `json:"password"`
}

type UpdateUserResponse struct {
	Message string `json:"message"`
}

type CreateTaskResponse struct {
	TaskID uint `json:"task_id"`
}

type TaskUpdateRequest struct {
	Title       *string            `json:"title"`
	Description *string            `json:"description"`
	StartDate   *time.Time         `json:"start_date"`
	EndDate     *time.Time         `json:"end_date"`
	Status      *models.TaskStatus `json:"status"`
}
