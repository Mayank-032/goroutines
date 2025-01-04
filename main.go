package main

import (
	"fmt"
	"lld/router"
	"log"
	"net/http"
)

func main() {
	r := http.NewServeMux()
	router.InitRoutes(r)

	var port = ":8000"
	fmt.Println("Starting to listen on PORT ", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Error: %v\n. server shutdown gracefully", err.Error())
		return
	}
}
