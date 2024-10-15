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

func InitializeDEE2Collector() *colly.Collector {
	dee2Collector := colly.NewCollector(
		colly.Async(true),
		colly.DetectCharset())

	// Limit the number of threads and set delay
	dee2Collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 16,              // Number of threads
		Delay:       1 * time.Second, // Delay between requests
	})
	/*
		dee2Collector.WithTransport(&http.Transport{
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
		featureEvent.PeriodStart, featureEvent.PeriodEnd, err = SplitDateRange(e.ChildText("div.container-center > div.row > div:nth-of-type(4)"))
		HugoDateHerrorCheck(err, featureEvent.FullName)
		featureEvent.RegistrationStart, featureEvent.RegistrationEnd, err = SplitDateRange(e.ChildText("div.container-center > div.row > div.row > div:nth-of-type(2)"))
		HugoDateHerrorCheck(err, featureEvent.FullName)
		featureEvent.ImpressionStart, featureEvent.ImpressionEnd, err = SplitDateRange(e.ChildText("div.container-center > div.row > div.row > div:nth-of-type(4)"))
		HugoDateHerrorCheck(err, featureEvent.FullName)
		featureEvent.EntryCount, err = strconv.Atoi(e.ChildText("div.container-center > div.row > div.row > div:nth-of-type(6)"))
		ConversionErrorCheck(err, featureEvent.FullName)
		featureEvent.ImpressionCount, err = strconv.Atoi(e.ChildText("div.container-center > div.row > div.row > div:nth-of-type(8)"))
		ConversionErrorCheck(err, featureEvent.FullName)
		featureEvent.InfoLink = e.ChildAttr("a.btn-large:nth-child(1)", "href")
		featureEvent.DetailLink = html.UnescapeString(manbowEventUrlPrefix + e.ChildAttr("a.btn:nth-child(2)", "href"))
		featureEvent.ListLink = html.UnescapeString(manbowEventUrlPrefix + e.ChildAttr("a.btn:nth-child(3)", "href"))

		/*
			featureEvent.InfoLink = html.UnescapeString(e.ChildAttr("a.btn-large:nth-child(1)" , "href"))
			featureEvent.DetailLink = html.UnescapeString(manbowEventUrlPrefix + e.ChildAttr("td:nth-of-type(7) a:nth-of-type(2)", "href"))
			featureEvent.ListLink = html.UnescapeString(manbowEventUrlPrefix + e.ChildAttr("td:nth-of-type(7) a:nth-of-type(3)", "href"))
			featureEvent.LastScrapeTime, err = GetHugoDateTime(time.Now().Format("2006/01/02 15:04:05"))
		*/
		HugoDateHerrorCheck(err, featureEvent.FullName)
		AddEvent(&featureEvent)
	})

	// scrape event listing
	dee2Collector.OnHTML("tr.event", func(e *colly.HTMLElement) {
		listEvent := Event{}
		listEvent.FullName = e.ChildText("td:nth-of-type(2) a")
		listEvent.Id, err = strconv.Atoi(e.ChildText("td:nth-of-type(1)"))
		ConversionErrorCheck(err, listEvent.FullName)
		if listEvent.Id != 142 {
			return
		}
		listEvent.RegistrationStart, listEvent.RegistrationEnd, err = SplitDateRange(e.ChildText("td:nth-of-type(3)"))
		HugoDateHerrorCheck(err, listEvent.FullName)
		listEvent.ImpressionStart, listEvent.ImpressionEnd, err = SplitDateRange(e.ChildText("td:nth-of-type(4)"))
		HugoDateHerrorCheck(err, listEvent.FullName)
		listEvent.PeriodStart = listEvent.RegistrationStart
		listEvent.PeriodEnd = listEvent.ImpressionEnd
		listEvent.EntryCount, err = strconv.Atoi(e.ChildText("td:nth-of-type(5)"))
		ConversionErrorCheck(err, listEvent.FullName)
		listEvent.ImpressionCount, err = strconv.Atoi(e.ChildText("td:nth-of-type(6)"))
		ConversionErrorCheck(err, listEvent.FullName)
		listEvent.InfoLink = html.UnescapeString(e.ChildAttr("td:nth-of-type(7) a", "href"))
		listEvent.DetailLink = html.UnescapeString(manbowEventUrlPrefix + e.ChildAttr("td:nth-of-type(7) a:nth-of-type(2)", "href"))
		listEvent.ListLink = html.UnescapeString(manbowEventUrlPrefix + e.ChildAttr("td:nth-of-type(7) a:nth-of-type(3)", "href"))
		listEvent.LastScrapeTime, err = GetHugoDateTime(time.Now().Format("2006/01/02 15:04:05"))
		HugoDateHerrorCheck(err, listEvent.FullName)
		AddEvent(&listEvent)
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
