package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

var modernEvents []Event

func InitializeLLCollector() *colly.Collector {
	llCollector := colly.NewCollector(
		colly.Async(),
		colly.DetectCharset())

	// Limit the number of threads and set delay
	llCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 16,              // Number of threads
		Delay:       1 * time.Second, // Delay between requests
	})

	// Set a timeout for requests
	llCollector.SetRequestTimeout(60 * time.Second)

	llCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
		fmt.Println("linkListCollector Visiting", r.URL.String())
		logger.Info().Msgf("linkListCollector Visiting %s", r.URL.String())
	})

	llCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("listLinkCollector Error:\t\t", err, r.StatusCode)
		logger.Error().Msgf("listLinkCollector Error: %s, %d", err, r.StatusCode)
	})

	llCollector.OnResponse(func(r *colly.Response) {
		/*
			fmt.Println("listLinkCollector Response received:\t", r.StatusCode)
			logger.Info().Msgf("listLinkCollector Response received: %d", r.StatusCode)
		*/
	})

	llCollector.OnHTML("#modern_list", func(e *colly.HTMLElement) {

		event, ok := e.Request.Ctx.GetAny("event").(Event)
		if ok {
			event.IsModern = true
			fmt.Printf("Modern Team List for event: %s\n", event.FullName)
			logger.Info().Str("eventName", event.FullName).Msg("Modern Team List")
		}
	})

	/*
		listLinkCollector.OnXML("", func(e *colly.XMLElement) {
			listLinkCollector.Visit(e.Attr("href"))
		})
	*/
	llCollector.OnScraped(func(r *colly.Response) {
		/*
			fmt.Println("listLinkCollector Scraped:\t\t", r.Request.URL)
			logger.Info().Msgf("listLinkCollector Scraped: %s", r.Request.URL)
		*/
	})

	return llCollector
}
