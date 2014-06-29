package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	ip := "0.0.0.0"
	port := os.Getenv("port")
	if len(port) == 0 {
		port = "8000"
	}
	http.Handle("/", http.FileServer(http.Dir("./")))
	fmt.Printf("Listening on %s:%s\n", ip, port)
	http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil)
}
