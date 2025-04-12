# 🕷️ Go Website Crawler CLI Tool

A command-line tool written in Go that accepts a website URL, checks its status, and recursively crawls the website to generate a site map. The tool supports exporting the site map in JSON (and optionally XML) format.

---

## 📌 Overview

This tool allows you to:

- Check if a website is reachable and view the HTTP status code and response time.
- Crawl all internal links within the same domain.
- Respect basic `robots.txt` rules.
- Avoid duplicate visits to the same URL.
- Limit crawl depth.
- Control request rate (rate limiting).
- Export the resulting site map as a JSON or XML file.

---

## ✨ Features

- ✅ Accept a URL as a command-line argument.
- 🔍 Perform an HTTP GET request to the URL and return:
  - HTTP status code (e.g., 200, 404, 500)
  - Response time
- 🧭 Recursively crawl the website:
  - Follow only internal links (same domain)
  - Avoid revisiting the same URL
  - Respect `robots.txt`
- 🌳 Generate and print a site map in a tree-like structure or grouped by depth.
- 📤 Export site map:
  - JSON format (`--json`)
  - XML format (`--xml`) *(optional)*
- 🚦 Support rate limiting between requests.
- 🎯 Control crawl depth using flags.
- 🧵 Use goroutines for concurrent crawling.

---

## 🛠️ Tech Stack

- **Language:** Go
- **Core Libraries:**
  - `net/http`
  - `golang.org/x/net/html`
- **CLI Libraries:**
  - [`cobra`](https://github.com/spf13/cobra) or [`urfave/cli`](https://github.com/urfave/cli)
- **Optional:**
  - `goquery` for HTML parsing

---

## 🚀 Getting Started

### Prerequisites

- [Go installed](https://golang.org/doc/install) (Go 1.18+ recommended)

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/go-website-crawler.git
cd go-website-crawler
```

### Run the Tool

```
go run main.go https://example.com
```

### Example Usage

```
go run main.go https://example.com
```

### 📤 Exporting Site Map

#### JSON

```
go run main.go https://example.com --json output.json
```

#### XML (Optional)

```
go run main.go https://example.com --xml output.xml
```

### ⚙️ Command-Line Options

```
Flag	Description
--------------------------

--json	Export the site map to a JSON file
--xml	Export the site map to an XML file (optional feature)
--depth	Limit the maximum crawl depth (e.g., --depth=2)
--rate	Apply rate limiting between requests (e.g., --rate=1s)
```

### 🧾 Sample Exported Site Map (JSON)

```json
{
  "site": "https://example.com",
  "status": "200 OK",
  "response_time_ms": 120,
  "site_map": {
    "/": {
      "children": {
        "/about": {},
        "/blog": {
          "children": {
            "/blog/post-1": {},
            "/blog/post-2": {}
          }
        },
        "/contact": {}
      }
    }
  }
}
```

### 📄 Sample robots.txt Support

The crawler will parse robots.txt and skip disallowed paths under User-agent: *. Advanced support (e.g., crawl-delay, per-agent rules) is not implemented.

### 🧪 Sample Exported Site Map (Tree View)

```
- https://example.com
  ├── /about
  ├── /blog
  │   ├── /blog/post-1
  │   └── /blog/post-2
  └── /contact
```