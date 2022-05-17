package main

import (
	"fmt"
	"net/http"
)

func hello(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "Yo!")
}
func main() {
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8080", nil)
}
