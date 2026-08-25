package main

import (
	"fmt"
	"net/http"
	"log"
)

func main() {
	http.HandleFunc("/",uploadHandler)
	http.HandleFunc("/upload",uploadStreamingHandler)
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}