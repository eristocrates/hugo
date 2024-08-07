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
var languageRegex = regexp.MustCompile(`Language\s*:\s*([^)]*)`)
var uploadUrlRegex = regexp.MustCompile(`\./upload/[^']*`)

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
var teamlistName = "td:nth-child(4) > a"
var teamlistNameLabel = "td:nth-child(4) > strong"
var teamlistProfileLink = "td:nth-child(4) > a"
var teamlistLeaderName = "td:nth-child(5)"
var teamlistLeaderCountry = "td:nth-child(5) > img.flag"
var teamlistMemberCount = "td:nth-child(6)"
var teamlistWorks = "td:nth-child(7)"
var teamlistMembers = "td:nth-child(8)"
var teamlistUpdate = "td:nth-child(9)"

var teamProfileTwitterButton = ".button-aqua"
var teamProfileTeamWebsiteButton = ".button-blue"

var fancyTitle = ".fancy-title"

var modernListSelectors = selectorSet{
	TeamList:      "#modern_list",
	FancyTitle:    fancyTitle,
	TeamElement:   "div.team_information",
	TeamName:      fmt.Sprintf("%s, %s, %s, %s", teamInfoName_a, teamInfoName_h3, teamInfoName_h2, modernListName_h3),
	FirstTeamName: fmt.Sprintf("%s, %s, %s, %s", teamInfoFirstName_a, teamInfoFirstName_h3, teamInfoFirstName_h2, modernListFirstName_h3),
}

var modernTeamlistSelectors = selectorSet{
	TeamRow:               teamlistRow,
	TeamId:                teamlistId,
	TeamEmblemSrc:         teamlistEmblemSrc,
	TeamListProfileLink:   teamlistProfileLink,
	TeamListName:          teamlistName,
	TeamListNameLabel:     teamlistNameLabel,
	TeamListLeaderName:    teamlistLeaderName,
	TeamListLeaderCountry: teamlistLeaderCountry,
	TeamListMemberCount:   teamlistMemberCount,
	TeamListWorks:         teamlistWorks,
	TeamListMembers:       teamlistMembers,
	TeamListUpdate:        teamlistUpdate,
}

var modernInfoListSelectors = selectorSet{
	PrimaryMenu: "div.container.clearfix > nav#primary-menu",
	MenuButtons: "li", //:nth-child(2)", // > a:nth-child(1)",
}

var teamProfileSelectors = selectorSet{
	SectionContent: "section#content",
	FancyTitle:     fancyTitle,
	TwitterButton:  teamProfileTwitterButton,
	WebsiteButton:  teamProfileTeamWebsiteButton,
	SongEntries:    "div.entry",
}

var teamProfileSectionHeaders = sectionHeaders{
	TeamProfile:     "Team Profile",
	TeamLeader:      "Team leader",
	Concept:         "Concept",
	Works:           "Works",
	RatioPoint:      "Ratio Point",
	TeamGenre:       "チームジャンル",
	TeamCommonality: "チームの共通点",
	TeamRaisonDetre: "チームを結成した理由",
	MemberList:      "Member List",
	Comment:         "Comment",
	RegistTime:      "Regist Time",
	LastUpdate:      "Last Update",
	TeamProfileEdit: "チームプロフィール編集",
}

// var songPageSelectors = selectorSet{}

/*
var modernEventXpaths = selectorSet{
	TeamList:    "//*[@id=\"modern_list\"]",
	TeamElement: "/html/body/div[2]/section[3]/div/div[3]/div[4]/div",
	TeamName:    "/html/body/div[2]/section[3]/div/div[3]/div[4]/div[1]/div[1]/div/div[2]/h3/a",
}
*/
