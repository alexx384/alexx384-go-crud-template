package db

import (
	"embed"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var dbMigrationFs embed.FS

func GetMigrationDriver() (source.Driver, error) {
	return iofs.New(dbMigrationFs, "migrations")
}
