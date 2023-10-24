# Go Web Application README

Welcome to the README for this Go web application. This document provides essential information about the project, its features, how to set it up, and how to use it.

## Table of Contents
- [Introduction](#introduction)
- [Tech Stack](#tech-stack)
- [Features](#features)
- [Usage](#usage)
- [Code Overview](#code-overview)
- [How to Set Up](#how-to-set-up)

## Introduction

This Go web application is designed to crawl web pages, check the existence and status of web URLs, and provide the crawled data to users. It is implemented with a focus on improving performance and response time by utilizing caching and concurrent crawling.

## Tech Stack

The application is built using the following technologies:

- **Go:** The primary programming language for the server-side logic.
- **HTML:** For rendering web pages.
- **Redis:** Used as a caching solution to store crawled data in real-time.
- **Postman:** A testing tool for verifying the functionality.

## Features

### Web URL Existence Check
- The application checks the existence of web URLs provided by users based on the HTTP status code returned.
- It supports retries, allowing users to customize the number of retry attempts.

### User Categories
- The application serves two categories of users: Paid and Free.
- For Paid customers, it assigns 5 concurrent crawler workers, while Free customers have 2 crawler workers. This parallel crawling reduces the overall crawling time.

### Caching Mechanism
- Before processing a request, the application checks if the requested data is present in the cache (Redis database).
- If a user has crawled a page within the last 60 minutes, the application fetches the data from the cache, reducing the crawling time to almost zero.

### Web Page Crawling
- The application uses the GoColly library to crawl web pages at a depth of 2.
- It maintains a visited array to store URLs that have already been crawled to reduce redundancy.

### Code Overview

The code includes the following key components:

- **CrawlHandler:** Handles the crawling requests from users, checks for cached data, and crawls the page if not found in the cache.
- **getPageFromCache:** Retrieves a page from the Redis cache based on the URL provided.
- **storePageInCache:** Stores the crawled page in the Redis cache.
- **crawlURL:** Initiates the crawling process, supporting retries and parallel crawling.
- **parseHTML:** Parses HTML content for URLs to crawl.
- **writeErrorToHTML:** Generates an error page if there is an issue with crawling.
- **writeCrawlTimeToHTML:** Generates a result page with crawl time and content.

## How to Set Up

To run this application, follow these steps:

1. Install Go on your system if not already installed.

2. Install the required Go packages using `go get`:
   ```
   go get github.com/go-redis/redis/v8
   go get github.com/gorilla/mux
   go get golang.org/x/net/html
   ```

3. Set up a Redis database and update the Redis client configuration in the code to match your Redis server settings.

4. Compile and run the Go application:
   ```
   go run main.go
   ```

5. Access the application via a web browser or tools  like Postman.

