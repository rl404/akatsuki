package api

import (
	"github.com/rl404/akatsuki/internal/delivery/grpc/schema"
	"github.com/rl404/akatsuki/internal/service"
)

func (api *API) animeFromService(animeSvc *service.Anime) *schema.Anime {
	var anime schema.Anime
	anime.Id = animeSvc.ID
	anime.Title = animeSvc.Title
	anime.AlternativeTitles = &schema.AlternativeTitles{
		Synonyms: animeSvc.AlternativeTitles.Synonyms,
		English:  animeSvc.AlternativeTitles.English,
		Japanese: animeSvc.AlternativeTitles.Japanese,
	}
	anime.Picture = animeSvc.Picture
	anime.StartDate = &schema.Date{
		Year:  int32(animeSvc.StartDate.Year),
		Month: int32(animeSvc.StartDate.Month),
		Day:   int32(animeSvc.StartDate.Day),
	}
	anime.EndDate = &schema.Date{
		Year:  int32(animeSvc.EndDate.Year),
		Month: int32(animeSvc.EndDate.Month),
		Day:   int32(animeSvc.EndDate.Day),
	}
	anime.Synopsis = animeSvc.Synopsis
	anime.Background = animeSvc.Background
	anime.Nsfw = animeSvc.NSFW
	anime.Type = api.typeFromEntity(animeSvc.Type)
	anime.Status = api.statusFromEntity(animeSvc.Status)
	anime.Episode = &schema.Episode{
		Count:    int32(animeSvc.Episode.Count),
		Duration: int32(animeSvc.Episode.Duration),
	}
	if animeSvc.Season.Season != "" {
		anime.Season = &schema.AnimeSeason{
			Season: api.seasonFromEntity(animeSvc.Season.Season),
			Year:   int32(animeSvc.Season.Year),
		}
	}
	if animeSvc.Broadcast.Day != "" {
		anime.Broadcast = &schema.Broadcast{
			Day:  api.dayFromEntity(animeSvc.Broadcast.Day),
			Time: animeSvc.Broadcast.Time,
		}
	}
	anime.Source = api.sourceFromEntity(animeSvc.Source)
	anime.Rating = api.ratingFromEntity(animeSvc.Rating)
	anime.Mean = animeSvc.Mean
	anime.Rank = int32(animeSvc.Rank)
	anime.Popularity = int32(animeSvc.Popularity)
	anime.Member = int32(animeSvc.Member)
	anime.Voter = int32(animeSvc.Voter)
	anime.Stats = &schema.Stats{
		Status: &schema.StatsStatus{
			Watching:  int32(animeSvc.Stats.Status.Watching),
			Completed: int32(animeSvc.Stats.Status.Completed),
			OnHold:    int32(animeSvc.Stats.Status.OnHold),
			Dropped:   int32(animeSvc.Stats.Status.Dropped),
			Planned:   int32(animeSvc.Stats.Status.Planned),
		},
	}
	anime.Pictures = animeSvc.Pictures
	anime.Genres = api.animeGenresFromService(animeSvc.Genres)
	anime.Related = api.animeRelatedFromService(animeSvc.Related)
	anime.Studios = api.animeStudiosFromService(animeSvc.Studios)
	return &anime
}

func (api *API) animeGenresFromService(animeGenres []service.AnimeGenre) []*schema.AnimeGenre {
	genres := make([]*schema.AnimeGenre, len(animeGenres))
	for i, g := range animeGenres {
		genres[i] = &schema.AnimeGenre{
			Id:   g.ID,
			Name: g.Name,
		}
	}
	return genres
}

func (api *API) animeRelatedFromService(animeRelated []service.AnimeRelated) []*schema.AnimeRelated {
	related := make([]*schema.AnimeRelated, len(animeRelated))
	for i, r := range animeRelated {
		related[i] = &schema.AnimeRelated{
			Id:       r.ID,
			Title:    r.Title,
			Picture:  r.Picture,
			Relation: api.relationFromEntity(r.Relation),
		}
	}
	return related
}

func (api *API) animeStudiosFromService(animeStudios []service.AnimeStudio) []*schema.AnimeStudio {
	studios := make([]*schema.AnimeStudio, len(animeStudios))
	for i, s := range animeStudios {
		studios[i] = &schema.AnimeStudio{
			Id:   s.ID,
			Name: s.Name,
		}
	}
	return studios
}

func (api *API) userAnimeRelationFromService(userAnimeRelation *service.UserAnimeRelation) *schema.UserAnimeRelation {
	nodes := make([]*schema.UserAnimeRelationNode, len(userAnimeRelation.Nodes))
	for i, n := range userAnimeRelation.Nodes {
		nodes[i] = &schema.UserAnimeRelationNode{
			AnimeId:          n.AnimeID,
			Title:            n.Title,
			Status:           api.statusFromEntity(n.Status),
			Score:            n.Score,
			Type:             api.typeFromEntity(n.Type),
			Source:           api.sourceFromEntity(n.Source),
			StartYear:        int32(n.StartYear),
			EpisodeCount:     int32(n.EpisodeCount),
			EpisodeDuration:  int32(n.EpisodeDuration),
			Season:           api.seasonFromEntity(n.Season),
			SeasonYear:       int32(n.SeasonYear),
			UserAnimeStatus:  api.userAnimeStatusFromEntity(n.UserAnimeStatus),
			UserAnimeScore:   int32(n.UserAnimeScore),
			UserEpisodeCount: int32(n.UserEpisodeCount),
		}
	}

	links := make([]*schema.UserAnimeRelationLink, len(userAnimeRelation.Links))
	for i, l := range userAnimeRelation.Links {
		links[i] = &schema.UserAnimeRelationLink{
			AnimeId1: l.AnimeID1,
			AnimeId2: l.AnimeID2,
			Relation: api.relationFromEntity(l.Relation),
		}
	}

	return &schema.UserAnimeRelation{
		Nodes: nodes,
		Links: links,
	}
}
