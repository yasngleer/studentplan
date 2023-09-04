package api

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yasngleer/studentplan/models"
	"github.com/yasngleer/studentplan/store"
)

func Getuserfromcontext(c echo.Context, store store.UserStore) (*models.User, error) {
	uid := c.Get("user_id").(string)
	iuid, _ := strconv.Atoi(uid)
	user, err := store.GetUserByID(uint(iuid))
	return user, err
}
