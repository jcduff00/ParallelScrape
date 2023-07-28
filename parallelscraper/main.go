package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type PollData struct {
	PollSource string `json:"poll_source"`
	Approve    string `json:"approve"`
	Disapprove string `json:"disapprove"`
}

func main() {
	var spreadsheetURL string

	flag.StringVar(&spreadsheetURL, "url", "", "URL of the spreadsheet to scrape")
	flag.Parse()

	if spreadsheetURL == "" {
		fmt.Println("Error: Please provide the URL of the spreadsheet using the -url flag.")
		return
	}

	c := colly.NewCollector()
	var data []PollData

	c.OnHTML("tr", func(e *colly.HTMLElement) {
		cells := e.DOM.Children()
		if cells.Length() == 3 {
			pollSource := strings.TrimSpace(cells.Eq(0).Text())
			approve := strings.TrimSpace(cells.Eq(1).Text())
			disapprove := strings.TrimSpace(cells.Eq(2).Text())

			pollData := PollData{
				PollSource: pollSource,
				Approve:    approve,
				Disapprove: disapprove,
			}

			data = append(data, pollData)
		}
	})

	if err := c.Visit(spreadsheetURL); err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("Error converting data to JSON:", err)
	}

	err = ioutil.WriteFile("poll_data.json", jsonData, 0644)
	if err != nil {
		log.Fatal("Error writing JSON data to file:", err)
	}

	fmt.Println("Data scraped successfully and saved in poll_data.json!")
}
