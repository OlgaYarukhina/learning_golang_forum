package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) GetPost(id int) (*Post, error) {
	post := &Post{}
	row := m.DB.QueryRow(`SELECT * FROM post WHERE ID = ?`, id)
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.User_id, &post.Created_at)
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

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Category, &s.User_id, &s.Created_at)
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

func (m *PostModel) Insert(header, description string, categoryId int, userId int, created_at time.Time) error {
	stmt := `INSERT INTO post (header, description,category_id,user_id, created_at)
    VALUES(?,?,?,?,current_date)`
	result, err := m.DB.Exec(stmt, header, description, categoryId, userId, created_at)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Println(id)
	return nil
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

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Category, &s.User_id, &s.Created_at)
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

func (m *PostModel) CreateLike(user_id, post_id int) error {
	likedOrNot := m.checkLike(user_id, post_id)
	switch likedOrNot {
	case true:
		stmt := "DELETE FROM likes WHERE user_id = ? AND post_id = ?"
		_, err := m.DB.Exec(stmt, user_id, post_id)
		if err != nil {
			return err
		}
	case false:

		stmt := `INSERT INTO likes (user_id, post_id)
    VALUES(?,?)`
		result, err := m.DB.Exec(stmt, user_id, post_id)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		log.Println(id)
		return nil
	}
	return nil
}

func (m *PostModel) checkLike(user_id, post_id int) bool {
	var id int
	// Query for a value based on a single row.
	if err := m.DB.QueryRow("SELECT id from likes where post_id = ? AND user_id = ?", post_id, user_id).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			return false
		}
		fmt.Println(err)
	}
	return true
}
