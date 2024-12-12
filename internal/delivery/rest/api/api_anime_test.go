package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/akatsuki/internal/delivery/rest/api"
	"github.com/rl404/akatsuki/internal/domain/anime/entity"
	"github.com/rl404/akatsuki/internal/service"
	"github.com/rl404/akatsuki/internal/utils"
	utils_test "github.com/rl404/akatsuki/tests/utils"
	"github.com/rl404/fairy/cache"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type animeTestSuite struct {
	suite.Suite
	db    *gorm.DB
	cache cache.Cacher
	api   *api.API
}

func TestAnime(t *testing.T) {
	suite.Run(t, new(animeTestSuite))
}

func (suite *animeTestSuite) SetupSuite() {
	cfg, err := utils_test.GetConfig()
	suite.Require().Nil(err)

	cache, err := utils_test.GetCache(cfg)
	suite.Require().Nil(err)

	ps, err := utils_test.GetPubsub(cfg)
	suite.Require().Nil(err)

	db, err := utils_test.GetDB(cfg)
	suite.Require().Nil(err)

	err = utils_test.Seed(db, utils_test.SeedComplete)
	suite.Require().Nil(err)

	service := utils_test.GetService(cfg, db, cache, ps)
	suite.cache, suite.db, suite.api = cache, db, api.New(service)
}

func (suite *animeTestSuite) TearDownSuite() {
	err := suite.cache.Close()
	suite.Require().Nil(err)

	err = utils_test.TruncateDB(suite.db)
	suite.Require().Nil(err)

	db, err := suite.db.DB()
	suite.Require().Nil(err)

	err = db.Close()
	suite.Require().Nil(err)
}

func (suite *animeTestSuite) TestGetAnime() {
	tests := []struct {
		name           string
		param          map[string]string
		expectedCode   int
		expectedReturn utils.Response
	}{
		{
			name:         "invalid-start-mean",
			param:        map[string]string{"start_mean": "random"},
			expectedCode: http.StatusBadRequest,
			expectedReturn: utils.Response{
				Status:  http.StatusBadRequest,
				Message: "invalid start_mean format",
			},
		},
		{
			name:         "invalid-end-mean",
			param:        map[string]string{"start_mean": "1.0", "end_mean": "random"},
			expectedCode: http.StatusBadRequest,
			expectedReturn: utils.Response{
				Status:  http.StatusBadRequest,
				Message: "invalid end_mean format",
			},
		},
		{
			name: "ok",
			param: map[string]string{
				"title":             "cowboy",
				"nsfw":              "false",
				"type":              "TV",
				"status":            "FINISHED",
				"season":            "SPRING",
				"season_year":       "1998",
				"start_airing_year": "1998",
				"end_airing_year":   "1999",
				"start_mean":        "8.0",
				"end_mean":          "9.0",
				"genre_id":          "1",
				"studio_id":         "14",
				"page":              "1",
				"limit":             "1",
			},
			expectedCode: http.StatusOK,
			expectedReturn: utils.Response{
				Status:  http.StatusOK,
				Message: "ok",
				Data: []service.Anime{{
					ID:    1,
					Title: "Cowboy Bebop",
					AlternativeTitles: service.AlternativeTitle{
						Synonyms: []string{},
						English:  "Cowboy Bebop",
						Japanese: "カウボーイビバップ",
					},
					Picture:    "https://api-cdn.myanimelist.net/images/anime/4/19644l.jpg",
					StartDate:  service.Date{Year: 1998, Month: 4, Day: 3},
					EndDate:    service.Date{Year: 1999, Month: 4, Day: 24},
					Synopsis:   "Crime is timeless. By the year 2071, humanity has expanded across the galaxy, filling the surface of other planets with settlements like those on Earth. These new societies are plagued by murder, drug use, and theft, and intergalactic outlaws are hunted by a growing number of tough bounty hunters.\n\nSpike Spiegel and Jet Black pursue criminals throughout space to make a humble living. Beneath his goofy and aloof demeanor, Spike is haunted by the weight of his violent past. Meanwhile, Jet manages his own troubled memories while taking care of Spike and the Bebop, their ship. The duo is joined by the beautiful con artist Faye Valentine, odd child Edward Wong Hau Pepelu Tivrusky IV, and Ein, a bioengineered Welsh Corgi.\n\nWhile developing bonds and working to catch a colorful cast of criminals, the Bebop crew's lives are disrupted by a menace from Spike's past. As a rival's maniacal plot continues to unravel, Spike must choose between life with his newfound family or revenge for his old wounds.\n\n[Written by MAL Rewrite]\n",
					Background: "When Cowboy Bebop first aired in spring of 1998 on TV Tokyo, only episodes 2, 3, 7-15, and 18 were broadcast, it was concluded with a recap special known as Yose Atsume Blues. This was due to anime censorship having increased following the big controversies over Evangelion, as a result most of the series was pulled from the air due to violent content. Satellite channel WOWOW picked up the series in the fall of that year and aired it in its entirety uncensored. Cowboy Bebop was not a ratings hit in Japan, but sold over 19,000 DVD units in the initial release run, and 81,000 overall. Protagonist Spike Spiegel won Best Male Character, and Megumi Hayashibara won Best Voice Actor for her role as Faye Valentine in the 1999 and 2000 Anime Grand Prix, respectively.\n\nCowboy Bebop's biggest influence has been in the United States, where it premiered on Adult Swim in 2001 with many reruns since. The show's heavy Western influence struck a chord with American viewers, where it became a \"gateway drug\" to anime aimed at adult audiences.",
					NSFW:       false,
					Type:       entity.TypeTV,
					Status:     entity.StatusFinished,
					Episode:    service.Episode{Count: 26, Duration: 1440},
					Season:     &service.Season{Season: entity.SeasonSpring, Year: 1998},
					Broadcast:  &service.Broadcast{Day: entity.DaySaturday, Time: "01:00"},
					Source:     entity.SourceOriginal,
					Rating:     entity.RatingR,
					Mean:       8.75,
					Rank:       40,
					Popularity: 43,
					Member:     1762458,
					Voter:      909364,
					Stats: service.Stats{
						Status: service.StatsStatus{
							Watching:  1012042,
							Completed: 1012042,
							OnHold:    101094,
							Dropped:   39629,
							Planned:   444653,
						},
					},
					Genres:  []service.AnimeGenre{},
					Related: []service.AnimeRelated{},
					Studios: []service.AnimeStudio{},
				}},
				Meta: service.Pagination{
					Page:  1,
					Total: 1,
					Limit: 1,
				},
			},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			url, err := url.Parse("/anime")
			suite.Require().Nil(err)

			query := url.Query()
			for k, v := range test.param {
				query.Add(k, v)
			}
			url.RawQuery = query.Encode()

			req, err := http.NewRequest(http.MethodGet, url.String(), nil)
			suite.Require().Nil(err)

			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/anime", suite.api.HandleGetAnime)
			r.ServeHTTP(recorder, req)

			expectedBody, err := json.Marshal(test.expectedReturn)
			suite.Require().Nil(err)

			suite.Equal(test.expectedCode, recorder.Code)
			suite.Equal(string(expectedBody), recorder.Body.String())
		})
	}
}

