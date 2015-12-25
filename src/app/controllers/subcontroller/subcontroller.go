package subcontroller

import (
	// "bytes"
	// "fmt"
	"github.com/drewweth/eta/src/app/models"
	"encoding/json"
	"github.com/drewweth/eta/src/app/controllers/helpers"
	"errors"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Page struct {
	Sub         string
	Posts				[]*models.Post
}

var querier *models.Querier

func Init(_querier *models.Querier) {
	querier = _querier
}

func ShowHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	subName := vars["sub"]

	posts, err := querier.GetPosts(subName)
	if err != nil {
		log.Println(err)
	}

	page := Page{subName, posts}

	startTime := time.Now()
	t, err := template.New("sub").ParseFiles("sub.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(startTime)
	log.Printf("Rendering: elapsed %f milliseconds\n", elapsed.Seconds() * 1000.0)
}

func InsertSubHandler(w http.ResponseWriter, req *http.Request){
	req.ParseForm()
	sub := req.PostFormValue("sub_name")
	err := querier.InsertSub(sub)
	if err != nil{
		friendlyErr := errors.New("Insert error. (" + err.Error() + ")")
		helpers.SendError(http.StatusBadRequest, friendlyErr, w)
		return
	}
	resp := models.GenericResponse{200, "Comment successful"}
	json.NewEncoder(w).Encode(resp)

}

func UpvotePostHandler(w http.ResponseWriter, req *http.Request){
}

func DownvotePostHandler(w http.ResponseWriter, req *http.Request){
}
