package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `
	insert into account (first_name, last_name, number, balance, created_at)
	values ($1, $2, $3, $4, $5)
	`
	_, err := s.db.Exec(query, account.FirstName, account.LastName, account.Number, account.Balance, account.Created_at)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := `select * from account where id = $1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("Account with id - %d wasn't found", id)
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query(`delete from account where id = $1`, id)
	return err
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `
	select * from account
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account  (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number int,
		balance int,
		created_at timestamp
	)`
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func NewPostgresStore() (*PostgresStore, error) {
	connString := "user=postgres dbname=postgres password=gobank sslmode=disable"

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := &Account{}
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.Created_at)
	return account, err
}
