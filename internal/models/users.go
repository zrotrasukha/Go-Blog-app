package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	//postgres
	insertQuery := `INSERT INTO users (name, email, hashed_password)
									VALUES ($1, $2, $3)`
	_, err = m.DB.Exec(insertQuery, name, email, hashedPassword)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && pqErr.Code.Name() == "unique_violation" {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, hashed_password FROM users WHERE email=$1`
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)

	fmt.Printf("User authenticated successfully with ID: %d\n", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("No user found with email: %s\n", email)
			return 0, ErrInvalidCredentials
		} else {
			fmt.Printf("Error querying user with email: %s, error: %v\n", email, err)
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			fmt.Println("Password mismatch for user with email:", email)
			return 0, err
		}
	}
	fmt.Printf("User authenticated successfully with ID: %d\n", id)
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
