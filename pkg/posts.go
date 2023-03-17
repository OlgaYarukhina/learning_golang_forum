package models

import (
	"database/sql"
	"errors"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Get(id int) (*Post, error) {

	statement := `SELECT id, header, description, created_at FROM post
    WHERE created_at > current_date AND id = ?`

	row := m.DB.QueryRow(statement, id)

	p := &Post{}

	err := row.Scan(&p.id, &p.header, &p.description, &p.created_at)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return p, nil
}
