package server

import (
	"crud/cmd/app/config"
	"crud/cmd/app/config/database"
	"crud/internal"
	"crud/internal/middleware"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // To allow file:// path in migration
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"log/slog"
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
	defer func(m *migrate.Migrate) {
		sourceErr, databaseError := m.Close()
		if sourceErr != nil {
			slog.Default().Warn("Source error during migration:", slog.String("error", sourceErr.Error()))
		}
		if databaseError != nil {
			slog.Default().Warn("Database error during migration:", slog.String("error", databaseError.Error()))
		}
	}(m)
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	} else {
		return err
	}
}

func ConfigureAppEngine(appConfig *config.Config, logLevelVar *slog.LevelVar) (*gin.Engine, *pgxpool.Pool, error) {
	logger := slog.Default()

	logger.Info("Starting server")

	appLogLevel, err := appConfig.App.ToSlogLevel()
	if err != nil {
		logger.Warn("Error converting app log level. Using default level",
			slog.String("error", err.Error()),
			slog.String("defaultLogLevel", slog.LevelInfo.String()))
		appLogLevel = slog.LevelInfo
	}
	//goland:noinspection GoDfaNilDereference
	logger.Info("Setting log level", slog.String("level", appLogLevel.String()))
	logLevelVar.Set(appLogLevel)

	dbPool, err := database.NewPool(appConfig.DB)
	if err != nil {
		logger.Error("Error connecting to database", slog.String("error", err.Error()))
		return nil, nil, err
	}

	if err = runDbMigration(dbPool); err != nil {
		dbPool.Close()
		logger.Error("Error running migration", slog.String("error", err.Error()))
		return nil, nil, err
	}

	if appConfig.App.IsAppInReleaseMode() {
		logger.Info("Running app in release mode")
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DebugPrintFunc = func(format string, v ...interface{}) {
		logger.Warn(fmt.Sprintf(format, v...))
	}
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.Debug("endpoint",
			slog.String("method", httpMethod),
			slog.String("path", absolutePath),
			slog.String("handler", handlerName),
			slog.Int("handlers", nuHandlers))
	}
	app := gin.New()
	app.Use(middleware.JSONLogMiddleware())
	app.Use(gin.Recovery())
	internal.SetupRouter(dbPool, app)

	return app, dbPool, nil
}
