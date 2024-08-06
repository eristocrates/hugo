package main

import "os"

// Song struct
type Song struct {
	SongId           int    `json:"songId"`
	SongPageLink     string `json:"songPageLink"`
	SongIsSpecial    bool   `json:"songIsSpecial"`
	SongSpecialTitle string `json:"songSpecialTitle"`
	SongJacket       string `json:"songJacket"`
	SongGenre        string `json:"songGenre"`
	SongTitle        string `json:"songTitle"`
	SongArtist       string `json:"songArtist"`
	SongRegistDate   string `json:"songRegistDate"`
	SongLastUpdate   string `json:"songLastUpdate"`
	SongCommentCount int    `json:"songCommentCount"`
}

type Concept struct {
	ConceptName  string `json:"conceptName"`
	ConceptImage string `json:"conceptImage"`
}

// Team struct
type Team struct {
	TeamId                  int           `json:"teamId"`
	TeamEmblemSrc           string        `json:"teamEmblemSrc"`
	TeamName                string        `json:"teamName"`
	TeamNameLabelRaw        []string      `json:"teamNameLabelRaw"`
	TeamIsRecruiting        bool          `json:"teamIsRecruiting"`
	TeamIsWithdrawn         bool          `json:"teamIsWithdrawn"`
	TeamIsDisqualified      bool          `json:"teamIsDisqualified"`
	TeamIsWarned            bool          `json:"teamIsWarned"`
	TeamProfileLink         string        `json:"teamProfileLink"`
	TeamLeaderName          string        `json:"teamLeaderName"`
	TeamLeaderCountryCode   string        `json:"teamLeaderCountryCode"`
	TeamLeaderCountryFlag   string        `json:"teamLeaderCountryFlag"`
	TeamLeaderLanguage      string        `json:"teamLeaderLanguage"`
	TeamMemberCount         int           `json:"teamMemberCount"`
	TeamReleasedWorksCount  int           `json:"teamReleasedWorksCount"`
	TeamDeclaredWorksCount  int           `json:"teamDeclaredWorksCount"`
	TeamMemberListRaw       string        `json:"teamMemberListRaw"`
	TeamMemberListProcessed []string      `json:"teamMemberListProcessed"`
	TeamMemberListIsCorrect bool          `json:"teamMemberListIsCorrect"`
	TeamLastUpdate          string        `json:"teamLastUpdate"`
	TeamTwitter             string        `json:"teamTwitter"`
	TeamWebsite             string        `json:"teamWebsite"`
	TeamConcepts            []Concept     `json:"teamConcepts"`
	RatioPoints             []pointValues `json:"ratioPoints"`
	TeamGenres              []string      `json:"teamGenres"`
	TeamCommonality         string        `json:"teamCommonality"`
	TeamRaisonDetre         string        `json:"teamRaisonDetre"`
	TeamComment             string        `json:"teamComment"`
	TeamRegistDate          string        `json:"teamRegistDate"`
	TestString              string        `json:"testString"`
	TestStringArray         []string      `json:"testStringArray"`
	Songs                   map[int]*Song `json:"songs"`
}

// Event struct
type Event struct {
	EventId           int           `json:"eventId"`
	FullName          string        `json:"fullName"`
	AbbrevName        string        `json:"abbrevName"`
	ShortName         string        `json:"shortName"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	BannerType1       string        `json:"bannerType1"`
	RegistrationStart string        `json:"registrationStart"`
	RegistrationEnd   string        `json:"registrationEnd"`
	ImpressionStart   string        `json:"impressionStart"`
	ImpressionEnd     string        `json:"impressionEnd"`
	PeriodStart       string        `json:"periodStart"`
	PeriodEnd         string        `json:"periodEnd"`
	EntryCount        int           `json:"entryCount"`
	ImpressionCount   int           `json:"impressionCount"`
	InfoLink          string        `json:"infoLink"`
	DetailLink        string        `json:"detailLink"`
	ListLink          string        `json:"listLink"`
	TeamListLink      string        `json:"teamListLink"`
	IsBof             bool          `json:"isBof"`
	LogoType1         string        `json:"logoType1"`
	LogoType2         string        `json:"logoType2"`
	LogoType3         string        `json:"logoType3"`
	LogoType4         string        `json:"logoType4"`
	LogoType5         string        `json:"logoType5"`
	TitleJpg          string        `json:"titleJpg"`
	BannerType2       string        `json:"bannerType2"`
	Video             string        `json:"video"`
	HeaderJpg         string        `json:"headerJpg"`
	HeaderPng         string        `json:"headerPng"`
	BackJpg           string        `json:"backJpg"`
	BackPng           string        `json:"backPng"`
	IsModern          bool          `json:"isModern"`
	IsPreModern       bool          `json:"isPreModern"`
	HasModernList     bool          `json:"hasModernList"`
	HasModernTeamList bool          `json:"hasModernTeamList"`
	TestString        string        `json:"testString"`
	TestStringArray   []string      `json:"testStringArray"`
	Teams             map[int]*Team `json:"teams"`
}

type CommaWriter struct {
	file      *os.File
	needComma bool
}

type CtxIds struct {
	EventId int
	TeamId  int
	SongId  int
}

type selectorSet struct {
	PrimaryMenu           string
	MenuButtons           string
	TeamList              string
	FancyTitle            string
	TeamElement           string
	TeamName              string
	TeamNameLabel         string
	FirstTeamName         string
	TeamRow               string
	TeamId                string
	TeamEmblemSrc         string
	TeamListProfileLink   string
	TeamListName          string
	TeamListNameLabel     string
	TeamListLeaderName    string
	TeamListLeaderCountry string
	TeamListMemberCount   string
	TeamListWorks         string
	TeamListMembers       string
	TeamListUpdate        string
	SectionContent        string
	TwitterButton         string
	WebsiteButton         string
	SongEntries           string
}

type sectionHeaders struct {
	TeamProfile     string
	TeamLeader      string
	Concept         string
	Works           string
	RatioPoint      string
	TeamGenre       string
	TeamCommonality string
	TeamRaisonDetre string
	MemberList      string
	Comment         string
	RegistTime      string
	LastUpdate      string
	TeamProfileEdit string
}

type pointValues struct {
	Name  string
	Value float64
	Type  string
}
