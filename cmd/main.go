package main

import (
	"net/http"

	"BookM/internal"
	"BookM/internal/api"
	"BookM/internal/storage/postgres"

	"log"

	"BookM/internal/storage/mongo"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func main() {
	//PostgreSQL
	url := "postgres://" + "postgres" + ":" + "postgres" + "@" + "localhost" + ":" + "5432" + "/" + "ttt" // todo вынеси все креды для постгри и монги в файл конфигруации и читай его при запуске
	conn, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Failed to init DB conf - %v", err)
	}
	postgresDB := postgres.New(conn)

	//mongo
	p := "mongodb://localhost:27017/"
	mongoDb, err := mongo.New(p)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoDb.Close()

	_, _ = postgresDB, mongoDb
	inter := internal.New(mongoDb)
	appi := api.New(inter) //подставить
	log.Fatal(http.ListenAndServe(":80", appi.Router()))
}
