package main

import (
	"fmt"
	"time"

	// "github.com/fatih/color"

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

	listLinkCollector := InitializeLLCollector()
	teamListLinkCollector := InitializeTLCollector()
	digitalEmergencyExitCollector := InitializeDEE2Collector(listLinkCollector, teamListLinkCollector)

	// Start scraping on DEE2 EVENT LIST Digital Emergency Exit 2 Event System
	digitalEmergencyExitCollector.Visit("https://manbow.nothing.sh/event/event.cgi/")

	digitalEmergencyExitCollector.Wait()
	listLinkCollector.Wait()
	teamListLinkCollector.Wait()

	logFile.Write([]byte("\n]"))     // Close the JSON array
	jsonLogFile.Write([]byte("\n]")) // Close the JSON array
	SaveEventsToFile(bofEvents)
	SaveEventsToFile(otherEvents)
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}
