# Go Web Application README

Welcome to the README for this Go web application. This document provides essential information about the project, its features, how to set it up, and how to use it.

## Table of Contents
- [Introduction](#introduction)
- [Tech Stack](#tech-stack)
- [Features](#features)
- [Usage](#Usage)
- [Code Overview](#code-overview)
- [How to Set Up](#how-to-set-up)

## Introduction

This single-page web application allows users to request the crawling of a URL. Users can enter the URL they want to crawl into a search bar and click on the "Crawl" button. The server will check if the URL has been crawled in the last 60 minutes. If the page is found in the server's cache, it will be read and returned to the user. If not, the server will crawl the URL in real-time and return the page. The maximum retries is given is 3, If an invalid user entered then it will check  3 times, and still if it not return anything then it will pop up message to retru with the correct url. It's implemented for 2 tyoes if userrs, Paid ones and non-paid ones. For paid ones 5 crawl APIs work concurrently and for unpaid users 2 crawl APi works, therefore paid user will have lesser crawl time.

I've implemented 
- Required
- Good To Have
- Great To Have

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

### Usage

This project utilizes several Go packages to enhance its functionality. Here's a brief overview of how these packages are used:

1. **`github.com/go-redis/redis/v8`**:
   - Purpose: This package is used for interacting with the Redis database.
   - Usage: It enables the application to connect to a Redis server and perform operations like storing and retrieving data from Redis. In this project, it is used to maintain a Redis cache, storing and retrieving crawled web pages in real-time to improve application performance.

2. **`github.com/gorilla/mux`**:
   - Purpose: Gorilla Mux is a powerful HTTP router for creating complex URL routing structures in Go applications.
   - Usage: In this project, Gorilla Mux is used to define routes and handlers for the web application. It allows for clean and flexible URL routing and helps in handling HTTP requests by specifying the logic for each route.

3. **`golang.org/x/net/html`**:
   - Purpose: This package is part of the Go standard library and is used for parsing HTML documents.
   - Usage: In this project, it is used to parse HTML content when crawling web pages. It allows the application to extract links and other data from HTML documents, making it essential for web scraping and content extraction tasks.

By incorporating these packages, this project benefits from advanced features, efficient routing, and HTML parsing capabilities, enhancing its overall functionality and performance.


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


# [Presentation Video](https://www.loom.com/share/01e431c08f9c40478f50aa101e1a6e73)

You can watch the presentation video for this Go web application by following this [LINK](https://www.loom.com/share/01e431c08f9c40478f50aa101e1a6e73). In the video, you'll get a detailed walkthrough of the application's features and functionality.