func (suite *animeTestSuite) TestGetAnimeByID() {
	tests := []struct {
		name           string
		param          string
		expectedCode   int
		expectedReturn utils.Response
	}{
		{
			name:         "invalid-id",
			param:        "random",
			expectedCode: http.StatusBadRequest,
			expectedReturn: utils.Response{
				Status:  http.StatusBadRequest,
				Message: "invalid anime id",
			},
		},
		{
			name:         "ok",
			param:        "1",
			expectedCode: http.StatusOK,
			expectedReturn: utils.Response{
				Status:  http.StatusOK,
				Message: "ok",
				Data: service.Anime{
					ID:    1,
					Title: "Cowboy Bebop",
					AlternativeTitles: service.AlternativeTitle{
						Synonyms: []string{},
						English:  "Cowboy Bebop",
						Japanese: "カウボーイビバップ",
					},
					Picture:    "https://api-cdn.myanimelist.net/images/anime/4/19644l.jpg",
					StartDate:  service.Date{Year: 1998, Month: 4, Day: 3},
					EndDate:    service.Date{Year: 1999, Month: 4, Day: 24},
					Synopsis:   "Crime is timeless. By the year 2071, humanity has expanded across the galaxy, filling the surface of other planets with settlements like those on Earth. These new societies are plagued by murder, drug use, and theft, and intergalactic outlaws are hunted by a growing number of tough bounty hunters.\n\nSpike Spiegel and Jet Black pursue criminals throughout space to make a humble living. Beneath his goofy and aloof demeanor, Spike is haunted by the weight of his violent past. Meanwhile, Jet manages his own troubled memories while taking care of Spike and the Bebop, their ship. The duo is joined by the beautiful con artist Faye Valentine, odd child Edward Wong Hau Pepelu Tivrusky IV, and Ein, a bioengineered Welsh Corgi.\n\nWhile developing bonds and working to catch a colorful cast of criminals, the Bebop crew's lives are disrupted by a menace from Spike's past. As a rival's maniacal plot continues to unravel, Spike must choose between life with his newfound family or revenge for his old wounds.\n\n[Written by MAL Rewrite]\n",
					Background: "When Cowboy Bebop first aired in spring of 1998 on TV Tokyo, only episodes 2, 3, 7-15, and 18 were broadcast, it was concluded with a recap special known as Yose Atsume Blues. This was due to anime censorship having increased following the big controversies over Evangelion, as a result most of the series was pulled from the air due to violent content. Satellite channel WOWOW picked up the series in the fall of that year and aired it in its entirety uncensored. Cowboy Bebop was not a ratings hit in Japan, but sold over 19,000 DVD units in the initial release run, and 81,000 overall. Protagonist Spike Spiegel won Best Male Character, and Megumi Hayashibara won Best Voice Actor for her role as Faye Valentine in the 1999 and 2000 Anime Grand Prix, respectively.\n\nCowboy Bebop's biggest influence has been in the United States, where it premiered on Adult Swim in 2001 with many reruns since. The show's heavy Western influence struck a chord with American viewers, where it became a \"gateway drug\" to anime aimed at adult audiences.",
					NSFW:       false,
					Type:       entity.TypeTV,
					Status:     entity.StatusFinished,
					Episode:    service.Episode{Count: 26, Duration: 1440},
					Season:     &service.Season{Season: entity.SeasonSpring, Year: 1998},
					Broadcast:  &service.Broadcast{Day: entity.DaySaturday, Time: "01:00"},
					Source:     entity.SourceOriginal,
					Rating:     entity.RatingR,
					Mean:       8.75,
					Rank:       40,
					Popularity: 43,
					Member:     1762458,
					Voter:      909364,
					Stats: service.Stats{
						Status: service.StatsStatus{
							Watching:  1012042,
							Completed: 1012042,
							OnHold:    101094,
							Dropped:   39629,
							Planned:   444653,
						},
					},
					Pictures: []string{
						"https://api-cdn.myanimelist.net/images/anime/7/3791l.jpg",
						"https://api-cdn.myanimelist.net/images/anime/12/19609l.jpg",
						"https://api-cdn.myanimelist.net/images/anime/4/19644l.jpg",
						"https://api-cdn.myanimelist.net/images/anime/4/22882l.jpg",
						"https://api-cdn.myanimelist.net/images/anime/12/22883l.jpg",
						"https://api-cdn.myanimelist.net/images/anime/6/22885l.jpg",
						"https://api-cdn.myanimelist.net/images/anime/11/53939l.jpg",
						"https://api-cdn.myanimelist.net/images/anime/2/78836l.jpg",
						"https://api-cdn.myanimelist.net/images/anime/12/81920l.jpg",
						"https://api-cdn.myanimelist.net/images/anime/1519/91488l.jpg",
						"https://api-cdn.myanimelist.net/images/anime/1764/109632l.jpg",
						"https://api-cdn.myanimelist.net/images/anime/1952/121242l.jpg",
						"https://api-cdn.myanimelist.net/images/anime/1623/122582l.jpg",
					},
					Genres: []service.AnimeGenre{
						{ID: 1, Name: "Action"},
						{ID: 24, Name: "Sci-Fi"},
						{ID: 29, Name: "Space"},
						{ID: 46, Name: "Award Winning"},
						{ID: 50, Name: "Adult Cast"},
					},
					Related: []service.AnimeRelated{
						{
							ID:       5,
							Title:    "Cowboy Bebop: Tengoku no Tobira",
							Picture:  "https://api-cdn.myanimelist.net/images/anime/1439/93480l.jpg",
							Relation: entity.RelationSideStory,
						},
						{
							ID:       17205,
							Title:    "Cowboy Bebop: Ein no Natsuyasumi",
							Picture:  "https://api-cdn.myanimelist.net/images/anime/1565/127387l.jpg",
							Relation: entity.RelationSideStory,
						},
						{
							ID:       4037,
							Title:    "Cowboy Bebop: Yose Atsume Blues",
							Picture:  "https://api-cdn.myanimelist.net/images/anime/10/54341l.jpg",
							Relation: entity.RelationSummary,
						},
					},
					Studios: []service.AnimeStudio{
						{ID: 14, Name: "Sunrise"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			req, err := http.NewRequest(http.MethodGet, "/anime/"+test.param, nil)
			suite.Require().Nil(err)

			recorder := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/anime/{animeID}", suite.api.HandleGetAnimeByID)
			r.ServeHTTP(recorder, req)

			expectedBody, err := json.Marshal(test.expectedReturn)
			suite.Require().Nil(err)

			suite.Equal(test.expectedCode, recorder.Code)
			suite.Equal(string(expectedBody), recorder.Body.String())
		})
	}
}
