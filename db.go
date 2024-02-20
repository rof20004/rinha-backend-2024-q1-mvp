package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var db *pgxpool.Pool

func init() {
	var (
		host    = os.Getenv("DATABASE_HOST")
		port    = os.Getenv("DATABASE_PORT")
		user    = os.Getenv("DATABASE_USER")
		pass    = os.Getenv("DATABASE_PASS")
		dbname  = os.Getenv("DATABASE_NAME")
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
	)

	log.Println("Connecting to database:", host, port)

	ctx := context.Background()

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("error while parsing '%s' postgresql database: %s\n", dbname, err.Error())
	}

	config.MinConns = 5
	config.MaxConns = 10

	log.Println("Opening connection pool to database:", config.ConnString())

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("error while opening '%s' postgresql database: %s\n", dbname, err.Error())
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("error while ping '%s' postgresql database: %s\n", dbname, err.Error())
	}

	db = pool
}
