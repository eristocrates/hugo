package main

import (
	"fmt"
	"regexp"
)

var bofRegex = regexp.MustCompile(`BOF[: ]?[0-9A-Z]+`)
var bmsofRegex = regexp.MustCompile(`BMS OF FIGHTERS[: ]?[0-9A-Z]+`)
var descriptionType1Regex = regexp.MustCompile(`-[^-]+-`)
var titleType1Regex = regexp.MustCompile(`(THE BMS OF FIGHTERS[^-]+)`)
var titleType2Regex = regexp.MustCompile(`(BOF[^-]+)`)
var titleType3Regex = regexp.MustCompile(`(BMS OF FIGHTERS[^-]+)`)

var manbowEventUrlPrefix = "https://manbow.nothing.sh/event/"

var teamInfoFirstName_a = "div.team_information:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(2) > h3:nth-child(1) > a:nth-child(1)"
var teamInfoFirstName_h3 = "div.team_information:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(2) > h3:nth-child(1)"
var teamInfoFirstName_h2 = "div.team_information:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(2) > h2:nth-child(1)"
var modernListFirstName_h3 = "#modern_list > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > h3:nth-child(1)"

var teamInfoName_a = "div.team_information > div:nth-child(2) > div:nth-child(1) > div:nth-child(2) > h3:nth-child(1) > a:nth-child(1)"
var teamInfoName_h3 = "div.team_information > div:nth-child(2) > div:nth-child(1) > div:nth-child(2) > h3:nth-child(1)"
var teamInfoName_h2 = "div.team_information > div:nth-child(2) > div:nth-child(1) > div:nth-child(2) > h2:nth-child(1)"
var modernListName_h3 = "#modern_list > div > div:nth-child(1) > div:nth-child(1) > h3:nth-child(1)"

var teamlistRow = "table > tbody > tr"
var teamlistId = "td:nth-child(1) > a"
var teamlistEmblemSrc = "td:nth-child(2) > a > img"
var teamlistName = "td:nth-child(4)"
var teamlistProfileLink = "td:nth-child(3) > a"

var modernListSelectors = selectorSet{
	TeamList:      "#modern_list",
	FancyTitle:    "#modern_list > div",
	TeamElement:   "div.team_information",
	TeamName:      fmt.Sprintf("%s, %s, %s, %s", teamInfoName_a, teamInfoName_h3, teamInfoName_h2, modernListName_h3),
	FirstTeamName: fmt.Sprintf("%s, %s, %s, %s", teamInfoFirstName_a, teamInfoFirstName_h3, teamInfoFirstName_h2, modernListFirstName_h3),
}

var modernTeamlistSelectors = selectorSet{
	TeamRow:             teamlistRow,
	TeamId:              teamlistId,
	TeamEmblemSrc:       teamlistEmblemSrc,
	TeamListProfileLink: teamlistProfileLink,
	TeamListName:        teamlistName,
}

/*
var modernEventXpaths = selectorSet{
	TeamList:    "//*[@id=\"modern_list\"]",
	TeamElement: "/html/body/div[2]/section[3]/div/div[3]/div[4]/div",
	TeamName:    "/html/body/div[2]/section[3]/div/div[3]/div[4]/div[1]/div[1]/div/div[2]/h3/a",
}
*/
