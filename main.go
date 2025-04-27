package main

import (
	"crud/internal"
	"crud/internal/config"
	"crud/internal/config/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // To allow file:// path in migration
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"os"
)

func runDbMigration(pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/repository/db/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	return m.Up()
}

func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	dbPool, err := database.NewPool(appConfig.DB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	if err = runDbMigration(dbPool); err != nil {
		fmt.Fprintf(os.Stderr, "Error running migration: %v\n", err)
		os.Exit(1)
	}

	app := gin.Default()
	internal.SetupRouter(dbPool, app)

	err = app.Run(":8080")
	if err != nil {
		fmt.Println("Something went wrong")
		return
	}
}
