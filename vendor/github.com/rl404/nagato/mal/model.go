package mal

import (
	"time"
)

// Anime is anime model.
type Anime struct {
	ID                     int                   `json:"id"`
	Title                  string                `json:"title"`
	MainPicture            Picture               `json:"main_picture"`
	AlternativeTitles      AlternativeTitles     `json:"alternative_titles"`
	StartDate              string                `json:"start_date"`
	EndDate                string                `json:"end_date"`
	Synopsis               string                `json:"synopsis"`
	Mean                   float64               `json:"mean"`
	Rank                   int                   `json:"rank"`
	Popularity             int                   `json:"popularity"`
	NumListUsers           int                   `json:"num_list_users"`
	NumScoringUsers        int                   `json:"num_scoring_users"`
	NSFW                   string                `json:"nsfw"`
	Genres                 []Genre               `json:"genres"`
	CreatedAt              time.Time             `json:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at"`
	MediaType              string                `json:"media_type"`
	Status                 string                `json:"status"`
	MyListStatus           MyAnimeListStatus     `json:"my_list_status"` // Need oauth2.
	NumEpisodes            int                   `json:"num_episodes"`
	StartSeason            Season                `json:"start_season"`
	Broadcast              Broadcast             `json:"broadcast"`
	Source                 string                `json:"source"`
	AverageEpisodeDuration int                   `json:"average_episode_duration"`
	Rating                 string                `json:"rating"`
	Studios                []Studio              `json:"studios"`
	Pictures               []Picture             `json:"pictures"`
	Background             string                `json:"background"`
	RelatedAnime           []RelatedAnime        `json:"related_anime"`
	RelatedManga           []RelatedManga        `json:"related_manga"` // Always empty.
	Recommendations        []AnimeRecommendation `json:"recommendations"`
	Statistics             Statistic             `json:"statistics"`
	NumFavorites           int                   `json:"num_favorites"`  // Undocumented.
	OpeningThemes          []ThemeSong           `json:"opening_themes"` // Undocumented.
	EndingThemes           []ThemeSong           `json:"ending_themes"`  // Undocumented.
	Videos                 []Video               `json:"videos"`         // Undocumented.
}

// Manga is manga model.
type Manga struct {
	ID                int                   `json:"id"`
	Title             string                `json:"title"`
	MainPicture       Picture               `json:"main_picture"`
	AlternativeTitles AlternativeTitles     `json:"alternative_titles"`
	StartDate         string                `json:"start_date"`
	EndDate           string                `json:"end_date"`
	Synopsis          string                `json:"synopsis"`
	Mean              float64               `json:"mean"`
	Rank              int                   `json:"rank"`
	Popularity        int                   `json:"popularity"`
	NumListUsers      int                   `json:"num_list_users"`
	NumScoringUsers   int                   `json:"num_scoring_users"`
	NSFW              string                `json:"nsfw"`
	Genres            []Genre               `json:"genres"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
	MediaType         string                `json:"media_type"`
	Status            string                `json:"status"`
	MyListStatus      MyMangaListStatus     `json:"my_list_status"` // Need oauth2.
	NumVolumes        int                   `json:"num_volumes"`
	NumChapters       int                   `json:"num_chapters"`
	Authors           []Author              `json:"authors"`
	Pictures          []Picture             `json:"pictures"`
	Background        string                `json:"background"`
	RelatedAnime      []RelatedAnime        `json:"related_anime"` // Always empty.
	RelatedManga      []RelatedManga        `json:"related_manga"`
	Recommendations   []MangaRecommendation `json:"recommendations"`
	Serialization     []Serialization       `json:"serialization"`
	NumFavorites      int                   `json:"num_favorites"` // Undocumented.
}

// Character is character model.
//
// Undocumented.
type Character struct {
	ID              int            `json:"id"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	AlternativeName string         `json:"alternative_name"`
	MainPicture     Picture        `json:"main_picture"` // Only contain `medium`.
	Biography       string         `json:"biography"`
	NumFavorites    int            `json:"num_favorites"`
	Animeography    []Animeography `json:"animeography"`
}

