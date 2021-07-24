package router

import (
	"eco4go/handlers"
	"eco4go/middlewares"
	"eco4go/templates"
	
	"log"
	"os"
	"os/signal"
	
	"github.com/labstack/echo/v4"
)

const address = "127.0.0.1:8080"

func newRoute() *echo.Echo {
	e := echo.New()
	middlewares.SetGeneralMiddlewares(e)

	tempFiles := []string{"templates/gohtml/month.gohtml"}
	templates.SetupTamplateConfig(tempFiles, e)

	e.GET("/month", handlers.ViewMonthComplection)
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
