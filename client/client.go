package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Sales struct {
	ID                           int     `json:"id,omitempty"`
	Data                         string  `json:"data"`
	Revenue                      float64 `json:"revenue"`
	SalesPerson 		      string  `json:"sales_person"`
}

func main() {
	s := Sales{
		Data: time.Now().Format("2006/01/02"),
		Revenue: 5000,
		SalesPerson: "Anna",
	}
	b, err := json.Marshal(&s)
	if err != nil {
		log.Print(err)
	}
	path := "add-sales"
	byf := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/"+path, byf)
	if err != nil {
		log.Print(err)
	}
	req.Header.Set("Content-Type", "application/json")
	defer req.Body.Close()
	client := &http.Client{Timeout: 10*time.Second}

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	answ, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}
	log.Print(string(answ))

}