// People is people model.
//
// Undocumented.
type People struct {
	ID               int      `json:"id"`
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	AlternativeNames []string `json:"alternative_names"`
	MainPicture      Picture  `json:"main_picture"` // Only contain `medium`.
	Birthday         string   `json:"birthday"`
	WebsiteURL       string   `json:"website_url"`
	More             string   `json:"more"`
	NumFavorites     int      `json:"num_favorites"`
}

// Picture is picture model.
type Picture struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
}

// AlternativeTitles is alternative title model.
type AlternativeTitles struct {
	Synonyms []string `json:"synonyms"`
	En       string   `json:"en"`
	Ja       string   `json:"ja"`
}

// Genre is genre model.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MyAnimeListStatus is my anime list status model.
type MyAnimeListStatus struct {
	Status             string    `json:"status"`
	Score              int       `json:"score"`
	NumEpisodesWatched int       `json:"num_episodes_watched"`
	IsRewatching       bool      `json:"is_rewatching"`
	StartDate          string    `json:"start_date"`
	FinishDate         string    `json:"finish_date"`
	Priority           int       `json:"priority"`
	NumTimesRewatched  int       `json:"num_times_rewatched"`
	RewatchValue       int       `json:"rewatch_value"`
	Tags               []string  `json:"tags"`
	Comments           string    `json:"comments"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// MyMangaListStatus is my manga list status model.
type MyMangaListStatus struct {
	Status          string    `json:"status"`
	Score           int       `json:"score"`
	NumVolumesRead  int       `json:"num_volumes_read"`
	NumChaptersRead int       `json:"num_chapters_read"`
	IsRereading     bool      `json:"is_rereading"`
	StartDate       string    `json:"start_date"`
	FinishDate      string    `json:"finish_date"`
	Priority        int       `json:"priority"`
	NumTimesReread  int       `json:"num_times_reread"`
	RereadValue     int       `json:"reread_value"`
	Tags            []string  `json:"tags"`
	Comments        string    `json:"comments"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Season is season model.
type Season struct {
	Year   int    `json:"year"`
	Season string `json:"season"`
}

// Broadcast is broadcast model.
type Broadcast struct {
	DayOfTheWeek string `json:"day_of_the_week"`
	StartTime    string `json:"start_time"`
}

// Studio is studio model.
type Studio struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// RelatedAnime is related anime model.
type RelatedAnime struct {
	Node                  Anime  `json:"node"`
	RelationType          string `json:"relation_type"`
	RelationTypeFormatted string `json:"relation_type_formatted"`
}

// RelatedManga is related manga model.
type RelatedManga struct {
	Node                  Manga  `json:"node"`
	RelationType          string `json:"relation_type"`
	RelationTypeFormatted string `json:"relation_type_formatted"`
}

// AnimeRecommendation is anime recommendation model.
type AnimeRecommendation struct {
	Node               Anime `json:"node"`
	NumRecommendations int   `json:"num_recommendations"`
}

// MangaRecommendation is manga recommendation model.
type MangaRecommendation struct {
	Node               Manga `json:"node"`
	NumRecommendations int   `json:"num_recommendations"`
}

// Statistic is statistic model.
type Statistic struct {
	NumListUsers int             `json:"num_list_users"`
	Status       StatisticStatus `json:"status"`
}

// StatisticStatus is statistic status model.
//
// Inconsistent data type.
// If value is > 0, api will return string.
// If value is 0, api will return int.
type StatisticStatus struct {
	Watching    interface{} `json:"watching"`
	Completed   interface{} `json:"completed"`
	OnHold      interface{} `json:"on_hold"`
	Dropped     interface{} `json:"dropped"`
	PlanToWatch interface{} `json:"plan_to_watch"`
}

// ThemeSong is theme song model.
//
// Undocumented.
type ThemeSong struct {
	ID      int    `json:"id"`
	AnimeID int    `json:"anime_id"`
	Text    string `json:"text"`
}

// Video is video model.
//
// Undocumented.
type Video struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
	Thumbnail string `json:"thumbnail"`
}

