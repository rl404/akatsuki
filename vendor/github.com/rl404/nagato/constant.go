package nagato

// NsfwType is nsfw types.
type NsfwType string

// Available nsfw types.
const (
	NsfwWhite NsfwType = "white"
	NsfwGray  NsfwType = "gray"
	NsfwBlack NsfwType = "black" // API never return this.
)

// MediaType is media types.
type MediaType string

// Available media types.
const (
	MediaUnknown MediaType = "unknown"

	// Anime.
	MediaTV      MediaType = "tv"
	MediaOVA     MediaType = "ova"
	MediaMovie   MediaType = "movie"
	MediaSpecial MediaType = "special"
	MediaONA     MediaType = "ona"
	MediaMusic   MediaType = "music"

	// Manga.
	MediaManga      MediaType = "manga"
	MediaNovel      MediaType = "novel"
	MediaOneShot    MediaType = "one_shot"
	MediaDoujinshi  MediaType = "doujinshi"
	MediaManhwa     MediaType = "manhwa"
	MediaManhua     MediaType = "manhua"
	MediaOEL        MediaType = "oel"
	MediaLightNovel MediaType = "light_novel" // Undocumented.
)

// StatusType is status types.
type StatusType string

// Available status types.
const (
	// Anime.
	StatusFinishedAiring  StatusType = "finished_airing"
	StatusCurrentlyAiring StatusType = "currently_airing"
	StatusNotYetAired     StatusType = "not_yet_aired"

	// Manga.
	StatusFinishedPublishing  StatusType = "finished"
	StatusCurrentlyPublishing StatusType = "currently_publishing"
	StatusNotYetPublished     StatusType = "not_yet_published"
	StatusOnHiatus            StatusType = "on_hiatus"    // Undocumented.
	StatusDiscontinued        StatusType = "discontinued" // Undocumented.
)

// UserAnimeStatusType is user anime status types.
type UserAnimeStatusType string

// Available user anime status types.
const (
	UserAnimeStatusWatching    UserAnimeStatusType = "watching"
	UserAnimeStatusCompleted   UserAnimeStatusType = "completed"
	UserAnimeStatusOnHold      UserAnimeStatusType = "on_hold"
	UserAnimeStatusDropped     UserAnimeStatusType = "dropped"
	UserAnimeStatusPlanToWatch UserAnimeStatusType = "plan_to_watch"
)

// UserMangaStatusType is user manga status types.
type UserMangaStatusType string

// Available user manga status types.
const (
	UserMangaStatusReading    UserMangaStatusType = "reading"
	UserMangaStatusCompleted  UserMangaStatusType = "completed"
	UserMangaStatusOnHold     UserMangaStatusType = "on_hold"
	UserMangaStatusDropped    UserMangaStatusType = "dropped"
	UserMangaStatusPlanToRead UserMangaStatusType = "plan_to_read"
)

// PriorityType is priority types.
type PriorityType int

// Available priority types.
const (
	PriorityLow PriorityType = iota
	PriorityMedium
	PriorityHigh
)

// RewatchValueType is rewatch value types.
type RewatchValueType int

// Available rewatch value types.
const (
	RewatchValueVeryLow RewatchValueType = iota + 1
	RewatchValueLow
	RewatchValueMedium
	RewatchValueHigh
	RewatchValueVeryHigh
)

// RereadValueType is reread value types.
type RereadValueType int

// Available reread value types.
const (
	RereadValueVeryLow RereadValueType = iota + 1
	RereadValueLow
	RereadValueMedium
	RereadValueHigh
	RereadValueVeryHigh
)

// SeasonType is season types.
type SeasonType string

// Available season types.
const (
	SeasonWinter SeasonType = "winter"
	SeasonSpring SeasonType = "spring"
	SeasonSummer SeasonType = "summer"
	SeasonFall   SeasonType = "fall"
)

// DayType is day types.
type DayType string

// Available day types.
const (
	DayMonday    DayType = "monday"
	DayTuesday   DayType = "tuesday"
	DayWednesday DayType = "wednesday"
	DayThursday  DayType = "thursday"
	DayFriday    DayType = "friday"
	DaySaturday  DayType = "saturday"
	DaySunday    DayType = "sunday"
	DayOther     DayType = "other"
)

