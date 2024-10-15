package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

		selectors = teamProfileSelectors
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
			downloadHtml, err := e.DOM.Find("blockquote p").Html()
			if err == nil {
				lines := newlineTabsRegex.Split(downloadHtml, -1)
				for _, line := range lines {
					if len(line) > 0 {
						song.DownloadRaw = append(song.DownloadRaw, line)
					}
				}
			}
			// downloadSection := e.DOM.Find("blockquote a").Text()
			// lines := newlineTabsRegex.Split(downloadSection, -1)
			var downloadSection []string

			// Find the <p> element within blockquote and get each line of text
			e.DOM.Find("blockquote p").Each(func(i int, selection *goquery.Selection) {
				selection.Contents().Each(func(j int, s *goquery.Selection) {
					// Check if we have an element node that is an <a> tag
					if goquery.NodeName(s) == "a" {
						href, exists := s.Attr("href")
						if exists {
							downloadSection = append(downloadSection, href)
							song.TestStringArray = append(song.TestStringArray, href)
						}
					} else {
						// For regular text nodes, just append the text
						text := strings.TrimSpace(s.Text())
						if text != "" {
							downloadSection = append(downloadSection, text)
							song.TestStringArray = append(song.TestStringArray, text)
						}
					}
				})
			})
			if len(downloadSection) == 1 && urlRegex.MatchString(downloadSection[0]) { // single line link
				song.DownloadProcessed = append(song.DownloadProcessed, DownloadLink{
					Url:  downloadSection[0],
					Desc: "",
					Tags: []LinkTag{{Id: Unlabeled, String: LinkCategory(Unlabeled).String()}},
				},
				)
			} else if len(downloadSection) > 1 {
				newLink := DownloadLink{Desc: "", Url: "", Tags: []LinkTag{}}
				var firstLineIsLink bool
				var checkEven bool

				if urlRegex.MatchString(downloadSection[0]) {
					firstLineIsLink = true
					checkEven = false
				} else {
					firstLineIsLink = false
					checkEven = true
				}
				for index, line := range downloadSection {
					if firstLineIsLink { // first line is a link
						if index == 0 {
							song.DownloadProcessed = append(song.DownloadProcessed, DownloadLink{
								Url:  downloadSection[0],
								Desc: "",
								Tags: []LinkTag{{Id: Unlabeled, String: LinkCategory(Unlabeled).String()}},
							})
							continue
						}
					}

					if (checkEven && index%2 == 0) || (!checkEven && index%2 != 0) {
						newLink.Tags = append(newLink.Tags, setTags(line)...)

						newLink.Desc = line
						if index != len(downloadSection)-1 {
							if !urlRegex.MatchString(downloadSection[index+1]) {
								newLink.Url = ""
								newLink.Tags = []LinkTag{{Id: Unlinked, String: LinkCategory(Unlinked).String()}}
								song.DownloadProcessed = append(song.DownloadProcessed, newLink)

								newLink.Tags = []LinkTag{}
								checkEven = !checkEven
							}
						} else if index == len(downloadSection)-1 {
							if !urlRegex.MatchString(downloadSection[index]) {
								newLink.Url = ""
								newLink.Tags = []LinkTag{{Id: Unlinked, String: LinkCategory(Unlinked).String()}}
								song.DownloadProcessed = append(song.DownloadProcessed, newLink)

								newLink.Tags = []LinkTag{}
								checkEven = !checkEven
							}
						}
					} else {
						// TODO check if line is a link
						newLink.Url = line
						song.DownloadProcessed = append(song.DownloadProcessed, newLink)

						newLink.Tags = []LinkTag{}
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

			sections := []string{}
			e.ForEach(selectors.FancyTitle, func(_ int, el *colly.HTMLElement) {

				section := strings.TrimSpace(el.Text)
				sections = append(sections, section)
				// TODO flag sections not tracked in teamProfileSectionHeaders

				if section == "Points" {
					parent := el.DOM.Parent()

					songPoint := PointValue{}
					// Iterate over the sibling elements at the same level
					parent.NextAll().Each(func(i int, s *goquery.Selection) {

						if s.HasClass("col-md-3") {
							songPoint.Name = s.Find("h5").Text()

							counterValue := s.Find(".counter").Text()
							songPoint.Value, _ = strconv.ParseFloat(counterValue, 32)
							songPoint.Desc = "value"
							song.Points = append(song.Points, songPoint)
							// Check if the sibling has the class "col_full"
							if s.HasClass("col_full") {
								return
							}
						}

					})
				}
				if section == "VOTE" {
					parent := el.DOM.Parent()

					votePoint := PointValue{}
					// Iterate over the sibling elements at the same level
					parent.NextAll().Each(func(i int, s *goquery.Selection) {

						if s.HasClass("col_half") {
							votePoint.Name = s.Find("h5").Text()

							counterValue := s.Find(".counter").Text()
							votePoint.Value, _ = strconv.ParseFloat(counterValue, 32)
							votePoint.Desc = "value"
							song.Votes = append(song.Votes, votePoint)
							// Check if the sibling has the class "col_full"
							if s.HasClass("col_full") {
								return
							}
						}
					})
				}
				if section == "Short Impression" {
					parent := el.DOM.Parent()

					// Iterate over the sibling elements at the same level
					parent.NextAll().Each(func(i int, s *goquery.Selection) {
						// fmt.Printf("parent text: '%s'\n", s.Text())

						if s.HasClass("col_full") {
							hasChildSpost := s.Find("div.spost").Length() > 0
							if hasChildSpost {
								shortImpression := ShortImpression{}
								s.Find("div.spost").Each(func(i int, t *goquery.Selection) {
									shortImpression.Points = parseToInt(t.Find("div.points_oneline").Text())
									shortImpression.UserName = strings.TrimSpace(t.Find("div.entry-title strong").Text())
									matches := inParensRegex.FindStringSubmatch(t.Find("div.entry-title small").Text())
									if len(matches) > 1 {
										shortImpression.UserId = matches[1]
									}
									shortImpression.CountryCode = t.Find("img.flag").AttrOr("title", "")
									shortImpression.CountryFlag = GetPrefixUrl(t.Find("img.flag").AttrOr("src", " "))
									// Assuming the string is stored in a variable called dateString
									dateString := t.Find("small").Text()

									// Find the match
									match := jpDateRegex.FindString(dateString)

									if len(match) > 1 {
										// song.TestString = match
										jpDate, err := ProcessJpDateString(match)
										if err == nil {
											shortImpression.Time, _ = GetHugoDateTime(jpDate)
										}
									}
									shortImpression.Content = strings.TrimSpace(t.Find("div.entry-title").Eq(1).Text())

									song.ShortImpressions = append(song.ShortImpressions, shortImpression)
								})
								return
							}
						}
					})

				}
				if section == "Long Impression" {
					parent := el.DOM.Parent()

					// Collect all the next sibling elements into an array
					var nextSiblings []*goquery.Selection
					// Iterate over the sibling elements at the same level
					parent.NextAll().Each(func(i int, s *goquery.Selection) {
						nextSiblings = append(nextSiblings, s)
					})
					// TODO figure out conditional arithmetic for tracking impressions, empty divs, the div for the button, and the response impressions for each div
					longImpression := LongImpression{}
					pointBreakdown := PointValue{}
					for i := 0; i < len(nextSiblings); i++ {
						if nextSiblings[i].HasClass("spost") && (nextSiblings[i].Find("div.entry-c")).Length() > 0 { // should be a header
							// song.TestStringArray = append(song.TestStringArray, nextSiblings[i].Text())
							//	song.TestString = nextSiblings[i].Find("nobr").Text()
							longImpression.PointsOverall = parseToInt(nextSiblings[i].Find("nobr").Text())
							longImpression.UserName = strings.TrimSpace(nextSiblings[i].Find("div.entry-title strong").Text())
							nextSiblings[i].Find("span").Each(func(_ int, span *goquery.Selection) {

								// Extract the text before ':' as pointBreakdown.Name
								fullText := span.Text()
								parts := strings.Split(fullText, ":")
								if len(parts) >= 2 {
									pointBreakdown.Name = strings.TrimSpace(parts[0])

									// Extract the number after ':' as pointBreakdown.Value
									valueText := strings.TrimSpace(parts[1])
									valueParts := strings.Split(valueText, " ")
									if len(valueParts) >= 2 {
										pointBreakdown.Value, err = strconv.ParseFloat(strings.TrimSpace(valueParts[0]), 64)
										if err != nil {
											fmt.Println("Error parsing point value:", err)
										}
									}
									pointBreakdown.Desc = "value"
								}

								// Add the pointBreakdown to the song's list of point breakdowns
								if pointBreakdown.Value != 0 {
									longImpression.PointBreakdown = append(longImpression.PointBreakdown, pointBreakdown)
								}
							})

							longImpression.CountryCode = nextSiblings[i].Find("img.flag").AttrOr("title", "")
							longImpression.CountryFlag = GetPrefixUrl(nextSiblings[i].Find("img.flag").AttrOr("src", " "))
							jpDate, err := ProcessJpDateString(nextSiblings[i].Find("ul.entry-meta li").Text())
							if err == nil {
								longImpression.Time, _ = GetHugoDateTime(jpDate)
							}

							matches := inParensRegex.FindStringSubmatch(nextSiblings[i].Find("small span").Text())
							if len(matches) > 1 {
								longImpression.UserId = matches[1]
							}
							i += 1 // should be content
							longImpression.Comment = strings.TrimSpace(nextSiblings[i].Find("p.event-desc-detail").Text())
							if strings.Contains(longImpression.Comment, "Impression is invalidity") {
								song.LongImpressions = append(song.LongImpressions, longImpression)
								continue
							}
							for nextSiblings[i+1].Find("div.entry-image").Length() > 0 { // check for reply header
								replyImpression := LongImpression{}
								i += 1 // should be reply header
								replyImpression.UserName = strings.TrimSpace(nextSiblings[i].Find("div.entry-title strong").Text())
								replyImpression.CountryCode = nextSiblings[i].Find("img.flag").AttrOr("title", "")
								replyImpression.CountryFlag = GetPrefixUrl(nextSiblings[i].Find("img.flag").AttrOr("src", " "))
								jpDate, err := ProcessJpDateString(nextSiblings[i].Find("ul.entry-meta li").Text())
								if err == nil {
									replyImpression.Time, _ = GetHugoDateTime(jpDate)
								}

								matches := inParensRegex.FindStringSubmatch(nextSiblings[i].Find("small span").Text())
								if len(matches) > 1 {
									replyImpression.UserId = matches[1]
								}
								i += 1 // should be content
								replyImpression.Comment = strings.TrimSpace(nextSiblings[i].Find("p.event-desc-detail").Text())
								longImpression.ResponseImpressions = append(longImpression.ResponseImpressions, replyImpression)

							}
							if nextSiblings[i+2].HasClass("center") {
								i += 2 // should be response button
								longImpression.ResponseButton = manbowEventUrlPrefix + strings.TrimSpace(nextSiblings[i].Find("a").AttrOr("href", " "))

								song.LongImpressions = append(song.LongImpressions, longImpression)
								i += 1
							}
						}
					}
				}
			})
			jpDate, err := ProcessJpDateString(e.ChildText("div.col_full:nth-child(16) > small:nth-child(1)"))
			if err == nil {
				song.LastVoteTime, _ = GetHugoDateTime(jpDate)
			}
			// song.TestStringArray = sections
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
