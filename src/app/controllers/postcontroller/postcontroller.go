package postcontroller

import (
	// "bytes"
	"fmt"
	"github.com/drewweth/eta/src/app/models"
"errors"
"encoding/json"
"github.com/drewweth/eta/src/app/controllers/helpers"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/schema"
)

type Page struct {
	Post					*models.Post
	SubName       string
	ID            string
	CommentCount int
	CommentHTML   template.HTML
	LinkIsNil			bool
}

var querier *models.Querier
var decoder *schema.Decoder


func Init(_querier *models.Querier) {
	querier = _querier
	decoder = schema.NewDecoder()

}

var funcMap = template.FuncMap{
	"renderComment": models.RenderComment,
	"markdown":      markDowner,
}

func markDowner(args ...interface{}) template.HTML {
	s := []byte(fmt.Sprintf("%s", args...))
	return template.HTML(s)
}

func ShowHandler(w http.ResponseWriter, req *http.Request) {
	// var pageTitle string
	var pageSubName string
	var pageID string
	var pageN int
	var pageHTML string
	var linkIsNil bool

	vars := mux.Vars(req)
	subName := vars["sub"]
	id := vars["id"]

	post, err := querier.GetPost(id)
	if err != nil {
		log.Println(err)
	}
	log.Println(post)

	htmlString, err := querier.GetCachedComments(id)
	if err != nil {
		log.Println(err)
	}

	if htmlString != "" { // If you have cached data
		log.Println("You have cached data")
		pageN = 0
		pageHTML = htmlString
	} else { // If you don't have cached data
		log.Println("You don't have cached data")
		buffer, n, err := querier.RenderCommentHTML(id)
		if err != nil {
			log.Println(err)
			// Do something here, error prone
		}
		pageN = n
		pageHTML = buffer.String()
	}

	pageSubName = subName
	pageID = id

	if post.Link != ""{
		linkIsNil = false
	}else{
		linkIsNil = true
	}

	page := Page{post, pageSubName, pageID, pageN, template.HTML(pageHTML), linkIsNil}

	startTime := time.Now()
	t, err := template.New("post").Funcs(funcMap).ParseFiles("post.html", "jquery.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(startTime)
	log.Printf("Rendering: elapsed %f milliseconds\n", elapsed.Seconds()*1000.0)
}

func parsePostParams(w http.ResponseWriter, req *http.Request)(models.PostReqParams, error){
	var params models.PostReqParams
	err := req.ParseForm()
	if err != nil {
		return params, err
	}
	err = decoder.Decode(&params, req.PostForm)
	// err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		return params, err
	}

	log.Println(params)

	if params.SubName == "" || params.Title == "" {
		friendlyErr := errors.New("Some params cannot be blank.")
		return params, friendlyErr
	}

	return params, nil
}

func InsertPostHandler(w http.ResponseWriter, req *http.Request){
	params, err := parsePostParams(w, req)

	if err != nil {
		friendlyErr := errors.New("Problem parsing parameters. (" + err.Error() + ")")
		helpers.SendError(http.StatusBadRequest, friendlyErr, w)
		return
	}
	err = querier.InsertPost(params.SubName, params.Title, params.Comment, params.Link)
	if err != nil {
		friendlyErr := errors.New("Insert error. (" + err.Error() + ")")
		helpers.SendError(http.StatusBadRequest, friendlyErr, w)
		return
	}

	resp := models.GenericResponse{200, "Post successful"}
	json.NewEncoder(w).Encode(resp)

}


func UpvotePostHandler(w http.ResponseWriter, req *http.Request){

}

func DownvotePostHandler(w http.ResponseWriter, req *http.Request){

}
