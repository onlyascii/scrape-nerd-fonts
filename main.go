package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	// Parse command line arguments
	url := flag.String("url", "", "URL to scrape for links")
	regexPattern := flag.String("regex", "", "Regex pattern to match links")
	flag.Parse()

	if *url == "" {
		fmt.Println("No URL provided")
		return
	}

	// Compile regex pattern
	var regex *regexp.Regexp
	var err error
	if *regexPattern != "" {
		regex, err = regexp.Compile(*regexPattern)
		if err != nil {
			fmt.Println("Invalid regex pattern:", err)
			return
		}
	}

	// Make HTTP request
	resp, err := http.Get(*url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// Extract links and download content if prefixed
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link := a.Val
					if regex != nil && regex.MatchString(link) {
						fmt.Println("Matched Regex:", link)
						downloadContent(link)
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
}

func downloadContent(link string) {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println("Error downloading link:", err)
		return
	}
	defer resp.Body.Close()

	// Create file
	fileName := strings.Split(link, "/")
	file, err := os.Create(fileName[len(fileName)-1])
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write content to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Downloaded:", link)
}
