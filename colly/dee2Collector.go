package main

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

var err error

func InitializeDEE2Collector(llCollector *colly.Collector) *colly.Collector {
	dee2Collector := colly.NewCollector(
		colly.Async(true),
		colly.DetectCharset())

	// Limit the number of threads and set delay
	dee2Collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 16,              // Number of threads
		Delay:       1 * time.Second, // Delay between requests
	})

	// Set a timeout for requests
	dee2Collector.SetRequestTimeout(60 * time.Second)

	dee2Collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")

		/*
			fmt.Println("digitalEmergenceExitCollector Visiting", r.URL.String())
			logger.Info().Msgf("digitalEmergenceExitCollector Visiting %s", r.URL.String())
		*/
	})

	dee2Collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("digitalEmergenceExitCollector Error:\t\t", err, r.StatusCode)
		logger.Error().Msgf("digitalEmergenceExitCollector Error: %s, %d", err, r.StatusCode)
	})

	dee2Collector.OnResponse(func(r *colly.Response) {
		/*
			fmt.Println("digitalEmergenceExitCollector Response received:\t", r.StatusCode)
			logger.Info().Msgf("digitalEmergenceExitCollector Response received: %d", r.StatusCode)
		*/
	})

	// scrape featured event
	dee2Collector.OnHTML(".col-8.carousel-stepstep", func(e *colly.HTMLElement) {
		featureEvent := Event{}
		idStr := e.ChildText("div.container-center > h2 > span.label-success")
		featureEvent.Id, _ = strconv.Atoi(idStr)
		featureEvent.FullName = strings.TrimSpace(e.ChildText("div.container-center > h2"))[len(idStr):]
		featureEvent.PeriodStart, featureEvent.PeriodEnd = SplitDateRange(e.ChildText("div.container-center > div.row > div:nth-of-type(4)"))
		featureEvent.RegistrationStart, featureEvent.RegistrationEnd = SplitDateRange(e.ChildText("div.container-center > div.row > div.row > div:nth-of-type(2)"))
		featureEvent.ImpressionStart, featureEvent.ImpressionEnd = SplitDateRange(e.ChildText("div.container-center > div.row > div.row > div:nth-of-type(4)"))
		featureEvent.EntryCount, err = strconv.Atoi(e.ChildText("div.container-center > div.row > div.row > div:nth-of-type(6)"))
		ConversionErrorCheck(err, featureEvent.FullName)
		featureEvent.ImpressionCount, err = strconv.Atoi(e.ChildText("div.container-center > div.row > div.row > div:nth-of-type(8)"))
		ConversionErrorCheck(err, featureEvent.FullName)

		AddEvent(&featureEvent)
	})

	// scrape event listing
	dee2Collector.OnHTML("tr.event", func(e *colly.HTMLElement) {
		listEvent := Event{}
		listEvent.FullName = e.ChildText("td:nth-of-type(2) a")
		listEvent.Id, err = strconv.Atoi(e.ChildText("td:nth-of-type(1)"))
		ConversionErrorCheck(err, listEvent.FullName)
		listEvent.RegistrationStart, listEvent.RegistrationEnd = SplitDateRange(e.ChildText("td:nth-of-type(3)"))
		listEvent.ImpressionStart, listEvent.ImpressionEnd = SplitDateRange(e.ChildText("td:nth-of-type(4)"))
		listEvent.PeriodStart = listEvent.RegistrationStart
		listEvent.PeriodEnd = listEvent.ImpressionEnd
		listEvent.EntryCount, err = strconv.Atoi(e.ChildText("td:nth-of-type(5)"))
		ConversionErrorCheck(err, listEvent.FullName)
		listEvent.ImpressionCount, err = strconv.Atoi(e.ChildText("td:nth-of-type(6)"))
		ConversionErrorCheck(err, listEvent.FullName)
		listEvent.InfoLink = html.UnescapeString(e.ChildAttr("td:nth-of-type(7) a", "href"))
		listEvent.DetailLink = html.UnescapeString(manbowEventUrlPrefix + e.ChildAttr("td:nth-of-type(7) a:nth-of-type(2)", "href"))
		listEvent.ListLink = html.UnescapeString(manbowEventUrlPrefix + e.ChildAttr("td:nth-of-type(7) a:nth-of-type(3)", "href"))
		AddEvent(&listEvent)
		if listEvent.IsBof {
			ctx := colly.NewContext()
			ctx.Put("event", listEvent)
			llCollector.Request("GET", listEvent.ListLink, nil, ctx, nil)

			llCollector.Visit(listEvent.ListLink)
		}
		// TODO Go to detailLink
		// TODO come up with some form of categorization for ListLink, and make a distinct collector for each type
		// TODO make functions for performaing specific tasks
		// TODO have each collector pass their category to the function as a branching flag for which elements to target
		// TODO selectorSet variables based on category flag
	})
	/*
		digitalEmergencyExitCollector.OnXML("", func(e *colly.XMLElement) {
		})
	*/
	dee2Collector.OnScraped(func(r *colly.Response) {
		/*
			fmt.Println("digitalEmergenceExitCollector Scraped", r.Request.URL)
			logger.Info().Msgf("digitalEmergenceExitCollector Scraped %s", r.Request.URL)
		*/
	})
	return dee2Collector
}
