package main

import (
	_ "embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"time"
)

type SiteSource struct {
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Labels   []string `json:"labels"`
	Language string   `json:"language"`
}

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Image       Image  `xml:"image"`
	Items       []Item `xml:"item"`
}

type Image struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	URL   string `xml:"url"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

//go:embed source.json
var sourceJSON string

// FetchRSSFeed fetches an RSS feed from a URL.
func FetchRSSFeed(url string) (*RSS, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	//read body
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch RSS feed: %s", resp.Status)
	}

	var feed RSS
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&feed)
	if err != nil {
		return nil, fmt.Errorf("failed to decode RSS feed: %w", err)
	}

	return &feed, nil
}

func main() {
	fmt.Println("Loading RSS sources...")

	var sources []SiteSource
	err := json.Unmarshal([]byte(sourceJSON), &sources)
	if err != nil {
		fmt.Printf("Error loading sources: %v\n", err)
		return
	}

	// Compute Week number
	currentTime := time.Now()
	year, week := currentTime.ISOWeek()
	fmt.Printf("Current Year: %d, Week: %d\n", year, week)

	// Open week number html file
	fmt.Printf("Opening week number HTML file for year %d, week %d...\n", year, week)

	fp, err := os.Create(fmt.Sprintf("week_%d.html", week))
	if err != nil {
		fmt.Printf("Error opening week number HTML file: %v\n", err)
		return
	}
	defer fp.Close()

	fp.WriteString(fmt.Sprintf("<html>\n<body>\n<h1>Week %d of Year %d</h1>\n", week, year))
	fp.WriteString("<ul>\n")

	for _, source := range sources {
		fmt.Printf("Loaded RSS source: %s (%s)\n", source.Name, source.URL)

		// Fetch RSS feed and process it

		feed, err := FetchRSSFeed(source.URL)
		if err != nil {
			fmt.Printf("Error fetching RSS feed for %s: %v\n", source.Name, err)
			continue
		}

		fmt.Printf("Fetched RSS feed: %s\n", feed.Channel.Title)
		for _, item := range feed.Channel.Items {
			// Check if pubDate is not older than 1 week
			pubDate, err := time.Parse(time.RFC1123, item.PubDate)
			if err == nil && time.Since(pubDate) > 7*24*time.Hour {
				continue
			}
			fp.WriteString(fmt.Sprintf(`<li><a href="%s">%s</a></li>`, item.Link, item.Title))
			fp.WriteString("\n")
		}
	}
	fp.WriteString("\n</ul></body>\n</html>")
}
