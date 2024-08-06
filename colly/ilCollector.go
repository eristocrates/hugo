package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

func InitializeILCollector() *colly.Collector {
	ilColector := colly.NewCollector(
		colly.Async(),
		colly.DetectCharset())

	// Limit the number of threads and set delay
	ilColector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 16,              // Number of threads
		Delay:       1 * time.Second, // Delay between requests
	})
	/*
		ilColector.WithTransport(&http.Transport{
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
	ilColector.SetRequestTimeout(60 * time.Second)

	ilColector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
		/*
			fmt.Println("linkListCollector Visiting", r.URL.String())
			logger.Info().Msgf("linkListCollector Visiting %s", r.URL.String())
		*/
	})

	ilColector.OnError(func(r *colly.Response, err error) {
		fmt.Println("infoLinkCollector Error:\t\t", err, r.StatusCode, r.Request.URL.String())
		logger.Error().Msgf("infoLinkCollector Error: %s, %d", err, r.StatusCode)
	})

	ilColector.OnResponse(func(r *colly.Response) {
		/*
			fmt.Println("infoLinkCollector Response received:\t", r.StatusCode)
			logger.Info().Msgf("infoLinkCollector Response received: %d", r.StatusCode)
		*/
	})

	ilColector.OnHTML(modernInfoListSelectors.PrimaryMenu, func(e *colly.HTMLElement) {
		eventId, ok := e.Request.Ctx.GetAny("eventId").(int)
		if ok {
			event := bofEvents[eventId]
			e.ForEachWithBreak(modernInfoListSelectors.MenuButtons, func(i int, h *colly.HTMLElement) bool {
				if h.Text == "Team" {
					event.TeamListLink = h.ChildAttr("a", "href")
					return false
				}
				return true
			})

		}
		/*
			ilCtx := colly.NewContext()
			ilCtx.Put("infoEvent", &infoEvent)

			tlCollector.Request("GET", infoEvent.TeamListLink, nil, ilCtx, nil)
			tlCollector.Visit(infoEvent.TeamListLink)
		*/
		// }
	})
	/*
		ilColector.OnXML(modernEventXpaths.TeamList, func(e *colly.XMLElement) {
			event, ok := e.Request.Ctx.GetAny("event").(*Event)
			if ok {
				event.IsModern = true
				selectors = modernEventXpaths
				event.TestOutput = e.ChildText(selectors.TeamName)
			}
		})
	*/
	ilColector.OnScraped(func(r *colly.Response) {
		/*
			fmt.Println("infoLinkCollector Scraped:\t\t", r.Request.URL)
			logger.Info().Msgf("infoLinkCollector Scraped: %s", r.Request.URL)
		*/
	})

	return ilColector
}
