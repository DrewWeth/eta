package commentcontroller

import (
	"encoding/json"
	"errors"
	"github.com/drewweth/eta/src/app/controllers/helpers"
	"github.com/drewweth/eta/src/app/models"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"github.com/gorilla/schema"
)

var querier *models.Querier
var decoder *schema.Decoder

func Init(_querier *models.Querier) {
	querier = _querier
	decoder = schema.NewDecoder()
}

func parseCommentParams(w http.ResponseWriter, req *http.Request) (models.AppDataReqParams, error) {
	var params models.AppDataReqParams
	err := req.ParseForm()
	if err != nil{
		return params, err
	}
	err = decoder.Decode(&params, req.PostForm)
	// err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		return params, err
	}

	log.Println(params)

	if params.PostID == "" || params.ParentID == "" || params.Comment == "" || params.APIToken == "" {
		friendlyErr := errors.New("Some params cannot be blank: ")
		return params, friendlyErr
	}

	return params, nil
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// MassInsertHandler adds 1000 comments to a post given the postID
func MassInsertHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// subName := vars["sub"]
	id := vars["id"]
	log.Println("ID", id)

	rand.Seed(time.Now().Unix())

	comments, _, _, err := querier.GetCommentsByPostID(id)
	if err != nil {
		log.Fatal(err)
	}
	if len(comments) > 0 {
		for i := 0; i < 1000; i++ {
			n := random(0, len(comments))
			log.Println(n)
			err := querier.InsertComment(id, comments[n].ID, "Test Data "+strconv.Itoa(i))
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		err := querier.InsertComment(id, id, "Root Test Data")
		if err != nil {
			log.Fatal(err)
		}
	}
}

// CreateCache creates an HTML representation of a post and stores it
func CreateCache(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// subName := vars["sub"]
	id := vars["id"]

	htmlBuffer, _, err := querier.RenderCommentHTML(id)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = querier.InsertCachedComment(id, htmlBuffer.String())
	if err != nil {
		log.Fatal(err)
	}
}

// PostHandler ccreates a new comment given a bunch of params
func PostHandler(w http.ResponseWriter, req *http.Request) {
	params, err := parseCommentParams(w, req)
	if err != nil {
		friendlyErr := errors.New("Problem parsing parameters. (" + err.Error() + ")")
		helpers.SendError(http.StatusBadRequest, friendlyErr, w)
		return
	}

	err = querier.InsertComment(params.PostID, params.ParentID, params.Comment)
	if err != nil {
		friendlyErr := errors.New("Insert error. (" + err.Error() + ")")
		helpers.SendError(http.StatusBadRequest, friendlyErr, w)
		return
	}

	resp := models.GenericResponse{200, "Comment successful"}
	json.NewEncoder(w).Encode(resp)
}

func UpvoteCommentHandler(w http.ResponseWriter, req *http.Request){

}

func DownvoteCommentHandler(w http.ResponseWriter, req *http.Request){

}
