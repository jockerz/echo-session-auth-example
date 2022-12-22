package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"

	"webapp/pkg"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	app := echo.New()

	// Create sesion auth
	pkg.CreateSessionAuth()

	// Required: Global Middlewares
	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &pkg.CustomContext{
				Context: c,
			}
			return next(cc)
		}
	})

	// Session middleware
	app.Use(pkg.SessionAuth.GetSessionMiddleware())

	// Session auth middleware
	app.Use(pkg.SessionAuth.AuthMiddlewareFunc)

	// Routes
	pkg.RegisterRoutes(app)

	if err := app.Start(":9090"); err != nil {
		panic(err)
	}
}
