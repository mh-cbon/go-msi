package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello, world\n")
}

func main() {

	fmt.Println("Starting....")

	// public views
	http.HandleFunc("/", HandleIndex)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
