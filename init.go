package main

//go get github.com/jackc/pgx/v5
import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

const schemaSQL = `
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_online BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    profile_picture TEXT
);

CREATE TABLE IF NOT EXISTS rooms (
    room_id TEXT PRIMARY KEY,
    last_message_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS events (
    event_id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL REFERENCES rooms(room_id),
    title TEXT NOT NULL,
    location TEXT,
    capacity INTEGER,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS room_user (
    room_id TEXT NOT NULL REFERENCES rooms(room_id),
    user_id TEXT NOT NULL REFERENCES users(user_id),
    PRIMARY KEY (room_id, user_id)
);

CREATE TABLE IF NOT EXISTS event_participation (
    event_id TEXT NOT NULL REFERENCES events(event_id),
    user_id TEXT NOT NULL REFERENCES users(user_id),
    PRIMARY KEY (event_id, user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    message_id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL REFERENCES rooms(room_id),
    sender_id TEXT NOT NULL REFERENCES users(user_id),
    content TEXT,
    message_type TEXT NOT NULL,
    media_url TEXT,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_read BOOLEAN DEFAULT FALSE
);
`

func main() {
	// Get connection string from environment variable or fallback
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		// Default connection string for local development
		connString = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	ctx := context.Background()

	// Connect to database
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("fail to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	// Execute schema creation
	_, err = conn.Exec(ctx, schemaSQL)
	if err != nil {
		log.Fatalf("fail to execute schema: %v\n", err)
	}

	fmt.Println("done")
}
