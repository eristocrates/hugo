package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

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
	/*
		llCollector.WithTransport(&http.Transport{
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
	llCollector.SetRequestTimeout(60 * time.Second)

	llCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
		/*
			fmt.Println("linkListCollector Visiting", r.URL.String())
			logger.Info().Msgf("linkListCollector Visiting %s", r.URL.String())
		*/
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

	llCollector.OnHTML(modernListSelectors.TeamList, func(e *colly.HTMLElement) {
		eventId, ok := e.Request.Ctx.GetAny("eventId").(int)
		if ok {
			event := bofEvents[eventId]
			event.HasModernList = true
			firstTeamName := e.ChildText(modernListSelectors.FirstTeamName)
			if firstTeamName != "" {
				event.IsModern = true
				selectors = modernListSelectors
			} else {
				event.IsPreModern = true
			}

			// testArray := e.ChildTexts(selectors.TeamName)
			/*
				if len(testArray) > 0 {
					var team Team
					e.ForEach(selectors.TeamName, func(i int, h *colly.HTMLElement) {
						team.TeamName = h.Text
						event.Teams = append(event.Teams, team)
					})
				}
			*/
			// TODO Further disambiguate between events with team_information and events without (modern vs premodern?)
			/*
				selectors = modernEventSelectors
				// TODO move to helperFunctions
				str := e.ChildText(selectors.FancyTitle)
				// Define a function to use as a delimiter
				isDelimiter := func(c rune) bool {
					return c == '\t' || c == '\n'
				}

				// Use strings.FieldsFunc with the custom delimiter function
				parts := strings.FieldsFunc(str, isDelimiter)
			*/
		}
	})
	/*
		llCollector.OnXML(modernEventXpaths.TeamList, func(e *colly.XMLElement) {
			event, ok := e.Request.Ctx.GetAny("event").(*Event)
			if ok {
				event.IsModern = true
				selectors = modernEventXpaths
				event.TestOutput = e.ChildText(selectors.TeamName)
			}
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
