package nagato

import (
	"strconv"
	"time"

	"github.com/rl404/nagato/mal"
)

func (c *Client) animeFieldsToStrs(fields ...AnimeField) []string {
	strs := make([]string, len(fields))
	for i, f := range fields {
		strs[i] = string(f)
	}
	return strs
}

func (c *Client) mangaFieldsToStrs(fields ...MangaField) []string {
	strs := make([]string, len(fields))
	for i, f := range fields {
		strs[i] = string(f)
	}
	return strs
}

func (c *Client) userFieldsToStrs(fields ...UserField) []string {
	strs := make([]string, len(fields))
	for i, f := range fields {
		strs[i] = string(f)
	}
	return strs
}

func (c *Client) animeToAnime(anime *mal.Anime) *Anime {
	if anime == nil {
		return nil
	}

	return &Anime{
		ID:    anime.ID,
		Title: anime.Title,
		MainPicture: Picture{
			Large:  anime.MainPicture.Large,
			Medium: anime.MainPicture.Medium,
		},
		AlternativeTitles: AlternativeTitles{
			Synonyms: anime.AlternativeTitles.Synonyms,
			English:  anime.AlternativeTitles.En,
			Japanese: anime.AlternativeTitles.Ja,
		},
		StartDate:       c.dateToDate(anime.StartDate),
		EndDate:         c.dateToDate(anime.EndDate),
		Synopsis:        anime.Synopsis,
		Mean:            anime.Mean,
		Rank:            anime.Rank,
		Popularity:      anime.Popularity,
		NumListUsers:    anime.NumListUsers,
		NumScoringUsers: anime.NumScoringUsers,
		NSFW:            c.nsfwToNsfw(anime.NSFW),
		Genres:          c.genreToGenre(anime.Genres),
		MediaType:       c.mediaToMedia(anime.MediaType),
		Status:          c.statusToStatus(anime.Status),
		MyListStatus: UserAnimeListStatus{
			Status:             c.listStatusToUserAnimeStatus(anime.MyListStatus.Status),
			Score:              anime.MyListStatus.Score,
			NumEpisodesWatched: anime.MyListStatus.NumEpisodesWatched,
			IsRewatching:       anime.MyListStatus.IsRewatching,
			StartDate:          c.dateToDate(anime.MyListStatus.StartDate),
			FinishDate:         c.dateToDate(anime.MyListStatus.FinishDate),
			Priority:           c.priorityToPriority(anime.MyListStatus.Priority),
			NumTimesRewatched:  anime.MyListStatus.NumTimesRewatched,
			RewatchValue:       c.rewatchToRewatch(anime.MyListStatus.RewatchValue),
			Tags:               anime.MyListStatus.Tags,
			Comments:           anime.MyListStatus.Comments,
			UpdatedAt:          anime.MyListStatus.UpdatedAt,
		},
		NumEpisodes: anime.NumEpisodes,
		StartSeason: Season{
			Season: c.seasonToSeason(anime.StartSeason.Season),
			Year:   anime.StartSeason.Year,
		},
		Broadcast: Broadcast{
			DayOfTheWeek: c.dayToDay(anime.Broadcast.DayOfTheWeek),
			StartTime:    anime.Broadcast.StartTime,
		},
		Source:                 c.sourceToSource(anime.Source),
		AverageEpisodeDuration: time.Duration(anime.AverageEpisodeDuration) * time.Second,
		Rating:                 c.ratingToRating(anime.Rating),
		Studios:                c.studioToStudio(anime.Studios),
		Pictures:               c.pictureToPicture(anime.Pictures),
		Background:             anime.Background,
		RelatedAnime:           c.relatedAnimeToRelatedAnime(anime.RelatedAnime),
		RelatedManga:           c.relatedMangaToRelatedManga(anime.RelatedManga),
		Recommendations:        c.animeRecommendationToAnimeRecommendation(anime.Recommendations),
		Statistics: Statistic{
			NumListUsers: anime.Statistics.NumListUsers,
			Status: StatisticStatus{
				Watching:    c.statisticStatusToStatisticStatus(anime.Statistics.Status.Watching),
				Completed:   c.statisticStatusToStatisticStatus(anime.Statistics.Status.Completed),
				OnHold:      c.statisticStatusToStatisticStatus(anime.Statistics.Status.OnHold),
				Dropped:     c.statisticStatusToStatisticStatus(anime.Statistics.Status.Dropped),
				PlanToWatch: c.statisticStatusToStatisticStatus(anime.Statistics.Status.PlanToWatch),
			},
		},
		NumFavorites:  anime.NumFavorites,
		OpeningThemes: c.songToSong(anime.OpeningThemes),
		EndingThemes:  c.songToSong(anime.EndingThemes),
		Videos:        c.videoToVideo(anime.Videos),
	}
}

