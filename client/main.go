package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	dolar, err := GetDolar()
	if err != nil {
		log.Printf("%v", err)
		return
	}
	err = SaveFile(dolar)
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("- 1 USD = %v BRL", dolar)
}

func GetDolar() (string, error) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		return "0", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*300))
	defer cancel()
	req = req.WithContext(ctx)
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return "0", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "0", err
	}
	return string(body), err
}

func SaveFile(dolar string) error {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	text := "Dólar:" + dolar
	_, err = f.Write([]byte(text)) // OUTROS DADOS
	if err != nil {
		panic(err)
	}
	return nil
}
