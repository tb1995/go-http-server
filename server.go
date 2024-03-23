package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port              string
	PublicDirectory   string
	NotFoundPagePath  string
	DefaultIndexPath string
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Println("Error loading configuration: ", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRequest(cfg))

	err = http.ListenAndServe(":"+cfg.Port, mux)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Println("Error starting server: ", err)
	}
}

func loadConfig() (*ServerConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	return &ServerConfig{
		Port:              os.Getenv("PORT"),
		PublicDirectory:   os.Getenv("PUBLIC_DIRECTORY_PATH"),
		NotFoundPagePath:  os.Getenv("NOT_FOUND_PAGE_PATH"),
		DefaultIndexPath: os.Getenv("DEFAULT_INDEX_PATH"),
	}, nil
}

func handleRequest(cfg *ServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        
		filePath := cfg.PublicDirectory + r.URL.Path

        if strings.Contains(r.URL.Path, "..") {
            fmt.Printf("🐟y behaviour")
            filePath = cfg.PublicDirectory + cfg.NotFoundPagePath
        } else if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			filePath = cfg.PublicDirectory + cfg.NotFoundPagePath
		} else if r.URL.Path == "/" {
			filePath = cfg.PublicDirectory + cfg.DefaultIndexPath
		}
		content, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error serving file:", err)
			return
		}
		fmt.Fprint(w, string(content))
	}
}