func (c *Client) animePagingToAnimeList(anime *mal.AnimePaging) []Anime {
	if anime == nil {
		return nil
	}

	res := make([]Anime, len(anime.Data))
	for i, a := range anime.Data {
		res[i] = *c.animeToAnime(&a.Node)
	}

	return res
}

func (c *Client) animeRankingPagingToAnimeList(anime *mal.AnimeRankingPaging) []Anime {
	if anime == nil {
		return nil
	}

	res := make([]Anime, len(anime.Data))
	for i, a := range anime.Data {
		res[i] = *c.animeToAnime(&a.Node)
	}

	return res
}

func (c *Client) seasonalAnimePagingToAnimeList(anime *mal.SeasonalAnimePaging) []Anime {
	if anime == nil {
		return nil
	}

	res := make([]Anime, len(anime.Data))
	for i, a := range anime.Data {
		res[i] = *c.animeToAnime(&a.Node)
	}

	return res
}

func (c *Client) suggestedAnimePagingToAnimeList(anime *mal.SuggestedAnimePaging) []Anime {
	if anime == nil {
		return nil
	}

	res := make([]Anime, len(anime.Data))
	for i, a := range anime.Data {
		res[i] = *c.animeToAnime(&a.Node)
	}

	return res
}

func (c *Client) nsfwToNsfw(str string) NsfwType {
	return map[string]NsfwType{
		"white": NsfwWhite,
		"gray":  NsfwGray,
		"black": NsfwBlack,
	}[str]
}

func (c *Client) genreToGenre(genres []mal.Genre) []Genre {
	res := make([]Genre, len(genres))
	for i, g := range genres {
		res[i] = Genre{
			ID:   g.ID,
			Name: g.Name,
		}
	}
	return res
}

func (c *Client) mediaToMedia(t string) MediaType {
	return map[string]MediaType{
		"unknown":     MediaUnknown,
		"tv":          MediaTV,
		"ova":         MediaOVA,
		"movie":       MediaMovie,
		"special":     MediaSpecial,
		"ona":         MediaONA,
		"music":       MediaMusic,
		"pv":          MediaPV,
		"cm":          MediaCM,
		"tv_special":  MediaTVSpecial,
		"manga":       MediaManga,
		"novel":       MediaNovel,
		"one_shot":    MediaOneShot,
		"doujinshi":   MediaDoujinshi,
		"manhwa":      MediaManhwa,
		"manhua":      MediaManhua,
		"oel":         MediaOEL,
		"light_novel": MediaLightNovel,
	}[t]
}

func (c *Client) statusToStatus(s string) StatusType {
	return map[string]StatusType{
		"finished_airing":      StatusFinishedAiring,
		"currently_airing":     StatusCurrentlyAiring,
		"not_yet_aired":        StatusNotYetAired,
		"finished":             StatusFinishedPublishing,
		"currently_publishing": StatusCurrentlyPublishing,
		"not_yet_published":    StatusNotYetPublished,
		"on_hiatus":            StatusOnHiatus,
		"discontinued":         StatusDiscontinued,
	}[s]
}

func (c *Client) listStatusToUserAnimeStatus(s string) UserAnimeStatusType {
	return map[string]UserAnimeStatusType{
		"watching":      UserAnimeStatusWatching,
		"completed":     UserAnimeStatusCompleted,
		"on_hold":       UserAnimeStatusOnHold,
		"dropped":       UserAnimeStatusDropped,
		"plan_to_watch": UserAnimeStatusPlanToWatch,
	}[s]
}

func (c *Client) listStatusToUserMangaStatus(s string) UserMangaStatusType {
	return map[string]UserMangaStatusType{
		"reading":      UserMangaStatusReading,
		"completed":    UserMangaStatusCompleted,
		"on_hold":      UserMangaStatusOnHold,
		"dropped":      UserMangaStatusDropped,
		"plan_to_read": UserMangaStatusPlanToRead,
	}[s]
}

