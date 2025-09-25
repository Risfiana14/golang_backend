package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	dsn := "host=localhost user=postgres password=12345678 dbname=Alumni_db port=5432 sslmode=disable"

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Gagal ping database:", err)
	}
	fmt.Println("Berhasil terhubung ke database PostgreSQL")
}