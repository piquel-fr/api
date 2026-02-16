package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/database/repository"
)

type StorageService interface {
	repository.Querier
	Close()
}

type databaseStorageService struct {
	repository.Queries
	connection *pgxpool.Pool
}

func NewDatabaseStorageService() StorageService {
	log.Printf("[Database] Attempting to connect to the database...\n")

	connection, err := pgxpool.New(context.Background(), config.Envs.DBURL)
	if err != nil {
		panic(err)
	}

	if err = connection.Ping(context.Background()); err != nil {
		panic(err)
	}

	queries := repository.New(connection)
	log.Printf("[Database] Successfully connected to the database!\n")

	return &databaseStorageService{
		connection: connection,
		Queries:    *queries,
	}
}

func (s *databaseStorageService) Close() {
	s.connection.Close()
}
