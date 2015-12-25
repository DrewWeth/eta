package main

import (
	"log"

	"github.com/codegangsta/negroni"
	"github.com/drewweth/eta/src/app/controllers/commentcontroller"
	"github.com/drewweth/eta/src/app/controllers/homecontroller"
	"github.com/drewweth/eta/src/app/controllers/postcontroller"
	"github.com/drewweth/eta/src/app/controllers/subcontroller"
	"github.com/drewweth/eta/src/app/models"
	"github.com/drewweth/eta/src/server"
	"github.com/gocql/gocql"
)

func connectToCassandra() *gocql.Session {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "eta"
	cluster.ProtoVersion = 4

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	return session
}

func main() {
	log.Println("Starting server...")

	session := connectToCassandra()
	querier := &models.Querier{session}
	defer session.Close()

	postcontroller.Init(querier)
	homecontroller.Init(querier)
	commentcontroller.Init(querier)
	subcontroller.Init(querier)

	router := server.NewServer()
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
