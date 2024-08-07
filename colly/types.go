package main

import "os"

// TODO add scraped time to event, team and song

type Event struct {
	Id                int           `json:"eventId"`
	FullName          string        `json:"eventFullName"`
	AbbrevName        string        `json:"eventAbbrevName"`
	ShortName         string        `json:"eventShortName"`
	Title             string        `json:"eventTitle"`
	Description       string        `json:"eventDescription"`
	BannerType1       string        `json:"eventBannerType1"`
	RegistrationStart string        `json:"eventRegistrationStart"`
	RegistrationEnd   string        `json:"eventRegistrationEnd"`
	ImpressionStart   string        `json:"eventImpressionStart"`
	ImpressionEnd     string        `json:"eventImpressionEnd"`
	PeriodStart       string        `json:"eventPeriodStart"`
	PeriodEnd         string        `json:"eventPeriodEnd"`
	EntryCount        int           `json:"eventEntryCount"`
	ImpressionCount   int           `json:"eventImpressionCount"`
	InfoLink          string        `json:"eventInfoLink"`
	DetailLink        string        `json:"eventDetailLink"`
	ListLink          string        `json:"eventListLink"`
	TeamListLink      string        `json:"eventTeamListLink"`
	IsBof             bool          `json:"eventIsBof"`
	LogoType1         string        `json:"eventLogoType1"`
	LogoType2         string        `json:"eventLogoType2"`
	LogoType3         string        `json:"eventLogoType3"`
	LogoType4         string        `json:"eventLogoType4"`
	LogoType5         string        `json:"eventLogoType5"`
	TitleJpg          string        `json:"eventTitleJpg"`
	BannerType2       string        `json:"eventBannerType2"`
	Video             string        `json:"eventVideo"`
	HeaderJpg         string        `json:"eventHeaderJpg"`
	HeaderPng         string        `json:"eventHeaderPng"`
	BackJpg           string        `json:"eventBackJpg"`
	BackPng           string        `json:"eventBackPng"`
	IsModern          bool          `json:"eventIsModern"`
	IsPreModern       bool          `json:"eventIsPreModern"`
	HasModernList     bool          `json:"eventHasModernList"`
	HasModernTeamList bool          `json:"eventHasModernTeamList"`
	LastScrapeTime    string        `json:"eventLastScrapeTime"`
	TestString        string        `json:"eventTestString"`
	TestStringArray   []string      `json:"eventTestStringArray"`
	Teams             map[int]*Team `json:"teams"`
}

type Team struct {
	Id                  int           `json:"teamId"`
	Emblem              string        `json:"teamEmblem"`
	Banner              string        `json:"teamBanner"`
	Name                string        `json:"teamName"`
	NameLabelRaw        []string      `json:"teamNameLabelRaw"`
	IsRecruiting        bool          `json:"teamIsRecruiting"`
	IsWithdrawn         bool          `json:"teamIsWithdrawn"`
	IsDisqualified      bool          `json:"teamIsDisqualified"`
	IsWarned            bool          `json:"teamIsWarned"`
	ProfileLink         string        `json:"teamProfileLink"`
	LeaderName          string        `json:"teamLeaderName"`
	LeaderCountryCode   string        `json:"teamLeaderCountryCode"`
	LeaderCountryFlag   string        `json:"teamLeaderCountryFlag"`
	LeaderLanguage      string        `json:"teamLeaderLanguage"`
	MemberCount         int           `json:"teamMemberCount"`
	ReleasedWorksCount  int           `json:"teamReleasedWorksCount"`
	DeclaredWorksCount  int           `json:"teamDeclaredWorksCount"`
	MemberListRaw       string        `json:"teamMemberListRaw"`
	MemberListProcessed []string      `json:"teamMemberListProcessed"`
	MemberListIsCorrect bool          `json:"teamMemberListIsCorrect"`
	LastUpdate          string        `json:"teamLastUpdate"`
	Twitter             string        `json:"teamTwitter"`
	Website             string        `json:"teamWebsite"`
	Concepts            []Concept     `json:"teamConcepts"`
	RatioPoints         []pointValue  `json:"ratioPoints"`
	Genres              []string      `json:"teamGenres"`
	Commonality         string        `json:"teamCommonality"`
	RaisonDetre         string        `json:"teamRaisonDetre"`
	Comment             string        `json:"teamComment"`
	RegistDate          string        `json:"teamRegistDate"`
	Impression          int           `json:"teamImpression"`
	Total               int           `json:"teamTotal"`
	Median              int           `json:"teamMedian"`
	LastScrapeTime      string        `json:"teamLastScrapeTime"`
	TestString          string        `json:"testString"`
	TestStringArray     []string      `json:"testStringArray"`
	Songs               map[int]*Song `json:"songs"`
}

type Song struct {
	Id                int            `json:"songId"`
	PageLink          string         `json:"songPageLink"`
	IsSpecial         bool           `json:"songIsSpecial"`
	SpecialTitle      string         `json:"songSpecialTitle"`
	Jacket            string         `json:"songJacket"`
	Header            string         `json:"songHeader"`
	Genre             string         `json:"songGenre"`
	Title             string         `json:"songTitle"`
	Artist            string         `json:"songArtist"`
	RegistDate        string         `json:"songRegistDate"`
	LastUpdate        string         `json:"songLastUpdate"`
	Keys              []string       `json:"songKeys"`
	CommentCount      int            `json:"songCommentCount"`
	Total             int            `json:"songTotal"`
	Median            int            `json:"songMedian"`
	Composition       string         `json:"songComposition"`
	LastScrapeTime    string         `json:"songLastScrapeTime"`
	Bpm               int            `json:"songBpm"`
	BpmLower          int            `json:"songBpmLower"`
	BpmUpper          int            `json:"songBpmUpper"`
	BpmAverage        int            `json:"songBpmAverage"`
	LevelLower        int            `json:"songLevelLower"`
	LevelUpper        int            `json:"songLevelUpper"`
	BgaStatus         string         `json:"songBgaStatus"`
	Youtube           string         `json:"songYoutube"`
	Size              int            `json:"songSize"`
	DownloadRaw       string         `json:"songDownloadRaw"`
	DownloadProcessed []downloadLink `json:"songDownloadProcessed"`
	TagsRaw           []string       `json:"songTagsRaw"`
	TagsProcessed     []tag          `json:"songTagsProcessed"`
	Soundcloud        string         `json:"songSoundcloud"`
	Bemuse            string         `json:"songBemuse"`
	Comment           string         `json:"songComment"`
	Points            []pointValue   `json:"songPoints"`
	Votes             []pointValue   `json:"songVotes"`
	LastVoteTime      string         `json:"songLastVoteTime"`
	TestString        string         `json:"songTestString"`
	TestStringArray   []string       `json:"songTestStringArray"`
}

type Concept struct {
	ConceptName  string `json:"conceptName"`
	ConceptImage string `json:"conceptImage"`
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

type pointValue struct {
	Name  string
	Value float64
	Desc  string
}

type downloadLink struct {
	Url  string
	Desc string
}

type tag struct {
	Name        string
	Translation string
	Category    string
}

type shortImpression struct {
	Points      int
	UserName    string
	UserId      string
	CountryCode string
	CountryFlag string
	Time        string
	Content     string
}

type longImpression struct {
	PointsOverall  int
	UserName       string
	CountryCode    string
	CountryFlag    string
	PointBreakdown []pointValue
	Time           string
	Comment        string
	ResponseButton string
	IsReply        bool
	Content        string
}
