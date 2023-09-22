package db

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func New(username, password, host, dbName string) (*DB, error) {
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", username, password, host, dbName)
	conn, err := sqlx.Open("pgx", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("[db.New] could not open connection to DB: %w", err)
	}
	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("[db.New] could not complete initial DB ping: %w", err)
	}
	return &DB{conn}, nil
}
