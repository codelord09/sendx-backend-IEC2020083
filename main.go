package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"golang.org/x/net/html"
)

var ctx = context.Background()
var redisClient *redis.Client
var crawledPages = make(map[string]CrawledPage)

const maxRetries = 3
const retryDelay = time.Second * 5

type CrawledPage struct {
	Content   string
	Timestamp time.Time
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/crawl", CrawlHandler).Methods("POST")
	r.HandleFunc("/", IndexHandler).Methods("GET")

	http.Handle("/", r)

	// Initialize the Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis-16568.c301.ap-south-1-1.ec2.cloud.redislabs.com:16568",
		Password: "F7yfvRdrYk5mjJGh1oqoa0v4702uhuBj", // no password set
		DB:       0,                                  // use default DB
		// Protocol: 3, // Protocol is not used here
	})

	log.Fatal(http.ListenAndServe(":5500", nil))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func CrawlHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.FormValue("url")
	paying := r.FormValue("paying") == "true"
	speedMultiplier := 1.0

	if paying {
		speedMultiplier = 5.0
	}

	// Check if the page exists in Redis cache
	cachedPage, err := getPageFromCache(url)
	if err == nil {
		// Page found in cache, return it
		writeCrawlTimeToHTML(w, r, 0, cachedPage.Content, true)
		return
	}

	startTime := time.Now()

	content, err := crawlURL(url, speedMultiplier)
	if err != nil {
		log.Println("Error: ", err)
		writeErrorToHTML(w, r, startTime)
		return
	}

	// Store the newly crawled page in the Redis cache
	storePageInCache(url, content)

	crawlTime := time.Since(startTime).Seconds()

	writeCrawlTimeToHTML(w, r, crawlTime, content, true)
}

func getPageFromCache(url string) (CrawledPage, error) {
	// Try to get the page from Redis cache
	cacheKey := "page:" + url
	val, err := redisClient.Get(ctx, cacheKey).Result()
	if err != nil && err == redis.Nil {
		// Key not found in the cache
		return CrawledPage{}, errors.New("Page not found in cache")
	} else if err != nil {
		// Error while accessing Redis
		log.Printf("Error accessing Redis: %v\n", err)
		return CrawledPage{}, err
	}

	// Parse the cached timestamp
	timestampStr, err := redisClient.Get(ctx, cacheKey+":timestamp").Result()
	if err != nil {
		log.Printf("Error accessing timestamp in Redis: %v\n", err)
		return CrawledPage{}, err
	}
	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		log.Printf("Error parsing timestamp from Redis: %v\n", err)
		return CrawledPage{}, err
	}

	return CrawledPage{
		Content:   val,
		Timestamp: timestamp,
	}, nil
}

func storePageInCache(url, content string) {
	cacheKey := "page:" + url
	// Store the page content
	redisClient.Set(ctx, cacheKey, content, 0)
	// Store the timestamp
	redisClient.Set(ctx, cacheKey+":timestamp", time.Now().Format(time.RFC3339), 0)
}

func crawlURL(url string, speedMultiplier float64) (string, error) {
	client := &http.Client{}

	for retry := 0; retry <= maxRetries; retry++ {
		resp, err := client.Get(url)
		if err != nil {
			if retry < maxRetries {
				log.Printf("Error: %v, Retrying...\n", err)
				time.Sleep(time.Duration(float64(retryDelay) / speedMultiplier))
				continue
			}
			return "", err
		}
		defer resp.Body.Close()

		content, err := parseHTML(resp.Body)
		if err != nil {
			if retry < maxRetries {
				log.Printf("Error: %v, Retrying...\n", err)
				time.Sleep(time.Duration(float64(retryDelay) / speedMultiplier))
				continue
			}
			return "", err
		}

		return content, nil
	}

	return "", errors.New("max retries exceeded")
}

func parseHTML(body io.Reader) (string, error) {
	// Parse the HTML content
	tokenizer := html.NewTokenizer(body)
	crawlingDetails := ""

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				// End of the HTML content, return the details and no error
				return crawlingDetails, nil
			}
			// There was an error during parsing, return it
			return crawlingDetails, err
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						decodedHref := html.EscapeString(attr.Val)
						crawlingDetails += fmt.Sprintf("<a href=\"%s\">%s</a><br>", decodedHref, decodedHref)
					}
				}
			}
		}
	}
}

func writeErrorToHTML(w http.ResponseWriter, r *http.Request, crawlTime time.Time) {
	file, err := os.Create("answer.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	errorContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Error</title>
		</head>
		<body>
			<h1>Unable To Process The Crawling Details At This Moment. Check the URL and try again</h1>
			<p>Crawl Time: %f seconds</p>
		</body>
		</html>
	`, time.Since(crawlTime).Seconds())

	_, err = file.WriteString(errorContent)
	if err != nil {
		log.Fatal(err)
	}

	http.ServeFile(w, r, "answer.html")
}

func writeCrawlTimeToHTML(w http.ResponseWriter, r *http.Request, crawlTime float64, content string, fetchedFromCache bool) {
	file, err := os.Create("answer.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var message string
	if fetchedFromCache {
		message = "Details were fetched from the cache."
	} else {
		message = fmt.Sprintf("Crawl Time: %f seconds", crawlTime)
	}

	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Crawling Result</title>
		</head>
		<body>
			<h1>Crawling Completed</h1>
			<p>%s</p>
			%s
		</body>
		</html>
	`, message, content)

	_, err = file.WriteString(htmlContent)
	if err != nil {
		log.Fatal(err)
	}

	http.ServeFile(w, r, "answer.html")
}
