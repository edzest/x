package main

import (
	"flag"
	"log"
	"os"

	"edzest.org/x/internal/server"
)

var portFlag = flag.String("port", "8080", "port to listen")

func main() {
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = *portFlag
	}

	addr := ":" + port

	srv := server.NewHTTPServer(addr)
	log.Printf("Starting server on address %s", addr)
	log.Fatal(srv.ListenAndServe())
}
