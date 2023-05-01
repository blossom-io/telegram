package postgres

import (
	"database/sql"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var driverName = "pgx"

type Postgres struct {
	DB *sql.DB
	Sq sq.StatementBuilderType
}

// New creates postgres.
func New(pgURL string) (*Postgres, error) {
	var (
		pg  Postgres
		err error
	)

	pg.Sq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar) // pgx supports only dollar format $1, $2, etc

	pg.DB, err = sql.Open(driverName, pgURL)
	if err != nil {
		log.Fatalln("postgres - New:", err)
	}

	result, err := pg.DB.Exec("SELECT now();")
	if err != nil {
		return nil, err
	}

	fmt.Println(result.RowsAffected())

	return &pg, nil
}

func (pg *Postgres) Close() error {
	return pg.DB.Close()
}
