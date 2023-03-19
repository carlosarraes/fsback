package db

import "database/sql"

type PostgresConn struct {
	DB *sql.DB
}

type User struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Progress  float64 `json:"progress"`
}

func (m *PostgresConn) GetUsers() ([]User, error) {
	rows, err := m.DB.Query("SELECT first_name, last_name, progress FROM data.user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.FirstName, &user.LastName, &user.Progress)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (m *PostgresConn) DeleteUser(lastName string) error {
	data, err := m.DB.Exec("DELETE FROM data.user WHERE last_name = $1", lastName)
	if err != nil {
		return err
	}
	status, _ := data.RowsAffected()
	if status == 0 {
		return err
	}

	return nil
}

func (m *PostgresConn) CreateUser(user User) error {
	data, err := m.DB.Exec("INSERT INTO data.user (first_name, last_name, progress) VALUES ($1, $2, $3)", user.FirstName, user.LastName, user.Progress)
	if err != nil {
		return err
	}
	status, _ := data.RowsAffected()
	if status == 0 {
		return err
	}

	return nil
}

func (m *PostgresConn) SumCheck() (float64, error) {
	sumCheck, err := m.DB.Query("SELECT sum(progress) FROM data.user")
	if err != nil {
		return 0, err
	}
	defer sumCheck.Close()
	var sum float64
	if sumCheck.Next() {
		err := sumCheck.Scan(&sum)
		if err != nil {
			return 0, err
		}
	}

	return sum, nil
}
