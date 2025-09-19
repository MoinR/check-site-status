# 🌍 Go Site Monitor

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-active-success.svg)]()

**Go Site Monitor** is a lightweight, open-source tool to monitor website availability in real time.
It reads a list of URLs, checks their status every 3 minutes, and updates a live dark-themed dashboard via WebSockets.

---

## ✨ Features

* 🔄 Monitor multiple URLs from a `sites.txt` file
* ⏱️ Automatic health checks every 3 minutes (configurable)
* 📡 Real-time updates via WebSockets
* 🌑 Dark-themed HTML dashboard (`#4a5568` background)
* ✅ Simple Go backend + HTML/JS frontend

---

## 🚀 Quick Start

### 1. Clone the repo

```bash
git clone https://github.com/<your-username>/go-site-monitor.git
cd go-site-monitor
```

### 2. Install dependencies

```bash
go mod tidy
go get github.com/gorilla/websocket
```

### 3. Add URLs to monitor

Create a `sites.txt` file with one URL per line:

```
https://google.com
https://github.com
https://openai.com
```

### 4. Run the server

```bash
go run main.go
```

Server starts at `:8080`.

### 5. Open the dashboard

Just open `index.html` in your browser → it connects to `ws://localhost:8080/ws` and shows live updates.

---

## 📷 Screenshot

<img width="1893" height="442" alt="image" src="https://github.com/user-attachments/assets/eef47c1c-18a3-431e-83e9-2ce0d98bcda8" />


---

## 🤝 Contributing

Contributions are welcome! 🎉

* Fork the repo
* Create a new branch (`feature/my-feature`)
* Commit changes
* Open a Pull Request

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).

---

