package main

import (
	"context"
	"database/sql"
	"encoding/json"
	util "goexpert/error"
	"io"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Usdbrl struct {
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
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", GetCotacao)
	http.ListenAndServe(":8080", nil)
}

func GetCotacao(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "cotacao" {
		w.WriteHeader(http.StatusNotFound)
	}

	ctxGet, cancelGet := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancelGet()

	request, err := http.NewRequestWithContext(
		ctxGet,
		"GET",
		"https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil,
	)

	util.ErrorHandler(err)

	response, err := http.DefaultClient.Do(request)
	util.ErrorHandler(err)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	util.ErrorHandler(err)

	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	util.ErrorHandler(err)

	err = PersistCotacao(cotacao)
	util.ErrorHandler(err)

	value, err := strconv.ParseFloat(cotacao.Usdbrl.Bid, 32)
	util.ErrorHandler(err)

	json.NewEncoder(w).Encode(value)
}

func PersistCotacao(cotacao Cotacao) error {
	db, err := sql.Open("sqlite3", "cotacao.db")
	util.ErrorHandler(err)
	defer db.Close()

	ctxDB, cancelDB := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancelDB()

	create_sts := `
	CREATE TABLE IF NOT EXISTS cotacoes(
		code TEXT NOT NULL,
		codein TEXT NOT NULL,
		name TEXT NOT NULL,
		high TEXT NOT NULL,
		low TEXT NOT NULL,
		varBid TEXT NOT NULL,
		pctChange TEXT NOT NULL,
		bid TEXT NOT NULL,
		ask TEXT NOT NULL,
		timestamp TEXT NOT NULL,
		create_date TEXT NOT NULL
	);
	`

	_, err = db.Exec(create_sts)
	util.ErrorHandler(err)

	insert_db, err := db.BeginTx(ctxDB, nil)
	util.ErrorHandler(err)

	insert_sts, err := insert_db.PrepareContext(
		ctxDB,
		"INSERT INTO cotacoes VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	)

	util.ErrorHandler(err)

	_, err = insert_sts.Exec(
		cotacao.Usdbrl.Code,
		cotacao.Usdbrl.Codein,
		cotacao.Usdbrl.Name,
		cotacao.Usdbrl.High,
		cotacao.Usdbrl.Low,
		cotacao.Usdbrl.VarBid,
		cotacao.Usdbrl.PctChange,
		cotacao.Usdbrl.Bid,
		cotacao.Usdbrl.Ask,
		cotacao.Usdbrl.Timestamp,
		cotacao.Usdbrl.CreateDate,
	)

	if err != nil {
		insert_db.Rollback()
	} else {
		insert_db.Commit()
	}

	return nil
}
