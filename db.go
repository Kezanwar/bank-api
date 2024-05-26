package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts([]*Account, error)
}

type PostGresDB struct {
	db *sql.DB
}

var db_user = os.Getenv("DB_USER")
var db_name = os.Getenv("DB_NAME")
var db_password = os.Getenv("DB_PASSWORD")

func NewPostgresDB() (*PostGresDB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", db_user, db_password, db_name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("postgres connected 🚀")

	return &PostGresDB{
		db: db,
	}, nil
}

func (s *PostGresDB) Init() error {
	return s.CreateAccountTable()
}

func (s *PostGresDB) CreateAccountTable() error {

	// _, errrrr := s.db.Exec("DROP TABLE Account")

	// if errrrr != nil {
	// 	return errrrr
	// }

	query := `CREATE TABLE if not exists Account (
    ID serial primary key,
	FirstName varchar(50),
    LastName varchar(50),
   	Number serial,
    Balance serial,
	CreatedAt timestamp
)`

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostGresDB) CreateAccount(acc *Account) error {

	query := `INSERT INTO Account 
	(FirstName, LastName, Number, Balance, CreatedAt)
	values ($1, $2, $3, $4, $5)
	;`

	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return err
	}

	// fmt.Printf("%+v\n", resp)

	return nil
}

func (s *PostGresDB) UpdateAccount(*Account) error {
	return nil
}

func (s *PostGresDB) DeleteAccountByID(id int) error {

	rows, findErr := s.db.Query(
		`SELECT * FROM Account where id = $1`, id,
	)

	if findErr != nil {
		return findErr
	}

	for rows.Next() {
		_, err := scanIntoAccount(rows)

		if err != nil {
			return err
		}
	}

	_, deleteErr := s.db.Query(
		`DELETE FROM Account where id = $1`, id,
	)

	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (s *PostGresDB) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query(
		`SELECT * FROM Account where id = $1`, id,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)

	}

	return nil, fmt.Errorf("Account %d not found", id)
}

func (s *PostGresDB) GetAccounts() ([]*Account, error) {

	rows, err := s.db.Query(
		`SELECT * FROM Account`,
	)

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

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := &Account{}
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	return account, err
}
