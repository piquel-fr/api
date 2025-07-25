package database

import (
	"context"
	"log"

	repository "github.com/piquel-fr/api/database/generated"
	"github.com/piquel-fr/api/services/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Queries *repository.Queries
var connection *pgxpool.Pool

func InitDatabase() {
	log.Printf("[Database] Attempting to connect to the database...\n")


    connection, err := pgxpool.New(context.Background(), config.Envs.DBURL)
	if err != nil {
		panic(err)
	}

	err = connection.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	Queries = repository.New(connection)
	log.Printf("[Database] Successfully connected to the database!\n")
}

func DeinitDatabase() {
    connection.Close()
}

