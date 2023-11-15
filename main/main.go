package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	showDate        bool
	showDescription bool
	siteURL         string
)

func init() {
	flag.BoolVar(&showDate, "date", true, "Include date in the output")
	flag.BoolVar(&showDescription, "description", true, "Include description in the output")
	flag.StringVar(&siteURL, "site", "https://thehackernews.com", "Specify the website URL to scrape (default: https://thehackernews.com)")
}

func scrapeWebsite(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".home-title").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		description := s.Next().Text()
		date := s.Next().Next().Text()

		var output []string
		output = append(output, title)
		if showDescription {
			output = append(output, fmt.Sprintf("\nDescription: %s", description))
		}
		if showDate {
			output = append(output, fmt.Sprintf("\nDate: %s", date))
		}

		fmt.Println(strings.Join(output, ""))
		fmt.Println("---")
	})
}

func main() {
	flag.Parse()

	fmt.Println("Select a website to scrape:")
	fmt.Println("1. https://thehackernews.com")
	fmt.Println("2. https://www.securityweek.com/")

	var choice int
	fmt.Print("Enter your choice (1 or 2): ")
	_, err := fmt.Scan(&choice)
	if err != nil {
		log.Fatal(err)
	}

	switch choice {
	case 1:
		siteURL = "https://thehackernews.com"
	case 2:
		siteURL = "https://www.securityweek.com/"
	default:
		log.Fatal("Invalid choice. Please enter 1 or 2.")
	}

	scrapeWebsite(siteURL)
}
