package main

import (
	"fmt"
	"net/http"

	"github.com/JLarky/goReactServerComponents/internal/routes"
)

func main() {
	r := routes.NewRouter()

	fmt.Println("Server starting on http://localhost:3333")
	http.ListenAndServe(":3333", r)
}
