package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database/repository"
)

var Queries *repository.Queries
var Connection *pgxpool.Pool

func InitDatabase() {
	log.Printf("[Database] Attempting to connect to the database...\n")

	conn, err := pgxpool.New(context.Background(), config.Envs.DBURL)
	if err != nil {
		panic(err)
	}

	Connection = conn

	if err = Connection.Ping(context.Background()); err != nil {
		panic(err)
	}

	Queries = repository.New(Connection)
	log.Printf("[Database] Successfully connected to the database!\n")
}
