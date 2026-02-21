package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

type Service struct{
	Pool *pgxpool.Pool
}

func New() (*Service, error){
	//get database URL from env variable
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == ""{
		log.Fatalf("DATABASE_URL env variable is not set")
	}

	//handle timeouts and cancellations(create background context)
	ctx := context.Background()

	//create connection pool
	//const pool = new Pool({ connectionString }) in Node
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil{
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	//ping db to verify connection
	err = pool.Ping(ctx)
	if err != nil{
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Successfully connected to the database")

	return &Service{
		Pool: pool,
	}, err
}


// Close gracefully shuts down the connection pool (used when the server stops)
func (s *Service) Close() {
	if s.Pool != nil{
		s.Pool.Close()
	}
}