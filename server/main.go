package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Usd struct {
	Usdbrl Usdbrl `json:"USDBRL"`
}
type Usdbrl struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Dolar struct {
	ID    int     `gorm:"primaryKey"`
	Code  string  `json:"code"`
	Value float64 `json:"value"`
	gorm.Model
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", handlerCotacao)
	http.ListenAndServe(":8080", mux)
}

func handlerCotacao(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dolar, err := GetDolar()
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte("0"))
		return
	}

	save, err := SaveDolar(dolar)
	status := "(SALVO)"
	if !save {
		status = "(NÃO SALVO)"
	}

	resultado_json, _ := json.Marshal(dolar)
	w.Write([]byte(resultado_json))
	log.Printf("- 1 USD = %v BRL %v", dolar, status)
	if err != nil {
		log.Printf("%v", err)
	}
}

func GetDolar() (float64, error) {
	req, err := http.NewRequest(http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*200))
	defer cancel()
	req = req.WithContext(ctx)
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var u Usd
	err = json.Unmarshal(body, &u)
	if err != nil {
		return 0, err
	}
	cambio, err := strconv.ParseFloat(u.Usdbrl.Bid, 64)
	if err != nil {
		return 0, err
	}
	return cambio, err
}

func SaveDolar(dolar float64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*10))
	defer cancel()

	db, err := gorm.Open(sqlite.Open("dolar.db"), &gorm.Config{})
	if err != nil {
		return false, err
	}
	// defer db.Close()
	db.AutoMigrate(&Dolar{})
	e := db.WithContext(ctx).Create(&Dolar{Code: "BRL", Value: dolar}).Error
	return e == nil, e
}
