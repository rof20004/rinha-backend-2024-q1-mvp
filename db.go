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
		connStr = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
	)

	log.Println("Connecting to database:", host, port)

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("error while opening '%s' postgresql database: %s\n", dbname, err.Error())
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("error while ping '%s' postgresql database: %s\n", dbname, err.Error())
	}

	db = pool
}
