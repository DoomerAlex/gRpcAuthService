package main

import (
	"errors"
	"flag"
	"fmt"
	// Библиотека для миграций
	"github.com/golang-migrate/migrate"
	// Драйвер для выполнения мишраций SQLite3
	_ "github.com/golang-migrate/migrate/database/sqlite3"
	// Драйвер для получения миграций из файлов
	_ "github.com/golang-migrate/migrate/source/file"
)

const (
	StoragePathParamName     = "storage-path"
	MigrationsPathParamName  = "migrations-path"
	MigrationsTableParamName = "migrations-table"
	DefaultMigrationTable    = "schema_migrations"
)

// Миграционная утилита
func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, StoragePathParamName, "", "path to the storage")
	flag.StringVar(&migrationsPath, MigrationsPathParamName, "", "path to migrations")
	flag.StringVar(&migrationsTable, MigrationsTableParamName, DefaultMigrationTable, "name of migrations table in BD")
	flag.Parse()

	if storagePath == "" {
		panic(StoragePathParamName + " is required")
	}
	if migrationsPath == "" {
		panic(MigrationsPathParamName + " is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