// Author is author model.
type Author struct {
	Node Person `json:"node"`
	Role string `json:"role"`
}

// Person is person model.
type Person struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Serialization is serialization model.
type Serialization struct {
	Node Magazine `json:"node"`
	Role string   `json:"role"` // Always empty.
}

// Magazine is magazine model.
type Magazine struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Paging is common paging model.
type Paging struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

// AnimePaging is anime paging model.
type AnimePaging struct {
	Data   []AnimePagingData `json:"data"`
	Paging Paging            `json:"paging"`
}

// AnimePagingData is anime paging data model.
type AnimePagingData struct {
	Node Anime `json:"node"`
}

// AnimeRankingPaging is anime ranking paging model.
type AnimeRankingPaging struct {
	Data   []AnimeRankingPagingData `json:"data"`
	Paging Paging                   `json:"paging"`
}

// AnimeRankingPagingData is anime ranking paging data model.
type AnimeRankingPagingData struct {
	Node    Anime        `json:"node"`
	Ranking AnimeRanking `json:"ranking"`
}

// AnimeRanking is anime ranking model.
type AnimeRanking struct {
	Rank         int `json:"rank"`
	PreviousRank int `json:"previous_rank"` // Always empty.
}

// SuggestedAnimePaging is suggested anime paging model.
type SuggestedAnimePaging struct {
	Data   []SuggestedAnimePagingData `json:"data"`
	Paging Paging                     `json:"paging"`
}

// SuggestedAnimePagingData is suggested anime paging data model.
type SuggestedAnimePagingData struct {
	Node Anime `json:"node"`
}

// SeasonalAnimePaging is seasonal anime paging model.
type SeasonalAnimePaging struct {
	Data   []SeasonalAnimePagingData `json:"data"`
	Paging Paging                    `json:"paging"`
}

// SeasonalAnimePagingData is seasonal anime paging data model.
type SeasonalAnimePagingData struct {
	Node Anime `json:"node"`
}

// MangaPaging is manga paging model.
type MangaPaging struct {
	Data   []MangaPagingData `json:"data"`
	Paging Paging            `json:"paging"`
}

// MangaPagingData is manga paging data model.
type MangaPagingData struct {
	Node Manga `json:"node"`
}

// MangaRankingPaging is manga ranking paging model.
type MangaRankingPaging struct {
	Data   []MangaRankingPagingData `json:"data"`
	Paging Paging                   `json:"paging"`
}

// MangaRankingPagingData is manga paging data model.
type MangaRankingPagingData struct {
	Node    Manga        `json:"node"`
	Ranking MangaRanking `json:"ranking"`
}

// MangaRanking is manga ranking model.
type MangaRanking struct {
	Rank         int `json:"rank"`
	PreviousRank int `json:"previous_rank"` // Always empty.
}

// UserAnimePaging is user anime paging model.
type UserAnimePaging struct {
	Data   []UserAnimePagingData `json:"data"`
	Paging Paging                `json:"paging"`
}

// UserAnimePagingData is user anime paging data model.
type UserAnimePagingData struct {
	Node       Anime             `json:"node"`
	ListStatus MyAnimeListStatus `json:"list_status"`
}

// UserMangaPaging is user manga paging model.
type UserMangaPaging struct {
	Data   []UserMangaPagingData `json:"data"`
	Paging Paging                `json:"paging"`
}

// UserMangaPagingData is user manga paging data model.
type UserMangaPagingData struct {
	Node       Manga             `json:"node"`
	ListStatus MyMangaListStatus `json:"list_status"`
}

// User is user model.
type User struct {
	ID              int                `json:"id"`
	Name            string             `json:"name"`
	Picture         string             `json:"picture"`
	Gender          string             `json:"gender"`
	Birthday        string             `json:"birthday"`
	Location        string             `json:"location"`
	JoinedAt        time.Time          `json:"joined_at"`
	AnimeStatistics UserAnimeStatistic `json:"anime_statistics"`
	TimeZone        string             `json:"time_zone"`
	IsSupporter     bool               `json:"is_supporter"`
}

