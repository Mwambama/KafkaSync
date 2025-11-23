package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
)

// Reusing existing config structs
type tomlConfig struct {
	Database databaseConfig `toml:"database"`
}

type databaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

// data structure for API response
type DownloadRecord struct {
	ID             int    `json:"id"`
	Filename       string `json:"filename"`
	RemoteLocation string `json:"remote_location"`
	Hash           string `json:"hash"`
	Status         string `json:"status"`
	DownloadedAt   string `json:"downloaded_at"`
}

var db *sql.DB

func init() {
	// Load the same config file for consumer
	var conf tomlConfig
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Connect to DB
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.DbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("❌ Database unreachable: %v", err)
	}
	log.Println("✅ API Server connected to Database")
}

// enableCORS allows the React app (on port 5173) to call this API (on port 8080)
func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func getDownloads(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w) // Enable access for React

	rows, err := db.Query("SELECT id, filename, remote_location, hash, status, downloaded_at FROM downloads ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var downloads []DownloadRecord
	for rows.Next() {
		var d DownloadRecord
		// Scan timestamps as strings for simplicity
		if err := rows.Scan(&d.ID, &d.Filename, &d.RemoteLocation, &d.Hash, &d.Status, &d.DownloadedAt); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		downloads = append(downloads, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(downloads)
}

func main() {
	http.HandleFunc("/api/downloads", getDownloads)

	log.Println("🚀 API Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