func (c *Client) priorityToPriority(p int) PriorityType {
	return map[int]PriorityType{
		0: PriorityLow,
		1: PriorityMedium,
		2: PriorityHigh,
	}[p]
}

func (c *Client) rewatchToRewatch(r int) RewatchValueType {
	return map[int]RewatchValueType{
		1: RewatchValueVeryLow,
		2: RewatchValueLow,
		3: RewatchValueMedium,
		4: RewatchValueHigh,
		5: RewatchValueVeryHigh,
	}[r]
}

func (c *Client) rereadToReread(r int) RereadValueType {
	return map[int]RereadValueType{
		1: RereadValueVeryLow,
		2: RereadValueLow,
		3: RereadValueMedium,
		4: RereadValueHigh,
		5: RereadValueVeryHigh,
	}[r]
}

func (c *Client) seasonToSeason(s string) SeasonType {
	return map[string]SeasonType{
		"winter": SeasonWinter,
		"spring": SeasonSpring,
		"summer": SeasonSummer,
		"fall":   SeasonFall,
	}[s]
}

func (c *Client) dayToDay(d string) DayType {
	return map[string]DayType{
		"monday":    DayMonday,
		"tuesday":   DayTuesday,
		"wednesday": DayWednesday,
		"thursday":  DayThursday,
		"friday":    DayFriday,
		"saturday":  DaySaturday,
		"sunday":    DaySunday,
		"other":     DayOther,
	}[d]
}

func (c *Client) sourceToSource(s string) SourceType {
	return map[string]SourceType{
		"other":         SourceOther,
		"original":      SourceOriginal,
		"manga":         SourceManga,
		"4_koma_manga":  Source4KomaManga,
		"web_manga":     SourceWebManga,
		"digital_manga": SourceDigitalManga,
		"novel":         SourceNovel,
		"light_novel":   SourceLightNovel,
		"visual_novel":  SourceVisualNovel,
		"game":          SourceGame,
		"card_game":     SourceCardGame,
		"book":          SourceBook,
		"picture_book":  SourcePictureBook,
		"radio":         SourceRadio,
		"music":         SourceMusic,
		"web_novel":     SourceWebNovel,
		"mixed_media":   SourceMixedMedia,
	}[s]
}

func (c *Client) ratingToRating(r string) RatingType {
	return map[string]RatingType{
		"g":     RatingG,
		"pg":    RatingPG,
		"pg_13": RatingPG13,
		"r":     RatingR,
		"r+":    RatingRPlus,
		"rx":    RatingRX,
	}[r]
}

func (c *Client) studioToStudio(studios []mal.Studio) []Studio {
	res := make([]Studio, len(studios))
	for i, s := range studios {
		res[i] = Studio{
			ID:   s.ID,
			Name: s.Name,
		}
	}
	return res
}

func (c *Client) pictureToPicture(pictures []mal.Picture) []Picture {
	res := make([]Picture, len(pictures))
	for i, p := range pictures {
		res[i] = Picture{
			Large:  p.Large,
			Medium: p.Medium,
		}
	}
	return res
}

func (c *Client) relatedAnimeToRelatedAnime(anime []mal.RelatedAnime) []RelatedAnime {
	res := make([]RelatedAnime, len(anime))
	for i, a := range anime {
		res[i] = RelatedAnime{
			Anime:        *c.animeToAnime(&a.Node),
			RelationType: c.relationToRelation(a.RelationType),
		}
	}
	return res
}

func (c *Client) relatedMangaToRelatedManga(manga []mal.RelatedManga) []RelatedManga {
	res := make([]RelatedManga, len(manga))
	for i, m := range manga {
		res[i] = RelatedManga{
			Manga:        *c.mangaToManga(&m.Node),
			RelationType: c.relationToRelation(m.RelationType),
		}
	}
	return res
}

func (c *Client) relationToRelation(r string) RelationType {
	return map[string]RelationType{
		"sequel":              RelationSequel,
		"prequel":             RelationPrequel,
		"alternative_setting": RelationAlternativeSetting,
		"alternative_version": RelationAlternativeVersion,
		"side_story":          RelationSideStory,
		"parent_story":        RelationParentStory,
		"summary":             RelationSummary,
		"full_story":          RelationFullStory,
		"spin_off":            RelationSpinOff,
		"other":               RelationOther,
		"character":           RelationCharacter,
		"adaptation":          RelationAdaptation,
	}[r]
}

