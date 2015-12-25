package models

import (
	"github.com/gocql/gocql"
	"log"
	"time"
	"bytes"
	"strconv"
)

type Sub struct {
	ID          string
	Name        string `json:"name"`
	Subscribers int    `json:"subscribers"`
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

func (querier Querier) GetPosts(subName string)([]*Post, error){
	posts, err := querier.QueryForPosts(subName)
	log.Println("Post count", len(posts))
	return posts, err
}

func renderSub(posts []*Post) *bytes.Buffer {
	var buffer bytes.Buffer

	for _, post := range posts {
		buffer.WriteString(`<div class="post" id="`+ post.ID + `"><span style='padding-right:10px;'>` + strconv.Itoa(post.Upvotes) + `</span><a href="/r/` + post.Sub + `/comments/` + post.ID + `">` + post.Title + `</a></div>`)
	}
	return &buffer
}

func (querier Querier)InsertSub(subName string)(error){
	if err := querier.Session.Query(`INSERT INTO subs (id, name, subscribers, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		gocql.TimeUUID(),
		subName,
		0,
		time.Now(),
		time.Now()).Exec(); err != nil {
		return err
	}
	return nil
}

func (querier Querier) GetSubs()([]*Sub, error){
	var id gocql.UUID
	var name string
	var subscribers int
	var createdAt time.Time
	var updatedAt time.Time

	var subs []*Sub

	iter := querier.Session.Query(`SELECT id, name, subscribers, created_at, updated_at FROM subs`).Iter()
	for iter.Scan(&id, &name, &subscribers, &createdAt, &updatedAt) {
				subs = append(subs, &Sub{
					id.String(),
					name,
					subscribers,
					createdAt,
					updatedAt,
				})
	}
	return subs, nil
}

func (querier Querier) QueryForPosts(subName string)([]*Post, error){
	var id gocql.UUID
	var title string
	var upvotes int
	var downvotes int
	var commentCount int
	var createdAt time.Time
	var updatedAt time.Time

	var posts []*Post

	iter := querier.Session.Query(`SELECT id, title, upvotes, downvotes, created_at, updated_at, comment_count FROM posts WHERE sub_id = ?`,
		subName).Iter()
	for iter.Scan(&id, &title, &upvotes, &downvotes, &createdAt, &updatedAt, &commentCount) {
				posts = append(posts, &Post{
					id.String(),
					title,
					"",
					"",
					"",
					subName,
					upvotes,
					downvotes,
					commentCount,
					createdAt,
					updatedAt,
				})
	}
	return posts, nil
}
