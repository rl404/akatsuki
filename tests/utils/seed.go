package utils_test

import (
	"github.com/rl404/akatsuki/internal/domain/anime/entity"
	animeSQL "github.com/rl404/akatsuki/internal/domain/anime/repository/sql"
	genreSQL "github.com/rl404/akatsuki/internal/domain/genre/repository/sql"
	studioSQL "github.com/rl404/akatsuki/internal/domain/studio/repository/sql"
	"gorm.io/gorm"
)

type seedTypes string

const (
	SeedComplete seedTypes = "complete"
)

var testData = map[seedTypes]struct {
	anime               []animeSQL.Anime
	animeGenres         []animeSQL.AnimeGenre
	animePictures       []animeSQL.AnimePicture
	animeRelateds       []animeSQL.AnimeRelated
	animeStudios        []animeSQL.AnimeStudio
	animeStatsHistories []animeSQL.AnimeStatsHistory
	genres              []genreSQL.Genre
	studios             []studioSQL.Studio
}{
	SeedComplete: {
		anime: []animeSQL.Anime{
			{
				ID:              1,
				Title:           "Cowboy Bebop",
				TitleSynonym:    "[]",
				TitleEnglish:    "Cowboy Bebop",
				TitleJapanese:   "カウボーイビバップ",
				Picture:         "https://api-cdn.myanimelist.net/images/anime/4/19644l.jpg",
				StartDay:        3,
				StartMonth:      4,
				StartYear:       1998,
				EndDay:          24,
				EndMonth:        4,
				EndYear:         1999,
				Synopsis:        "Crime is timeless. By the year 2071, humanity has expanded across the galaxy, filling the surface of other planets with settlements like those on Earth. These new societies are plagued by murder, drug use, and theft, and intergalactic outlaws are hunted by a growing number of tough bounty hunters.\n\nSpike Spiegel and Jet Black pursue criminals throughout space to make a humble living. Beneath his goofy and aloof demeanor, Spike is haunted by the weight of his violent past. Meanwhile, Jet manages his own troubled memories while taking care of Spike and the Bebop, their ship. The duo is joined by the beautiful con artist Faye Valentine, odd child Edward Wong Hau Pepelu Tivrusky IV, and Ein, a bioengineered Welsh Corgi.\n\nWhile developing bonds and working to catch a colorful cast of criminals, the Bebop crew's lives are disrupted by a menace from Spike's past. As a rival's maniacal plot continues to unravel, Spike must choose between life with his newfound family or revenge for his old wounds.\n\n[Written by MAL Rewrite]\n",
				Background:      "When Cowboy Bebop first aired in spring of 1998 on TV Tokyo, only episodes 2, 3, 7-15, and 18 were broadcast, it was concluded with a recap special known as Yose Atsume Blues. This was due to anime censorship having increased following the big controversies over Evangelion, as a result most of the series was pulled from the air due to violent content. Satellite channel WOWOW picked up the series in the fall of that year and aired it in its entirety uncensored. Cowboy Bebop was not a ratings hit in Japan, but sold over 19,000 DVD units in the initial release run, and 81,000 overall. Protagonist Spike Spiegel won Best Male Character, and Megumi Hayashibara won Best Voice Actor for her role as Faye Valentine in the 1999 and 2000 Anime Grand Prix, respectively.\n\nCowboy Bebop's biggest influence has been in the United States, where it premiered on Adult Swim in 2001 with many reruns since. The show's heavy Western influence struck a chord with American viewers, where it became a \"gateway drug\" to anime aimed at adult audiences.",
				NSFW:            false,
				Type:            entity.TypeTV,
				Status:          entity.StatusFinished,
				Episode:         26,
				EpisodeDuration: 1440,
				Season:          entity.SeasonSpring,
				SeasonYear:      1998,
				BroadcastDay:    entity.DaySaturday,
				BroadcastTime:   "01:00",
				Source:          entity.SourceOriginal,
				Rating:          entity.RatingR,
				Mean:            8.75,
				Rank:            40,
				Popularity:      43,
				Member:          1762458,
				Voter:           909364,
				UserWatching:    1012042,
				UserCompleted:   1012042,
				UserOnHold:      101094,
				UserDropped:     39629,
				UserPlanned:     444653,
			},
			{
				ID:      5,
				Title:   "Cowboy Bebop: Tengoku no Tobira",
				Picture: "https://api-cdn.myanimelist.net/images/anime/1439/93480l.jpg",
			},
			{
				ID:      17205,
				Title:   "Cowboy Bebop: Ein no Natsuyasumi",
				Picture: "https://api-cdn.myanimelist.net/images/anime/1565/127387l.jpg",
			},
			{
				ID:      4037,
				Title:   "Cowboy Bebop: Yose Atsume Blues",
				Picture: "https://api-cdn.myanimelist.net/images/anime/10/54341l.jpg",
			},
		},
		animeGenres: []animeSQL.AnimeGenre{
			{AnimeID: 1, GenreID: 29},
			{AnimeID: 1, GenreID: 46},
			{AnimeID: 1, GenreID: 50},
			{AnimeID: 1, GenreID: 1},
			{AnimeID: 1, GenreID: 24},
		},
		animePictures: []animeSQL.AnimePicture{
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/7/3791l.jpg"},
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/12/19609l.jpg"},
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/4/19644l.jpg"},
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/4/22882l.jpg"},
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/12/22883l.jpg"},
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/6/22885l.jpg"},
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/11/53939l.jpg"},
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/2/78836l.jpg"},
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/12/81920l.jpg"},
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/1519/91488l.jpg"},
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/1764/109632l.jpg"},
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/1952/121242l.jpg"},
			{AnimeID: 1, URL: "https://api-cdn.myanimelist.net/images/anime/1623/122582l.jpg"},
		},
		animeRelateds: []animeSQL.AnimeRelated{
			{AnimeID1: 1, AnimeID2: 5, Relation: entity.RelationSideStory},
			{AnimeID1: 1, AnimeID2: 17205, Relation: entity.RelationSideStory},
			{AnimeID1: 1, AnimeID2: 4037, Relation: entity.RelationSummary},
		},
		animeStudios: []animeSQL.AnimeStudio{
			{AnimeID: 1, StudioID: 14},
		},
		animeStatsHistories: []animeSQL.AnimeStatsHistory{
			{
				ID:            1,
				AnimeID:       1,
				Mean:          8.75,
				Rank:          40,
				Popularity:    43,
				Member:        1762458,
				Voter:         909364,
				UserWatching:  1012042,
				UserCompleted: 1012042,
				UserOnHold:    101094,
				UserDropped:   39629,
				UserPlanned:   444653,
			},
		},
		genres: []genreSQL.Genre{
			{ID: 1, Name: "Action"},
			{ID: 24, Name: "Sci-Fi"},
			{ID: 29, Name: "Space"},
			{ID: 46, Name: "Award Winning"},
			{ID: 50, Name: "Adult Cast"},
		},
		studios: []studioSQL.Studio{
			{ID: 14, Name: "Sunrise"},
		},
	},
}

