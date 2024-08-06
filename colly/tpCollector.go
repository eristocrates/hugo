package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func InitializeTPCollector() *colly.Collector {
	tpCollector := colly.NewCollector(
		colly.Async(),
		colly.DetectCharset())

	// Limit the number of threads and set delay
	tpCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 16,              // Number of threads
		Delay:       1 * time.Second, // Delay between requests
	})
	/*
		tpCollector.WithTransport(&http.Transport{
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
	tpCollector.SetRequestTimeout(60 * time.Second)

	tpCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
		/*
			fmt.Println("teamLinkListCollector Visiting", r.URL.String())
			logger.Info().Msgf("teamLinkListCollector Visiting %s", r.URL.String())
		*/
	})

	tpCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("teamProfileCollector Error:\t\t", err, r.StatusCode, r.Request.URL.String())
		logger.Error().Msgf("teamProfileCollector Error: %s, %d", err, r.StatusCode)
	})

	tpCollector.OnResponse(func(r *colly.Response) {
		/*
			fmt.Println("teamProfileCollector Response received:\t", r.StatusCode)
			logger.Info().Msgf("teamProfileCollector Response received: %d", r.StatusCode)
		*/
	})

	tpCollector.OnHTML(teamProfileSelectors.SectionContent, func(e *colly.HTMLElement) {
		selectors = teamProfileSelectors
		ids, ok := e.Request.Ctx.GetAny("eventTeamIds").(CtxIds)
		if ok {
			team := bofEvents[ids.EventId].Teams[ids.TeamId]
			// TODO consider scraping warning message box
			team.TeamTwitter = e.ChildAttr(selectors.TwitterButton, "href")
			team.TeamWebsite = e.ChildAttr(selectors.WebsiteButton, "href")

			e.ForEach(selectors.FancyTitle, func(_ int, el *colly.HTMLElement) {
				section := strings.TrimSpace(el.Text)
				// TODO flag sections not tracked in teamProfileSectionHeaders

				if section == teamProfileSectionHeaders.TeamLeader {
					team.TeamLeaderLanguage = GetLanguage(el.DOM.Next().Text())
				}

				if section == teamProfileSectionHeaders.Concept {
					concept := Concept{}
					el.DOM.NextAll().Each(func(_ int, s *goquery.Selection) {

						concept.ConceptName = strings.TrimSpace(s.Text())
						imgSrc, exists := s.Find("img").Attr("src")
						if exists {
							concept.ConceptImage = GetPrefixUrl(imgSrc)
						}
						team.TeamConcepts = append(team.TeamConcepts, concept)
					})
				}

				if section == teamProfileSectionHeaders.RatioPoint {
					// Traverse up to the parent element
					parent := el.DOM.Parent()

					ratioPoint := pointValues{}
					// Iterate over the sibling elements at the same level
					parent.NextAll().Each(func(i int, s *goquery.Selection) {

						if s.HasClass("col_one_fourth") {
							ratioPoint.Name = s.Find("h5").Text()

							counterValue := s.Find(".counter").Text()
							if strings.Contains(counterValue, "x") {
								ratioPoint.Value, _ = strconv.ParseFloat(strings.Split(counterValue, "x")[1], 64)
								ratioPoint.Value = math.Round(ratioPoint.Value*10) / 10
								ratioPoint.Type = "multiplier"
							} else {
								ratioPoint.Value, _ = strconv.ParseFloat(counterValue, 32)
								ratioPoint.Type = "value"
							}
							team.RatioPoints = append(team.RatioPoints, ratioPoint)
							// Check if the sibling has the class "col_last"
							if s.HasClass("col_last") {
								return
							}
						}

					})
				}

				if section == teamProfileSectionHeaders.TeamGenre {
					genreString := strings.TrimSpace(el.DOM.Next().Text())
					genres := strings.FieldsFunc(genreString, func(r rune) bool {
						return r == '・' || r == '/'
					})

					for i, genre := range genres {
						genres[i] = strings.TrimSpace(genre)
					}

					team.TeamGenres = genres

				}

				if section == teamProfileSectionHeaders.TeamCommonality {
					team.TeamCommonality = strings.TrimSpace(el.DOM.Next().Text())
				}

				if section == teamProfileSectionHeaders.TeamRaisonDetre {
					team.TeamRaisonDetre = strings.TrimSpace(el.DOM.Next().Text())
				}

				if section == teamProfileSectionHeaders.Comment {
					team.TeamComment = strings.TrimSpace(el.DOM.Next().Text())
				}

				if section == teamProfileSectionHeaders.RegistTime {
					team.TeamRegistDate, err = GetHugoDateTime(strings.TrimSpace(el.DOM.Next().Text()))
					HugoDateHerrorCheck(err, team.TeamName)
				}

			})

			team.Songs = make(map[int]*Song)
			e.ForEach(selectors.SongEntries, func(_ int, el *colly.HTMLElement) {
				song := Song{}
				song.SongPageLink = fmt.Sprintf("%s%s", manbowEventUrlPrefix, el.ChildAttr("a", "href"))
				song.SongId, err = GetIdFromURL(song.SongPageLink, "song")
				if err != nil {
					logger.Error().Err(err).Msgf("Error extracting song id from url: %s", song.SongPageLink)
				}
				song.SongSpecialTitle = strings.TrimSpace(el.ChildText("div.sale-flash"))
				song.SongIsSpecial = song.SongSpecialTitle != ""
				song.SongJacket = GetPrefixUrl(el.ChildAttr("img", "src"))
				song.SongGenre = strings.TrimSpace(el.ChildText("div.entry-title > small"))
				song.SongTitle = strings.TrimSpace(el.ChildText("div.entry-title > h2 > a"))
				song.SongArtist = strings.TrimSpace(el.ChildText("div.entry-title > h5"))
				song.SongRegistDate, err = GetHugoDateTime(strings.TrimSpace(el.ChildText("ul.entry-meta > li:nth-child(2)")))
				HugoDateHerrorCheck(err, song.SongTitle)
				song.SongLastUpdate, err = GetHugoDateTime(strings.TrimSpace(el.ChildText("ul.entry-meta > li:nth-child(3)")))
				HugoDateHerrorCheck(err, song.SongTitle)
				song.SongCommentCount, err = strconv.Atoi(strings.TrimSpace(el.ChildText("ul.entry-meta > li:nth-child(4)")))
				ConversionErrorCheck(err, song.SongTitle)

				team.Songs[song.SongId] = &song
			})

		}
	})

	/*
		tpCollector.OnXML(modernEventXpaths.TeamList, func(e *colly.XMLElement) {
			event, ok := e.Request.Ctx.GetAny("event").(*Event)
			if ok {
				event.IsModern = true
				selectors = modernEventXpaths
				event.TestOutput = e.ChildText(selectors.TeamName)
			}
		})
	*/
	tpCollector.OnScraped(func(r *colly.Response) {
		/*
			fmt.Println("teamProfileCollector Scraped:\t\t", r.Request.URL)
			logger.Info().Msgf("teamProfileCollector Scraped: %s", r.Request.URL)
		*/
	})

	return tpCollector
}
