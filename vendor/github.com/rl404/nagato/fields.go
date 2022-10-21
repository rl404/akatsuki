package nagato

import (
	"fmt"
	"strings"
)

// AnimeField is anime fields.
type AnimeField string

// Available anime fields.
//
// id, title, main_picture are returned by default.
const (
	AnimeFieldAlternativeTitles      AnimeField = "alternative_titles"
	AnimeFieldStartDate              AnimeField = "start_date"
	AnimeFieldEndDate                AnimeField = "end_date"
	AnimeFieldSynopsis               AnimeField = "synopsis"
	AnimeFieldMean                   AnimeField = "mean"
	AnimeFieldRank                   AnimeField = "rank"
	AnimeFieldPopularity             AnimeField = "popularity"
	AnimeFieldNumListUsers           AnimeField = "num_list_users"
	AnimeFieldNumScoringUsers        AnimeField = "num_scoring_users"
	AnimeFieldNSFW                   AnimeField = "nsfw"
	AnimeFieldGenres                 AnimeField = "genres"
	AnimeFieldMediaType              AnimeField = "media_type"
	AnimeFieldStatus                 AnimeField = "status"
	AnimeFieldMyListStatus           AnimeField = "my_list_status" // Need oauth2.
	AnimeFieldNumEpisodes            AnimeField = "num_episodes"
	AnimeFieldStartSeason            AnimeField = "start_season"
	AnimeFieldBroadcast              AnimeField = "broadcast"
	AnimeFieldSource                 AnimeField = "source"
	AnimeFieldAverageEpisodeDuration AnimeField = "average_episode_duration"
	AnimeFieldRating                 AnimeField = "rating"
	AnimeFieldStudios                AnimeField = "studios"
	AnimeFieldPictures               AnimeField = "pictures"       // Can't include this field in a list.
	AnimeFieldBackground             AnimeField = "background"     // Can't include this field in a list.
	AnimeFieldStatistics             AnimeField = "statistics"     // Can't include this field in a list.
	AnimeFieldNumFavorites           AnimeField = "num_favorites"  // Undocumented.
	AnimeFieldOpeningThemes          AnimeField = "opening_themes" // Undocumented & can't include this field in a list.
	AnimeFieldEndingThemes           AnimeField = "ending_themes"  // Undocumented & can't include this field in a list.
	AnimeFieldVideos                 AnimeField = "videos"         // Undocumented & can't include this field in a list.
)

// AnimeFieldRelatedAnime is to generate related anime field.
//
// Can't include this field in a list.
func AnimeFieldRelatedAnime(fields ...AnimeField) AnimeField {
	nested := ""
	if len(fields) > 0 {
		strs := make([]string, len(fields))
		for i, f := range fields {
			strs[i] = string(f)
		}
		nested = fmt.Sprintf("{%s}", strings.Join(strs, ","))
	}
	return AnimeField("related_anime" + nested)
}

// AnimeFieldRelatedManga is to generate related manga field.
//
// Related manga always return empty from API.
//
// Can't include this field in a list.
func AnimeFieldRelatedManga(fields ...MangaField) AnimeField {
	nested := ""
	if len(fields) > 0 {
		strs := make([]string, len(fields))
		for i, f := range fields {
			strs[i] = string(f)
		}
		nested = fmt.Sprintf("{%s}", strings.Join(strs, ","))
	}
	return AnimeField("related_manga" + nested)
}

// AnimeFieldRecommendations is to generate recommendations field.
//
// Can't include this field in a list.
func AnimeFieldRecommendations(fields ...AnimeField) AnimeField {
	nested := ""
	if len(fields) > 0 {
		strs := make([]string, len(fields))
		for i, f := range fields {
			strs[i] = string(f)
		}
		nested = fmt.Sprintf("{%s}", strings.Join(strs, ","))
	}
	return AnimeField("recommendations" + nested)
}

// UserAnimeField is user anime fields.
type UserAnimeField string

// Available user anime fields.
//
// status, score, num_episodes_watched, is_rewatching,
// start_date, finish_date, priority are returned by default.
const (
	UserAnimeNumTimesRewatched UserAnimeField = "num_times_rewatched"
	UserAnimeRewatchValue      UserAnimeField = "rewatch_value"
	UserAnimeTags              UserAnimeField = "tags"
	UserAnimeComments          UserAnimeField = "comments"
)

