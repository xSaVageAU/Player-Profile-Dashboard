package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"playerprofile/internal/config"
	"playerprofile/internal/handlers"
	"playerprofile/internal/nats"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Printf("Warning: config.yaml not found, using defaults or config.example.yaml: %v", err)
		cfg, err = config.LoadConfig("config.example.yaml")
		if err != nil {
			log.Fatalf("Fatal: Could not load any config file: %v", err)
		}
	}

	// 2. Connect to NATS
	natsClient, err := nats.Connect(cfg.NATS)
	if err != nil {
		log.Printf("Warning: NATS connection failed: %v", err)
	} else {
		defer natsClient.Close()
	}

	// 3. Parse templates
	tmpl, err := template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// 4. Static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 5. Routes
	http.HandleFunc("/", handlers.ProfileHandler(tmpl))

	// 6. Start Server
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on http://localhost%s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
