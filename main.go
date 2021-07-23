// eco4go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
	"log"
	"net/http"
	"os"
	"os/signal"
	"eco4go/dbmod"
)

const address = "127.0.0.1:8080"

func allDataHandler(c echo.Context) error {
	sales := dbmod.Sales{}
	var ms []dbmod.Sales
	jSales, err := sales.GetAllDataJSON()
	if err != nil {
		return err
	}
	json.Unmarshal(jSales, &ms)
	return c.JSON(http.StatusOK, ms)
}

func addDaySales(c echo.Context) error {
	now := time.Now().Format("2006/01/02")
	sales := dbmod.Sales{}
	err := c.Bind(&sales)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("error while binding addDaySales request: %v", err))
	}
	err = sales.InsertDailySales()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("error while insert addDaySales into db: %v", err))
	}
	log.Printf("sales for %v is added", now)
	return c.String(http.StatusCreated, "we got your sales!")
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

// simple middleware
func AddServerToHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "eco4go/0.1")
		return next(c)
	}
}

func OpenEcho() {
	e := echo.New()

	e.Use(AddServerToHeader)
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

	e.GET("/", allDataHandler)
	e.POST("/daysales", addDaySales)

	chWait := make(chan os.Signal, 1)
	signal.Notify(chWait, os.Interrupt)
	go func() {
		e.Logger.Fatal(e.Start(address))
	}()
	<-chWait
	log.Println("echo has been stopped")
}

func main() {
	dbmod.OpenDB()
	defer dbmod.CloseDB()
	OpenEcho()
}
