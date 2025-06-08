package main

import (
	"auth/config"
	"auth/internal/controller"
	"auth/internal/repo/pg"
	"auth/internal/server"
	"auth/internal/service"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	pgDB := connectToDB()
	userStorage := pg.NewUserStorage(pgDB)
	authService := service.NewAuthService(userStorage)
	authController := controller.NewAuthController(authService)

	server := server.NewServer(
		*config.NewConfig(),
		*authController,
	)
	server.Start()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	var counts int64
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
