package dbrepo

import (
	"database/sql"
	"errors"
)

type TestDBRepo struct{}

func (m *TestDBRepo) Connection() *sql.DB {
	return nil
}

func (m *TestDBRepo) GetUsers() ([]User, error) {
	testUser := User{
		FirstName: "John",
		LastName:  "Doe",
		Progress:  42.0,
	}

	users := []User{testUser}
	return users, nil
}

func (m *TestDBRepo) DeleteUser(lastName string) error {
	if lastName == "NonExisting" {
		return errors.New("no rows affected")
	}
	return nil
}

func (m *TestDBRepo) CreateUser(user User) error {
	return nil
}

func (m *TestDBRepo) SumCheck(progress float64) error {
	if progress > 1.00 {
		return errors.New("progress is greater than 1.00")
	}
	return nil
}
