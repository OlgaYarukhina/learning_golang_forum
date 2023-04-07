package models

import (
	"database/sql"
	"fmt"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) GetPost(id int) (*Post, error) {
	post := &Post{}
	row := m.DB.QueryRow(`SELECT * FROM post WHERE ID = ?`, id)
	err := row.Scan(&post.ID, &post.Header, &post.Description, &post.Category_id, &post.User_id, &post.Created_at)
	if err != nil {
		if err == sql.ErrNoRows {
			return post, fmt.Errorf("post not found")
		}
		return post, err
	}
	return post, nil
}

func (m *PostModel) getCountLikesByPostId(id *int) (int, error) {
	var id_post int
	stmt := `SELECT id FROM likes WHERE post_id = ?`

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

func (m *PostModel) Latest() ([]*Post, error) {

	stmt := `SELECT * FROM post ORDER BY created_at ASC LIMIT 200`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		s := &Post{}

		err = rows.Scan(&s.ID, &s.Header, &s.Description, &s.Category_id, &s.User_id, &s.Created_at)
		s.Like, err = m.getCountLikesByPostId(&s.ID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
func (m *PostModel) GetPostsByCategory(id int) ([]*Post, error) {

	stmt := `SELECT * FROM post WHERE id=?`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		s := &Post{}

		err = rows.Scan(&s.ID, &s.Header, &s.Description, &s.Category_id, &s.User_id, &s.Created_at)
		s.Like, err = m.getCountLikesByPostId(&s.ID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
