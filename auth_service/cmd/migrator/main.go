package main

import (
	"errors"
	"flag"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.Parse()

	if storagePath == "" {
		storagePath = os.Getenv("STORAGEPATH")
	}

	if migrationsPath == "" {
		migrationsPath = os.Getenv("MIGRATIONSPATH")
	}

	if storagePath == "" {
		panic("storage path is empty!")
	}

	if migrationsPath == "" {
		panic("Migrations path is empty!")
	}

	m, err := migrate.New("file://"+migrationsPath, "sqlite3://"+storagePath)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migrations to apply...")
		}

		panic(err)
	}

	log.Println("Migrations applied!")
}