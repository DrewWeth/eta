// create keyspace eta with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
// create table eta.comments(id UUID, post_id text, content text, upvotes int, downvotes int, created_at timestamp, updated_at timestamp, PRIMARY KEY(id));
// create index on eta.comments(post_id);

package main

import (
	"os"
	"encoding/json"
	"sort"
	"time"
	"fmt"
	"github.com/gocql/gocql"
	"log"
)

var session *gocql.Session

type Comment struct { 
        ID       string  `json:"id"` 
        ParentID string  `json:"-"` 
        Content     string  `json:"name"` 
	Upvotes		int	`json:"upvotes"`
	Downvotes	int 	`json:"downvotes"`
	CreatedAt	time.Time
	UpdatedAt	time.Time	`json:"updated_at"`
        Children []*Comment `json:"children"` 
}

type Sub struct {
	ID	string
	Name      string `json:"name"`
	Subscribers int `json:"subscribers"`
	UpdatedAt	time.Time
	CreatedAt	time.Time
}

type Post struct {
	ID		string	`json:"title"`
	Title		string `json:"title"`
	Link		string	`json:"link"`
	Comment      	string `json:"content"`
	PostToken  	string `json:"-"`
	Sub 		string `json:"sub"`
	Upvotes 	int	`json:"upvotes"`
	Downvotes 	int	`json:"downvotes"`
	UpdatedAt	time.Time
	CreatedAt	time.Time
}

type User struct {
	ID           string        `json:"id"`
	Username     string        `json:"username"`
	Email		string		`json:"email"`
	PasswordHash string        `json:"-"`
	APIToken     string        `json:"api_token"`
	UpdatedAt	time.Time
	CreatedAt    time.Time     `json:"created_at"`
	Subs []Sub `json:"subs"`
}

type UpvoteSorter []*Comment
func (a UpvoteSorter) Len() int           { return len(a) }
func (a UpvoteSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UpvoteSorter) Less(i, j int) bool { return !(a[i].Upvotes < a[j].Upvotes) }

func GetComments(postID string) ([]byte, int, error) {
	var roots []*Comment
	
	comments, commentLookup, childrenLookup, err := getCommentsByPostID(postID)
	if err != nil{
		fmt.Println(err)
		return nil, 0, err
	}
	
	for _, node := range comments { // for all comments
		if node.ParentID == postID{
			roots = append(roots, node)
		}
		for _, childID := range childrenLookup[node.ID] { // get the ID's of their children
			node.Children = append(node.Children, commentLookup[childID]) // add the child to the list via lookup
		}
	}

	log.Println("Root", roots)

	bytes, _ := json.MarshalIndent(roots, "", "\t") //formated output
	//bytes, _ := json.Marshal(roots)
	return bytes, len(comments), nil
}

func getCommentsByPostID(postIDToQuery string)([]*Comment, map[string]*Comment, map[string][]string, error){
	var id gocql.UUID
	var content string
	var parentID string
	var upvotes int
	var downvotes int
	var updatedAt time.Time

	commentLookup := make(map[string]*Comment)
	childrenLookup := make(map[string][]string)
	var comments []*Comment	

	iter := session.Query(`SELECT id, parent_id, content, upvotes, downvotes, updated_at FROM comments WHERE post_id = ?`, postIDToQuery).Iter()
	for iter.Scan(&id, &parentID, &content, &upvotes, &downvotes, &updatedAt) {
		comment := &Comment{id.String(), parentID, content, upvotes, downvotes, time.Now(), updatedAt, nil}
		comments = append(comments, comment)
		commentLookup[id.String()] = comment
	}

	sort.Sort(UpvoteSorter(comments))

	for _, node := range comments {
		childrenLookup[node.ParentID] = append(childrenLookup[node.ParentID], node.ID)
	}

	if err := iter.Close(); err != nil {
		return nil, nil, nil, err
	}

	log.Println(commentLookup)
	log.Println(childrenLookup)

	return comments, commentLookup, childrenLookup, nil
}

func InsertComment(postID, parentID, content string) error{
	if err := session.Query(`INSERT INTO comments (id, post_id, parent_id, content, upvotes, downvotes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		gocql.TimeUUID(),
		postID,
		parentID,
		content,
		0,
		0,
		time.Now(),
		time.Now()).Exec(); err != nil {
		return err
	}
	return nil
}

func checkErr(err error){
	if err != nil{
		log.Println(err)
	}
}


func main() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "eta"
	cluster.ProtoVersion = 4

	log.Println(cluster)
	_session, err := cluster.CreateSession()
	session = _session
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	
	
	if len(os.Args) > 1 {
		err = InsertComment(os.Args[1], os.Args[2], os.Args[3])
		checkErr(err)
	} 

	startTime := time.Now()
	bytes, n, err := GetComments("abc")	
	if err != nil{
		log.Println(err)
		return
	}
	fmt.Println(string(bytes), n)
	elapsed := time.Since(startTime)
	log.Printf("trace end: elapsed %f milliseconds\n", elapsed.Seconds() * 1000.0)

}
