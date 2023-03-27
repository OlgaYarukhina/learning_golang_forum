package models

import (
	"database/sql"
	"log"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(header, description string, categoryId, userId int) error {
	stmt := `INSERT INTO post (header, description, category_id, user_id, created_at)
    VALUES(?,?,?,?, current_date)`

	result, err := m.DB.Exec(stmt, header, description, categoryId, userId)
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


// Get one user posts
func (m *PostModel) GetUserPosts(userID int) []*Post {
	var allPosts []*Post

	rows, err := m.DB.Query(`SELECT id, header, description FROM post WHERE user_id = ?`, userID)
	if err != nil {
		log.Println("Bad request")
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		post := &Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {
			log.Println("Something wrong with db")
			log.Println(err)
			return nil
		}
		allPosts = append(allPosts, post)
	}

	fromLastPost := []*Post{}
	for i:= len(allPosts)-1; i >= 0; i-- {
		fromLastPost = append(fromLastPost, allPosts[i])
	}

	return fromLastPost
}