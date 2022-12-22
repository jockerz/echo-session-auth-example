package pkg

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type LoginParam struct {
	Username string `form:"username" json:"username" query:"username"`
	Remember bool   `form:"remember" json:"remember" query:"remember"`
}

var TplLogin = `<html>
<head><title>Login</title></head>
<body>
	<p>This is login page</p>
	<p>Error: %v</p>
	<form method="post">
		<label>Username</label><input type="text" name="username"/><br/>
		<label>Remember</label><input type="checkbox" name="remember" value="1"/><br/>
		<button type="submit">Login</button>
	</form>
</body>
</html>`
var TplProtected = `<html>
<head><title>Protected</title></head>
<body>
	<p>This is protected page</p>
	<ul>
		<li><a href="/login">Login</a></li>
		<li><a href="/logout">Logout</a></li>
	</ul>
</body>
</html>`
var TplFreshlyProtected = `<html>
<head><title>Freshly Protected</title></head>
<body>
	<p>This is freshly protected page</p>
	<ul>
		<li><a href="/login">Login</a></li>
		<li><a href="/logout">Logout</a></li>
	</ul>
</body>
</html>`

func Login(c echo.Context) error {
	ctx := c.(*CustomContext)
	var args LoginParam
	errMsg := "nil"

	if err := ctx.Bind(&args); err != nil {
		return err
	}

	if ctx.Request().Method == http.MethodPost {
		for _, u := range Users {
			if strings.EqualFold(u.Username, args.Username) {
				// User is found
				strID := fmt.Sprintf("%v", u.ID)
				// Authenticate User ID
				SessionAuth.Login(ctx, strID, true, args.Remember)

				// Redirect to protected page
				return ctx.Redirect(http.StatusFound, "/")
			}
		}
		errMsg = "Invalid user"

	}
	return ctx.HTML(http.StatusOK, fmt.Sprintf(TplLogin, errMsg))
}

func Logout(c echo.Context) error {
	ctx := c.(*CustomContext)
	SessionAuth.Logout(ctx)
	return ctx.Redirect(http.StatusFound, "/login")
}

func ProtectedPage(c echo.Context) error {
	ctx := c.(*CustomContext)
	SessionAuth.LoginRequired(ctx)
	return ctx.HTML(http.StatusOK, TplProtected)
}

func FreshOnlyProtectedPage(c echo.Context) error {
	ctx := c.(*CustomContext)
	SessionAuth.FreshLoginRequired(ctx)
	return ctx.HTML(http.StatusOK, TplFreshlyProtected)
}

func RegisterRoutes(app *echo.Echo) {
	app.GET("/login", Login)
	app.POST("/login", Login)
	app.GET("/logout", Logout)
	app.GET("/", ProtectedPage)
	app.GET("/fresh", FreshOnlyProtectedPage)
}
