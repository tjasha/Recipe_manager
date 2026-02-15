package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	//defined port
	port := "8080"

	// Test API endpoint
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello Test test"})
	})

	// Serve React files
	distDir := "../frontend/dist"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		path := filepath.Join(distDir, r.URL.Path)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			http.ServeFile(w, r, path)
			return
		}

		// if other file doesn't exist, open index.html
		http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
	})

	log.Printf("Server running at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
