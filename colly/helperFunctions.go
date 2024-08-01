package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

func (w *CommaWriter) Write(p []byte) (n int, err error) {
	if w.needComma {
		if _, err := w.file.Write([]byte(",\n")); err != nil {
			return 0, err
		}
	}
	n, err = w.file.Write(p)
	w.needComma = true
	return
}

func SaveEventsToFile(events map[int]*Event) {
	for _, event := range events {
		var filePath string
		if event.IsBof {
			filePath = fmt.Sprintf("../logs/bof/event%d.json", event.EventId)
		} else {
			filePath = fmt.Sprintf("../logs/other/event%d.json", event.EventId)
		}

		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err, fmt.Sprintf("Could not create file for event %d", event.EventId))
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetEscapeHTML(false) // This line prevents HTML escaping
		encoder.SetIndent("", "  ")
		err = encoder.Encode(event)
		if err != nil {
			log.Fatal(err, fmt.Sprintf("Could not encode event%d to JSON", event.EventId))
		}
	}
}

func BofCheck(event *Event) bool {
	if strings.Contains(event.FullName, "BOF") || strings.Contains(event.FullName, "BMS OF FIGHTERS") {

		if strings.Contains(event.FullName, "BOF") {
			match := bofRegex.FindString(event.FullName)
			if match != "" {
				event.AbbrevName = match
				logger.Info().Msg(fmt.Sprintf("BOF AbbrevName found: %s\n", event.AbbrevName))
			}
		}

		if strings.Contains(event.FullName, "BMS OF FIGHTERS") {
			match := bmsofRegex.FindString(event.FullName)
			if match != "" {
				event.AbbrevName = strings.Replace(match, "BMS OF FIGHTERS", "BOF", 1)
				logger.Info().Msg(fmt.Sprintf("BMS OF FIGHTERS AbbrevName found: %s\n", event.AbbrevName))
			}
		}

		// start title logic
		matches := titleType1Regex.FindStringSubmatch(event.FullName)
		if len(matches) > 1 {
			event.Title = strings.TrimSpace(matches[1])

			logger.Info().Msg(fmt.Sprintf("Title Type 1 found: %s\n", event.Title))
		}

		match := descriptionType1Regex.FindString(event.FullName)
		if match != "" {
			event.Description = strings.ReplaceAll(match, "THE BMS OF FIGHTERS ", "")

			logger.Info().Msg(fmt.Sprintf("Dashed Phrase found: %s\n", event.AbbrevName))
		}

		if event.Title == "" {
			matches = titleType2Regex.FindStringSubmatch(event.FullName)
			if len(matches) > 1 {
				event.Title = strings.TrimSpace(matches[1])
				logger.Info().Msg(fmt.Sprintf("Title Type 2 found: %s\n", event.Title))
			}
		}
		if event.Title == "" {
			matches = titleType3Regex.FindStringSubmatch(event.FullName)
			if len(matches) > 1 {
				event.Title = strings.TrimSpace(matches[1])
				logger.Info().Msg(fmt.Sprintf("Title Type 3 found: %s\n", event.Title))
			}
		}

		// Apply Kanji Descriptions
		if event.Title == "" && event.Description == "" {
			event.Title, event.Description = SplitKanji(event.FullName)
		} else if HasKanji(event.FullName) {
			event.Title, event.Description = SplitKanji(event.FullName)
		} else if event.Description == "" || strings.Contains(event.Description, "BOF") {
			_, event.Description = SplitKanji(event.FullName)
		}

		if strings.Contains(event.Title, " preliminary skirmish") && event.Description == "" {
			event.Title = strings.Replace(event.Title, " preliminary skirmish", "", 1)
			event.Description = "preliminary skirmish"
		}

		// handle special cases not worth the effort of regexing
		// TODO evaluade if these are maybe worth the effort
		if event.EventId == 88 {
			event.Description = "konzertsaal"
		}
		if event.EventId == 66 {
			event.Title = "THE BMS OF FIGHTERS 2010"
		}

		event.ShortName = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(event.AbbrevName, ":", ""), " ", ""))

		if event.EventId == 104 {
			event.FullName = "大血戦 -THE BMS OF FIGHTERS ULTIMATE-"
			event.Title = "大血戦"
			event.Description = "-THE BMS OF FIGHTERS ULTIMATE-"
			event.ShortName = "bofu_daikessen"
		}

		if strings.Contains(event.FullName, "FOON") {
			event.IsBof = false
			return false
		} else {
			event.IsBof = true
		}

		// TODO consider replacing this blind logic ilCollector
		event.LogoType1 = fmt.Sprintf("https://www.bmsoffighters.net/%s/img/%s.png", event.ShortName, event.ShortName)
		event.LogoType2 = fmt.Sprintf("https://www.bmsoffighters.net/%s/img/logo.png", event.ShortName)
		event.LogoType3 = fmt.Sprintf("https://www.bmsoffighters.net/%s/img/%s_logo.png", event.ShortName, event.ShortName)
		event.LogoType4 = fmt.Sprintf("https://www.bmsoffighters.net/%s/img/%s.png", event.ShortName, strings.ToUpper(event.ShortName))
		event.LogoType5 = fmt.Sprintf("https://www.bmsoffighters.net/%s/index_files/%s_2.png", event.ShortName, event.ShortName)
		event.BannerType2 = fmt.Sprintf("https://www.bmsoffighters.net/%s/img/banner.jpg", event.ShortName)
		event.TitleJpg = fmt.Sprintf("https://www.bmsoffighters.net/%s/img/title.jpg", event.ShortName)

		event.Video = fmt.Sprintf("https://www.bmsoffighters.net/%s/img/%s.mp4", event.ShortName, event.ShortName)
		event.HeaderJpg = fmt.Sprintf("https://www.bmsoffighters.net/%s/img/header.jpg", event.ShortName)
		event.HeaderPng = fmt.Sprintf("https://www.bmsoffighters.net/%s/img/header.png", event.ShortName)
		event.BackJpg = fmt.Sprintf("https://www.bmsoffighters.net/%s/img/back.jpg", event.ShortName)
		event.BackPng = fmt.Sprintf("https://www.bmsoffighters.net/%s/img/back.png", event.ShortName)

		if event.InfoLink == "http://www.bmsoffighters.net/" && !strings.Contains(event.FullName, "preliminary") {

			event.InfoLink = fmt.Sprintf("https://www.bmsoffighters.net/%s/index.html", event.ShortName)
		}
		// event.TeamListLink = fmt.Sprintf("https://manbow.nothing.sh/event/event_teamprofile.cgi?event=%d", event.EventId)
		return true

	}
	event.Title = event.FullName
	event.IsBof = false
	return false
}