// SourceType is source types.
type SourceType string

// Available source types.
const (
	SourceOther        SourceType = "other"
	SourceOriginal     SourceType = "original"
	SourceManga        SourceType = "manga"
	Source4KomaManga   SourceType = "4_koma_manga"
	SourceWebManga     SourceType = "web_manga"
	SourceDigitalManga SourceType = "digital_manga"
	SourceNovel        SourceType = "novel"
	SourceLightNovel   SourceType = "light_novel"
	SourceVisualNovel  SourceType = "visual_novel"
	SourceGame         SourceType = "game"
	SourceCardGame     SourceType = "card_game"
	SourceBook         SourceType = "book"
	SourcePictureBook  SourceType = "picture_book"
	SourceRadio        SourceType = "radio"
	SourceMusic        SourceType = "music"
	SourceWebNovel     SourceType = "web_novel"   // Undocumented.
	SourceMixedMedia   SourceType = "mixed_media" // Undocumented.
)

// RatingType is rating types.
type RatingType string

// Available rating types.
const (
	RatingG     RatingType = "g"
	RatingPG    RatingType = "pg"
	RatingPG13  RatingType = "pg_13"
	RatingR     RatingType = "r"
	RatingRPlus RatingType = "r+"
	RatingRX    RatingType = "rx"
)

// RelationType is relation types.
type RelationType string

// Available relation types.
const (
	RelationSequel             RelationType = "sequel"
	RelationPrequel            RelationType = "prequel"
	RelationAlternativeSetting RelationType = "alternative_setting"
	RelationAlternativeVersion RelationType = "alternative_version"
	RelationSideStory          RelationType = "side_story"
	RelationParentStory        RelationType = "parent_story"
	RelationSummary            RelationType = "summary"
	RelationFullStory          RelationType = "full_story"
	RelationSpinOff            RelationType = "spin_off"  // Undocumented.
	RelationOther              RelationType = "other"     // Undocumented.
	RelationCharacter          RelationType = "character" // Undocumented.
)

// RankingType is ranking types.
type RankingType string

// Available ranking types.
const (
	RankingAll          RankingType = "all"
	RankingByPopularity RankingType = "bypopularity"
	RankingFavorite     RankingType = "favorite"

	// Anime.
	RankingAiring   RankingType = "airing"
	RankingUpcoming RankingType = "upcoming"
	RankingTV       RankingType = "tv"
	RankingOVA      RankingType = "ova"
	RankingMovie    RankingType = "movie"
	RankingSpecial  RankingType = "special"

	// Manga.
	RankingManga   RankingType = "manga"
	RankingNovel   RankingType = "novels"
	RankingOneShot RankingType = "oneshots"
	RankingDoujin  RankingType = "doujin"
	RankingManhwa  RankingType = "manhwa"
	RankingManhua  RankingType = "manhua"
)

// SeasonalAnimeSortType is seasonal anime sort types.
type SeasonalAnimeSortType string

// Available seasonal anime sort types.
const (
	SeasonalAnimeSortByScore        SeasonalAnimeSortType = "anime_score" // Not working.
	SeasonalAnimeSortByNumListUsers SeasonalAnimeSortType = "anime_num_list_users"
)

// UserAnimeSortType is user anime sort types.
type UserAnimeSortType string

// Available user anime sort types.
const (
	UserAnimeSortScore     UserAnimeSortType = "list_score"
	UserAnimeSortUpdatedAt UserAnimeSortType = "list_updated_at"
	UserAnimeSortTitle     UserAnimeSortType = "anime_title"
	UserAnimeSortStartDate UserAnimeSortType = "anime_start_date"
	UserAnimeSortID        UserAnimeSortType = "anime_id" // Not working.
)

// UserMangaSortType is user manga sort types.
type UserMangaSortType string

// Available user manga sort types.
const (
	UserMangaSortScore     UserMangaSortType = "list_score"
	UserMangaSortUpdatedAt UserMangaSortType = "list_updated_at"
	UserMangaSortTitle     UserMangaSortType = "manga_title"
	UserMangaSortStartDate UserMangaSortType = "manga_start_date"
	UserMangaSortID        UserMangaSortType = "manga_id" // Not working.
)
