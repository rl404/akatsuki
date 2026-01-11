package api

import (
	"github.com/rl404/akatsuki/internal/delivery/grpc/schema"
	"github.com/rl404/akatsuki/internal/domain/anime/entity"
	userAnimeEntity "github.com/rl404/akatsuki/internal/domain/user_anime/entity"
)

func (api *API) typeFromEntity(_type entity.Type) schema.Type {
	switch _type {
	case entity.TypeTV:
		return schema.Type_TYPE_TV
	case entity.TypeOVA:
		return schema.Type_TYPE_OVA
	case entity.TypeONA:
		return schema.Type_TYPE_ONA
	case entity.TypeMovie:
		return schema.Type_TYPE_MOVIE
	case entity.TypeSpecial:
		return schema.Type_TYPE_SPECIAL
	case entity.TypeMusic:
		return schema.Type_TYPE_MUSIC
	case entity.TypeCM:
		return schema.Type_TYPE_CM
	case entity.TypePV:
		return schema.Type_TYPE_PV
	case entity.TypeTVSpecial:
		return schema.Type_TYPE_TV_SPECIAL
	default:
		return schema.Type_TYPE_UNKNOWN
	}
}

func (api *API) statusFromEntity(status entity.Status) schema.Status {
	switch status {
	case entity.StatusFinished:
		return schema.Status_STATUS_FINISHED
	case entity.StatusReleasing:
		return schema.Status_STATUS_RELEASING
	case entity.StatusNotYet:
		return schema.Status_STATUS_NOT_YET
	default:
		return schema.Status_STATUS_UNKNOWN
	}
}

func (api *API) seasonFromEntity(season entity.Season) schema.Season {
	switch season {
	case entity.SeasonWinter:
		return schema.Season_SEASON_WINTER
	case entity.SeasonSpring:
		return schema.Season_SEASON_SPRING
	case entity.SeasonSummer:
		return schema.Season_SEASON_SUMMER
	case entity.SeasonFall:
		return schema.Season_SEASON_FALL
	default:
		return schema.Season_SEASON_UNKNOWN
	}
}

func (api *API) dayFromEntity(day entity.Day) schema.Day {
	switch day {
	case entity.DayMonday:
		return schema.Day_DAY_MONDAY
	case entity.DayTuesday:
		return schema.Day_DAY_TUESDAY
	case entity.DayWednesday:
		return schema.Day_DAY_WEDNESDAY
	case entity.DayThursday:
		return schema.Day_DAY_THURSDAY
	case entity.DayFriday:
		return schema.Day_DAY_FRIDAY
	case entity.DaySaturday:
		return schema.Day_DAY_SATURDAY
	case entity.DaySunday:
		return schema.Day_DAY_SUNDAY
	case entity.DayOther:
		return schema.Day_DAY_OTHER
	default:
		return schema.Day_DAY_UNKNOWN
	}
}

func (api *API) sourceFromEntity(source entity.Source) schema.Source {
	switch source {
	case entity.SourceOriginal:
		return schema.Source_SOURCE_ORIGINAL
	case entity.SourceManga:
		return schema.Source_SOURCE_MANGA
	case entity.Source4Koma:
		return schema.Source_SOURCE_4_KOMA_MANGA
	case entity.SourceWebManga:
		return schema.Source_SOURCE_WEB_MANGA
	case entity.SourceDigitalManga:
		return schema.Source_SOURCE_DIGITAL_MANGA
	case entity.SourceNovel:
		return schema.Source_SOURCE_NOVEL
	case entity.SourceLightNovel:
		return schema.Source_SOURCE_LIGHT_NOVEL
	case entity.SourceVisualNovel:
		return schema.Source_SOURCE_VISUAL_NOVEL
	case entity.SourceGame:
		return schema.Source_SOURCE_GAME
	case entity.SourceCardGame:
		return schema.Source_SOURCE_CARD_GAME
	case entity.SourceBook:
		return schema.Source_SOURCE_BOOK
	case entity.SourcePictureBook:
		return schema.Source_SOURCE_PICTURE_BOOK
	case entity.SourceRadio:
		return schema.Source_SOURCE_RADIO
	case entity.SourceMusic:
		return schema.Source_SOURCE_MUSIC
	case entity.SourceOther:
		return schema.Source_SOURCE_OTHER
	case entity.SourceWebNovel:
		return schema.Source_SOURCE_WEB_NOVEL
	case entity.SourceMixedMedia:
		return schema.Source_SOURCE_MIXED_MEDIA
	default:
		return schema.Source_SOURCE_UNKNOWN
	}
}

func (api *API) ratingFromEntity(rating entity.Rating) schema.Rating {
	switch rating {
	case entity.RatingG:
		return schema.Rating_RATING_G
	case entity.RatingPG:
		return schema.Rating_RATING_PG
	case entity.RatingPG13:
		return schema.Rating_RATING_PG13
	case entity.RatingR:
		return schema.Rating_RATING_R
	case entity.RatingRPlus:
		return schema.Rating_RATING_R_PLUS
	case entity.RatingRX:
		return schema.Rating_RATING_RX
	default:
		return schema.Rating_RATING_UNKNOWN
	}
}

func (api *API) relationFromEntity(relation entity.Relation) schema.Relation {
	switch relation {
	case entity.RelationSequel:
		return schema.Relation_RELATION_SEQUEL
	case entity.RelationPrequel:
		return schema.Relation_RELATION_PREQUEL
	case entity.RelationAlternativeSetting:
		return schema.Relation_RELATION_ALTERNATIVE_SETTING
	case entity.RelationAlternativeVersion:
		return schema.Relation_RELATION_ALTERNATIVE_VERSION
	case entity.RelationSideStory:
		return schema.Relation_RELATION_SIDE_STORY
	case entity.RelationParentStory:
		return schema.Relation_RELATION_PARENT_STORY
	case entity.RelationSummary:
		return schema.Relation_RELATION_SUMMARY
	case entity.RelationFullStory:
		return schema.Relation_RELATION_FULL_STORY
	case entity.RelationSpinOff:
		return schema.Relation_RELATION_SPIN_OFF
	case entity.RelationAdaptation:
		return schema.Relation_RELATION_ADAPTATION
	case entity.RelationCharacter:
		return schema.Relation_RELATION_CHARACTER
	case entity.RelationOther:
		return schema.Relation_RELATION_OTHER
	default:
		return schema.Relation_RELATION_UNKNOWN
	}
}

func (api *API) userAnimeStatusFromEntity(status userAnimeEntity.Status) schema.UserAnimeStatus {
	switch status {
	case userAnimeEntity.StatusWatching:
		return schema.UserAnimeStatus_USER_ANIME_STATUS_WATCHING
	case userAnimeEntity.StatusCompleted:
		return schema.UserAnimeStatus_USER_ANIME_STATUS_COMPLETED
	case userAnimeEntity.StatusOnHold:
		return schema.UserAnimeStatus_USER_ANIME_STATUS_ON_HOLD
	case userAnimeEntity.StatusDropped:
		return schema.UserAnimeStatus_USER_ANIME_STATUS_DROPPED
	case userAnimeEntity.StatusPlanned:
		return schema.UserAnimeStatus_USER_ANIME_STATUS_PLANNED
	default:
		return schema.UserAnimeStatus_USER_ANIME_STATUS_UNKNOWN
	}
}
