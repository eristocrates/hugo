package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

func InitializeSPCollector() *colly.Collector {
	spCollector := colly.NewCollector(
		colly.Async(),
		colly.DetectCharset())

	// Limit the number of threads and set delay
	spCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 16,              // Number of threads
		Delay:       1 * time.Second, // Delay between requests
	})
	/*
		spCollector.WithTransport(&http.Transport{
			MaxIdleConnsPerHost: 10,
			Proxy:               http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		})
	*/

	// Set a timeout for requests
	spCollector.SetRequestTimeout(60 * time.Second)

	spCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
		/*
			fmt.Println("teamLinkListCollector Visiting", r.URL.String())
			logger.Info().Msgf("teamLinkListCollector Visiting %s", r.URL.String())
		*/
	})

	spCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("songPageCollector Error:\t\t", err, r.StatusCode, r.Request.URL.String())
		logger.Error().Msgf("songPageCollector Error: %s, %d", err, r.StatusCode)
	})

	spCollector.OnResponse(func(r *colly.Response) {
		/*
			fmt.Println("songPageCollector Response received:\t", r.StatusCode)
			logger.Info().Msgf("songPageCollector Response received: %d", r.StatusCode)
		*/
	})

	spCollector.OnHTML("", func(e *colly.HTMLElement) {
		ids, ok := e.Request.Ctx.GetAny("eventTeamIds").(CtxIds)
		if ok {
			song := bofEvents[ids.EventId].Teams[ids.TeamId].Songs[ids.SongId]
		}
	})

	/*
		spCollector.OnXML(modernEventXpaths.TeamList, func(e *colly.XMLElement) {
			event, ok := e.Request.Ctx.GetAny("event").(*Event)
			if ok {
				event.IsModern = true
				selectors = modernEventXpaths
				event.TestOutput = e.ChildText(selectors.TeamName)
			}
		})
	*/
	spCollector.OnScraped(func(r *colly.Response) {
		/*
			fmt.Println("songPageCollector Scraped:\t\t", r.Request.URL)
			logger.Info().Msgf("songPageCollector Scraped: %s", r.Request.URL)
		*/
	})

	return spCollector
}
