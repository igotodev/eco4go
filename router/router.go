package router

import (
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"os/signal"
	"eco4go/handlers"
	"eco4go/middlewares"
	"eco4go/templates"
)

const address = "127.0.0.1:8080"

func newRoute() *echo.Echo {
	e := echo.New()
	middlewares.SetGeneralMiddlewares(e)
	//month.gohtml
	tempFiles := []string{"templates/gohtml/index.gohtml"}
	templates.SetupTamplateConfig(tempFiles, e)
	e.Static("/static", "static/css")

	e.GET("/month", handlers.ViewMonthCompletion)
	e.GET("/all-json", handlers.AllDataJSON)
	e.POST("/add-sales", handlers.AddDaySales)

	return e
}

func OpenEcho() {
	e := newRoute()
	chWait := make(chan os.Signal, 1)
	signal.Notify(chWait, os.Interrupt)
	go func() {
		e.Logger.Fatal(e.Start(address))
	}()
	<-chWait
	log.Println("echo has been stopped")
}
