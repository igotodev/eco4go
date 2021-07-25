package handlers

import (
	"strconv"
	"eco4go/dbmod"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type AtSales struct {
	Id                     int     `json:"id"`
	MonthCompletion        float64 `json:"month_completion"`
	MonthCompletionPercent float64 `json:"month_completion_percent"`
	Target                 float64 `json:"target"`
	Month                  string  `json:"month"`
	DaysInMonth            int     `json:"days_in_month"`
	CurrentDate 		   string  `json:"current_date"`
	NeedInDay 			   float64 `json:"need_in_day"`
}

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
	daySales, err := strconv.ParseFloat(c.FormValue("add_sales"), 64)
	if err != nil {
		return c.String(http.StatusInternalServerError,
			fmt.Sprintf("error while prase float from form value: %v", err))
	}
	sales := dbmod.Sales{
		SalesPerson: "Anna",
		Revenue: daySales,
	}
	/*
	err := c.Bind(&sales)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("error while binding addDaySales request: %v", err))
	}
	 */
	err = sales.InsertDailySales()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("error while insert addDaySales into db: %v", err))
	}
	log.Printf("sales for %v is added", now)
	err = c.Redirect(http.StatusSeeOther, "/month")
	if err != nil {
		log.Printf("error while redirect: %v", err)
	}
	return nil
}

func ViewMonthCompletion(c echo.Context) error {
	sales := dbmod.Sales{}
	mc, err := sales.MonthCompletion()
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error while call MonthCompletion: %v", err))
	}
	currentDay, _ := strconv.Atoi(time.Now().Format("02"))
	inPr := mc / float64(10000 / 100)
	needInDay := (float64(10000) - mc) / float64(daysInCurrentMonth() - currentDay)
	atSales :=AtSales{
		DaysInMonth: daysInCurrentMonth(), 
		Target: 10000, 
		MonthCompletionPercent: inPr, 
		MonthCompletion: mc, 
		Month: time.Now().Month().String(), 
		CurrentDate: time.Now().Format("2006/01/02"),
		NeedInDay: needInDay,
	}
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "month", atSales)
}

func daysInCurrentMonth() int {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	daysInMonth, err := strconv.Atoi(lastOfMonth.Format("02"))
	if err != nil {
		log.Fatalf("internal error while call dayInMonth: %v", err)
	}
	return daysInMonth
}