func (c *Client) animeRecommendationToAnimeRecommendation(anime []mal.AnimeRecommendation) []AnimeRecommendation {
	res := make([]AnimeRecommendation, len(anime))
	for i, a := range anime {
		res[i] = AnimeRecommendation{
			Anime:              *c.animeToAnime(&a.Node),
			NumRecommendations: a.NumRecommendations,
		}
	}
	return res
}

func (c *Client) mangaRecommendationToMangaRecommendation(manga []mal.MangaRecommendation) []MangaRecommendation {
	res := make([]MangaRecommendation, len(manga))
	for i, m := range manga {
		res[i] = MangaRecommendation{
			Manga:              *c.mangaToManga(&m.Node),
			NumRecommendations: m.NumRecommendations,
		}
	}
	return res
}

func (c *Client) statisticStatusToStatisticStatus(s interface{}) int {
	if v, ok := s.(int); ok {
		return v
	}

	if v, ok := s.(string); ok {
		vv, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return vv
	}

	return 0
}

func (c *Client) songToSong(songs []mal.ThemeSong) []ThemeSong {
	res := make([]ThemeSong, len(songs))
	for i, s := range songs {
		res[i] = ThemeSong{
			ID:   s.ID,
			Name: s.Text,
		}
	}
	return res
}

func (c *Client) videoToVideo(videos []mal.Video) []Video {
	res := make([]Video, len(videos))
	for i, v := range videos {
		res[i] = Video{
			ID:        v.ID,
			Title:     v.Title,
			URL:       v.URL,
			Thumbnail: v.Thumbnail,
		}
	}
	return res
}

func (c *Client) mangaToManga(manga *mal.Manga) *Manga {
	if manga == nil {
		return nil
	}

	return &Manga{
		ID:    manga.ID,
		Title: manga.Title,
		MainPicture: Picture{
			Large:  manga.MainPicture.Large,
			Medium: manga.MainPicture.Medium,
		},
		AlternativeTitles: AlternativeTitles{
			Synonyms: manga.AlternativeTitles.Synonyms,
			English:  manga.AlternativeTitles.En,
			Japanese: manga.AlternativeTitles.Ja,
		},
		StartDate:       c.dateToDate(manga.StartDate),
		EndDate:         c.dateToDate(manga.EndDate),
		Synopsis:        manga.Synopsis,
		Mean:            manga.Mean,
		Rank:            manga.Rank,
		Popularity:      manga.Popularity,
		NumListUsers:    manga.NumListUsers,
		NumScoringUsers: manga.NumScoringUsers,
		NSFW:            c.nsfwToNsfw(manga.NSFW),
		Genres:          c.genreToGenre(manga.Genres),
		MediaType:       c.mediaToMedia(manga.MediaType),
		Status:          c.statusToStatus(manga.Status),
		MyListStatus: UserMangaListStatus{
			Status:          c.listStatusToUserMangaStatus(manga.MyListStatus.Status),
			Score:           manga.MyListStatus.Score,
			NumVolumesRead:  manga.MyListStatus.NumVolumesRead,
			NumChaptersRead: manga.MyListStatus.NumChaptersRead,
			IsRereading:     manga.MyListStatus.IsRereading,
			StartDate:       c.dateToDate(manga.MyListStatus.StartDate),
			FinishDate:      c.dateToDate(manga.MyListStatus.FinishDate),
			Priority:        c.priorityToPriority(manga.MyListStatus.Priority),
			NumTimesReread:  manga.MyListStatus.NumTimesReread,
			RereadValue:     c.rereadToReread(manga.MyListStatus.RereadValue),
			Tags:            manga.MyListStatus.Tags,
			Comments:        manga.MyListStatus.Comments,
			UpdatedAt:       manga.MyListStatus.UpdatedAt,
		},
		NumVolumes:      manga.NumVolumes,
		NumChapters:     manga.NumChapters,
		Authors:         c.authorToAuthor(manga.Authors),
		Pictures:        c.pictureToPicture(manga.Pictures),
		Background:      manga.Background,
		RelatedAnime:    c.relatedAnimeToRelatedAnime(manga.RelatedAnime),
		RelatedManga:    c.relatedMangaToRelatedManga(manga.RelatedManga),
		Recommendations: c.mangaRecommendationToMangaRecommendation(manga.Recommendations),
		Serialization:   c.serializationToSerialization(manga.Serialization),
		NumFavorites:    manga.NumFavorites,
	}
}

