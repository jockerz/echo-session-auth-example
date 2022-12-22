package pkg

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	sessionauth "github.com/jockerz/echo-session-auth"
	"github.com/labstack/echo/v4"
)

type (
	CustomContext struct {
		echo.Context
		User interface{}
	}

	User struct {
		Username string
		ID       int
	}
)

var (
	Users = []*User{
		&User{"First", 1},
		&User{"Second", 2},
	}
	Config = sessionauth.MakeConfig(
		[]byte("changeme"),
		"/login",                // UnAuthRedirect
		[]string{"favicon.ico"}, // Excluded path
		[]*regexp.Regexp{},      // Exlucede regex path
	)
	SessionAuth *sessionauth.SessionAuth
)

func GetuserByID(c echo.Context, UserID interface{}) error {
	ctx := c.(*CustomContext)

	var uid int
	uid, err := strconv.Atoi(fmt.Sprintf("%v", UserID))

	if err != nil {
		return err
	}

	for _, u := range Users {
		if u.ID == uid {
			// User is found
			ctx.User = u
			return nil
		}
	}
	return errors.New("user not found")
}

func CreateSessionAuth() {
	sa, err := sessionauth.Create(Config, GetuserByID)
	if err != nil {
		panic(err)
	}
	SessionAuth = sa
}
