package main

import (
	"fmt"
	"time"

	// "github.com/fatih/color"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog"
	// "github.com/vbauerster/mpb/v8"
	// "github.com/vbauerster/mpb/v8/decor"
)

var logger zerolog.Logger
var bofEvents map[int]*Event
var otherEvents map[int]*Event
var selectors selectorSet

func main() {
	start := time.Now()
	bofEvents = make(map[int]*Event)
	otherEvents = make(map[int]*Event)

	// logger := InitializeLogger()

	teamListLinkCollector := InitializeTLCollector()
	infoLinkCollector := InitializeILCollector()
	listLinkCollector := InitializeLLCollector()
	digitalEmergencyExitCollector := InitializeDEE2Collector()

	// Start scraping on DEE2 EVENT LIST Digital Emergency Exit 2 Event System
	digitalEmergencyExitCollector.Visit("https://manbow.nothing.sh/event/event.cgi/")
	digitalEmergencyExitCollector.Wait()

	// TODO remove this once boftt is added to main event page
	boftt := Event{
		EventId:       146,
		FullName:      "BOF:TT [THE BMS OF FIGHTERS : TT -Sonata for the 20th Ceremony-]",
		HasModernList: true,
		InfoLink:      "https://www.bmsoffighters.net/boftt/index.html",
		ListLink:      "https://manbow.nothing.sh/event/event.cgi?action=List_def&event=146",
	}

	AddEvent(&boftt)

	for id, event := range bofEvents {
		ctx := colly.NewContext()
		ctx.Put("eventId", id)
		infoLinkCollector.Request("GET", event.InfoLink, nil, ctx, nil)
	}
	infoLinkCollector.Wait()

	for id, event := range bofEvents {
		ctx := colly.NewContext()
		ctx.Put("eventId", id)
		teamListLinkCollector.Request("GET", event.TeamListLink, nil, ctx, nil)
	}
	teamListLinkCollector.Wait()

	for id, event := range bofEvents {
		ctx := colly.NewContext()
		ctx.Put("eventId", id)
		listLinkCollector.Request("GET", event.ListLink, nil, ctx, nil)
	}
	listLinkCollector.Wait()

	logFile.Write([]byte("\n]"))     // Close the JSON array
	jsonLogFile.Write([]byte("\n]")) // Close the JSON array
	SaveEventsToFile(bofEvents)
	SaveEventsToFile(otherEvents)
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}
