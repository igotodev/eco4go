package handlers

import (
	"eco4go/dbmod"

	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"

	"github.com/labstack/echo/v4"
)

func AllDataJSON(c echo.Context) error {
	sales := dbmod.Sales{}
	var ms []dbmod.Sales
	jSales, err := sales.GetAllDataJSON()
	if err != nil {
		return err
	}
	json.Unmarshal(jSales, &ms)
	return c.JSON(http.StatusOK, ms)
}

func AddDaySales(c echo.Context) error {
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

func ViewMonthComplection(c echo.Context) error {
	sales := dbmod.Sales{}
	mc, err := sales.MonthCompletion()
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "month", mc)
}
