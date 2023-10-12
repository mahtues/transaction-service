package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("starting avocado application...")

	mux := http.NewServeMux()

	if err := http.ListenAndServe(":8000", mux); err != nil {
		fmt.Println("failed to start server")
	}
}
