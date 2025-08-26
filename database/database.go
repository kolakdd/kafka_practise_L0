package database

import (
	"kafkapractisel0/services"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func InitDB(e services.Env) (db *sqlx.DB) {
	dsn := e.GetDatabaseDSN()
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	configDB(db)
	err = db.Ping()
	if err != nil {
		log.Fatal("Error ping to database: ", err)
	}
	return db
}

func configDB(db *sqlx.DB) {
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
}
