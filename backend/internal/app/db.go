package app

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	db_host     = os.Getenv("LDAP_AUTH_POSTGRES_HOST")
	db_port     = os.Getenv("LDAP_AUTH_POSTGRES_PORT")
	db_user     = os.Getenv("LDAP_AUTH_POSTGRES_USER")
	db_password = os.Getenv("LDAP_AUTH_POSTGRES_PASSWORD")
	db_dbname   = os.Getenv("LDAP_AUTH_POSTGRES_DB")
)

var db *sql.DB

func SetupDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db_host, db_port, db_user, db_password, db_dbname)

	time.Sleep(4 * time.Second)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
}

func CloseDB() {
	db.Close()
}