// AnimeFieldUserStatus is to generate user status field.
//
// Used only for getting user anime list.
func AnimeFieldUserStatus(fields ...UserAnimeField) AnimeField {
	nested := ""
	if len(fields) > 0 {
		strs := make([]string, len(fields))
		for i, f := range fields {
			strs[i] = string(f)
		}
		nested = fmt.Sprintf("{%s}", strings.Join(strs, ","))
	}
	return AnimeField("list_status" + nested)
}

// MangaField is manga fields.
type MangaField string

// Available manga fields.
//
// id, title, main_picture are returned by default.
const (
	MangaFieldAlternativeTitles MangaField = "alternative_titles"
	MangaFieldStartDate         MangaField = "start_date"
	MangaFieldEndDate           MangaField = "end_date"
	MangaFieldSynopsis          MangaField = "synopsis"
	MangaFieldMean              MangaField = "mean"
	MangaFieldRank              MangaField = "rank"
	MangaFieldPopularity        MangaField = "popularity"
	MangaFieldNumListUsers      MangaField = "num_list_users"
	MangaFieldNumScoringUsers   MangaField = "num_scoring_users"
	MangaFieldNSFW              MangaField = "nsfw"
	MangaFieldGenres            MangaField = "genres"
	MangaFieldMediaType         MangaField = "media_type"
	MangaFieldStatus            MangaField = "status"
	MangaFieldMyListStatus      MangaField = "my_list_status" // Need oauth2.
	MangaFieldNumVolumes        MangaField = "num_volumes"
	MangaFieldNumChapters       MangaField = "num_chapters"
	MangaFieldAuthors           MangaField = "authors{first_name,last_name}"
	MangaFieldPictures          MangaField = "pictures"      // Can't include this field in a list.
	MangaFieldBackground        MangaField = "background"    // Can't include this field in a list.
	MangaFieldSerialization     MangaField = "serialization" // Can't include this field in a list.
	MangaFieldNumFavorites      MangaField = "num_favorites" // Undocumented.
)

// MangaFieldRelatedAnime is to generate related anime field.
//
// Related anime always return empty from API.
//
// Can't include this field in a list.
func MangaFieldRelatedAnime(fields ...AnimeField) MangaField {
	nested := ""
	if len(fields) > 0 {
		strs := make([]string, len(fields))
		for i, f := range fields {
			strs[i] = string(f)
		}
		nested = fmt.Sprintf("{%s}", strings.Join(strs, ","))
	}
	return MangaField("related_anime" + nested)
}

// MangaFieldRelatedManga is to generate related manga field.
//
// Can't include this field in a list.
func MangaFieldRelatedManga(fields ...MangaField) MangaField {
	nested := ""
	if len(fields) > 0 {
		strs := make([]string, len(fields))
		for i, f := range fields {
			strs[i] = string(f)
		}
		nested = fmt.Sprintf("{%s}", strings.Join(strs, ","))
	}
	return MangaField("related_manga" + nested)
}

// MangaFieldRecommendations is to generate recommendations field.
//
// Can't include this field in a list.
func MangaFieldRecommendations(fields ...MangaField) MangaField {
	nested := ""
	if len(fields) > 0 {
		strs := make([]string, len(fields))
		for i, f := range fields {
			strs[i] = string(f)
		}
		nested = fmt.Sprintf("{%s}", strings.Join(strs, ","))
	}
	return MangaField("recommendations" + nested)
}

// UserMangaField is user manga fields.
type UserMangaField string

// Available user manga fields.
//
// status, score, num_volumes_read, num_chapters_read, is_rereading,
// start_date, finish_date, priority are returned by default.
const (
	UserMangaNumTimesReread UserMangaField = "num_times_reread"
	UserMangaRereadValue    UserMangaField = "reread_value"
	UserMangaTags           UserMangaField = "tags"
	UserMangaComments       UserMangaField = "comments"
)

// MangaFieldUserStatus is to generate user status field.
//
// Used only for getting user manga list.
func MangaFieldUserStatus(fields ...UserMangaField) MangaField {
	nested := ""
	if len(fields) > 0 {
		strs := make([]string, len(fields))
		for i, f := range fields {
			strs[i] = string(f)
		}
		nested = fmt.Sprintf("{%s}", strings.Join(strs, ","))
	}
	return MangaField("list_status" + nested)
}
