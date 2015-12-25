package homecontroller

import (
	// "fmt"
	"github.com/drewweth/eta/src/app/models"
	"html/template"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Subs []*models.Sub
}

var querier *models.Querier

func Init(_querier *models.Querier) {
	querier = _querier
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	subs, err := querier.GetSubs()
	if err != nil{
		log.Println(err)
	}
	log.Println("Sub count:", len(subs))
	page := Page{"Home", subs}

	t, err := template.New("index").ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, page)
	if err != nil {
		log.Fatal(err)
	}
}
