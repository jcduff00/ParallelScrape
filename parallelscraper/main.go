package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type PollData struct {
	PollSource string
	Approve    string
	Disapprove string
}

func main() {
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

	spreadsheetURL := "https://docs.google.com/spreadsheets/d/1Szc6lW9wwq5M0ZnJ2rk7K-Y5SFaCJMMlwKFY9B-hPUg/edit#gid=0"
	if err := c.Visit(spreadsheetURL); err != nil {
		log.Fatal(err)
	}

	csvData := "Poll Source,Approve,Disapprove\n"
	for _, poll := range data {
		row := fmt.Sprintf("%s,%s,%s\n", poll.PollSource, poll.Approve, poll.Disapprove)
		csvData += row
	}

	err := ioutil.WriteFile("poll_data.csv", []byte(csvData), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Let's take a look at this data!")
}
