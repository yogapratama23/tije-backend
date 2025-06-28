package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func RunMigration(conn *pgx.Conn) {
	stdLibConn := stdlib.OpenDB(*conn.Config())
	driver, err := postgres.WithInstance(stdLibConn, &postgres.Config{})
	if err != nil {
		panic(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Running DB Migrations")

	err = m.Up()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Migration completed")
}
