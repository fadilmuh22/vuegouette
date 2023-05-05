package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", viper.GetString("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}
