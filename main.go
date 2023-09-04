package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/yasngleer/studentplan/api"
	"github.com/yasngleer/studentplan/middleware"
	"github.com/yasngleer/studentplan/store"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:secret@tcp(db:3306)/stu?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	userStore := store.NewUserStore(db)
	taskStore := store.NewTaskStore(db)
	e := echo.New()

	handler := api.NewAPI(userStore, taskStore)
	e.POST("/register", handler.RegisterUser)
	e.POST("/login", handler.LoginUser)
	a := e.Group("")
	a.Use(middleware.JWTMiddleware())
	a.PUT("/user", handler.UpdateUser)
	a.GET("/tasks", handler.ListTasks)
	a.POST("/tasks", handler.CreateTask)

	a.GET("/tasks/:id", handler.GetTask)
	a.PUT("/tasks/:id", handler.UpdateTask)

	a.DELETE("/tasks/:id", handler.DeleteTask)
	e.Start(":8080")

}
