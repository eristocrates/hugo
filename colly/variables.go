package main

import "regexp"

var bofRegex = regexp.MustCompile(`BOF[: ]?[0-9A-Z]+`)
var bmsofRegex = regexp.MustCompile(`BMS OF FIGHTERS[: ]?[0-9A-Z]+`)
var descriptionType1Regex = regexp.MustCompile(`-[^-]+-`)
var titleType1Regex = regexp.MustCompile(`(THE BMS OF FIGHTERS[^-]+)`)
var titleType2Regex = regexp.MustCompile(`(BOF[^-]+)`)
var titleType3Regex = regexp.MustCompile(`(BMS OF FIGHTERS[^-]+)`)

var manbowEventUrlPrefix = "https://manbow.nothing.sh/event/"

var modernEventSelectors = selectorSet{
	TeamList:    "#modern_list",
	FancyTitle:  "#modern_list > div",
	TeamElement: "div.team_information",
	TeamName:    "div.row div.col-sm-12 div.fancy-title.title-dotted-border.title-center h3 adiv.team_information:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(2) > h3:nth-child(1) > a:nth-child(1)",
}

/*
var modernEventXpaths = selectorSet{
	TeamList:    "//*[@id=\"modern_list\"]",
	TeamElement: "/html/body/div[2]/section[3]/div/div[3]/div[4]/div",
	TeamName:    "/html/body/div[2]/section[3]/div/div[3]/div[4]/div[1]/div[1]/div/div[2]/h3/a",
}
*/