func (c *Client) authorToAuthor(authors []mal.Author) []Author {
	res := make([]Author, len(authors))
	for i, a := range authors {
		res[i] = Author{
			Person: Person{
				ID:        a.Node.ID,
				FirstName: a.Node.FirstName,
				LastName:  a.Node.LastName,
			},
			Role: a.Role,
		}
	}
	return res
}

func (c *Client) serializationToSerialization(serialization []mal.Serialization) []Serialization {
	res := make([]Serialization, len(serialization))
	for i, s := range serialization {
		res[i] = Serialization{
			Magazine: Magazine{
				ID:   s.Node.ID,
				Name: s.Node.Name,
			},
			Role: s.Role,
		}
	}
	return res
}

func (c *Client) mangaPagingToMangaList(manga *mal.MangaPaging) []Manga {
	if manga == nil {
		return nil
	}

	res := make([]Manga, len(manga.Data))
	for i, m := range manga.Data {
		res[i] = *c.mangaToManga(&m.Node)
	}

	return res
}

func (c *Client) mangaRankingPagingToMangaList(manga *mal.MangaRankingPaging) []Manga {
	if manga == nil {
		return nil
	}

	res := make([]Manga, len(manga.Data))
	for i, m := range manga.Data {
		res[i] = *c.mangaToManga(&m.Node)
	}

	return res
}

func (c *Client) listStatusToUserAnimeListStatus(a mal.MyAnimeListStatus) *UserAnimeListStatus {
	return &UserAnimeListStatus{
		Status:             c.listStatusToUserAnimeStatus(a.Status),
		Score:              a.Score,
		NumEpisodesWatched: a.NumEpisodesWatched,
		IsRewatching:       a.IsRewatching,
		StartDate:          c.dateToDate(a.StartDate),
		FinishDate:         c.dateToDate(a.FinishDate),
		Priority:           c.priorityToPriority(a.Priority),
		NumTimesRewatched:  a.NumTimesRewatched,
		RewatchValue:       c.rewatchToRewatch(a.RewatchValue),
		Tags:               a.Tags,
		Comments:           a.Comments,
		UpdatedAt:          a.UpdatedAt,
	}
}

func (c *Client) userAnimePagingToUserAnimeList(anime *mal.UserAnimePaging) []UserAnime {
	if anime == nil {
		return nil
	}

	res := make([]UserAnime, len(anime.Data))
	for i, a := range anime.Data {
		res[i] = UserAnime{
			Anime:  *c.animeToAnime(&a.Node),
			Status: *c.listStatusToUserAnimeListStatus(a.ListStatus),
		}
	}

	return res
}

func (c *Client) listStatusToUserMangaListStatus(m mal.MyMangaListStatus) *UserMangaListStatus {
	return &UserMangaListStatus{
		Status:          c.listStatusToUserMangaStatus(m.Status),
		Score:           m.Score,
		NumVolumesRead:  m.NumVolumesRead,
		NumChaptersRead: m.NumChaptersRead,
		IsRereading:     m.IsRereading,
		StartDate:       c.dateToDate(m.StartDate),
		FinishDate:      c.dateToDate(m.FinishDate),
		Priority:        c.priorityToPriority(m.Priority),
		NumTimesReread:  m.NumTimesReread,
		RereadValue:     c.rereadToReread(m.RereadValue),
		Tags:            m.Tags,
		Comments:        m.Comments,
		UpdatedAt:       m.UpdatedAt,
	}
}

func (c *Client) userMangaPagingToUserMangaList(manga *mal.UserMangaPaging) []UserManga {
	if manga == nil {
		return nil
	}

	res := make([]UserManga, len(manga.Data))
	for i, m := range manga.Data {
		res[i] = UserManga{
			Manga:  *c.mangaToManga(&m.Node),
			Status: *c.listStatusToUserMangaListStatus(m.ListStatus),
		}
	}

	return res
}

