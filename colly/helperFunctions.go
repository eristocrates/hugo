package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
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
func HugoDateHerrorCheck(err error, eventName string) {
	if err != nil {
		logger.Error().Msg(fmt.Sprintf("Error converting Date Time for %s to iso layout: %s\n", eventName, err))
	}
}

func GetHugoDateTime(dateStr string) (string, error) {
	dateLayouts := []string{
		"2006/01/02",
		"2006/01/02 15:04",
		"2006/01/02 15:04:05",
	}
	isoLayout := "2006-01-02T15:04:05+09:00"

	var parsedTime time.Time
	var err error

	for _, layout := range dateLayouts {
		parsedTime, err = time.Parse(layout, dateStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return "", err
	}

	return parsedTime.Format(isoLayout), nil
}

func SplitDateRange(dateRange string) (startDate string, endDate string, err error) {
	dates := strings.Split(dateRange, "～")
	if len(dates) != 2 {
		return "", "", fmt.Errorf("invalid date range format")
	}

	startDate, err = GetHugoDateTime(dates[0])
	if err != nil {
		return "", "", fmt.Errorf("error converting start date: %w", err)
	}

	endDate, err = GetHugoDateTime(dates[1])
	if err != nil {
		return "", "", fmt.Errorf("error converting end date: %w", err)
	}

	return startDate, endDate, nil
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

func splitMembers(input string, memberCount int) ([]string, bool) {
	/*

		// Second pass
		result = secondPassSplit(input)
		if len(result) == memberCount {
			return result, true
		}
	*/

	// fourth pass
	fourthPassResult := fourthPassSplit(input)
	if len(fourthPassResult) == memberCount {
		return fourthPassResult, true
	}

	var thirdPassResult []string
	// third pass
	for _, element := range fourthPassResult {
		midResult := thirdPassSplit(element)
		thirdPassResult = append(thirdPassResult, midResult...)
	}
	// result = thirdPassSplit(input)
	if len(thirdPassResult) == memberCount {
		return thirdPassResult, true
	}
	/*
		// First pass
		result = firstPassSplit(input)
		if len(result) == memberCount {
			return result, true
		}
	*/

	return thirdPassResult, false
}

func firstPassSplit(input string) []string {
	re := regexp.MustCompile(`(.*?）)`)
	matches := re.FindAllString(input, -1)

	var result []string
	for _, match := range matches {
		trimmed := strings.TrimSpace(match)
		if trimmed != "" && !strings.ContainsRune(trimmed, '\uFFFD') {
			result = append(result, trimmed)
		}
	}

	// Handle any remaining text after the last "）"
	lastIndex := strings.LastIndex(input, "）")
	if lastIndex != -1 && lastIndex < len(input)-1 {
		remaining := strings.TrimSpace(input[lastIndex+1:])
		if remaining != "" && !strings.ContainsRune(remaining, '\uFFFD') {
			result = append(result, remaining)
		}
	}

	return result
}

func secondPassSplit(input string) []string {
	re := regexp.MustCompile(`(.*?]),?\s*`)
	matches := re.FindAllString(input, -1)

	var result []string
	for _, match := range matches {
		trimmed := strings.TrimRight(strings.TrimSpace(match), ", ")
		if trimmed != "" && !strings.ContainsRune(trimmed, '\uFFFD') {
			result = append(result, trimmed)
		}
	}

	// Handle any remaining text after the last "]"
	lastIndex := strings.LastIndex(input, "]")
	if lastIndex != -1 && lastIndex < len(input)-1 {
		remaining := strings.TrimSpace(input[lastIndex+1:])
		if remaining != "" && !strings.ContainsRune(remaining, '\uFFFD') {
			result = append(result, remaining)
		}
	}

	return result
}

func thirdPassSplit(input string) []string {
	var result []string
	var current []rune

	for i, r := range input {
		current = append(current, r)

		if r == ')' || r == '）' {
			shouldSplit := false
			for j := i + 1; j < len(input); j++ {
				nextRune := rune(input[j])
				if !unicode.IsSpace(nextRune) {
					if nextRune != '(' && nextRune != '（' {
						shouldSplit = true
					}
					break
				}
			}
			if shouldSplit {
				result = append(result, string(current))
				current = nil
			}
		}
	}

	if len(current) > 0 {
		result = append(result, string(current))
	}

	var cleanedResult []string
	for i := 0; i < len(result); i++ {
		cleaned := strings.TrimLeftFunc(result[i], func(r rune) bool {
			return unicode.IsSpace(r) || r == ',' || r == '、' || r == ';'
		})
		if cleaned != "" {
			cleanedResult = append(cleanedResult, cleaned)
		}
	}

	return cleanedResult
}

func fourthPassSplit(input string) []string {
	var result []string
	var current []rune

	for _, r := range input {
		current = append(current, r)
		if r == '】' || r == ']' {
			result = append(result, string(current))
			current = nil
		}
	}

	if len(current) > 0 {
		result = append(result, string(current))
	}

	// Clean elements after the first one
	for i := 1; i < len(result); i++ {
		result[i] = strings.TrimLeftFunc(result[i], func(r rune) bool {
			return cleanPrefix(r)
		})
	}

	return result
}

func cleanPrefix(r rune) bool {
	return unicode.IsSpace(r) || r == ',' || r == '、' || r == ';' || r == '/' || r == '•' || r == '1' || r == '2' || r == '3' || r == '4' || r == '5' || r == '6' || r == '7' || r == '8' || r == '9' || r == '0' || r == '.' || r == '✦'

}

func ProcessTeamNameLabel(team *Team) {
	team.TeamIsRecruiting = false
	team.TeamIsWithdrawn = false
	team.TeamIsDisqualified = false
	team.TeamIsWarned = false
	for _, label := range team.TeamNameLabelRaw {
		recruiting := "チームメンバー募集中！"
		withdrawn := "チーム辞退"
		disqualified := "失格"
		warned := "チーム規定違反警告"
		if label == recruiting {
			team.TeamIsRecruiting = true
		}
		if label == withdrawn {
			team.TeamIsWithdrawn = true
		}
		if label == disqualified {
			team.TeamIsDisqualified = true
		}
		if label == warned {
			team.TeamIsWarned = true
		}
	}

}

func GetLanguage(str string) string {
	matches := languageRegex.FindStringSubmatch(str)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""

}

func GetPrefixUrl(url string) string {
	return strings.Replace(url, "./", manbowEventUrlPrefix, 1)
}
func GetIdFromURL(urlStr string, param string) (int, error) {
	var idString string
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return 0, err
	}

	// Extract the "num" query parameter
	if param == "song" {
		idString = parsedURL.Query().Get("num")
	}
	return strconv.Atoi(idString)
}
