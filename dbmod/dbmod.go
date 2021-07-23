package dbmod

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

type Sales struct {
	ID                           int     `json:"id,omitempty"`
	Data                         string  `json:"data"`
	Revenue                      float64 `json:"revenue"`
	SalesPerson 				 string  `json:"sales_person"`
}

type Auth struct {
	Login string
	Password string
	Cookie string
	CreationUserTime string
	Invite bool
}

var myDB *sql.DB

const dsn = "user:password@tcp(localhost:3306)/database" // example

var (
	mu sync.Mutex
	wg sync.WaitGroup
)


func (sales *Sales) InsertDailySales() error {
	//mu.Lock()
	//defer mu.Unlock()
	dataNow := time.Now().Format("2006-01-02")

	_, err := myDB.Exec("INSERT INTO `monthly_sales` (`data`, `revenue`, `sales_person`) " +
		"VALUES (?, ?, ?);", dataNow, sales.Revenue, sales.SalesPerson)
	if err != nil {
		return err
	}
	return nil
}

func (sales *Sales) SelectBetween(fromData string, toData string) ([]Sales, error) {
	//mu.Lock()
	//defer mu.Unlock()
	result, err := myDB.Query("SELECT * FROM `monthly_sales` WHERE `data` BETWEEN ? AND ?;", fromData, toData)
	if err != nil {
		return []Sales{}, err
	}

	allSales := []Sales{}
	for result.Next() {
		result.Scan(&sales.ID, &sales.Data, &sales.Revenue, &sales.SalesPerson)
		allSales = append(allSales, *sales)
	}
	result.Close()
	return allSales, nil
}

func (sales *Sales) SelectThisMonth() ([]Sales, error) {
	//mu.Lock()
	//defer mu.Unlock()
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	startThisMonth := year + "-" + month + "-" + "01"
	dateNow := time.Now().Format("2006-01-02")
	result, err := myDB.Query("SELECT * FROM `monthly_sales` WHERE `data` BETWEEN ? AND ?;", startThisMonth, dateNow)
	if err != nil {
		return []Sales{}, err
	}
	allSales := []Sales{}
	for result.Next() {
		result.Scan(&sales.ID, &sales.Data, &sales.Revenue, &sales.SalesPerson)
		allSales = append(allSales, *sales)
	}
	result.Close()
	return allSales, nil
}

func (sales *Sales) GetAllDataJSON() ([]byte, error) {
	//mu.Lock()
	//defer mu.Unlock()
	//dateNow := time.Now().Format("2006-01-02")
	result, err := myDB.Query("SELECT * FROM `monthly_sales`;")
	if err != nil {
		return nil, err
	}

	allSales := []Sales{}
	for result.Next() {
		result.Scan(&sales.ID, &sales.Data, &sales.Revenue, &sales.SalesPerson)
		allSales = append(allSales, *sales)
	}
	result.Close()

	jd, err := json.Marshal(allSales)
	return jd, nil
}

func (sales *Sales) MonthCompletion() (monthCompl float64, err error) {
	s, err := sales.SelectThisMonth()
	if err != nil {
		return 0, err
	}
	var monthRevenue float64
	for _, v := range s {
		monthRevenue += v.Revenue
	}
	return monthRevenue, nil
}

func CheckUser(login string, pass string) (ok bool, err error) {
	res, err := myDB.Query("SELECT `invite` FROM `auth4` WHERE `login`=? && `password`=?;", login, pass)
	if err != nil {
		return false, fmt.Errorf("error during database query execution: %v", err)
	}
	defer res.Close()

	au := Auth{}
	for res.Next() {
		res.Scan(&au.Invite)
	}
	return au.Invite, nil
}

func CheckCookie(login string, pass string) (cookie string, err error) {
	res, err := myDB.Query("SELECT `cookie` FROM `auth4` WHERE `login`=? && `password`=?;", login, pass)
	if err != nil {
		return "", fmt.Errorf("error during database query execution: %v", err)
	}
	defer res.Close()

	au := Auth{}
	for res.Next() {
		res.Scan(&au.Cookie)
	}
	return au.Cookie, nil
}

func UpdateCookie(login string, pass string, c string) (bool, error) {
	_, err := myDB.Exec("UPDATE `auth4` SET `cookie`=? WHERE `login`=? && `password`=?;", c, login, pass)
	if err != nil {
		return false, fmt.Errorf("error during database query execution: %v", err)
	}
	return true, nil
}

func OpenDB() {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	myDB = db
}

func CloseDB() {
	myDB.Close()
}

