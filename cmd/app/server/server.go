package server

import (
	"context"
	"crud/cmd/app/config"
	"crud/cmd/app/config/database"
	"crud/internal"
	"crud/internal/middleware"
	"crud/internal/repository/db"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // To allow file:// path in migration
	"github.com/hellofresh/health-go/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"time"
)

func runDbMigration(pool *pgxpool.Pool) error {
	migrationDriver, err := db.GetMigrationDriver()
	if err != nil {
		return err
	}
	sqlDB := stdlib.OpenDBFromPool(pool)
	conn, err := sqlDB.Conn(context.Background())
	defer func(conn *sql.Conn) {
		err := conn.Close()
		if err != nil {
			slog.Default().Warn("failed to close database connection", slog.String("error", err.Error()))
		}
	}(conn)
	driver, err := postgres.WithConnection(context.Background(), conn, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", migrationDriver, "postgres", driver)
	if err != nil {
		return err
	}
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

func setupHealthCheck(app *gin.Engine, pool *pgxpool.Pool) error {
	healthcheck, err := health.New(health.WithSystemInfo(), health.WithComponent(health.Component{
		Name:    "crud",
		Version: "v1.0.0",
	}))
	if err != nil {
		return err
	}
	err = healthcheck.Register(health.Config{
		Name:      "database",
		Timeout:   time.Second * 2,
		SkipOnErr: false,
		Check: func(ctx context.Context) error {
			return pool.Ping(ctx)
		},
	})
	if err != nil {
		return err
	}
	app.GET("/status", func(c *gin.Context) {
		healthcheck.HandlerFunc(c.Writer, c.Request)
	})
	return nil
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
	if err = setupHealthCheck(app, dbPool); err != nil {
		logger.Error("Error setting up health check", slog.String("error", err.Error()))
		return nil, nil, err
	}
	app.Use(middleware.JSONLogMiddleware())
	app.Use(gin.Recovery())
	internal.SetupRouter(dbPool, app)
	return app, dbPool, nil
}
