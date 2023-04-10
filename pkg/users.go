package models

import (
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

//получаем все данные о пользователе по его логину
func (m *UserModel) GetUserByUsername(username string) (User, error) {
	stmt := `SELECT * FROM user WHERE username = ?`
	s := User{}
	row := m.DB.QueryRow(stmt, username)
	err := row.Scan(&s.ID, &s.Username, &s.Password, &s.Email, &s.Created_at)
	if err != nil {
		return s, nil
	}

	return s, err
}

// Insert - Метод для создания
func (m *UserModel) Insert(username, password, email string, created_at time.Time) error {

	stmt := `INSERT INTO user (username, password, email, created_at)
    VALUES(?, ?, ?, current_date)`

	_, err := m.DB.Exec(stmt, username, password, email, created_at)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) CheckUser(email string) (User, error) {
	stmt := `SELECT * FROM user WHERE email = ?`
	s := User{}
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&s.ID, &s.Username, &s.Password, &s.Email, &s.Created_at)
	if err != nil {
		return s, nil
	}

	return s, err

}
