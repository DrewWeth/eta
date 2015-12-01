package main

import (
	"log"

	"github.com/codegangsta/negroni"
	"github.com/drewweth/eta/src/server"
)

func main() {
	log.Println("Starting server...")

	router := server.NewServer()

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
