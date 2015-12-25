// create keyspace eta with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
// create table eta.comments(id UUID, post_id text, content text, upvotes int, downvotes int, created_at timestamp, updated_at timestamp, PRIMARY KEY(id));
// create index on eta.comments(post_id);

// create table eta.caches(id UUID, post_id text, content text, upvotes int, downvotes int, created_at timestamp, updated_at timestamp, PRIMARY KEY(id));
// create index on eta.caches(post_id);

create table eta.posts(id UUID, sub_id text, title text, link text, comment text, upvotes int, downvotes int, comment_count int, created_at timestamp, updated_at timestamp, PRIMARY KEY(id, comment_count, upvotes, downvotes));
create index on eta.posts(sub_id);

// create table eta.subs(id UUID, name text, subscribers int, created_at timestamp, updated_at timestamp, PRIMARY KEY(id));
// create index on eta.subs(name);

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"os"
	"sort"
	"time"
)

type Comment struct {
	ID        string `json:"id"`
	ParentID  string `json:"-"`
	Content   string `json:"name"`
	Upvotes   int    `json:"upvotes"`
	Downvotes int    `json:"downvotes"`
	CreatedAt time.Time
	UpdatedAt time.Time  `json:"updated_at"`
	Children  []*Comment `json:"children"`
}

type Sub struct {
	ID          string
	Name        string `json:"name"`
	Subscribers int    `json:"subscribers"`
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

type Post struct {
	ID        string `json:"title"`
	Title     string `json:"title"`
	Link      string `json:"link"`
	Comment   string `json:"content"`
	PostToken string `json:"-"`
	Sub       string `json:"sub"`
	Upvotes   int    `json:"upvotes"`
	Downvotes int    `json:"downvotes"`
	UpdatedAt time.Time
	CreatedAt time.Time
}

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	APIToken     string `json:"api_token"`
	UpdatedAt    time.Time
	CreatedAt    time.Time `json:"created_at"`
	Subs         []Sub     `json:"subs"`
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {

	if len(os.Args) > 1 {
		err = InsertComment(os.Args[1], os.Args[2], os.Args[3])
		checkErr(err)
	}

	startTime := time.Now()

	fmt.Println(string(bytes), n)
	elapsed := time.Since(startTime)
	log.Printf("trace end: elapsed %f milliseconds\n", elapsed.Seconds()*1000.0)

}
