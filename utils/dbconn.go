package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var (
	db   *sql.DB
	once sync.Once
)

func ConnectDB() *sql.DB {
	once.Do(func() {

		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️ .env file tidak ditemukan, pakai environment system")
		}

		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		sslmode := os.Getenv("DB_SSLMODE")

		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user,
			password,
			host,
			port,
			name,
			sslmode,
		)

		conn, err := sql.Open("pgx", dsn)
		if err != nil {
			log.Fatal("❌ Gagal open DB:", err)
		}

		// test koneksi
		if err = conn.Ping(); err != nil {
			log.Fatal("❌ Gagal connect DB:", err)
		}

		// optional: connection pool tuning
		conn.SetMaxOpenConns(25)
		conn.SetMaxIdleConns(10)

		db = conn
		log.Println("✅ PostgreSQL connected")
	})

	return db
}