func (c *Client) userToUser(user *mal.User) *User {
	return &User{
		ID:       user.ID,
		Name:     user.Name,
		Picture:  user.Picture,
		Gender:   user.Gender,
		Birthday: user.Birthday,
		Location: user.Location,
		JoinedAt: user.JoinedAt,
		AnimeStatistics: UserAnimeStatistic{
			WatchingCount:    user.AnimeStatistics.NumItemsWatching,
			CompletedCount:   user.AnimeStatistics.NumItemsCompleted,
			OnHoldCount:      user.AnimeStatistics.NumItemsOnHold,
			DroppedCount:     user.AnimeStatistics.NumItemsDropped,
			PlanToWatchCount: user.AnimeStatistics.NumItemsPlanToWatch,
			TotalCount:       user.AnimeStatistics.NumItems,
			WatchedDays:      user.AnimeStatistics.NumdaysWatched,
			WatchingDays:     user.AnimeStatistics.NumdaysWatching,
			CompletedDays:    user.AnimeStatistics.NumdaysCompleted,
			OnHoldDays:       user.AnimeStatistics.NumdaysOnHold,
			DroppedDays:      user.AnimeStatistics.NumdaysDropped,
			TotalDays:        user.AnimeStatistics.NumDays,
			Episode:          user.AnimeStatistics.NumEpisodes,
			RewatchedTimes:   user.AnimeStatistics.NumTimesRewatched,
			MeanScore:        user.AnimeStatistics.MeanScore,
		},
		TimeZone:    user.TimeZone,
		IsSupporter: user.IsSupporter,
	}
}

func (c *Client) forumCategoriesToForumCategories(data []mal.ForumBoardCategory) []ForumBoardCategory {
	res := make([]ForumBoardCategory, len(data))
	for i, cc := range data {
		res[i] = ForumBoardCategory{
			Title:  cc.Title,
			Boards: c.forumBoardsToForumBoards(cc.Boards),
		}
	}
	return res
}

func (c *Client) forumBoardsToForumBoards(data []mal.ForumBoard) []ForumBoard {
	res := make([]ForumBoard, len(data))
	for i, d := range data {
		res[i] = ForumBoard{
			ID:          d.ID,
			Title:       d.Title,
			Description: d.Description,
			Subboards:   c.forumSubboardsToForumSubboards(d.Subboards),
		}
	}
	return res
}

func (c *Client) forumSubboardsToForumSubboards(data []mal.ForumSubboard) []ForumSubboard {
	res := make([]ForumSubboard, len(data))
	for i, d := range data {
		res[i] = ForumSubboard{
			ID:    d.ID,
			Title: d.Title,
		}
	}
	return res
}

func (c *Client) forumTopicsToForumTopic(data []mal.ForumTopic) []ForumTopic {
	res := make([]ForumTopic, len(data))
	for i, d := range data {
		res[i] = ForumTopic{
			ID:                d.ID,
			Title:             d.Title,
			CreatedAt:         d.CreatedAt,
			CreatedBy:         d.CreatedBy.Name,
			PostCount:         d.NumberOfPosts,
			LastPostCreatedAt: d.LastPostCreatedAt,
			LastPostCreatedBy: d.LastPostCreatedBy.Name,
			IsLocked:          d.IsLocked,
		}
	}
	return res
}

func (c *Client) forumTopicDetailsToForumTopicDetails(d mal.ForumTopicData) *ForumTopicDetail {
	return &ForumTopicDetail{
		Title: d.Title,
		Posts: c.forumPostsToForumPosts(d.Posts),
		Poll: ForumPoll{
			ID:       d.Poll.ID,
			Question: d.Poll.Question,
			IsClosed: d.Poll.Closed,
			Options:  c.forumPollOptionsToForumPollOptions(d.Poll.Options),
		},
	}
}

func (c *Client) forumPostsToForumPosts(data []mal.ForumPost) []ForumPost {
	res := make([]ForumPost, len(data))
	for i, d := range data {
		res[i] = ForumPost{
			ID:        d.ID,
			Number:    d.Number,
			CreatedAt: d.CreatedAt,
			CreatedBy: d.CreatedBy.Name,
			Body:      d.Body,
			Signature: d.Signature,
		}
	}
	return res
}

func (c *Client) forumPollOptionsToForumPollOptions(data []mal.ForumPollOption) []ForumPollOption {
	res := make([]ForumPollOption, len(data))
	for i, d := range data {
		res[i] = ForumPollOption{
			ID:    d.ID,
			Text:  d.Text,
			Votes: d.Votes,
		}
	}
	return res
}
