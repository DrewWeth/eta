package models

import (
	"github.com/gocql/gocql"
	// "encoding/json"
	"time"
	"log"
)

type Post struct {
	ID        string `json:"title"`
	Title     string `json:"title"`
	Link      string `json:"link"`
	Comment   string `json:"content"`
	PostToken string `json:"-"`
	Sub       string `json:"sub"`
	Upvotes   int    `json:"upvotes"`
	Downvotes int    `json:"downvotes"`
	CommentCount	int	`json:"comment_count"` // Redundency
	UpdatedAt time.Time
	CreatedAt time.Time
}

type PostReqParams struct{
	SubName	string `schema:"sub_name"`
	Title	string `schema:"title"`
	Comment string `schema:"comment"`
	Link string `schema:"link"`
}

func (querier *Querier) GetPost(postID string)(*Post, error){
	var id gocql.UUID
	var comment, title, link, sub string
	var upvotes int
	var downvotes int
	var commentCount int
	var createdAt time.Time
	var updatedAt time.Time

	if err := querier.Session.Query(`SELECT id, title, link, sub_id, comment, upvotes, downvotes, comment_count, created_at, updated_at FROM posts WHERE id = ? LIMIT 1`,
		postID).Consistency(gocql.One).Scan(&id, &title, &link, &sub, &comment, &upvotes, &downvotes, &commentCount, &createdAt, &updatedAt); err != nil {
		log.Println(err)
		return nil, err
	}
	return &Post{id.String(), title, link, comment, "", sub, upvotes, downvotes, commentCount, updatedAt, createdAt}, nil
}

func (querier *Querier) InsertPost(subName, title, comment, link string) error {
	if err := querier.Session.Query(`INSERT INTO posts (id, sub_id, title, comment, link, upvotes, downvotes, comment_count, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		gocql.TimeUUID(),
		subName,
		title,
		comment,
		link,
		0,
		0,
		0,
		time.Now(),
		time.Now(),
		).Exec(); err != nil {
		return err
	}
	return nil
}
