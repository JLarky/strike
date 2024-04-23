package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/JLarky/strike/internal/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}

	r := routes.NewRouter()

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, r)
}
