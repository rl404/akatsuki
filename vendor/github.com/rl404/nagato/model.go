package nagato

import (
	"time"
)

// Anime is anime model.
type Anime struct {
	ID                     int
	Title                  string
	MainPicture            Picture
	AlternativeTitles      AlternativeTitles
	StartDate              Date
	EndDate                Date
	Synopsis               string
	Mean                   float64
	Rank                   int
	Popularity             int
	NumListUsers           int
	NumScoringUsers        int
	NSFW                   NsfwType
	Genres                 []Genre
	MediaType              MediaType
	Status                 StatusType
	MyListStatus           UserAnimeListStatus // Need oauth.
	NumEpisodes            int
	StartSeason            Season
	Broadcast              Broadcast
	Source                 SourceType
	AverageEpisodeDuration time.Duration
	Rating                 RatingType
	Studios                []Studio
	Pictures               []Picture
	Background             string
	RelatedAnime           []RelatedAnime
	RelatedManga           []RelatedManga // Always empty.
	Recommendations        []AnimeRecommendation
	Statistics             Statistic
	NumFavorites           int         // Undocumented.
	OpeningThemes          []ThemeSong // Undocumented.
	EndingThemes           []ThemeSong // Undocumented.
	Videos                 []Video     // Undocumented.
}

// Manga is manga model.
type Manga struct {
	ID                int
	Title             string
	MainPicture       Picture
	AlternativeTitles AlternativeTitles
	StartDate         Date
	EndDate           Date
	Synopsis          string
	Mean              float64
	Rank              int
	Popularity        int
	NumListUsers      int
	NumScoringUsers   int
	NSFW              NsfwType
	Genres            []Genre
	MediaType         MediaType
	Status            StatusType
	MyListStatus      UserMangaListStatus // Need oauth2.
	NumVolumes        int
	NumChapters       int
	Authors           []Author
	Pictures          []Picture
	Background        string
	RelatedAnime      []RelatedAnime // Always empty.
	RelatedManga      []RelatedManga
	Recommendations   []MangaRecommendation
	Serialization     []Serialization // Undocumented.
}

// Picture is picture model.
type Picture struct {
	Large  string
	Medium string
}

// AlternativeTitles is alternative title model.
type AlternativeTitles struct {
	Synonyms []string
	English  string
	Japanese string
}

// Date is date model.
type Date struct {
	Year  int `validate:"gte=0"`
	Month int `validate:"gte=0,lte=12"`
	Day   int `validate:"gte=0,lte=31"`
}

// Genre is genre model.
type Genre struct {
	ID   int
	Name string
}

// UserAnimeListStatus is user anime list status model.
type UserAnimeListStatus struct {
	Status             UserAnimeStatusType
	Score              int
	NumEpisodesWatched int
	IsRewatching       bool
	StartDate          Date
	FinishDate         Date
	Priority           PriorityType
	NumTimesRewatched  int
	RewatchValue       RewatchValueType
	Tags               []string
	Comments           string
	UpdatedAt          time.Time
}

// UserMangaListStatus is my manga list status model.
type UserMangaListStatus struct {
	Status          UserMangaStatusType
	Score           int
	NumVolumesRead  int
	NumChaptersRead int
	IsRereading     bool
	StartDate       Date
	FinishDate      Date
	Priority        PriorityType
	NumTimesReread  int
	RereadValue     RereadValueType
	Tags            []string
	Comments        string
	UpdatedAt       time.Time
}

// Season is season model.
type Season struct {
	Year   int
	Season SeasonType
}

// Broadcast is broadcast model.
type Broadcast struct {
	DayOfTheWeek DayType
	StartTime    string
}

// Studio is studio model.
type Studio struct {
	ID   int
	Name string
}

// RelatedAnime is related anime model.
type RelatedAnime struct {
	Anime        Anime
	RelationType RelationType
}

// RelatedManga is related manga model.
type RelatedManga struct {
	Manga        Manga
	RelationType RelationType
}

// AnimeRecommendation is anime recommendation model.
type AnimeRecommendation struct {
	Anime              Anime
	NumRecommendations int
}

// Statistic is statistic model.
type Statistic struct {
	NumListUsers int
	Status       StatisticStatus
}

// StatisticStatus is statistic status model.
type StatisticStatus struct {
	Watching    int
	Completed   int
	OnHold      int
	Dropped     int
	PlanToWatch int
}

// ThemeSong is theme song model.
//
// Undocumented.
type ThemeSong struct {
	ID   int
	Name string
}

// Video is video model.
//
// Undocumented.
type Video struct {
	ID        int
	Title     string
	URL       string
	Thumbnail string
}

// Author is author model.
type Author struct {
	Person Person
	Role   string
}

// Person is person model.
type Person struct {
	ID        int
	FirstName string
	LastName  string
}

// Serialization is serialization model.
type Serialization struct {
	Magazine Magazine
	Role     string // Always empty.
}

// Magazine is magazine model.
type Magazine struct {
	ID   int
	Name string
}

// MangaRecommendation is manga recommendation model.
type MangaRecommendation struct {
	Manga              Manga
	NumRecommendations int
}

// UserAnime is user anime model.
type UserAnime struct {
	Anime  Anime
	Status UserAnimeListStatus
}

// UserManga is user manga model.
type UserManga struct {
	Manga  Manga
	Status UserMangaListStatus
}

// User is user model.
type User struct {
	ID              int
	Name            string
	Picture         string
	Gender          string
	Birthday        string
	Location        string
	JoinedAt        time.Time
	AnimeStatistics UserAnimeStatistic
	TimeZone        string
	IsSupporter     bool
}

// UserAnimeStatistic is user anime statistic.
type UserAnimeStatistic struct {
	WatchingCount    int
	CompletedCount   int
	OnHoldCount      int
	DroppedCount     int
	PlanToWatchCount int
	TotalCount       int
	WatchedDays      float64
	WatchingDays     float64
	CompletedDays    float64
	OnHoldDays       float64
	DroppedDays      float64
	TotalDays        float64
	Episode          int
	RewatchedTimes   int
	MeanScore        float64
}

// ForumBoardCategory is forum board category model.
type ForumBoardCategory struct {
	Title  string
	Boards []ForumBoard
}

// ForumBoard is forum board model.
type ForumBoard struct {
	ID          int
	Title       string
	Description string
	Subboards   []ForumSubboard
}

// ForumSubboard is forum subboard model.
type ForumSubboard struct {
	ID    int
	Title string
}

// ForumTopic is forum topic model.
type ForumTopic struct {
	ID                int
	Title             string
	CreatedAt         time.Time
	CreatedBy         string
	PostCount         int
	LastPostCreatedAt time.Time
	LastPostCreatedBy string
	IsLocked          bool
}

// ForumTopicDetail is forum topic detail model.
type ForumTopicDetail struct {
	Title string
	Posts []ForumPost
	Poll  ForumPoll
}

// ForumPost is forum post model.
type ForumPost struct {
	ID        int
	Number    int
	CreatedAt time.Time
	CreatedBy string
	Body      string
	Signature string
}

// ForumPoll is forum poll model.
type ForumPoll struct {
	ID       int
	Question string
	IsClosed bool
	Options  []ForumPollOption
}

// ForumPollOption is forum poll option.
type ForumPollOption struct {
	ID    int
	Text  string
	Votes int
}
