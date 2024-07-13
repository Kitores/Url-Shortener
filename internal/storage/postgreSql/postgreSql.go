package postgreSql

import (
	"JustTesting/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"sync"
)

type PostgreSqlx struct {
	db *sqlx.DB
}

type Storage struct {
	db *sql.DB
}

type Postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *PostgreSqlx
	pgOnce     sync.Once
)

// Use singleton template to make sure that only have one connection pool
func NewPG(connString string) (*PostgreSqlx, error) {
	pgOnce.Do(func() {
		db, err := sqlx.Connect("pgx", connString)
		if err != nil {
			fmt.Errorf("unable to create conection pool: %w", err)
		}
		pgInstance = &PostgreSqlx{db: db}
		if pgInstance.db == nil {
			log.Fatal("unable to create connection")
		}
		fmt.Println(pgInstance)
	})
	return pgInstance, nil
}
func (pg *Postgres) Ping(ctx context.Context, connString string) error {
	return pg.db.Ping(ctx)
}
func (pg *PostgreSqlx) SaveURL(urlToSave string, alias string) error {
	//var id int64
	const functionName = "storage/postgreSql/SaveURL()"
	query := fmt.Sprintf("INSERT INTO urls(alias,url) VALUES('%s','%s')", alias, urlToSave)
	_, err := pg.db.Exec(query)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Unique constraint violation
				return fmt.Errorf("%s: %w", functionName, storage.ErrUrlExists)
			} else {
				return fmt.Errorf("unexpected error while saving url %s: %w", urlToSave, err)
			}
		} else {
			fmt.Errorf("unexpected error while saving url %s: %w", urlToSave, err)
		}
	}
	//id, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("%s unable to get last inserted id: %w", functionName, err)
	}
	//TODO: handling error UrlExistErr
	return err
}
func (pg *PostgreSqlx) GetAlias(url string) (alias string, err error) {
	const functionName = "storage/postgreSql/GetAlias()"
	query := fmt.Sprintf("SELECT alias FROM urls WHERE url='%s'", url)
	err = pg.db.QueryRow(query).Scan(&alias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrAliasNotFound
		}
		return "", fmt.Errorf("%s: %w", functionName, err)
	}
	return
}

func (pg *PostgreSqlx) GetURL(alias string) (string, error) {
	var url string
	const functionName = "storage/postgreSql/GetURL()"
	//query := fmt.Sprintf("SELECT url FROM urls WHERE alias='%s'", alias)
	query2 := "SELECT url FROM urls WHERE alias=$1"
	err := pg.db.QueryRow(query2, alias).Scan(&url)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlNotFound
		}
		return "", fmt.Errorf("%s: %w", functionName, err)
	}
	return url, nil
}

func InitStorage() (*Storage, error) {
	var err error
	var db *sql.DB
	connStr := "postgres://postgres:angelo4ek@localhost:5432/test1?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecting to PostgreSQL")
		return nil, err
	}
	if err = db.Ping(); err != nil {
		fmt.Println("Error pinging PostgreSQL")
		return nil, err
	}
	fmt.Println("The database is connected")
	return &Storage{db: db}, nil
}

func (pg *PostgreSqlx) Delete(url string) error {
	functionName := "storage/postgreSql/Delete()"
	var err error
	deleteStmt := `DELETE FROM urls WHERE url = $1`
	_, err = pg.db.Exec(deleteStmt, url)
	if err != nil {
		fmt.Errorf("%s: %w", functionName, err)
	}
	return err
}
func (pg *PostgreSqlx) DeleteRange(left int64, right int64) error {
	functionName := "storage/postgreSql/Delete()"
	var err error
	for i := left; i <= right; i++ {
		deleteStmt := `DELETE FROM urls WHERE id = $1`
		_, err = pg.db.Exec(deleteStmt, i)
		if err != nil {
			fmt.Errorf("%s: %w", functionName, err)
		}
	}
	return err
}
func Update(strg *Postgres) {
	var err error
	updateStmt := `UPDATE urls SET id=1 WHERE id=6	`
	_, err = strg.db.Exec(context.Background(), updateStmt)
	if err != nil {
		fmt.Println(err)
	}
}
