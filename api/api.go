package api

import (
	"net/http"
	"time"

	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yasngleer/studentplan/models"
	"github.com/yasngleer/studentplan/store"
	"github.com/yasngleer/studentplan/utils"
	"golang.org/x/crypto/bcrypt"
)

type API struct {
	userStore *store.UserStore
	taskStore *store.TaskStore
}

func NewAPI(userStore *store.UserStore, taskStore *store.TaskStore) *API {

	return &API{
		userStore: userStore,
		taskStore: taskStore,
	}

}

// Handlers for user

func (api *API) RegisterUser(c echo.Context) error {
	rr := new(RegisterRequest)
	err := c.Bind(rr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	passbytes, err := bcrypt.GenerateFromPassword([]byte(rr.Password), 14)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := &models.User{
		Username:     rr.Username,
		Email:        rr.Email,
		PasswordHash: string(passbytes),
	}

	if err := api.userStore.CreateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, RegisterResponse{
		UserID: user.ID,
	})
}

func (api *API) LoginUser(c echo.Context) error {
	r := new(LoginRequest)
	err := c.Bind(r)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user, err := api.userStore.GetUserByUsername(r.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(r.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Password Mismatch")
	}
	token := utils.GenToken(strconv.Itoa(int(user.ID)))
	return c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})
}

func (api *API) UpdateUser(c echo.Context) error {
	r := new(UpdateUserRequest)
	err := c.Bind(r)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())

	}
	if r.Password == "" {
		return c.JSON(http.StatusBadRequest, "Empty password")
	}
	user, err := Getuserfromcontext(c, *api.userStore)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	passhash, _ := bcrypt.GenerateFromPassword([]byte(r.Password), 14)
	user.PasswordHash = string(passhash)
	api.userStore.UpdateUser(user)
	return c.JSON(200, UpdateUserResponse{
		Message: "Success",
	})
}

// Handlers for task

func (api *API) ListTasks(c echo.Context) error {
	date1 := new(time.Time)
	date2 := new(time.Time)

	switch date := c.QueryParam("date"); date {
	case "day":
		start := time.Now()
		*date1 = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Now().Location())
		*date2 = time.Now()
	case "week":
		start := time.Now().AddDate(0, 0, -7)
		*date1 = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Now().Location())
		*date2 = time.Now()
	case "month":
		start := time.Now().AddDate(0, 1, 0)
		*date1 = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Now().Location())
		*date2 = time.Now()
	default:
		date1 = nil
		date2 = nil
	}
	user, err := Getuserfromcontext(c, *api.userStore)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	tasks, err := api.taskStore.GetTasksByUser(user.ID, date1, date2)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(200, tasks)

}

func (api *API) CreateTask(c echo.Context) error {
	task := new(models.Task)
	err := c.Bind(task)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user, err := Getuserfromcontext(c, *api.userStore)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	task.UserID = user.ID
	err = api.taskStore.CreateTask(task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(200, task)
}

func (api *API) GetTask(c echo.Context) error {
	taskid := c.Param("id")
	user, err := Getuserfromcontext(c, *api.userStore)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	itaskid, _ := strconv.Atoi(taskid)
	task, err := api.taskStore.GetTaskByID(uint(itaskid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if task.UserID != user.ID {
		return c.JSON(http.StatusUnauthorized, "Dont have access to task")
	}
	return c.JSON(200, task)
}

func (api *API) UpdateTask(c echo.Context) error {
	r := new(TaskUpdateRequest)
	err := c.Bind(r)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	taskid := c.Param("id")
	itaskid, _ := strconv.Atoi(taskid)
	user, err := Getuserfromcontext(c, *api.userStore)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	task, err := api.taskStore.GetTaskByID(uint(itaskid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if task.UserID != user.ID {
		return c.JSON(http.StatusUnauthorized, "Dont have access to task")
	}

	if r.Title != nil {
		task.Title = *r.Title
	}
	if r.Description != nil {
		task.Description = *r.Description
	}
	if r.Status != nil {
		task.Status = *r.Status
	}
	if r.StartDate != nil {
		task.StartDate = *r.StartDate
	}
	if r.EndDate != nil {
		task.EndDate = *r.EndDate
	}
	api.taskStore.UpdateTask(task)
	return c.JSON(200, task)
}

//todo check if user has access
func (api *API) DeleteTask(c echo.Context) error {
	taskid := c.Param("id")
	itaskid, _ := strconv.Atoi(taskid)
	user, err := Getuserfromcontext(c, *api.userStore)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	task, err := api.taskStore.GetTaskByID(uint(itaskid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if task.UserID != user.ID {
		return c.JSON(http.StatusUnauthorized, "Dont have access to task")
	}
	err = api.taskStore.DeleteTask(uint(itaskid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "")
}
