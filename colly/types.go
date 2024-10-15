package main

import "os"

type Event struct {
	Id                int           `json:"eventId"`
	FullName          string        `json:"eventFullName"`
	AbbrevName        string        `json:"eventAbbrevName"`
	ShortName         string        `json:"eventShortName"`
	Title             string        `json:"eventTitle"`
	Description       string        `json:"eventDescription"`
	Banner            string        `json:"eventBanner"`
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
	Logo              string        `json:"eventLogo"`
	TitleJpg          string        `json:"eventTitleJpg"`
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
	RatioPoints         []PointValue  `json:"ratioPoints"`
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
	Id                    int               `json:"songId"`
	PageLink              string            `json:"songPageLink"`
	IsSpecial             bool              `json:"songIsSpecial"`
	SpecialTitle          string            `json:"songSpecialTitle"`
	Jacket                string            `json:"songJacket"`
	Header                string            `json:"songHeader"`
	Genre                 string            `json:"songGenre"`
	Title                 string            `json:"songTitle"`
	Artist                string            `json:"songArtist"`
	RegistDate            string            `json:"songRegistDate"`
	LastUpdate            string            `json:"songLastUpdate"`
	Keys                  []string          `json:"songKeys"`
	CommentCount          int               `json:"songCommentCount"`
	Total                 int               `json:"songTotal"`
	Median                int               `json:"songMedian"`
	Composition           string            `json:"songComposition"`
	LastScrapeTime        string            `json:"songLastScrapeTime"`
	Bpm                   int               `json:"songBpm"`
	BpmLower              int               `json:"songBpmLower"`
	BpmUpper              int               `json:"songBpmUpper"`
	BpmAverage            int               `json:"songBpmAverage"`
	LevelLower            int               `json:"songLevelLower"`
	LevelUpper            int               `json:"songLevelUpper"`
	BgaStatus             []string          `json:"songBgaStatus"`
	Youtube               string            `json:"songYoutube"`
	Size                  int               `json:"songSize"`
	DownloadRaw           []string          `json:"songDownloadRaw"`
	DownloadProcessed     []DownloadLink    `json:"songDownloadProcessed"`
	TagsRaw               []string          `json:"songTagsRaw"`
	TagsProcessed         []Tag             `json:"songTagsProcessed"`
	Soundcloud            string            `json:"songSoundcloud"`
	Bemuse                string            `json:"songBemuse"`
	Comment               string            `json:"songComment"`
	Points                []PointValue      `json:"songPoints"`
	Votes                 []PointValue      `json:"songVotes"`
	LastVoteTime          string            `json:"songLastVoteTime"`
	ProdEnv               string            `json:"songProdEnv"`
	ShortImpressionButton string            `json:"songShortImpressionButton"`
	LongImpressionButton  string            `json:"songLongImpressionButton"`
	ShortImpressions      []ShortImpression `json:"songShortImpressions"`
	LongImpressions       []LongImpression  `json:"songLongImpressions"`
	TestString            string            `json:"songTestString"`
	TestStringArray       []string          `json:"songTestStringArray"`
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

type PointValue struct {
	Name  string
	Value float64
	Desc  string
}

type DownloadLink struct {
	Url  string
	Desc string
	Tags []LinkTag
}

type Tag struct {
	Name     string
	Category string
}

type ShortImpression struct {
	Points      int
	UserName    string
	UserId      string
	CountryCode string
	CountryFlag string
	Time        string
	Content     string
}

type LongImpression struct {
	PointsOverall       int
	UserName            string
	CountryCode         string
	CountryFlag         string
	UserId              string
	PointBreakdown      []PointValue
	Time                string
	ResponseButton      string
	Comment             string
	ResponseImpressions []LongImpression
}

type LinkCategory int

const (
	Unlabeled LinkCategory = iota
	Larger
	Smaller
	DP
	PMS
	Latest
	Prior
	Beatoraja
	LR2
	HQ
	LQ
	BGA
	Google
	OneDrive
	Chinese
	Unlinked
	Untracked
)

// String method returns the name of the Status
func (s LinkCategory) String() string {
	switch s {
	case Unlabeled:
		return "Unlabeled"
	case Larger:
		return "Larger"
	case Smaller:
		return "Smaller"
	case DP:
		return "DP"
	case PMS:
		return "PMS"
	case Untracked:
		return "Untracked"
	case Latest:
		return "Latest"
	case Prior:
		return "Prior"
	case Beatoraja:
		return "Beatoraja"
	case LR2:
		return "LR2"
	case HQ:
		return "HQ"
	case LQ:
		return "LQ"
	case BGA:
		return "BGA"
	case Google:
		return "Google"
	case OneDrive:
		return "OneDrive"
	case Chinese:
		return "Chinese"
	case Unlinked:
		return "Unlinked"
	default:
		return "Not a category"
	}
}

type LinkTag struct {
	Id     LinkCategory
	String string
}
