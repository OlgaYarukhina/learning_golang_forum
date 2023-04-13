package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Latest() ([]*Post, error) {
	stmt := `SELECT * FROM post ORDER BY created_at ASC LIMIT 200`
	rows, err := m.DB.Query(stmt)

	defer rows.Close()
	var posts []*Post
	for rows.Next() {
		s := &Post{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.User_id, &s.Created_at)
		s.Category_name, err = m.getCategoryRelation(&s.ID)
		s.Like, s.Dislike, err = m.getCountLikesByPostId(&s.ID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}

	return posts, nil
}

func (m *PostModel) getCategoryRelation(post_id *int) ([]string, error) {
	stmt := `SELECT category_id FROM categoryPostRelation WHERE post_id = ?`
	var array_of_names []string
	var category_id int
	rows, err := m.DB.Query(stmt, post_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&category_id)
		name, err := m.getNameOfCategoryById(&category_id)
		if err != nil {
			return nil, err
		}
		array_of_names = append(array_of_names, name)
	}

	return array_of_names, nil
}

func (m *PostModel) getNameOfCategoryById(category_id *int) (string, error) {
	var category_name string
	// Query for a value based on a single row.
	if err := m.DB.QueryRow("SELECT name from category where id = ?", category_id).Scan(&category_name); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
		}
		fmt.Println(err)
	}
	return category_name, nil
}

func (m *PostModel) GetPost(id int) (*Post, error) {
	post := &Post{}
	row := m.DB.QueryRow(`SELECT * FROM post WHERE ID = ?`, id)
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.User_id, &post.Created_at)
	post.Category_name, err = m.getCategoryRelation(&post.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return post, fmt.Errorf("post not found")
		}
		return post, err
	}
	return post, nil
}

func (m *PostModel) getCountLikesByPostId(id *int) (int, int, error) {
	var id_post int
	stmt := `SELECT id FROM like WHERE post_id = ? AND type_of = "like"`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()
	var likes []int
	var dislikes []int
	stmt_dislike := `SELECT id FROM like WHERE post_id = ? AND type_of = "dislike"`
	rows_dis, err := m.DB.Query(stmt_dislike, id)
	for rows_dis.Next() {
		var c int
		err = rows_dis.Scan(&id_post)
		dislikes = append(dislikes, c)
	}
	for rows.Next() {
		var s int
		err = rows.Scan(&id_post)
		if err != nil {
			return 0, 0, err
		}
		likes = append(likes, s)
	}
	if err = rows.Err(); err != nil {
		return 0, 0, err
	}

	return len(likes), len(dislikes), nil
}

func (m *PostModel) Insert(header, description string, userId int, created_at time.Time, category []string) error {
	stmt := `INSERT INTO post (header, description, user_id, created_at)
    VALUES(?,?,?, current_date)`

	result, err := m.DB.Exec(stmt, header, description, userId, created_at)

	id, err := result.LastInsertId()
	for _, category_id := range category {
		cat_id, err := strconv.Atoi(category_id)
		stmt := `INSERT INTO categoryPostRelation (post_id,category_id) VALUES (?,?)`
		_, err = m.DB.Exec(stmt, id, cat_id)
		if err != nil {
			fmt.Println(err)
		}
	}
	if err != nil {
		return err
	}

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
	for i := len(allPosts) - 1; i >= 0; i-- {
		fromLastPost = append(fromLastPost, allPosts[i])
	}

	return fromLastPost
}

func (m *PostModel) GetPostsByCategory(id int) ([]*Post, error) {
	stmt := `SELECT * FROM categoryPostRelation WHERE category_id=?`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categoryPostRelation []*CategoryPostRelation
	var postIds []int
	for rows.Next() {
		s := &CategoryPostRelation{}
		err = rows.Scan(&s.ID, &s.Post_id, &s.Category_id)
		if err != nil {
			return nil, err
		}
		categoryPostRelation = append(categoryPostRelation, s)
		postIds = append(postIds, s.Post_id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	posts, err := getPostsByIDList(m.DB, postIds)

	return posts, nil
}

func getPostsByIDList(db *sql.DB, idList []int) ([]*Post, error) {
	// Create a string of question marks corresponding to the number of IDs.
	placeholders := strings.Trim(strings.Repeat(",?", len(idList)), ",")

	// Prepare the SQL statement with placeholders for the IDs.
	stmt, err := db.Prepare(fmt.Sprintf("SELECT id, header, description, user_id, created_at FROM post WHERE id IN (%s)", placeholders))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Convert the slice of IDs to a variadic argument.
	args := make([]interface{}, len(idList))
	for i, id := range idList {
		args[i] = id
	}

	// Create a slice to hold the retrieved posts.
	posts := make([]*Post, 0, len(idList))

	// Query the database with the list of IDs.
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the result set and add each post to the slice.
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.User_id, &post.Created_at); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	// Check for any errors during iteration.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) CreateLike(user_id, post_id int, type_of_like string) error {
	likedOrNot := m.checkLike(user_id, post_id)
	switch likedOrNot {
	case true:
		stmt := "DELETE FROM like WHERE user_id = ? AND post_id = ? AND type_of = ?"
		result, err := m.DB.Exec(stmt, user_id, post_id, type_of_like)
		res, err := result.RowsAffected()
		if res == int64(0) {
			stmt := "DELETE FROM like WHERE user_id = ? AND post_id = ?"
			_, err := m.DB.Exec(stmt, user_id, post_id, type_of_like)
			if err != nil {

				return err
			}
			m.CreateLike(user_id, post_id, type_of_like)
		}
		if err != nil {
			return err
		}
	case false:
		stmt := `INSERT INTO like(post_id, user_id,type_of)
    VALUES(?,?,?)`
		_, err := m.DB.Exec(stmt, post_id, user_id, type_of_like)

		if err != nil {
			return err
		}

	}
	return nil
}

func (m *PostModel) checkLike(user_id, post_id int) bool {
	var id int
	// Query for a value based on a single row.
	if err := m.DB.QueryRow("SELECT id from like where post_id = ? AND user_id = ?", post_id, user_id).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			return false
		}
		fmt.Println(err)
	}
	return true
}
