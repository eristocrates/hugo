package main

import "os"

// Song struct
type Song struct {
	SongID   string `json:"song_id"`
	SongName string `json:"song_name"`
}

// Team struct
type Team struct {
	TeamId                int      `json:"teamId"`
	TeamEmblemSrc         string   `json:"teamEmblemSrc"`
	TeamName              string   `json:"teamName"`
	TeamProfileLink       string   `json:"teamProfileLink"`
	TeamLeaderName        string   `json:"teamLeaderName"`
	TeamLeaderCountryCode string   `json:"teamLeaderCountryCode"`
	TeamLeaderCountryFlag string   `json:"teamLeaderCountryFlag"`
	TeamMemberCount       int      `json:"teamMemberCount"`
	TeamReleasedWorkCount int      `json:"teamReleasedWorkCount"`
	TeamDeclaredWorkCount int      `json:"teamDeclaredWorkCount"`
	TeamMemberListRaw     string   `json:"teamMemberListRaw"`
	TeamUpdate            string   `json:"teamUpdate"`
	TestString            string   `json:"testString"`
	TestStringArray       []string `json:"testStringArray"`
	Songs                 []Song   `json:"songs"`
}

// Event struct
type Event struct {
	EventId           int      `json:"eventId"`
	FullName          string   `json:"fullName"`
	AbbrevName        string   `json:"abbrevName"`
	ShortName         string   `json:"shortName"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	BannerType1       string   `json:"bannerType1"`
	RegistrationStart string   `json:"registrationStart"`
	RegistrationEnd   string   `json:"registrationEnd"`
	ImpressionStart   string   `json:"impressionStart"`
	ImpressionEnd     string   `json:"impressionEnd"`
	PeriodStart       string   `json:"periodStart"`
	PeriodEnd         string   `json:"periodEnd"`
	EntryCount        int      `json:"entryCount"`
	ImpressionCount   int      `json:"impressionCount"`
	InfoLink          string   `json:"infoLink"`
	DetailLink        string   `json:"detailLink"`
	ListLink          string   `json:"listLink"`
	TeamListLink      string   `json:"teamListLink"`
	IsBof             bool     `json:"isBof"`
	LogoType1         string   `json:"logoType1"`
	LogoType2         string   `json:"logoType2"`
	LogoType3         string   `json:"logoType3"`
	LogoType4         string   `json:"logoType4"`
	LogoType5         string   `json:"logoType5"`
	TitleJpg          string   `json:"titleJpg"`
	BannerType2       string   `json:"bannerType2"`
	Video             string   `json:"video"`
	HeaderJpg         string   `json:"headerJpg"`
	HeaderPng         string   `json:"headerPng"`
	BackJpg           string   `json:"backJpg"`
	BackPng           string   `json:"backPng"`
	IsModern          bool     `json:"isModern"`
	IsPreModern       bool     `json:"isPreModern"`
	HasModernList     bool     `json:"hasModernList"`
	HasModernTeamList bool     `json:"hasModernTeamList"`
	TestString        string   `json:"testString"`
	TestStringArray   []string `json:"testStringArray"`
	Teams             []Team   `json:"teams"`
}

type CommaWriter struct {
	file      *os.File
	needComma bool
}

type selectorSet struct {
	PrimaryMenu           string
	MenuButtons           string
	TeamList              string
	FancyTitle            string
	TeamElement           string
	TeamName              string
	FirstTeamName         string
	TeamRow               string
	TeamId                string
	TeamEmblemSrc         string
	TeamListProfileLink   string
	TeamListName          string
	TeamListLeaderName    string
	TeamListLeaderCountry string
}