func AddEvent(event *Event) {

	event.BannerType1 = fmt.Sprintf("%simages/%d.jpg", manbowEventUrlPrefix, event.EventId)
	if BofCheck(event) {
		bofEvents[event.EventId] = event
		logger.Info().Msgf("Added BOF event: %s (ID: %d)", event.FullName, event.EventId)
	} else {
		otherEvents[event.EventId] = event
		logger.Info().Msgf("Added other event: %s (ID: %d)", event.FullName, event.EventId)
	}
}
func ConversionErrorCheck(err error, eventName string) {
	if err != nil {
		logger.Error().Msg(fmt.Sprintf("Error converting event ID for %s to int: %s\n", eventName, err))
	}
}
func SplitDateRange(dateRange string) (startDate string, endDate string) {
	layout := "2006/01/02"
	isoLayout := "2006-01-02T15:04:05+09:00"

	start := strings.Split(dateRange, "～")[0]
	end := strings.Split(dateRange, "～")[1]

	startTime, _ := time.Parse(layout, start)
	endTime, _ := time.Parse(layout, end)

	startDate = startTime.Format(isoLayout)
	endDate = endTime.Format(isoLayout)

	return
}

func IsKanji(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana)
}
func SplitKanji(s string) (nonCJK, cjk string) {
	var nonCJKBuilder, cjkBuilder strings.Builder
	for _, r := range s {
		if IsKanji(r) {
			cjkBuilder.WriteRune(r)
		} else {
			nonCJKBuilder.WriteRune(r)
		}
	}
	return strings.TrimSpace(nonCJKBuilder.String()), strings.TrimSpace(cjkBuilder.String())
}

func HasKanji(str string) bool {
	for _, r := range str {
		if IsKanji(r) {
			return true
		}
	}
	return false
}
