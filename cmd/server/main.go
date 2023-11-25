package main

import (
	"log"
	"os"

	"edzest.org/x/internal/server"
)

func main() {
	var port string
	if os.Getenv("PORT") == "" {
		port = ":8080"
	} else {
		port = ":" + os.Getenv("PORT")
	}
	srv := server.NewHTTPServer(port)
	log.Fatal(srv.ListenAndServe())
}
