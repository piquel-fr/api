package database

import (
	"context"
	"log"

	"github.com/PiquelChips/piquel.fr/services/config"
	repository "github.com/PiquelChips/piquel.fr/database/generated"
	"github.com/jackc/pgx/v5"
)

var Queries *repository.Queries
var connection *pgx.Conn

func InitDatabase() {

	log.Printf("[Database] Attempting to connect to the database...\n")

	connection, err := pgx.Connect(context.Background(), config.Envs.DB_URL)
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
    connection.Close(context.Background())
}

