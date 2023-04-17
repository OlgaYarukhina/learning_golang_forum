package models

import (
	"database/sql"
	"fmt"
	"time"
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
		s.Like, err = m.getLikeCountsByCommentId(&s.ID)
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

func (m *CommentModel) Insert(comment string, post_id int, user_id int, created_at time.Time) error {
	stmt := `INSERT INTO comments (comment, post_id, user_id, created_at)
    VALUES(?, ?, ?, current_date)`
	_, err := m.DB.Exec(stmt, comment, post_id, user_id, created_at)
	if err != nil {
		return err
	}
	return nil
}

func (m *CommentModel) CreateLike(comment_id, user_id int) error {
	likedOrNot := m.checkLike(user_id, comment_id)
	switch likedOrNot {
	case true:
		stmt := "DELETE FROM comment_like WHERE comment_id = ? AND user_id = ?"
		_, err := m.DB.Exec(stmt, comment_id, user_id)

		if err != nil {
			return err
		}
	case false:
		stmt := `INSERT INTO comment_like(comment_id, user_id)
    VALUES(?,?)`
		_, err := m.DB.Exec(stmt, comment_id, user_id)

		if err != nil {
			return err
		}

	}
	return nil
}

func (m *CommentModel) checkLike(user_id, comment_id int) bool {
	var id int
	// Query for a value based on a single row.
	if err := m.DB.QueryRow("SELECT id from comment_like where comment_id = ? AND user_id = ?", comment_id, user_id).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			return false
		}
		fmt.Println(err)
	}
	return true
}

func (m *CommentModel) getLikeCountsByCommentId(id *int) (int, error) {
	var id_post int
	stmt := `SELECT id FROM comment_like WHERE comment_id = ?`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var likes []int
	for rows.Next() {
		var s int
		err = rows.Scan(&id_post)
		if err != nil {
			return 0, err
		}
		likes = append(likes, s)
	}
	if err = rows.Err(); err != nil {
		return 0, err
	}

	return len(likes), nil
}
