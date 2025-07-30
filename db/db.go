package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var Conn *pgx.Conn

func Init() {
	_ = godotenv.Load()
	url := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	log.Println("Connected to PostgreSQL")
	Conn = conn
}
