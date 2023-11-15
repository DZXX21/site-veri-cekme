package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/net/html"
)

var inTitle, inDescription, inDateTime bool

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Choose a website to fetch information from:")
	fmt.Println("1: The Hacker News")
	fmt.Println("2: Security Week")
	fmt.Print("Enter choice (1 or 2): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var url string
	switch choice {
	case "1":
		url = "https://thehackernews.com/"
	case "2":
		url = "https://www.securityweek.com/"
	default:
		fmt.Println("Invalid choice")
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching the page:", err)
		return
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)

	filterChoice := getFilterChoice()

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			if choice == "1" {
				if t.Data == "h2" && hasClass(t.Attr, "home-title") {
					inTitle = true
				} else if t.Data == "div" && hasClass(t.Attr, "home-desc") {
					inDescription = true
				} else if t.Data == "span" && hasClass(t.Attr, "h-datetime") {
					inDateTime = true
				}
			} else if choice == "2" {
				if t.Data == "h2" && hasClass(t.Attr, "zox-s-title3") {
					inTitle = true
				} else if t.Data == "p" && hasClass(t.Attr, "zox-s-graph") {
					inDescription = true
				}
			}
		case tt == html.TextToken:
			text := string(z.Text())

			switch filterChoice {
			case "1":
				filterTitles(text)
			case "2":
				filterDescriptions(text)
			case "3":
				filterDateTime(text)
			case "4":
				filterTitles(text)
				filterDescriptions(text)
				filterDateTime(text)
			default:
				fmt.Println("Invalid choice")
				return
			}
		}
	}
}

func hasClass(attrs []html.Attribute, className string) bool {
	for _, attr := range attrs {
		if attr.Key == "class" && strings.Contains(attr.Val, className) {
			return true
		}
	}
	return false
}

func getFilterChoice() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Choose what to filter:")
	fmt.Println("1: Titles")
	fmt.Println("2: Descriptions")
	fmt.Println("3: Date/Time")
	fmt.Println("4: All")
	fmt.Print("Enter choice (1, 2, 3, or 4): ")

	filterChoice, _ := reader.ReadString('\n')
	return strings.TrimSpace(filterChoice)
}

func filterTitles(text string) {
	if inTitle {
		fmt.Print("Title: ")
		color.Green(text) // Başlığı yeşil renkte göster
		inTitle = false
	}
}

func filterDescriptions(text string) {
	if inDescription {
		fmt.Print("Description: ")
		color.Cyan(text) // Açıklamayı mavi renkte göster
		inDescription = false
	}
}

func filterDateTime(text string) {
	if inDateTime {
		fmt.Print("Date/Time: ")
		color.Yellow(text) // Tarih ve saat bilgisini sarı renkte göster
		inDateTime = false
	}
}
