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

// TODO reorder type properties to match the order seen on the website, and reorder the functions in the collectors in the same order

var logger zerolog.Logger
var bofEvents map[int]*Event
var otherEvents map[int]*Event
var selectors selectorSet

func main() {
	start := time.Now()
	bofEvents = make(map[int]*Event)
	otherEvents = make(map[int]*Event)

	// logger := InitializeLogger()

	digitalEmergencyExitCollector := InitializeDEE2Collector()
	infoLinkCollector := InitializeILCollector()
	teamListLinkCollector := InitializeTLCollector()
	teamProfileCollector := InitializeTPCollector()
	songPageCollector := InitializeSPCollector()
	// listLinkCollector := InitializeLLCollector()

	// Start scraping on DEE2 EVENT LIST Digital Emergency Exit 2 Event System
	digitalEmergencyExitCollector.Visit("https://manbow.nothing.sh/event/event.cgi/")
	digitalEmergencyExitCollector.Wait()

	// TODO remove this once boftt is added to main event page
	/*
		boftt := Event{
			Id:                146,
			FullName:          "BOF:TT [THE BMS OF FIGHTERS : TT -Sonata for the 20th Ceremony-]",
			HasModernList:     true,
			InfoLink:          "https://www.bmsoffighters.net/boftt/index.html",
			ListLink:          "https://manbow.nothing.sh/event/event.cgi?action=List_def&event=146",
			RegistrationStart: "2024-07-01 00:00:00",
			RegistrationEnd:   "2024-09-07 00:00:00",
			ImpressionStart:   "2024-09-01 00:00:00",
			ImpressionEnd:     "2024-10-04 00:00:00",
			PeriodStart:       "2024-07-01 00:00:00",
			PeriodEnd:         "2024-10-04 00:00:00",
		}

		AddEvent(&boftt)
	*/

	for id, event := range bofEvents {
		if id == 142 || id == 146 {
			ctx := colly.NewContext()
			ctx.Put("eventId", id)
			infoLinkCollector.Request("GET", event.InfoLink, nil, ctx, nil)
		}
	}
	infoLinkCollector.Wait()

	for id, event := range bofEvents {
		ctx := colly.NewContext()
		ctx.Put("eventId", id)
		teamListLinkCollector.Request("GET", event.TeamListLink, nil, ctx, nil)
	}
	teamListLinkCollector.Wait()

	for eventId, event := range bofEvents {
		for teamId, team := range event.Teams {

			ctx := colly.NewContext()
			ids := CtxIds{
				EventId: eventId,
				TeamId:  teamId,
			}
			ctx.Put("eventTeamIds", ids)
			teamProfileCollector.Request("GET", team.ProfileLink, nil, ctx, nil)
		}
	}
	teamProfileCollector.Wait()

	for eventId, event := range bofEvents {
		for teamId, team := range event.Teams {
			for songId, song := range team.Songs {

				ctx := colly.NewContext()
				ids := CtxIds{
					EventId: eventId,
					TeamId:  teamId,
					SongId:  songId,
				}

				ctx.Put("eventTeamSongIds", ids)
				songPageCollector.Request("GET", song.PageLink, nil, ctx, nil)
			}
		}
	}
	songPageCollector.Wait()

	/*
		for id, event := range bofEvents {
			ctx := colly.NewContext()
			ctx.Put("eventId", id)
			listLinkCollector.Request("GET", event.ListLink, nil, ctx, nil)
		}
		listLinkCollector.Wait()
	*/
	// TODO download files
	// TODO automate extraction per every download link type (i assume classic player pack (LR2, ruv-it) vs modern player pack (Beatoraja, qwilight))
	// TODO parse bms files for difficulty levels, average bpm, bpm per song type, maybe notes per second?
	// TODO automate repackaging into a single zip file for classic and modern
	// TODO automate upload to mega/google drive/pixeldrain/etc
	// TODO feature that link prominently to hugo event page
	// TODO run mock test on bofnt based on song RegistDates

	logFile.Write([]byte("\n]"))     // Close the JSON array
	jsonLogFile.Write([]byte("\n]")) // Close the JSON array
	SaveEventsToFile(bofEvents)
	SaveEventsToFile(otherEvents)
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}