func Seed(db *gorm.DB, types ...seedTypes) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	for _, t := range types {
		if len(testData[t].anime) > 0 {
			if err := tx.Create(testData[t].anime).Error; err != nil {
				return err
			}
		}

		if len(testData[t].animeGenres) > 0 {
			if err := tx.Create(testData[t].animeGenres).Error; err != nil {
				return err
			}
		}

		if len(testData[t].animePictures) > 0 {
			if err := tx.Create(testData[t].animePictures).Error; err != nil {
				return err
			}
		}

		if len(testData[t].animeRelateds) > 0 {
			if err := tx.Create(testData[t].animeRelateds).Error; err != nil {
				return err
			}
		}

		if len(testData[t].animeStudios) > 0 {
			if err := tx.Create(testData[t].animeStudios).Error; err != nil {
				return err
			}
		}

		if len(testData[t].animeStatsHistories) > 0 {
			if err := tx.Create(testData[t].animeStatsHistories).Error; err != nil {
				return err
			}
		}

		if len(testData[t].genres) > 0 {
			if err := tx.Create(testData[t].genres).Error; err != nil {
				return err
			}
		}

		if len(testData[t].studios) > 0 {
			if err := tx.Create(testData[t].studios).Error; err != nil {
				return err
			}
		}
	}

	return tx.Commit().Error
}
