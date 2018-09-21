package main

import (
	"fmt"
	"net/http"

	rc "./robchess"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hey there")
}

func main() {
	var posTest rc.Position

	fmt.Println(posTest)
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)
}
