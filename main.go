package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type SiteStatus struct {
	URL    string `json:"url"`
	Status string `json:"status"`
	Time   string `json:"time"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.RWMutex
)

func checkURL(url string) string {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "DOWN"
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return "UP"
	}
	return fmt.Sprintf("STATUS %d", resp.StatusCode)
}

func readURLs(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}
	return urls, scanner.Err()
}

// marshalJSON is a helper function for testing
func marshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func broadcastStatus(urls []string) {
	for {
		var results []SiteStatus
		for _, url := range urls {
			results = append(results, SiteStatus{
				URL:    url,
				Status: checkURL(url),
				Time:   time.Now().Format(time.RFC3339),
			})
		}

		data, err := json.Marshal(results)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		log.Printf("Broadcasting status update to %d clients", len(clients))
		
		clientsMutex.RLock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Error writing to client: %v", err)
				// Remove failed client
				clientsMutex.RUnlock()
				clientsMutex.Lock()
				delete(clients, client)
				client.Close()
				clientsMutex.Unlock()
				clientsMutex.RLock()
			}
		}
		clientsMutex.RUnlock()

		time.Sleep(10 * time.Second)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer func() {
		clientsMutex.Lock()
		delete(clients, ws)
		clientsMutex.Unlock()
		ws.Close()
		log.Println("Client disconnected")
	}()

	clientsMutex.Lock()
	clients[ws] = true
	clientsMutex.Unlock()
	log.Printf("New client connected. Total clients: %d", len(clients))

	// Keep connection alive and listen for close
	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			break
		}
	}
}

func main() {
	urls, err := readURLs("sites.txt")
	if err != nil {
		log.Fatalf("Error reading URLs: %v", err)
	}

	if len(urls) == 0 {
		log.Fatal("No URLs found in sites.txt")
	}

	log.Printf("Monitoring %d URLs", len(urls))

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	go broadcastStatus(urls)

	log.Println("🚀 Server started at :8080")
	log.Println("📊 Open http://localhost:8080 in your browser")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
