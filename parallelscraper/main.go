package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type PollData struct {
	Source          string `json:"source"`
	DatesAdminister string `json:"dates_administered"`
	DatesUpdated    string `json:"dates_updated"`
	JoeBiden        string `json:"joe_biden"`
	RFKJr           string `json:"robert_f_kennedy_jr"`
	Marianne        string `json:"marianne_williamson"`
	OtherUndecided  string `json:"other_undecided"`
	Margin          string `json:"margin"`
}

func main() {
	var url string

	flag.StringVar(&url, "url", "", "URL of the webpage to scrape")
	flag.Parse()

	if url == "" {
		fmt.Println("Error: Please provide the URL of the webpage using the -url flag.")
		return
	}

	c := colly.NewCollector()
	var data []PollData

	c.OnHTML(".wikitable tbody tr", func(e *colly.HTMLElement) {
		if e.Index == 0 {
			return
		}

		source := strings.TrimSpace(e.ChildText("td:nth-child(1)"))
		datesAdminister := strings.TrimSpace(e.ChildText("td:nth-child(2)"))
		datesUpdated := strings.TrimSpace(e.ChildText("td:nth-child(3)"))
		joeBiden := strings.TrimSpace(e.ChildText("td:nth-child(4)"))
		rfkJr := strings.TrimSpace(e.ChildText("td:nth-child(5)"))
		marianne := strings.TrimSpace(e.ChildText("td:nth-child(6)"))
		otherUndecided := strings.TrimSpace(e.ChildText("td:nth-child(7)"))
		margin := strings.TrimSpace(e.ChildText("td:nth-child(8)"))

		pollData := PollData{
			Source:          source,
			DatesAdminister: datesAdminister,
			DatesUpdated:    datesUpdated,
			JoeBiden:        joeBiden,
			RFKJr:           rfkJr,
			Marianne:        marianne,
			OtherUndecided:  otherUndecided,
			Margin:          margin,
		}

		data = append(data, pollData)
	})

	c.OnScraped(func(r *colly.Response) {
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Fatal("Error converting data to JSON:", err)
		}

		file, err := os.Create("data.json")
		if err != nil {
			log.Fatal("Error creating JSON file:", err)
		}
		defer file.Close()

		_, err = file.Write(jsonData)
		if err != nil {
			log.Fatal("Error writing JSON data to file:", err)
		}

		fmt.Println("Data scraped successfully and saved in data.json!")
	})

	if err := c.Visit(url); err != nil {
		log.Fatal(err)
	}
}
