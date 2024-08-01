package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func InitializeTLCollector() *colly.Collector {
	tlCollector := colly.NewCollector(
		colly.Async(),
		colly.DetectCharset())

	// Limit the number of threads and set delay
	tlCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 16,              // Number of threads
		Delay:       1 * time.Second, // Delay between requests
	})
	/*
		tlCollector.WithTransport(&http.Transport{
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
	tlCollector.SetRequestTimeout(60 * time.Second)

	tlCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
		/*
			fmt.Println("teamLinkListCollector Visiting", r.URL.String())
			logger.Info().Msgf("teamLinkListCollector Visiting %s", r.URL.String())
		*/
	})

	tlCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("teamListLinkCollector Error:\t\t", err, r.StatusCode)
		logger.Error().Msgf("teamListLinkCollector Error: %s, %d", err, r.StatusCode)
	})

	tlCollector.OnResponse(func(r *colly.Response) {
		/*
			fmt.Println("teamListLinkCollector Response received:\t", r.StatusCode)
			logger.Info().Msgf("teamListLinkCollector Response received: %d", r.StatusCode)
		*/
	})

	tlCollector.OnHTML("#teamlist", func(e *colly.HTMLElement) {

		eventId, ok := e.Request.Ctx.GetAny("eventId").(int)
		if ok {
			event := bofEvents[eventId]
			selectors = modernTeamlistSelectors

			team := Team{}
			e.ForEach(selectors.TeamRow, func(i int, s *colly.HTMLElement) {

				team.TeamId, err = strconv.Atoi(s.ChildText(selectors.TeamId))
				team.TeamEmblemSrc = strings.Replace((s.ChildAttr(selectors.TeamEmblemSrc, "src")), "./", manbowEventUrlPrefix, 1)
				team.TeamName = s.ChildText(selectors.TeamListName)
				team.TeamProfileLink = fmt.Sprintf("%s%s", manbowEventUrlPrefix, s.ChildAttr(selectors.TeamListProfileLink, "href"))
				ConversionErrorCheck(err, event.ShortName)

				event.Teams = append(event.Teams, team)
			})

		}
	})
	/*
		tlCollector.OnXML(modernEventXpaths.TeamList, func(e *colly.XMLElement) {
			event, ok := e.Request.Ctx.GetAny("event").(*Event)
			if ok {
				event.IsModern = true
				selectors = modernEventXpaths
				event.TestOutput = e.ChildText(selectors.TeamName)
			}
		})
	*/
	tlCollector.OnScraped(func(r *colly.Response) {
		/*
			fmt.Println("teamListLinkCollector Scraped:\t\t", r.Request.URL)
			logger.Info().Msgf("teamListLinkCollector Scraped: %s", r.Request.URL)
		*/
	})

	return tlCollector
}
