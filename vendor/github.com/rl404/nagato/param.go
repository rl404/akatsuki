package nagato

type idParam struct {
	ID int `validate:"gt=0"`
}

type usernameParam struct {
	Username string `validate:"required" mod:"trim,default=@me"`
}

// GetAnimeListParam is get anime list param model.
type GetAnimeListParam struct {
	Query  string `validate:"required,gte=3,lte=64" mod:"trim"`
	NSFW   bool   ``
	Limit  int    `validate:"gt=0,lte=100" mod:"default=100"`
	Offset int    `validate:"gte=0"`
}

// GetAnimeRankingParam is get anime ranking param.
type GetAnimeRankingParam struct {
	RankingType RankingType `validate:"oneof=all airing upcoming tv ova movie special bypopularity favorite" mod:"trim,default=all"`
	NSFW        bool        ``
	Limit       int         `validate:"gt=0,lte=500" mod:"default=100"`
	Offset      int         `validate:"gte=0"`
}

// GetSeasonalAnimeParam is get seasonal anime param.
type GetSeasonalAnimeParam struct {
	Year   int                   `validate:"gt=0"`
	Season SeasonType            `validate:"required,oneof=winter spring summer fall" mod:"trim"`
	NSFW   bool                  ``
	Sort   SeasonalAnimeSortType `validate:"oneof=anime_num_list_users anime_score" mod:"trim,default=anime_num_list_users"`
	Limit  int                   `validate:"gt=0,lte=500" mod:"default=100"`
	Offset int                   `validate:"gte=0"`
}

// GetSuggestedAnimeParam is get suggested anime param.
type GetSuggestedAnimeParam struct {
	NSFW   bool ``
	Limit  int  `validate:"gt=0,lte=100" mod:"default=100"`
	Offset int  `validate:"gte=0"`
}

// GetMangaListParam is get manga list param model.
type GetMangaListParam struct {
	Query  string `validate:"required,gte=3,lte=64" mod:"trim"`
	NSFW   bool   ``
	Limit  int    `validate:"gt=0,lte=100" mod:"default=100"`
	Offset int    `validate:"gte=0"`
}

// GetMangaRankingParam is get manga ranking param.
type GetMangaRankingParam struct {
	RankingType RankingType `validate:"oneof=all manga novels oneshots doujin manhwa manhua bypopularity favorite" mod:"trim,default=all"`
	NSFW        bool        ``
	Limit       int         `validate:"gt=0,lte=500" mod:"default=100"`
	Offset      int         `validate:"gte=0"`
}

// GetUserAnimeListParam is get user anime list param.
type GetUserAnimeListParam struct {
	Username string              `validate:"required,gte=3,lte=64" mod:"trim"`
	Status   UserAnimeStatusType `validate:"oneof='' watching completed on_hold dropped plan_to_watch" mod:"trim"`
	NSFW     bool                ``
	Sort     UserAnimeSortType   `validate:"oneof=list_score list_updated_at anime_title anime_start_date anime_id" mod:"trim,default=anime_title"`
	Limit    int                 `validate:"gt=0,lte=1000" mod:"default=100"`
	Offset   int                 `validate:"gte=0"`
}

// GetUserMangaListParam is get user manga list param.
type GetUserMangaListParam struct {
	Username string              `validate:"required,gte=3,lte=64" mod:"trim"`
	Status   UserMangaStatusType `validate:"oneof='' reading completed on_hold dropped plan_to_read" mod:"trim"`
	NSFW     bool                ``
	Sort     UserMangaSortType   `validate:"oneof=list_score list_updated_at manga_title manga_start_date manga_id" mod:"trim,default=manga_title"`
	Limit    int                 `validate:"gt=0,lte=1000" mod:"default=100"`
	Offset   int                 `validate:"gte=0"`
}

// UpdateMyAnimeListStatusParam is update my anime list status param.
type UpdateMyAnimeListStatusParam struct {
	ID             int                 `validate:"required,gt=0"`
	Status         UserAnimeStatusType `validate:"oneof='' watching completed on_hold dropped plan_to_watch" mod:"trim,default=plan_to_watch"`
	IsRewatching   bool                ``
	Score          int                 `validate:"gte=0,lte=10"`
	Episode        int                 `validate:"gte=0"`
	Priority       PriorityType        `validate:"gte=0,lte=2"`
	RewatchedTimes int                 `validate:"gte=0"`
	RewatchValue   RewatchValueType    `validate:"gte=0,lte=5"`
	Tags           []string            ``
	Comment        string              `mod:"trim"`
	StartDate      Date                ``
	FinishDate     Date                ``
}

// UpdateMyMangaListStatusParam is update my manga list status param.
type UpdateMyMangaListStatusParam struct {
	ID          int                 `validate:"required,gt=0"`
	Status      UserMangaStatusType `validate:"oneof='' reading completed on_hold dropped plan_to_read" mod:"trim,default=plan_to_read"`
	IsRereading bool                ``
	Score       int                 `validate:"gte=0,lte=10"`
	Volume      int                 `validate:"gte=0"`
	Chapter     int                 `validate:"gte=0"`
	Priority    PriorityType        `validate:"gte=0,lte=2"`
	RereadTimes int                 `validate:"gte=0"`
	RereadValue RereadValueType     `validate:"gte=0,lte=5"`
	Tags        []string            ``
	Comment     string              `mod:"trim"`
	StartDate   Date                ``
	FinishDate  Date                ``
}

// GetForumTopicsParam is get forum topics param.
type GetForumTopicsParam struct {
	BoardID       int                `validate:"gte=0"`
	SubboardID    int                `validate:"gte=0"`
	Query         string             `mod:"trim"`
	TopicUsername string             `mod:"trim"`
	Username      string             `mod:"trim"`
	Sort          ForumTopicSortType `validate:"oneof=recent" mod:"trim,default=recent"`
	Limit         int                `validate:"gt=0,lte=100" mod:"default=100"`
	Offset        int                `validate:"gte=0"`
}

// GetForumTopicDetailsParam is get forum topic details param.
type GetForumTopicDetailsParam struct {
	ID     int `validate:"gt=0"`
	Limit  int `validate:"gt=0,lte=100" mod:"default=100"`
	Offset int `validate:"gte=0"`
}