// UserAnimeStatistic is user anime statistic model.
type UserAnimeStatistic struct {
	NumItemsWatching    int     `json:"num_items_watching"`
	NumItemsCompleted   int     `json:"num_items_completed"`
	NumItemsOnHold      int     `json:"num_items_on_hold"`
	NumItemsDropped     int     `json:"num_items_dropped"`
	NumItemsPlanToWatch int     `json:"num_items_plan_to_watch"`
	NumItems            int     `json:"num_items"`
	NumdaysWatched      float64 `json:"num_days_watched"`
	NumdaysWatching     float64 `json:"num_days_watching"`
	NumdaysCompleted    float64 `json:"num_days_completed"`
	NumdaysOnHold       float64 `json:"num_days_on_hold"`
	NumdaysDropped      float64 `json:"num_days_dropped"`
	NumDays             float64 `json:"num_days"`
	NumEpisodes         int     `json:"num_episodes"`
	NumTimesRewatched   int     `json:"num_times_rewatched"`
	MeanScore           float64 `json:"mean_score"`
}

// ForumBoardCategories is forum board categories model.
type ForumBoardCategories struct {
	Categories []ForumBoardCategory `json:"categories"`
}

// ForumBoardCategory is forum board category model.
type ForumBoardCategory struct {
	Title  string       `json:"title"`
	Boards []ForumBoard `json:"boards"`
}

// ForumBoard is forum board model.
type ForumBoard struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Subboards   []ForumSubboard `json:"subboards"`
}

// ForumSubboard is forum subboard model.
type ForumSubboard struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// ForumTopicPaging is forum topic paging model.
type ForumTopicPaging struct {
	Data   []ForumTopic `json:"data"`
	Paging Paging       `json:"paging"`
}

// ForumTopic is forum topic model.
type ForumTopic struct {
	ID                int                 `json:"id"`
	Title             string              `json:"title"`
	CreatedAt         time.Time           `json:"created_at"`
	CreatedBy         ForumTopicCreatedBy `json:"created_by"`
	NumberOfPosts     int                 `json:"number_of_posts"`
	LastPostCreatedAt time.Time           `json:"last_post_created_at"`
	LastPostCreatedBy ForumTopicCreatedBy `json:"last_post_created_by"`
	IsLocked          bool                `json:"is_locked"`
}

// ForumTopicCreatedBy is forum topic created by model.
type ForumTopicCreatedBy struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ForumTopicDetailPaging is forum topic detail paging model.
type ForumTopicDetailPaging struct {
	Data   ForumTopicData `json:"data"`
	Paging Paging         `json:"paging"`
}

// ForumTopicData is forum topic detail model.
type ForumTopicData struct {
	Title string      `json:"title"`
	Posts []ForumPost `json:"posts"`
	Poll  ForumPoll   `json:"poll"`
}

// ForumPost is forum post model.
type ForumPost struct {
	ID        int                `json:"id"`
	Number    int                `json:"number"`
	CreatedAt time.Time          `json:"created_at"`
	CreatedBy ForumPostCreatedBy `json:"created_by"`
	Body      string             `json:"body"`
	Signature string             `json:"signature"`
}

// ForumPostCreatedBy is forum post created by model.
type ForumPostCreatedBy struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ForumAvatar string `json:"forum_avator"` // Typo from API.
}

// ForumPoll is forum poll model.
type ForumPoll struct {
	ID       int               `json:"id"`
	Question string            `json:"question"`
	Close    bool              `json:"close"`
	Options  []ForumPollOption `json:"options"`
}

// ForumPollOption is forum poll option model.
type ForumPollOption struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

// Animeography is animeography model.
//
// Undocumented.
type Animeography struct {
	Node Anime  `json:"node"`
	Role string `json:"role"`
}
