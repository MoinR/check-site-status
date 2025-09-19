package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

var clients = make(map[*websocket.Conn]bool)

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
		urls = append(urls, scanner.Text())
	}
	return urls, scanner.Err()
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

		data, _ := json.Marshal(results)
		fmt.Println(string(data))
		for client := range clients {
			client.WriteMessage(websocket.TextMessage, data)
		}

		time.Sleep(10 * time.Second)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	select {}
}

func main() {
	urls, err := readURLs("sites.txt")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/ws", handleConnections)

	go broadcastStatus(urls)

	fmt.Println("🚀 Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
