package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/buraksekili/broadcast/handlers"
)

func main() {
	flag.Parse()

	a, err := handlers.NewHTTPAgent()
	if err != nil {
		log.Fatal(err)
	}

	m := http.NewServeMux()
	m.HandleFunc("/", a.Broadcast)

	log.Fatal(http.ListenAndServe(":8932", m))
}
