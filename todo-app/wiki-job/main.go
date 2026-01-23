package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Task struct {
	Title string `json:"title"`
	State int    `json:"state"`
}

func main() {
	wikiUrl := getWikiLink()
	data := Task{
		Title: "Read " + wikiUrl,
		State: 0,
	}
	jsonData, _ := json.Marshal(data)
	backendUrl := os.Getenv("BACKEND_URL")
	fallbackURL := "http://todo-backend-svc"

	if backendUrl == "" {
		log.Printf("env BACKEND_URL was unset, using default: %q", fallbackURL)
		backendUrl = fallbackURL
	}
	log.Print("Started Job")

	resp, err := http.Post(
		backendUrl+"/tasks", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to connect to backend: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Fatalf("Backend returned non-success code: %d", resp.StatusCode)
	}

	log.Printf("Successfully added todo. Status: %s", resp.Status)
}

func getWikiLink() string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://en.wikipedia.org/wiki/Special:Random", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	// Wikipedia requires a descriptive User-Agent
	req.Header.Set("User-Agent", "MyK3sTodoBot/1.0 (contact: your@email.com)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error getWikiLink:", err)
		return ""
	}
	defer resp.Body.Close()

	// Now it should follow the redirect to the final article URL
	return resp.Request.URL.String()
}
