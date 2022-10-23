package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ConfigDB struct {
	Host string
	Port int
	Name string
}

func Connect(cnf ConfigDB) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s", cnf.Host, cnf.Port, cnf.Name)
	db, err := sqlx.Connect("mongodb", dsn)

	return db, err
}
