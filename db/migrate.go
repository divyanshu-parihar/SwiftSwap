package db

import (
	"context"
	"fmt"
	"log"

	"github.com/uptrace/bun"
)

func RunMigrations(db *bun.DB) error {
	ctx := context.Background()

	// Create the users table if it doesn't exist
	_, err := db.NewCreateTable().Model((*Transaction)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	return nil
}
func main() {
	db := StartServer()
	if err := RunMigrations(db); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
}
