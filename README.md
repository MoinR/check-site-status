# 🌍 Go Site Monitor

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![Test Coverage](https://img.shields.io/badge/coverage-40%25-yellow.svg)]()

**Go Site Monitor** is a lightweight, open-source tool to monitor website availability in real time.
It reads a list of URLs, checks their status every 10 seconds, and updates a beautiful live dashboard via WebSockets.

---

## ✨ Features

* 🔄 Monitor multiple URLs from a `sites.txt` file
* ⏱️ Automatic health checks every 10 seconds
* 📡 Real-time updates via WebSockets
* 🎨 Modern gradient-themed dashboard with statistics
* 🔌 Auto-reconnection on disconnect (up to 10 retries)
* 📊 Live statistics: Total, Online, Offline, Warning counts
* 🎯 Status badges: UP (green), DOWN (red), WARNING (yellow)
* 📱 Fully responsive design for mobile and desktop
* 🚀 Zero external dependencies in frontend (pure vanilla JS)
* 🔒 Thread-safe client management
* 💬 Comment support in sites.txt (lines starting with #)
* ✅ Comprehensive test coverage (40%)

---

## 🚀 Quick Start

### 1. Clone the repo

```bash
git clone https://github.com/MoinR/check-site-status.git
cd check-site-status
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Add URLs to monitor

Edit `sites.txt` file with one URL per line:

```
# Example sites to monitor
https://google.com
https://github.com
https://openai.com

# Lines starting with # are comments
```

### 4. Run the server

```bash
go run main.go
```

Server starts at `:8080`.

### 5. Open the dashboard

Open your browser and navigate to **http://localhost:8080**

The dashboard will automatically connect via WebSocket and display live status updates!

---

## 📷 Screenshot

<img width="1280" height="600" alt="Modern Site Monitor Dashboard" src="https://github.com/user-attachments/assets/03ba03ef-e732-4189-ba32-65afabc5ca8d" />

---

## 🧪 Testing

Run the test suite:

```bash
go test -v
```

Check test coverage:

```bash
go test -cover
```

---

## 🏗️ Architecture

- **Backend**: Go with Gorilla WebSocket library
- **Frontend**: Pure HTML/CSS/JavaScript (no frameworks)
- **Communication**: WebSocket for real-time updates
- **Deployment**: Single binary, no database required

---

## 🤝 Contributing

Contributions are welcome! 🎉

* Fork the repo
* Create a new branch (`feature/my-feature`)
* Commit changes
* Run tests: `go test`
* Open a Pull Request

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).

---

