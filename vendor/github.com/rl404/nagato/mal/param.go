package mal

import (
	"net/url"
	"strconv"
)

// GetAnimeListParam is get anime list param.
type GetAnimeListParam struct {
	Query  string
	Nsfw   bool
	Limit  int
	Offset int
}

// GetAnimeRankingParam is get anime ranking param.
type GetAnimeRankingParam struct {
	RankingType string
	Nsfw        bool
	Limit       int
	Offset      int
}

// GetSeasonalAnimeParam is get seasonal anime param.
type GetSeasonalAnimeParam struct {
	Year   int
	Season string
	Nsfw   bool
	Sort   string
	Limit  int
	Offset int
}

// GetSuggestedAnimeParam is get suggested anime param.
type GetSuggestedAnimeParam struct {
	Nsfw   bool
	Limit  int
	Offset int
}

// GetMangaListParam is get manga list param.
type GetMangaListParam struct {
	Query  string
	Nsfw   bool
	Limit  int
	Offset int
}

// GetMangaRankingParam is get manga ranking param.
type GetMangaRankingParam struct {
	RankingType string
	Nsfw        bool
	Limit       int
	Offset      int
}

// UpdateMyAnimeListStatusParam is update my anime list status param.
type UpdateMyAnimeListStatusParam struct {
	ID                 int
	Status             string
	IsRewatching       bool
	Score              int
	NumWatchedEpisodes int
	Priority           int
	NumTimesRewatched  int
	RewatchValue       int
	Tags               string
	Comments           string
	StartDate          string // Undocumented.
	FinishDate         string // Undocumented.
}

func (p *UpdateMyAnimeListStatusParam) encode() []byte {
	data := url.Values{}
	data.Set("status", p.Status)
	data.Set("is_rewatching", strconv.FormatBool(p.IsRewatching))
	data.Set("score", strconv.Itoa(p.Score))
	data.Set("num_watched_episodes", strconv.Itoa(p.NumWatchedEpisodes))
	data.Set("priority", strconv.Itoa(p.Priority))
	data.Set("num_times_rewatched", strconv.Itoa(p.NumTimesRewatched))
	data.Set("rewatch_value", strconv.Itoa(p.RewatchValue))
	data.Set("tags", p.Tags)
	data.Set("comments", p.Comments)
	data.Set("start_date", p.StartDate)
	data.Set("finish_date", p.FinishDate)
	return []byte(data.Encode())
}

// GetUserAnimeListParam is get user anime list param.
type GetUserAnimeListParam struct {
	Username string
	Status   string
	Nsfw     bool
	Sort     string
	Limit    int
	Offset   int
}

// UpdateMyMangaListStatusParam is update my manga list status param.
type UpdateMyMangaListStatusParam struct {
	ID              int
	Status          string
	IsRereading     bool
	Score           int
	NumVolumesRead  int
	NumChaptersRead int
	Priority        int
	NumTimesReread  int
	RereadValue     int
	Tags            string
	Comments        string
	StartDate       string // Undocumented.
	FinishDate      string // Undocumented.
}

func (p *UpdateMyMangaListStatusParam) encode() []byte {
	data := url.Values{}
	data.Set("status", p.Status)
	data.Set("is_rereading", strconv.FormatBool(p.IsRereading))
	data.Set("score", strconv.Itoa(p.Score))
	data.Set("num_volumes_read", strconv.Itoa(p.NumVolumesRead))
	data.Set("num_chapters_read", strconv.Itoa(p.NumChaptersRead))
	data.Set("priority", strconv.Itoa(p.Priority))
	data.Set("num_times_reread", strconv.Itoa(p.NumTimesReread))
	data.Set("reread_value", strconv.Itoa(p.RereadValue))
	data.Set("tags", p.Tags)
	data.Set("comments", p.Comments)
	data.Set("start_date", p.StartDate)
	data.Set("finish_date", p.FinishDate)
	return []byte(data.Encode())
}

// GetUserMangaListParam is get user manga list param.
type GetUserMangaListParam struct {
	Username string
	Status   string
	Nsfw     bool
	Sort     string
	Limit    int
	Offset   int
}

// GetForumTopicsParam is get forum topics param.
type GetForumTopicsParam struct {
	BoardID       int
	SubboardID    int
	Query         string
	TopicUsername string
	Username      string
	Sort          string // Only `recent` for now.
	Limit         int
	Offset        int
}

// GetForumTopicDetailsParam is get forum topic details param.
type GetForumTopicDetailsParam struct {
	ID     int
	Limit  int
	Offset int
}
