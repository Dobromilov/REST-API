package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL) // подключаюсь к postgres и делаю ping
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(25) // max кол-во подключений одновременно
	db.SetMaxIdleConns(5)  // max соеден. в простое

	return db, nil
}
