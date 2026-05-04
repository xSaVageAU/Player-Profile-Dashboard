package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"playerprofile/internal/handlers"
)

func main() {
	// Parse templates
	// Using a layout pattern: layout.html + partials
	tmpl, err := template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", handlers.ProfileHandler(tmpl))

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
