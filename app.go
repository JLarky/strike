package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/JLarky/goReactServerComponents/internal/routes"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}

	r := routes.NewRouter()

	fmt.Printf("Server starting on http://localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
