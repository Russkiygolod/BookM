package main

// todo файлы запуска приложений называются main.go
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
	// todo убери эти комментарии, логика довольна простая, не требует пояснений дополнительных
	//  БД  //
	// выбор бд обуславливается постановкой задачи
	//PostgreSQL
	url := "postgres://" + "postgres" + ":" + "postgres" + "@" + "localhost" + ":" + "5432" + "/" + "ttt" // todo вынеси все креды для постгри и монги в файл конфигруации и читай его при запуске
	conn, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Failed to init DB conf - %v", err)
	}
	db2 := postgres.New(conn) // todo дай осмысленное название переменной (postgresDB например)

	//mongo
	p := "mongodb://localhost:27017/"
	mongoDb, err := mongo.New(p)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoDb.Close()

	// сервер //

	_, _ = db2, mongoDb //выбор БД todo это конечно лютый треш, ну да ладно
	// todo переименуй переменную или название пакета, она у тебя также называется  как имя пакета, (a просто или appi)
	apii := internal.New(db2)
	appi := api.New(apii)                     //подставить
	http.ListenAndServe(":80", appi.Router()) // todo не обрабатываешь ошибку, которая возвращается методом, посмотри в инете как обрабатывают http.ListenAndServe
}
