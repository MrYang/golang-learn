package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Init(database string) {
	log.Println("init db")
	var err error
	db, err = sql.Open("mysql", database)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
}
