package repository

import (
	"database/sql"

	"github.com/carlosarraes/fsback/repository/dbrepo"
)

type DbRepo interface {
	Connection() *sql.DB
	GetUsers() ([]dbrepo.User, error)
	DeleteUser(lastName string) error
	CreateUser(user dbrepo.User) error
	SumCheck(progress float64) error
}
