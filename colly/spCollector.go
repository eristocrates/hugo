package main

import (
	"fmt"
	"strconv"
	"strings"
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

	spCollector.OnHTML(teamProfileSelectors.SectionContent, func(e *colly.HTMLElement) {
		ids, ok := e.Request.Ctx.GetAny("eventTeamSongIds").(CtxIds)
		if ok {
			song := bofEvents[ids.EventId].Teams[ids.TeamId].Songs[ids.SongId]
			uploadUrl, err := GetUploadPath(e.ChildAttr("div.section", "style"))
			if err == nil {
				song.Header = GetPrefixUrl(uploadUrl)
			}
			song.Jacket = GetPrefixUrl(e.ChildAttr(".moreinfo-header > img:nth-child(1)", "src"))
			bmsInfos := e.ChildTexts(".bmsinfo")
			if len(bmsInfos) > 0 {
				song.Keys = strings.Split(bmsInfos[0], " ")
			}
			bpmString := e.ChildText(".table > tbody:nth-child(1) > tr:nth-child(7) > td:nth-child(4)")
			if strings.Contains(bpmString, "～") {
				song.BpmLower, err = strconv.Atoi(strings.Split(bpmString, "～")[0])
				ConversionErrorCheck(err, song.Title)
				song.BpmUpper, err = strconv.Atoi(strings.Split(bpmString, "～")[1])
				ConversionErrorCheck(err, song.Title)
				song.BpmAverage = (song.BpmLower + song.BpmUpper) / 2
			} else {
				song.Bpm, err = strconv.Atoi(bpmString)
				ConversionErrorCheck(err, song.Title)
			}
			levelString := e.ChildText(".table > tbody:nth-child(1) > tr:nth-child(7) > td:nth-child(2)")
			if strings.Contains(levelString, "～") {
				song.LevelLower, err = strconv.Atoi(strings.TrimLeft(strings.Split(levelString, "～")[0], "★x"))
				ConversionErrorCheck(err, song.Title)
				song.LevelUpper, err = strconv.Atoi(strings.TrimLeft(strings.Split(levelString, "～")[1], "★x"))
				ConversionErrorCheck(err, song.Title)
			}

			song.BgaStatus = strings.Split(e.ChildText(".table > tbody:nth-child(1) > tr:nth-child(6) > td:nth-child(4)"), "・")
			song.Youtube = e.ChildAttr("div.col_one_third > iframe", "src")
			downloadHtml, err := e.DOM.Find("blockquote").Html()
			if err == nil {
				lines := newlineTabsRegex.Split(downloadHtml, -1)
				for _, line := range lines {
					if len(line) > 0 {
						song.DownloadRaw = append(song.DownloadRaw, line)
					}
				}
			}

			tags := e.ChildTexts("div.col_full div.bmsinfo2")
			if len(tags) > 0 {
				tagsRaw := strings.Replace(tags[0], "TAG : ", "", 1)
				song.TagsRaw = strings.Split(tagsRaw, " ")
				for _, tag := range song.TagsRaw {
					var processedTag = Tag{}
					if strings.Contains(tag, "-") {
						splitTag := strings.Split(tag, "-")
						processedTag.Category = splitTag[0]
						processedTag.Name = splitTag[1]
					}
					song.TagsProcessed = append(song.TagsProcessed, processedTag)
				}
			}

			song.Soundcloud = e.ChildText("div.m_audition")
			song.Bemuse = e.ChildAttr("div.bmson-iframe-content > iframe", "src")
			song.Comment = e.ChildText("div.col_full:nth-child(4) > div:nth-child(9) > p:nth-child(2)")
			sizeString := e.ChildText(".table > tbody:nth-child(1) > tr:nth-child(6) > td:nth-child(2)")
			if strings.HasSuffix(sizeString, "KB") {
				sizeString = strings.TrimSuffix(sizeString, "KB")
				sizeInt, err := strconv.Atoi(sizeString)
				if err != nil {
					ConversionErrorCheck(err, song.Title)
				}
				song.Size = sizeInt
			}
			// TODO maybe make this an array by splitting at the dot?
			song.Composition = e.ChildText(".table > tbody:nth-child(1) > tr:nth-child(3) > td:nth-child(3)")
			song.ProdEnv = e.ChildText(".seisakukankyo")
			longImpressionButtonString := e.ChildAttr(".button-aqua", "href")
			song.LongImpressionButton = fmt.Sprintf("%s%s", manbowEventUrlPrefix, longImpressionButtonString)
			shortImpressionButtonString := e.ChildAttr(".button-blue", "href")
			song.ShortImpressionButton = fmt.Sprintf("%s%s", song.PageLink, shortImpressionButtonString)
			// TODO get points
			// TODO get vote
			// TODO get short impressions
			// TODO get long impressions

			// TODO process download links
			song.LastScrapeTime, err = GetHugoDateTime(time.Now().Format("2006/01/02 15:04:05"))
			HugoDateHerrorCheck(err, song.Title)
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
