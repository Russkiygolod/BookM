package main

import (
	"BookM/pkg/api"
	"BookM/pkg/storage/postgres"
	"net/http"

	"BookM/pkg/storage/mongo"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func main() {
	//  БД  //
	// выбор бд обуславливается постановкой задачи
	//PostgreSQL
	url := "postgres://" + "postgres" + ":" + "postgres" + "@" + "localhost" + ":" + "5432" + "/" + "ttt"
	conn, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Failed to init DB conf - %v", err)
	}
	db2 := postgres.New(conn)

	//mongo
	p := "mongodb://localhost:27017/"
	mongoDb, err := mongo.New(p)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoDb.Close()

	// сервер //
	_, _ = db2, mongoDb //выбор БД
	api := api.New(db2) //подставить
	http.ListenAndServe(":80", api.Router())
}
