package models

import (
	"database/sql"
	"fmt"
)

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) GetUserById(id *int) string {
	var username string
	// Query for a value based on a single row.
	if err := m.DB.QueryRow("SELECT username from user where id = ?", id).Scan(&username); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
		}
		fmt.Println(err)
	}
	return username
}
func (m *CommentModel) Comments(id int) ([]*Comment, error) {
	stmt := `SELECT * FROM comments WHERE post_id = ? ORDER BY created_at ASC LIMIT 200`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []*Comment

	for rows.Next() {
		s := &Comment{}
		err = rows.Scan(&s.ID, &s.Comment, &s.Post_id, &s.User_id, &s.Created_at)
		s.Username = m.GetUserById(&s.User_id)
		if err != nil {
			return nil, err
		}
		comments = append(comments, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
