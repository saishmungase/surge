# Surge ⚡

A lightning-fast, concurrent API load testing and benchmarking CLI tool built in Go. 

Surge is designed to test the limits of your distributed systems and APIs. It bypasses the single-threaded limitations of traditional scripts by leveraging Go's Goroutines to hammer your endpoints with highly concurrent HTTP traffic, calculating precise metrics like p99 latency and Requests Per Second (RPS).

## 🚀 Installation

### Option 1: Using Go (Recommended)
If you already have Go installed on your machine, you can install Surge globally in one command:
```bash
go install [github.com/saishmungase/surge@latest](https://github.com/saishmungase/surge@latest)
```

### Option 2: Pre-compiled Binaries
Don't have Go installed? No problem. 
Head over to the [Releases page](https://github.com/saishmungase/surge/releases) and download the executable for your operating system (Windows, Mac, or Linux). Extract it and run it directly!

## 🛠️ Usage

Surge uses a simple `attack` command to start the load test.

**Basic GET Request Load Test:**
Hit an API with 500 requests, keeping 50 requests in-flight concurrently.
```bash
surge attack --url [https://api.github.com/zen](https://api.github.com/zen) --requests 500 --concurrency 50
```

**Testing POST Routes with JSON Payloads:**
Surge fully supports custom HTTP methods and JSON bodies for testing complex endpoints.
```bash
surge attack --url [https://api.example.com/data](https://api.example.com/data) \
  --method POST \
  --body '{"title": "test", "userId": 1}' \
  --requests 100 \
  --concurrency 10
```

## 📊 Metrics Output

Once the attack completes, Surge provides a clean terminal dashboard with critical network metrics:
* **Total Time:** The total duration of the test.
* **Requests/sec (RPS):** The exact throughput of your server.
* **Status Codes:** A breakdown of successful (2xx) vs failed (5xx) requests.
* **Latency Profiling:** Identifies your fastest request, p50 (median) latency, p99 latency, and the slowest request to help spot performance bottlenecks.

---
*Built to test backend infrastructure and APIs.*
