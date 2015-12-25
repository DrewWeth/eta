package models

import (
	// "encoding/json"
	"bytes"
	"github.com/gocql/gocql"
	"log"
	"sort"
	"strconv"
	"time"
	"html/template"
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

type UpvoteSorter []*Comment

func (a UpvoteSorter) Len() int           { return len(a) }
func (a UpvoteSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UpvoteSorter) Less(i, j int) bool { return !(a[i].Upvotes < a[j].Upvotes) }

func (querier *Querier) InsertComment(postID, parentID, content string) error {
	if err := querier.Session.Query(`INSERT INTO comments (id, post_id, parent_id, content, upvotes, downvotes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
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
	// err := querier.UpCommentCount(postID)
	// if err !=nil{
	// 	log.Println("UpCommentCount encountered error on id:", postID)
	// 	return err
	// }
	return nil
}

func (querier *Querier) UpCommentCount(postID string) error{
	if err := querier.Session.Query(`UPDATE posts SET comment_count = comment_count + 1 WHERE id = ?`,
		postID).Exec(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (querier *Querier) GetComments(postID string) ([]*Comment, int, error) {
	var roots []*Comment

	comments, commentLookup, childrenLookup, err := querier.GetCommentsByPostID(postID)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	for _, node := range comments { // for all comments
		if node.ParentID == postID {
			roots = append(roots, node)
		}
		for _, childID := range childrenLookup[node.ID] { // get the ID's of their children
			node.Children = append(node.Children, commentLookup[childID]) // add the child to the list via lookup
		}
	}
	// log.Println("Roots", roots)

	// bytes, _ := json.MarshalIndent(roots, "", "\t") //formated output
	//bytes, _ := json.Marshal(roots)
	return roots, len(comments), nil
}

func (querier Querier) GetCachedComments(postIDToQuery string) (string, error) {
	var id gocql.UUID
	var content string
	var upvotes int
	var downvotes int
	var createdAt time.Time
	var updatedAt time.Time

	if err := querier.Session.Query(`SELECT id, content, upvotes, downvotes, created_at, updated_at FROM caches WHERE post_id = ? LIMIT 1`,
		postIDToQuery).Consistency(gocql.One).Scan(&id, &content, &upvotes, &downvotes, &createdAt, &updatedAt); err != nil {
		log.Println(err)
		return "", err
	}
	return content, nil
}

func (querier Querier) InsertCachedComment(id, htmlString string) error {
	if err := querier.Session.Query(`INSERT INTO caches (id, post_id, content) VALUES (?, ?, ?)`,
		gocql.TimeUUID(),
		id,
		htmlString).Exec(); err != nil {
		return err
	}
	return nil
}

func (querier Querier) RenderCommentHTML(id string) (bytes.Buffer, int, error) {
	roots, n, err := querier.GetComments(id)
	if err != nil {
		log.Println(err)
		var emptyBuffer bytes.Buffer
		return emptyBuffer, 0, err
	}

	var buffer bytes.Buffer
	startTime := time.Now()

	for _, node := range roots {
		RenderComment(node, &buffer, 0)
	}
	elapsed := time.Since(startTime)
	log.Printf("Recursive tree rendering: elapsed %f milliseconds\n", elapsed.Seconds() * 1000.0)
	return buffer, n, nil
}

func RenderComment(comment *Comment, buffer *bytes.Buffer, depth int) *bytes.Buffer {
	buffer.WriteString("<div class='comment' id='" + comment.ID + "'><div><b>Username</b> "+strconv.Itoa(comment.Upvotes)+"</div><div>" + template.HTMLEscapeString(comment.Content) + "</div><div><a class='reply'>Reply</a></div></div>")

	if len(comment.Children) == 0 {
		var buff bytes.Buffer
		return &buff
	}

	for i, node := range comment.Children {
		if i < 6 {
			if i == 0 {
				if depth%2 == 1 {
					buffer.WriteString("<div style='margin-left:30px;background:#eee'>")
				} else {
					buffer.WriteString("<div style='margin-left:30px;background-color: rgb(247, 247, 248);'>")
				}
			}
			RenderComment(node, buffer, depth+1)

		}else if i == 6{
			buffer.WriteString("<div><a href=''>"+ strconv.Itoa(len(comment.Children) - i)+" more comments</a>")
		}
	}
	buffer.WriteString("</div>")

	return buffer
}

func (querier Querier) GetCommentsByPostID(postIDToQuery string) ([]*Comment, map[string]*Comment, map[string][]string, error) {
	var id gocql.UUID
	var idAsString string
	var content string
	var parentID string
	var upvotes int
	var downvotes int
	var updatedAt time.Time

	commentLookup := make(map[string]*Comment)
	childrenLookup := make(map[string][]string)
	var comments []*Comment

	startTime := time.Now()
	iter := querier.Session.Query(`SELECT id, parent_id, content, upvotes, downvotes, updated_at FROM comments WHERE post_id = ?`, postIDToQuery).Iter()
	elapsed := time.Since(startTime)
	log.Printf("Querying: elapsed %f milliseconds\n", elapsed.Seconds()*1000.0)

	startTime = time.Now() // Timing building the tree
	for iter.Scan(&id, &parentID, &content, &upvotes, &downvotes, &updatedAt) {
		idAsString = id.String()
		comment := &Comment{idAsString, parentID, content, upvotes, downvotes, time.Now(), updatedAt, nil}
		comments = append(comments, comment)
		commentLookup[idAsString] = comment
	}

	sort.Sort(UpvoteSorter(comments))

	for _, node := range comments {
		childrenLookup[node.ParentID] = append(childrenLookup[node.ParentID], node.ID)
	}
	elapsed = time.Since(startTime)
	log.Printf("Building the tree: elapsed %f milliseconds\n", elapsed.Seconds()*1000.0)

	if err := iter.Close(); err != nil {
		return nil, nil, nil, err
	}

	return comments, commentLookup, childrenLookup, nil
}
