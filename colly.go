package main

import (
	"fmt"
	"strings"

	// "github.com/fatih/color"

	"encoding/json"
	"log"
	"os"
	"runtime/debug"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog"
	// "github.com/vbauerster/mpb/v8"
	// "github.com/vbauerster/mpb/v8/decor"
)

// Song struct
type Song struct {
	SongID   string `json:"song_id"`
	SongName string `json:"song_name"`
}

// Team struct
type Team struct {
	TeamID   string `json:"team_id"`
	TeamName string `json:"team_name"`
	Songs    []Song `json:"songs"`
}

// Event struct
type Event struct {
	Id          string `json:"id"`
	FullName    string `json:"fullName"`
	AbbrevName  string `json:"abbrevName"`
	Name        string `json:"shortName"`
	Dates       string `json:"dates"`
	ResultDates string `json:"resultDates"`
	EntryCount  string `json:"entryCount"`
	PlayCount   string `json:"playCount"`
	Links       string `json:"links"`
	Teams       []Team `json:"teams"`
}

type commaWriter struct {
	file      *os.File
	needComma bool
}

func (w *commaWriter) Write(p []byte) (n int, err error) {
	if w.needComma {
		if _, err := w.file.Write([]byte(",\n")); err != nil {
			return 0, err
		}
	}
	n, err = w.file.Write(p)
	w.needComma = true
	return
}

func saveEventsToFile(events []Event) {
	// TODO save each event to a separate file
	file, err := os.Create("events.json")
	if err != nil {
		log.Fatal(err, "Could not create file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(events)
	if err != nil {
		log.Fatal(err, "Could not encode events to JSON")
	}
}

func main() {
	var events []Event

	buildInfo, _ := debug.ReadBuildInfo()

	logFilePath := "logs/colly.log"
	jsonLogFilePath := "logs/colly.json"
	logFile, err := os.OpenFile(
		logFilePath,
		// os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		os.O_TRUNC|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}
	jsonLogFile, err := os.OpenFile(
		jsonLogFilePath,
		// os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		os.O_TRUNC|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}

	lcw := &commaWriter{file: logFile}
	jcw := &commaWriter{file: jsonLogFile}
	multi := zerolog.MultiLevelWriter(lcw, jcw)

	logger := zerolog.New(multi).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()

	logFile.Write([]byte("[\n"))     // Open the JSON array
	jsonLogFile.Write([]byte("[\n")) // Open the JSON array

	defer logFile.Close()
	defer jsonLogFile.Close()

	// Instantiate default collector
	digitalEmergencyExit := colly.NewCollector(
		colly.DetectCharset(),
	// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
	// colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	)

	// On every a element which has href attribute call callback
	digitalEmergencyExit.OnHTML("tr.event", func(e *colly.HTMLElement) {
		event := Event{}
		fullName := e.ChildText("td:nth-of-type(2) a")
		bof := "BMS OF FIGHTERS"
		if strings.Contains(fullName, bof) {
			event.Id = e.ChildText("td:nth-of-type(1)")
			event.FullName = e.ChildText("td:nth-of-type(2) a")
			event.Dates = e.ChildText("td:nth-of-type(3)")
			event.ResultDates = e.ChildText("td:nth-of-type(4)")
			event.EntryCount = e.ChildText("td:nth-of-type(5)")
			event.PlayCount = e.ChildText("td:nth-of-type(6)")
			event.Links = e.ChildAttr("td:nth-of-type(7) a", "href")
			logger.Info().Msg(fmt.Sprintf("BOF added: %s\n", event))
		}
		/*
		 */

		// logger.Info().Msg(fmt.Sprintf("Event ID: %s\nEvent Name: %s\nEvent Dates: %s\nResult Dates: %s\nEntry Count: %s\nPlay Count: %s\nLinks: %s\n", eventID, eventName, eventDates, resultDates, entryCount, playCount, links))
	})

	// Before making a request print "Visiting ..."
	digitalEmergencyExit.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on DEE2 EVENT LIST Digital Emergency Exit 2 Event System
	digitalEmergencyExit.Visit("https://manbow.nothing.sh/event/event.cgi/")
	logFile.Write([]byte("\n]"))     // Close the JSON array
	jsonLogFile.Write([]byte("\n]")) // Close the JSON array
	saveEventsToFile(events)

}
