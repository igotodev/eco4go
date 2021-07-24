package middlewares

import (
	"eco4go/dbmod"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
)

func addServerToHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "eco4go/0.1")
		return next(c)
	}
}

func AuthValid(name, pass string, c echo.Context) (bool, error) {
	ok, err := dbmod.CheckUser(name, pass)
	if err != nil {
		log.Print("something went wrong: ", err)
		return false, err
	}
	if ok {
		return true, nil
	}
	log.Printf("something went wrong, ok is %v", ok)
	return false, nil
}

func SetGeneralMiddlewares(e *echo.Echo) {
	e.Use(addServerToHeader)
	e.Use(middleware.Recover())
	e.Use(middleware.Timeout())
	//e.Use(middleware.BasicAuth(AuthValid))
	//e.Use(CheckCookie)

	c := color.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: "${time_custom} from " + c.Blue("${remote_ip}") +
			" to " + c.Green("${host}${path}") +
			" with code " + c.BlueBg("${status}") + "\n",
		CustomTimeFormat: "2006/01/02 15:04:05",
	}))
}
