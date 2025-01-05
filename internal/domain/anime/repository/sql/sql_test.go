package sql_test

import (
	"context"
	"database/sql/driver"
	_errors "errors"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rl404/akatsuki/internal/domain/anime/entity"
	"github.com/rl404/akatsuki/internal/domain/anime/repository/sql"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type testSuite struct {
	suite.Suite
	db     *gorm.DB
	dbMock sqlmock.Sqlmock
}

func TestSQL(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) SetupSuite() {
	db, mock, err := sqlmock.New()
	suite.Require().Nil(err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	suite.Require().Nil(err)

	suite.db, suite.dbMock = gormDB, mock
}

func (suite *testSuite) TearDownSuite() {
	db, err := suite.db.DB()
	require.Nil(suite.T(), err)
	db.Close()
}

func (suite *testSuite) TestGet() {
	ctx := context.Background()
	errDummy := _errors.New("dummy error")
	boolTrue := true

	tests := []struct {
		name             string
		param            entity.GetRequest
		query            string
		queryArgs        []driver.Value
		queryReturn      []*sqlmock.Rows
		queryError       error
		queryCountCalled bool
		queryCount       string
		queryCountArgs   []driver.Value
		queryCountReturn []*sqlmock.Rows
		queryCountError  error
		expectedData     []*entity.Anime
		expectedTotal    int
		expectedCode     int
		expectedError    error
	}{
		{
			name:          "err-select",
			param:         entity.GetRequest{},
			query:         `SELECT * FROM "anime" WHERE "anime"."deleted_at" IS NULL ORDER BY rank = 0 nulls last, rank asc LIMIT $1`,
			queryArgs:     []driver.Value{0},
			queryReturn:   []*sqlmock.Rows{},
			queryError:    errDummy,
			expectedData:  nil,
			expectedTotal: 0,
			expectedCode:  http.StatusInternalServerError,
			expectedError: errors.ErrInternalDB,
		},
		{
			name:             "err-count",
			param:            entity.GetRequest{Limit: 1},
			query:            `SELECT * FROM "anime" WHERE "anime"."deleted_at" IS NULL ORDER BY rank = 0 nulls last, rank asc LIMIT $1`,
			queryArgs:        []driver.Value{1},
			queryReturn:      []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)},
			queryError:       nil,
			queryCountCalled: true,
			queryCount:       `SELECT count(*) FROM "anime" WHERE "anime"."deleted_at" IS NULL`,
			queryCountArgs:   []driver.Value{},
			queryCountReturn: nil,
			queryCountError:  errDummy,
			expectedData:     nil,
			expectedTotal:    0,
			expectedCode:     http.StatusInternalServerError,
			expectedError:    errors.ErrInternalDB,
		},
		{
			name: "ok",
			param: entity.GetRequest{
				Title:           "title",
				NSFW:            &boolTrue,
				Type:            entity.TypeMovie,
				Status:          entity.StatusFinished,
				Season:          entity.SeasonFall,
				SeasonYear:      2024,
				StartMean:       1,
				EndMean:         2,
				StartAiringYear: 2000,
				EndAiringYear:   2001,
				GenreID:         3,
				StudioID:        4,
				Limit:           1,
			},
			query: `SELECT "anime"."id","anime"."title","anime"."title_synonym","anime"."title_english","anime"."title_japanese","anime"."picture","anime"."start_day","anime"."start_month","anime"."start_year","anime"."end_day","anime"."end_month","anime"."end_year","anime"."synopsis","anime"."nsfw","anime"."type","anime"."status","anime"."episode","anime"."episode_duration","anime"."season","anime"."season_year","anime"."broadcast_day","anime"."broadcast_time","anime"."source","anime"."rating","anime"."background","anime"."mean","anime"."rank","anime"."popularity","anime"."member","anime"."voter","anime"."user_watching","anime"."user_completed","anime"."user_on_hold","anime"."user_dropped","anime"."user_planned","anime"."created_at","anime"."updated_at","anime"."deleted_at" FROM "anime" join (SELECT "anime_id" FROM "anime_genre" WHERE genre_id = $1) ag on ag.anime_id = id join (SELECT "anime_id" FROM "anime_studio" WHERE studio_id = $2) ast on ast.anime_id = id WHERE (title ilike $3 or title_synonym ilike $4 or title_english ilike $5 or title_japanese ilike $6) AND nsfw = $7 AND type = $8 AND status = $9 AND season = $10 AND season_year = $11 AND mean >= $12 AND mean <= $13 AND start_year >= $14 AND start_year <= $15 AND "anime"."deleted_at" IS NULL ORDER BY rank = 0 nulls last, rank asc LIMIT $16`,
			queryArgs: []driver.Value{
				3,
				4,
				"%title%",
				"%title%",
				"%title%",
				"%title%",
				&boolTrue,
				entity.TypeMovie,
				entity.StatusFinished,
				entity.SeasonFall,
				2024,
				1.0,
				2.0,
				2000,
				2001,
				1,
			},
			queryReturn:      []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)},
			queryError:       nil,
			queryCountCalled: true,
			queryCount:       `SELECT count(*) FROM "anime" join (SELECT "anime_id" FROM "anime_genre" WHERE genre_id = $1) ag on ag.anime_id = id join (SELECT "anime_id" FROM "anime_studio" WHERE studio_id = $2) ast on ast.anime_id = id WHERE (title ilike $3 or title_synonym ilike $4 or title_english ilike $5 or title_japanese ilike $6) AND nsfw = $7 AND type = $8 AND status = $9 AND season = $10 AND season_year = $11 AND mean >= $12 AND mean <= $13 AND start_year >= $14 AND start_year <= $15 AND "anime"."deleted_at" IS NULL`,
			queryCountArgs: []driver.Value{
				3,
				4,
				"%title%",
				"%title%",
				"%title%",
				"%title%",
				&boolTrue,
				entity.TypeMovie,
				entity.StatusFinished,
				entity.SeasonFall,
				2024,
				1.0,
				2.0,
				2000,
				2001,
			},
			queryCountReturn: []*sqlmock.Rows{sqlmock.NewRows([]string{"count"}).AddRow(1)},
			queryCountError:  nil,
			expectedData:     []*entity.Anime{{ID: 1}},
			expectedTotal:    1,
			expectedCode:     http.StatusOK,
			expectedError:    nil,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.dbMock.ExpectQuery(regexp.QuoteMeta(test.query)).
				WithArgs(test.queryArgs...).
				WillReturnRows(test.queryReturn...).
				WillReturnError(test.queryError)

			if test.queryCountCalled {
				suite.dbMock.ExpectQuery(regexp.QuoteMeta(test.queryCount)).
					WithArgs(test.queryCountArgs...).
					WillReturnRows(test.queryCountReturn...).
					WillReturnError(test.queryCountError)
			}

			sql := sql.New(suite.db, 0, 0, 0)

			data, total, code, err := sql.Get(ctx, test.param)
			suite.Equal(test.expectedData, data)
			suite.Equal(test.expectedTotal, total)
			suite.Equal(test.expectedCode, code)
			suite.ErrorIs(test.expectedError, err)
			suite.Nil(suite.dbMock.ExpectationsWereMet())
		})
	}
}

func (suite *testSuite) TestGetByID() {
	ctx := context.Background()
	errDummy := _errors.New("dummy error")

	tests := []struct {
		name                    string
		param                   int64
		queryAnime              string
		queryAnimeArgs          []driver.Value
		queryAnimeReturn        []*sqlmock.Rows
		queryAnimeError         error
		queryAnimeGenreCalled   bool
		queryAnimeGenre         string
		queryAnimeGenreArgs     []driver.Value
		queryAnimeGenreReturn   []*sqlmock.Rows
		queryAnimeGenreError    error
		queryAnimePictureCalled bool
		queryAnimePicture       string
		queryAnimePictureArgs   []driver.Value
		queryAnimePictureReturn []*sqlmock.Rows
		queryAnimePictureError  error
		queryAnimeRelatedCalled bool
		queryAnimeRelated       string
		queryAnimeRelatedArgs   []driver.Value
		queryAnimeRelatedReturn []*sqlmock.Rows
		queryAnimeRelatedError  error
		queryAnimeStudioCalled  bool
		queryAnimeStudio        string
		queryAnimeStudioArgs    []driver.Value
		queryAnimeStudioReturn  []*sqlmock.Rows
		queryAnimeStudioError   error
		expectedData            *entity.Anime
		expectedCode            int
		expectedError           error
	}{
		{
			name:             "error-anime-not-found",
			param:            1,
			queryAnime:       `SELECT * FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			queryAnimeArgs:   []driver.Value{1, 1},
			queryAnimeReturn: nil,
			queryAnimeError:  gorm.ErrRecordNotFound,
			expectedData:     nil,
			expectedCode:     http.StatusNotFound,
			expectedError:    errors.ErrAnimeNotFound,
		},
		{
			name:             "error-anime",
			param:            1,
			queryAnime:       `SELECT * FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			queryAnimeArgs:   []driver.Value{1, 1},
			queryAnimeReturn: nil,
			queryAnimeError:  errDummy,
			expectedData:     nil,
			expectedCode:     http.StatusInternalServerError,
			expectedError:    errors.ErrInternalDB,
		},
		{
			name:                  "error-anime-genre",
			param:                 1,
			queryAnime:            `SELECT * FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			queryAnimeArgs:        []driver.Value{1, 1},
			queryAnimeReturn:      []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)},
			queryAnimeError:       nil,
			queryAnimeGenreCalled: true,
			queryAnimeGenre:       `SELECT * FROM "anime_genre" WHERE anime_id = $1`,
			queryAnimeGenreArgs:   []driver.Value{1},
			queryAnimeGenreReturn: nil,
			queryAnimeGenreError:  errDummy,
			expectedData:          nil,
			expectedCode:          http.StatusInternalServerError,
			expectedError:         errors.ErrInternalDB,
		},
		{
			name:                    "error-anime-picture",
			param:                   1,
			queryAnime:              `SELECT * FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			queryAnimeArgs:          []driver.Value{1, 1},
			queryAnimeReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)},
			queryAnimeError:         nil,
			queryAnimeGenreCalled:   true,
			queryAnimeGenre:         `SELECT * FROM "anime_genre" WHERE anime_id = $1`,
			queryAnimeGenreArgs:     []driver.Value{1},
			queryAnimeGenreReturn:   []*sqlmock.Rows{sqlmock.NewRows([]string{"genre_id"}).AddRow(2)},
			queryAnimeGenreError:    nil,
			queryAnimePictureCalled: true,
			queryAnimePicture:       `SELECT * FROM "anime_picture" WHERE anime_id = $1`,
			queryAnimePictureArgs:   []driver.Value{1},
			queryAnimePictureReturn: nil,
			queryAnimePictureError:  errDummy,
			expectedData:            nil,
			expectedCode:            http.StatusInternalServerError,
			expectedError:           errors.ErrInternalDB,
		},
		{
			name:                    "error-anime-related",
			param:                   1,
			queryAnime:              `SELECT * FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			queryAnimeArgs:          []driver.Value{1, 1},
			queryAnimeReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)},
			queryAnimeError:         nil,
			queryAnimeGenreCalled:   true,
			queryAnimeGenre:         `SELECT * FROM "anime_genre" WHERE anime_id = $1`,
			queryAnimeGenreArgs:     []driver.Value{1},
			queryAnimeGenreReturn:   []*sqlmock.Rows{sqlmock.NewRows([]string{"genre_id"}).AddRow(2)},
			queryAnimeGenreError:    nil,
			queryAnimePictureCalled: true,
			queryAnimePicture:       `SELECT * FROM "anime_picture" WHERE anime_id = $1`,
			queryAnimePictureArgs:   []driver.Value{1},
			queryAnimePictureReturn: []*sqlmock.Rows{sqlmock.NewRows([]string{"url"}).AddRow("url")},
			queryAnimePictureError:  nil,
			queryAnimeRelatedCalled: true,
			queryAnimeRelated:       `SELECT * FROM "anime_related" WHERE anime_id1 = $1`,
			queryAnimeRelatedArgs:   []driver.Value{1},
			queryAnimeRelatedReturn: nil,
			queryAnimeRelatedError:  errDummy,
			expectedData:            nil,
			expectedCode:            http.StatusInternalServerError,
			expectedError:           errors.ErrInternalDB,
		},
		{
			name:                    "error-anime-studio",
			param:                   1,
			queryAnime:              `SELECT * FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			queryAnimeArgs:          []driver.Value{1, 1},
			queryAnimeReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)},
			queryAnimeError:         nil,
			queryAnimeGenreCalled:   true,
			queryAnimeGenre:         `SELECT * FROM "anime_genre" WHERE anime_id = $1`,
			queryAnimeGenreArgs:     []driver.Value{1},
			queryAnimeGenreReturn:   []*sqlmock.Rows{sqlmock.NewRows([]string{"genre_id"}).AddRow(2)},
			queryAnimeGenreError:    nil,
			queryAnimePictureCalled: true,
			queryAnimePicture:       `SELECT * FROM "anime_picture" WHERE anime_id = $1`,
			queryAnimePictureArgs:   []driver.Value{1},
			queryAnimePictureReturn: []*sqlmock.Rows{sqlmock.NewRows([]string{"url"}).AddRow("url")},
			queryAnimePictureError:  nil,
			queryAnimeRelatedCalled: true,
			queryAnimeRelated:       `SELECT * FROM "anime_related" WHERE anime_id1 = $1`,
			queryAnimeRelatedArgs:   []driver.Value{1},
			queryAnimeRelatedReturn: []*sqlmock.Rows{sqlmock.NewRows([]string{"anime_id2", "relation"}).AddRow(3, "SEQUEL")},
			queryAnimeRelatedError:  nil,
			queryAnimeStudioCalled:  true,
			queryAnimeStudio:        `SELECT * FROM "anime_studio" WHERE anime_id = $1`,
			queryAnimeStudioArgs:    []driver.Value{1},
			queryAnimeStudioReturn:  nil,
			queryAnimeStudioError:   errDummy,
			expectedData:            nil,
			expectedCode:            http.StatusInternalServerError,
			expectedError:           errors.ErrInternalDB,
		},
		{
			name:                    "ok",
			param:                   1,
			queryAnime:              `SELECT * FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			queryAnimeArgs:          []driver.Value{1, 1},
			queryAnimeReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)},
			queryAnimeError:         nil,
			queryAnimeGenreCalled:   true,
			queryAnimeGenre:         `SELECT * FROM "anime_genre" WHERE anime_id = $1`,
			queryAnimeGenreArgs:     []driver.Value{1},
			queryAnimeGenreReturn:   []*sqlmock.Rows{sqlmock.NewRows([]string{"genre_id"}).AddRow(2)},
			queryAnimeGenreError:    nil,
			queryAnimePictureCalled: true,
			queryAnimePicture:       `SELECT * FROM "anime_picture" WHERE anime_id = $1`,
			queryAnimePictureArgs:   []driver.Value{1},
			queryAnimePictureReturn: []*sqlmock.Rows{sqlmock.NewRows([]string{"url"}).AddRow("url")},
			queryAnimePictureError:  nil,
			queryAnimeRelatedCalled: true,
			queryAnimeRelated:       `SELECT * FROM "anime_related" WHERE anime_id1 = $1`,
			queryAnimeRelatedArgs:   []driver.Value{1},
			queryAnimeRelatedReturn: []*sqlmock.Rows{sqlmock.NewRows([]string{"anime_id2", "relation"}).AddRow(3, "SEQUEL")},
			queryAnimeRelatedError:  nil,
			queryAnimeStudioCalled:  true,
			queryAnimeStudio:        `SELECT * FROM "anime_studio" WHERE anime_id = $1`,
			queryAnimeStudioArgs:    []driver.Value{1},
			queryAnimeStudioReturn:  []*sqlmock.Rows{sqlmock.NewRows([]string{"studio_id"}).AddRow(3)},
			queryAnimeStudioError:   nil,
			expectedData: &entity.Anime{
				ID:        1,
				GenreIDs:  []int64{2},
				Pictures:  []string{"url"},
				Related:   []entity.Related{{ID: 3, Relation: entity.RelationSequel}},
				StudioIDs: []int64{3},
			},
			expectedCode:  http.StatusOK,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.dbMock.ExpectQuery(regexp.QuoteMeta(test.queryAnime)).
				WithArgs(test.queryAnimeArgs...).
				WillReturnRows(test.queryAnimeReturn...).
				WillReturnError(test.queryAnimeError)

			if test.queryAnimeGenreCalled {
				suite.dbMock.ExpectQuery(regexp.QuoteMeta(test.queryAnimeGenre)).
					WithArgs(test.queryAnimeGenreArgs...).
					WillReturnRows(test.queryAnimeGenreReturn...).
					WillReturnError(test.queryAnimeGenreError)
			}

			if test.queryAnimePictureCalled {
				suite.dbMock.ExpectQuery(regexp.QuoteMeta(test.queryAnimePicture)).
					WithArgs(test.queryAnimePictureArgs...).
					WillReturnRows(test.queryAnimePictureReturn...).
					WillReturnError(test.queryAnimePictureError)
			}

			if test.queryAnimeRelatedCalled {
				suite.dbMock.ExpectQuery(regexp.QuoteMeta(test.queryAnimeRelated)).
					WithArgs(test.queryAnimeRelatedArgs...).
					WillReturnRows(test.queryAnimeRelatedReturn...).
					WillReturnError(test.queryAnimeRelatedError)
			}

			if test.queryAnimeStudioCalled {
				suite.dbMock.ExpectQuery(regexp.QuoteMeta(test.queryAnimeStudio)).
					WithArgs(test.queryAnimeStudioArgs...).
					WillReturnRows(test.queryAnimeStudioReturn...).
					WillReturnError(test.queryAnimeStudioError)
			}

			sql := sql.New(suite.db, 0, 0, 0)

			data, code, err := sql.GetByID(ctx, test.param)
			suite.Equal(test.expectedData, data)
			suite.Equal(test.expectedCode, code)
			suite.ErrorIs(test.expectedError, err)
			suite.Nil(suite.dbMock.ExpectationsWereMet())
		})
	}
}

func (suite *testSuite) TestGetByIDs() {
	ctx := context.Background()
	errDummy := _errors.New("dummy error")

	tests := []struct {
		name          string
		param         []int64
		query         string
		queryArgs     []driver.Value
		queryReturn   []*sqlmock.Rows
		queryError    error
		expectedData  []*entity.Anime
		expectedCode  int
		expectedError error
	}{
		{
			name:          "error",
			param:         []int64{1},
			query:         `SELECT * FROM "anime" WHERE id in ($1) AND "anime"."deleted_at" IS NULL`,
			queryArgs:     []driver.Value{1},
			queryReturn:   []*sqlmock.Rows{},
			queryError:    errDummy,
			expectedData:  nil,
			expectedCode:  http.StatusInternalServerError,
			expectedError: errors.ErrInternalDB,
		},
		{
			name:          "ok",
			param:         []int64{1},
			query:         `SELECT * FROM "anime" WHERE id in ($1) AND "anime"."deleted_at" IS NULL`,
			queryArgs:     []driver.Value{1},
			queryReturn:   []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)},
			queryError:    nil,
			expectedData:  []*entity.Anime{{ID: 1}},
			expectedCode:  http.StatusOK,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.dbMock.ExpectQuery(regexp.QuoteMeta(test.query)).
				WithArgs(test.queryArgs...).
				WillReturnRows(test.queryReturn...).
				WillReturnError(test.queryError)

			sql := sql.New(suite.db, 0, 0, 0)

			data, code, err := sql.GetByIDs(ctx, test.param)
			suite.Equal(test.expectedData, data)
			suite.Equal(test.expectedCode, code)
			suite.ErrorIs(test.expectedError, err)
			suite.Nil(suite.dbMock.ExpectationsWereMet())
		})
	}
}

func (suite *testSuite) TestUpdate() {
	ctx := context.Background()
	errDummy := _errors.New("dummy error")
	now := time.Now()

	anime := entity.Anime{
		ID:    1,
		Title: "title",
		AlternativeTitle: entity.AlternativeTitle{
			Synonyms: []string{},
		},
		GenreIDs:  []int64{2},
		Pictures:  []string{"www"},
		Related:   []entity.Related{{ID: 3, Relation: entity.RelationFullStory}},
		StudioIDs: []int64{4},
	}

	tests := []struct {
		name                     string
		param                    entity.Anime
		beginError               error
		selectCalled             bool
		selectQuery              string
		selectQueryArgs          []driver.Value
		selectQueryReturn        []*sqlmock.Rows
		selectQueryError         error
		saveCalled               bool
		saveQuery                string
		saveQueryArgs            []driver.Value
		saveQueryResult          driver.Result
		saveQueryError           error
		deleteGenreCalled        bool
		deleteGenreQuery         string
		deleteGenreQueryArgs     []driver.Value
		deleteGenreQueryResult   driver.Result
		deleteGenreQueryError    error
		createGenreCalled        bool
		createGenreQuery         string
		createGenreQueryArgs     []driver.Value
		createGenreQueryResult   driver.Result
		createGenreQueryError    error
		deletePictureCalled      bool
		deletePictureQuery       string
		deletePictureQueryArgs   []driver.Value
		deletePictureQueryResult driver.Result
		deletePictureQueryError  error
		createPictureCalled      bool
		createPictureQuery       string
		createPictureQueryArgs   []driver.Value
		createPictureQueryResult driver.Result
		createPictureQueryError  error
		deleteRelatedCalled      bool
		deleteRelatedQuery       string
		deleteRelatedQueryArgs   []driver.Value
		deleteRelatedQueryResult driver.Result
		deleteRelatedQueryError  error
		createRelatedCalled      bool
		createRelatedQuery       string
		createRelatedQueryArgs   []driver.Value
		createRelatedQueryResult driver.Result
		createRelatedQueryError  error
		deleteStudioCalled       bool
		deleteStudioQuery        string
		deleteStudioQueryArgs    []driver.Value
		deleteStudioQueryResult  driver.Result
		deleteStudioQueryError   error
		createStudioCalled       bool
		createStudioQuery        string
		createStudioQueryArgs    []driver.Value
		createStudioQueryResult  driver.Result
		createStudioQueryError   error
		createHistoryCalled      bool
		createHistoryQuery       string
		createHistoryQueryArgs   []driver.Value
		createHistoryQueryReturn []*sqlmock.Rows
		createHistoryQueryError  error
		rollbackCalled           bool
		commitCalled             bool
		commitError              error
		expectedCode             int
		expectedError            error
	}{
		{
			name:          "error-begin",
			param:         anime,
			beginError:    errDummy,
			expectedCode:  http.StatusInternalServerError,
			expectedError: errors.ErrInternalDB,
		},
		{
			name:              "error-select",
			param:             anime,
			selectCalled:      true,
			selectQuery:       `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:   []driver.Value{1, 1},
			selectQueryReturn: []*sqlmock.Rows{},
			selectQueryError:  errDummy,
			rollbackCalled:    true,
			expectedCode:      http.StatusInternalServerError,
			expectedError:     errors.ErrInternalDB,
		},
		{
			name:              "error-save",
			param:             anime,
			selectCalled:      true,
			selectQuery:       `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:   []driver.Value{1, 1},
			selectQueryReturn: []*sqlmock.Rows{sqlmock.NewRows([]string{"created_at"}).AddRow(&now)},
			selectQueryError:  nil,
			saveCalled:        true,
			saveQuery:         `UPDATE "anime" SET "title"=$1,"title_synonym"=$2,"title_english"=$3,"title_japanese"=$4,"picture"=$5,"start_day"=$6,"start_month"=$7,"start_year"=$8,"end_day"=$9,"end_month"=$10,"end_year"=$11,"synopsis"=$12,"nsfw"=$13,"type"=$14,"status"=$15,"episode"=$16,"episode_duration"=$17,"season"=$18,"season_year"=$19,"broadcast_day"=$20,"broadcast_time"=$21,"source"=$22,"rating"=$23,"background"=$24,"mean"=$25,"rank"=$26,"popularity"=$27,"member"=$28,"voter"=$29,"user_watching"=$30,"user_completed"=$31,"user_on_hold"=$32,"user_dropped"=$33,"user_planned"=$34,"created_at"=$35,"updated_at"=$36,"deleted_at"=$37 WHERE "anime"."deleted_at" IS NULL AND "id" = $38`,
			saveQueryArgs:     []driver.Value{anime.Title, "[]", "", "", "", 0, 0, 0, 0, 0, 0, "", false, "", "", 0, 0, "", 0, "", "", "", "", "", 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, now, sqlmock.AnyArg(), nil, 1},
			saveQueryError:    errDummy,
			rollbackCalled:    true,
			expectedCode:      http.StatusInternalServerError,
			expectedError:     errors.ErrInternalDB,
		},
		{
			name:                  "error-delete-genre",
			param:                 anime,
			selectCalled:          true,
			selectQuery:           `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:       []driver.Value{1, 1},
			selectQueryReturn:     []*sqlmock.Rows{sqlmock.NewRows([]string{"created_at"}).AddRow(&now)},
			selectQueryError:      nil,
			saveCalled:            true,
			saveQuery:             `UPDATE "anime" SET "title"=$1,"title_synonym"=$2,"title_english"=$3,"title_japanese"=$4,"picture"=$5,"start_day"=$6,"start_month"=$7,"start_year"=$8,"end_day"=$9,"end_month"=$10,"end_year"=$11,"synopsis"=$12,"nsfw"=$13,"type"=$14,"status"=$15,"episode"=$16,"episode_duration"=$17,"season"=$18,"season_year"=$19,"broadcast_day"=$20,"broadcast_time"=$21,"source"=$22,"rating"=$23,"background"=$24,"mean"=$25,"rank"=$26,"popularity"=$27,"member"=$28,"voter"=$29,"user_watching"=$30,"user_completed"=$31,"user_on_hold"=$32,"user_dropped"=$33,"user_planned"=$34,"created_at"=$35,"updated_at"=$36,"deleted_at"=$37 WHERE "anime"."deleted_at" IS NULL AND "id" = $38`,
			saveQueryArgs:         []driver.Value{anime.Title, "[]", "", "", "", 0, 0, 0, 0, 0, 0, "", false, "", "", 0, 0, "", 0, "", "", "", "", "", 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, now, sqlmock.AnyArg(), nil, 1},
			saveQueryResult:       sqlmock.NewResult(0, 1),
			saveQueryError:        nil,
			deleteGenreCalled:     true,
			deleteGenreQuery:      `DELETE FROM "anime_genre" WHERE anime_id = $1`,
			deleteGenreQueryArgs:  []driver.Value{1},
			deleteGenreQueryError: errDummy,
			rollbackCalled:        true,
			expectedCode:          http.StatusInternalServerError,
			expectedError:         errors.ErrInternalDB,
		},
		{
			name:                   "error-create-genre",
			param:                  anime,
			selectCalled:           true,
			selectQuery:            `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:        []driver.Value{1, 1},
			selectQueryReturn:      []*sqlmock.Rows{sqlmock.NewRows([]string{"created_at"}).AddRow(&now)},
			selectQueryError:       nil,
			saveCalled:             true,
			saveQuery:              `UPDATE "anime" SET "title"=$1,"title_synonym"=$2,"title_english"=$3,"title_japanese"=$4,"picture"=$5,"start_day"=$6,"start_month"=$7,"start_year"=$8,"end_day"=$9,"end_month"=$10,"end_year"=$11,"synopsis"=$12,"nsfw"=$13,"type"=$14,"status"=$15,"episode"=$16,"episode_duration"=$17,"season"=$18,"season_year"=$19,"broadcast_day"=$20,"broadcast_time"=$21,"source"=$22,"rating"=$23,"background"=$24,"mean"=$25,"rank"=$26,"popularity"=$27,"member"=$28,"voter"=$29,"user_watching"=$30,"user_completed"=$31,"user_on_hold"=$32,"user_dropped"=$33,"user_planned"=$34,"created_at"=$35,"updated_at"=$36,"deleted_at"=$37 WHERE "anime"."deleted_at" IS NULL AND "id" = $38`,
			saveQueryArgs:          []driver.Value{anime.Title, "[]", "", "", "", 0, 0, 0, 0, 0, 0, "", false, "", "", 0, 0, "", 0, "", "", "", "", "", 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, now, sqlmock.AnyArg(), nil, 1},
			saveQueryResult:        sqlmock.NewResult(0, 1),
			saveQueryError:         nil,
			deleteGenreCalled:      true,
			deleteGenreQuery:       `DELETE FROM "anime_genre" WHERE anime_id = $1`,
			deleteGenreQueryArgs:   []driver.Value{1},
			deleteGenreQueryResult: sqlmock.NewResult(0, 1),
			deleteGenreQueryError:  nil,
			createGenreCalled:      true,
			createGenreQuery:       `INSERT INTO "anime_genre" ("anime_id","genre_id") VALUES ($1,$2)`,
			createGenreQueryArgs:   []driver.Value{1, 2},
			createGenreQueryError:  errDummy,
			rollbackCalled:         true,
			expectedCode:           http.StatusInternalServerError,
			expectedError:          errors.ErrInternalDB,
		},
		{
			name:                     "error-delete-picture",
			param:                    anime,
			selectCalled:             true,
			selectQuery:              `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:          []driver.Value{1, 1},
			selectQueryReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"created_at"}).AddRow(&now)},
			selectQueryError:         nil,
			saveCalled:               true,
			saveQuery:                `UPDATE "anime" SET "title"=$1,"title_synonym"=$2,"title_english"=$3,"title_japanese"=$4,"picture"=$5,"start_day"=$6,"start_month"=$7,"start_year"=$8,"end_day"=$9,"end_month"=$10,"end_year"=$11,"synopsis"=$12,"nsfw"=$13,"type"=$14,"status"=$15,"episode"=$16,"episode_duration"=$17,"season"=$18,"season_year"=$19,"broadcast_day"=$20,"broadcast_time"=$21,"source"=$22,"rating"=$23,"background"=$24,"mean"=$25,"rank"=$26,"popularity"=$27,"member"=$28,"voter"=$29,"user_watching"=$30,"user_completed"=$31,"user_on_hold"=$32,"user_dropped"=$33,"user_planned"=$34,"created_at"=$35,"updated_at"=$36,"deleted_at"=$37 WHERE "anime"."deleted_at" IS NULL AND "id" = $38`,
			saveQueryArgs:            []driver.Value{anime.Title, "[]", "", "", "", 0, 0, 0, 0, 0, 0, "", false, "", "", 0, 0, "", 0, "", "", "", "", "", 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, now, sqlmock.AnyArg(), nil, 1},
			saveQueryResult:          sqlmock.NewResult(0, 1),
			saveQueryError:           nil,
			deleteGenreCalled:        true,
			deleteGenreQuery:         `DELETE FROM "anime_genre" WHERE anime_id = $1`,
			deleteGenreQueryArgs:     []driver.Value{1},
			deleteGenreQueryResult:   sqlmock.NewResult(0, 1),
			deleteGenreQueryError:    nil,
			createGenreCalled:        true,
			createGenreQuery:         `INSERT INTO "anime_genre" ("anime_id","genre_id") VALUES ($1,$2)`,
			createGenreQueryArgs:     []driver.Value{1, 2},
			createGenreQueryResult:   sqlmock.NewResult(0, 1),
			createGenreQueryError:    nil,
			deletePictureCalled:      true,
			deletePictureQuery:       `DELETE FROM "anime_picture" WHERE anime_id = $1`,
			deletePictureQueryArgs:   []driver.Value{1},
			deletePictureQueryResult: sqlmock.NewResult(0, 1),
			deletePictureQueryError:  errDummy,
			rollbackCalled:           true,
			expectedCode:             http.StatusInternalServerError,
			expectedError:            errors.ErrInternalDB,
		},
		{
			name:                     "error-create-picture",
			param:                    anime,
			selectCalled:             true,
			selectQuery:              `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:          []driver.Value{1, 1},
			selectQueryReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"created_at"}).AddRow(&now)},
			selectQueryError:         nil,
			saveCalled:               true,
			saveQuery:                `UPDATE "anime" SET "title"=$1,"title_synonym"=$2,"title_english"=$3,"title_japanese"=$4,"picture"=$5,"start_day"=$6,"start_month"=$7,"start_year"=$8,"end_day"=$9,"end_month"=$10,"end_year"=$11,"synopsis"=$12,"nsfw"=$13,"type"=$14,"status"=$15,"episode"=$16,"episode_duration"=$17,"season"=$18,"season_year"=$19,"broadcast_day"=$20,"broadcast_time"=$21,"source"=$22,"rating"=$23,"background"=$24,"mean"=$25,"rank"=$26,"popularity"=$27,"member"=$28,"voter"=$29,"user_watching"=$30,"user_completed"=$31,"user_on_hold"=$32,"user_dropped"=$33,"user_planned"=$34,"created_at"=$35,"updated_at"=$36,"deleted_at"=$37 WHERE "anime"."deleted_at" IS NULL AND "id" = $38`,
			saveQueryArgs:            []driver.Value{anime.Title, "[]", "", "", "", 0, 0, 0, 0, 0, 0, "", false, "", "", 0, 0, "", 0, "", "", "", "", "", 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, now, sqlmock.AnyArg(), nil, 1},
			saveQueryResult:          sqlmock.NewResult(0, 1),
			saveQueryError:           nil,
			deleteGenreCalled:        true,
			deleteGenreQuery:         `DELETE FROM "anime_genre" WHERE anime_id = $1`,
			deleteGenreQueryArgs:     []driver.Value{1},
			deleteGenreQueryResult:   sqlmock.NewResult(0, 1),
			deleteGenreQueryError:    nil,
			createGenreCalled:        true,
			createGenreQuery:         `INSERT INTO "anime_genre" ("anime_id","genre_id") VALUES ($1,$2)`,
			createGenreQueryArgs:     []driver.Value{1, 2},
			createGenreQueryResult:   sqlmock.NewResult(0, 1),
			createGenreQueryError:    nil,
			deletePictureCalled:      true,
			deletePictureQuery:       `DELETE FROM "anime_picture" WHERE anime_id = $1`,
			deletePictureQueryArgs:   []driver.Value{1},
			deletePictureQueryResult: sqlmock.NewResult(0, 1),
			deletePictureQueryError:  nil,
			createPictureCalled:      true,
			createPictureQuery:       `INSERT INTO "anime_picture" ("anime_id","url") VALUES ($1,$2)`,
			createPictureQueryArgs:   []driver.Value{1, "www"},
			createPictureQueryResult: sqlmock.NewResult(0, 1),
			createPictureQueryError:  errDummy,
			rollbackCalled:           true,
			expectedCode:             http.StatusInternalServerError,
			expectedError:            errors.ErrInternalDB,
		},
		{
			name:                     "error-delete-related",
			param:                    anime,
			selectCalled:             true,
			selectQuery:              `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:          []driver.Value{1, 1},
			selectQueryReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"created_at"}).AddRow(&now)},
			selectQueryError:         nil,
			saveCalled:               true,
			saveQuery:                `UPDATE "anime" SET "title"=$1,"title_synonym"=$2,"title_english"=$3,"title_japanese"=$4,"picture"=$5,"start_day"=$6,"start_month"=$7,"start_year"=$8,"end_day"=$9,"end_month"=$10,"end_year"=$11,"synopsis"=$12,"nsfw"=$13,"type"=$14,"status"=$15,"episode"=$16,"episode_duration"=$17,"season"=$18,"season_year"=$19,"broadcast_day"=$20,"broadcast_time"=$21,"source"=$22,"rating"=$23,"background"=$24,"mean"=$25,"rank"=$26,"popularity"=$27,"member"=$28,"voter"=$29,"user_watching"=$30,"user_completed"=$31,"user_on_hold"=$32,"user_dropped"=$33,"user_planned"=$34,"created_at"=$35,"updated_at"=$36,"deleted_at"=$37 WHERE "anime"."deleted_at" IS NULL AND "id" = $38`,
			saveQueryArgs:            []driver.Value{anime.Title, "[]", "", "", "", 0, 0, 0, 0, 0, 0, "", false, "", "", 0, 0, "", 0, "", "", "", "", "", 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, now, sqlmock.AnyArg(), nil, 1},
			saveQueryResult:          sqlmock.NewResult(0, 1),
			saveQueryError:           nil,
			deleteGenreCalled:        true,
			deleteGenreQuery:         `DELETE FROM "anime_genre" WHERE anime_id = $1`,
			deleteGenreQueryArgs:     []driver.Value{1},
			deleteGenreQueryResult:   sqlmock.NewResult(0, 1),
			deleteGenreQueryError:    nil,
			createGenreCalled:        true,
			createGenreQuery:         `INSERT INTO "anime_genre" ("anime_id","genre_id") VALUES ($1,$2)`,
			createGenreQueryArgs:     []driver.Value{1, 2},
			createGenreQueryResult:   sqlmock.NewResult(0, 1),
			createGenreQueryError:    nil,
			deletePictureCalled:      true,
			deletePictureQuery:       `DELETE FROM "anime_picture" WHERE anime_id = $1`,
			deletePictureQueryArgs:   []driver.Value{1},
			deletePictureQueryResult: sqlmock.NewResult(0, 1),
			deletePictureQueryError:  nil,
			createPictureCalled:      true,
			createPictureQuery:       `INSERT INTO "anime_picture" ("anime_id","url") VALUES ($1,$2)`,
			createPictureQueryArgs:   []driver.Value{1, "www"},
			createPictureQueryResult: sqlmock.NewResult(0, 1),
			createPictureQueryError:  nil,
			deleteRelatedCalled:      true,
			deleteRelatedQuery:       `DELETE FROM "anime_related" WHERE anime_id1 = $1`,
			deleteRelatedQueryArgs:   []driver.Value{1},
			deleteRelatedQueryResult: sqlmock.NewResult(0, 1),
			deleteRelatedQueryError:  errDummy,
			rollbackCalled:           true,
			expectedCode:             http.StatusInternalServerError,
			expectedError:            errors.ErrInternalDB,
		},
		{
			name:                     "error-create-related",
			param:                    anime,
			selectCalled:             true,
			selectQuery:              `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:          []driver.Value{1, 1},
			selectQueryReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"created_at"}).AddRow(&now)},
			selectQueryError:         nil,
			saveCalled:               true,
			saveQuery:                `UPDATE "anime" SET "title"=$1,"title_synonym"=$2,"title_english"=$3,"title_japanese"=$4,"picture"=$5,"start_day"=$6,"start_month"=$7,"start_year"=$8,"end_day"=$9,"end_month"=$10,"end_year"=$11,"synopsis"=$12,"nsfw"=$13,"type"=$14,"status"=$15,"episode"=$16,"episode_duration"=$17,"season"=$18,"season_year"=$19,"broadcast_day"=$20,"broadcast_time"=$21,"source"=$22,"rating"=$23,"background"=$24,"mean"=$25,"rank"=$26,"popularity"=$27,"member"=$28,"voter"=$29,"user_watching"=$30,"user_completed"=$31,"user_on_hold"=$32,"user_dropped"=$33,"user_planned"=$34,"created_at"=$35,"updated_at"=$36,"deleted_at"=$37 WHERE "anime"."deleted_at" IS NULL AND "id" = $38`,
			saveQueryArgs:            []driver.Value{anime.Title, "[]", "", "", "", 0, 0, 0, 0, 0, 0, "", false, "", "", 0, 0, "", 0, "", "", "", "", "", 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, now, sqlmock.AnyArg(), nil, 1},
			saveQueryResult:          sqlmock.NewResult(0, 1),
			saveQueryError:           nil,
			deleteGenreCalled:        true,
			deleteGenreQuery:         `DELETE FROM "anime_genre" WHERE anime_id = $1`,
			deleteGenreQueryArgs:     []driver.Value{1},
			deleteGenreQueryResult:   sqlmock.NewResult(0, 1),
			deleteGenreQueryError:    nil,
			createGenreCalled:        true,
			createGenreQuery:         `INSERT INTO "anime_genre" ("anime_id","genre_id") VALUES ($1,$2)`,
			createGenreQueryArgs:     []driver.Value{1, 2},
			createGenreQueryResult:   sqlmock.NewResult(0, 1),
			createGenreQueryError:    nil,
			deletePictureCalled:      true,
			deletePictureQuery:       `DELETE FROM "anime_picture" WHERE anime_id = $1`,
			deletePictureQueryArgs:   []driver.Value{1},
			deletePictureQueryResult: sqlmock.NewResult(0, 1),
			deletePictureQueryError:  nil,
			createPictureCalled:      true,
			createPictureQuery:       `INSERT INTO "anime_picture" ("anime_id","url") VALUES ($1,$2)`,
			createPictureQueryArgs:   []driver.Value{1, "www"},
			createPictureQueryResult: sqlmock.NewResult(0, 1),
			createPictureQueryError:  nil,
			deleteRelatedCalled:      true,
			deleteRelatedQuery:       `DELETE FROM "anime_related" WHERE anime_id1 = $1`,
			deleteRelatedQueryArgs:   []driver.Value{1},
			deleteRelatedQueryResult: sqlmock.NewResult(0, 1),
			deleteRelatedQueryError:  nil,
			createRelatedCalled:      true,
			createRelatedQuery:       `INSERT INTO "anime_related" ("anime_id1","anime_id2","relation") VALUES ($1,$2,$3)`,
			createRelatedQueryArgs:   []driver.Value{1, 3, "FULL_STORY"},
			createRelatedQueryResult: sqlmock.NewResult(0, 1),
			createRelatedQueryError:  errDummy,
			rollbackCalled:           true,
			expectedCode:             http.StatusInternalServerError,
			expectedError:            errors.ErrInternalDB,
		},
		{
			name:                     "error-delete-studio",
			param:                    anime,
			selectCalled:             true,
			selectQuery:              `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:          []driver.Value{1, 1},
			selectQueryReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"created_at"}).AddRow(&now)},
			selectQueryError:         nil,
			saveCalled:               true,
			saveQuery:                `UPDATE "anime" SET "title"=$1,"title_synonym"=$2,"title_english"=$3,"title_japanese"=$4,"picture"=$5,"start_day"=$6,"start_month"=$7,"start_year"=$8,"end_day"=$9,"end_month"=$10,"end_year"=$11,"synopsis"=$12,"nsfw"=$13,"type"=$14,"status"=$15,"episode"=$16,"episode_duration"=$17,"season"=$18,"season_year"=$19,"broadcast_day"=$20,"broadcast_time"=$21,"source"=$22,"rating"=$23,"background"=$24,"mean"=$25,"rank"=$26,"popularity"=$27,"member"=$28,"voter"=$29,"user_watching"=$30,"user_completed"=$31,"user_on_hold"=$32,"user_dropped"=$33,"user_planned"=$34,"created_at"=$35,"updated_at"=$36,"deleted_at"=$37 WHERE "anime"."deleted_at" IS NULL AND "id" = $38`,
			saveQueryArgs:            []driver.Value{anime.Title, "[]", "", "", "", 0, 0, 0, 0, 0, 0, "", false, "", "", 0, 0, "", 0, "", "", "", "", "", 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, now, sqlmock.AnyArg(), nil, 1},
			saveQueryResult:          sqlmock.NewResult(0, 1),
			saveQueryError:           nil,
			deleteGenreCalled:        true,
			deleteGenreQuery:         `DELETE FROM "anime_genre" WHERE anime_id = $1`,
			deleteGenreQueryArgs:     []driver.Value{1},
			deleteGenreQueryResult:   sqlmock.NewResult(0, 1),
			deleteGenreQueryError:    nil,
			createGenreCalled:        true,
			createGenreQuery:         `INSERT INTO "anime_genre" ("anime_id","genre_id") VALUES ($1,$2)`,
			createGenreQueryArgs:     []driver.Value{1, 2},
			createGenreQueryResult:   sqlmock.NewResult(0, 1),
			createGenreQueryError:    nil,
			deletePictureCalled:      true,
			deletePictureQuery:       `DELETE FROM "anime_picture" WHERE anime_id = $1`,
			deletePictureQueryArgs:   []driver.Value{1},
			deletePictureQueryResult: sqlmock.NewResult(0, 1),
			deletePictureQueryError:  nil,
			createPictureCalled:      true,
			createPictureQuery:       `INSERT INTO "anime_picture" ("anime_id","url") VALUES ($1,$2)`,
			createPictureQueryArgs:   []driver.Value{1, "www"},
			createPictureQueryResult: sqlmock.NewResult(0, 1),
			createPictureQueryError:  nil,
			deleteRelatedCalled:      true,
			deleteRelatedQuery:       `DELETE FROM "anime_related" WHERE anime_id1 = $1`,
			deleteRelatedQueryArgs:   []driver.Value{1},
			deleteRelatedQueryResult: sqlmock.NewResult(0, 1),
			deleteRelatedQueryError:  nil,
			createRelatedCalled:      true,
			createRelatedQuery:       `INSERT INTO "anime_related" ("anime_id1","anime_id2","relation") VALUES ($1,$2,$3)`,
			createRelatedQueryArgs:   []driver.Value{1, 3, "FULL_STORY"},
			createRelatedQueryResult: sqlmock.NewResult(0, 1),
			createRelatedQueryError:  nil,
			deleteStudioCalled:       true,
			deleteStudioQuery:        `DELETE FROM "anime_studio" WHERE anime_id = $1`,
			deleteStudioQueryArgs:    []driver.Value{1},
			deleteStudioQueryResult:  sqlmock.NewResult(0, 1),
			deleteStudioQueryError:   errDummy,
			rollbackCalled:           true,
			expectedCode:             http.StatusInternalServerError,
			expectedError:            errors.ErrInternalDB,
		},
		{
			name:                     "error-create-studio",
			param:                    anime,
			selectCalled:             true,
			selectQuery:              `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:          []driver.Value{1, 1},
			selectQueryReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"created_at"}).AddRow(&now)},
			selectQueryError:         nil,
			saveCalled:               true,
			saveQuery:                `UPDATE "anime" SET "title"=$1,"title_synonym"=$2,"title_english"=$3,"title_japanese"=$4,"picture"=$5,"start_day"=$6,"start_month"=$7,"start_year"=$8,"end_day"=$9,"end_month"=$10,"end_year"=$11,"synopsis"=$12,"nsfw"=$13,"type"=$14,"status"=$15,"episode"=$16,"episode_duration"=$17,"season"=$18,"season_year"=$19,"broadcast_day"=$20,"broadcast_time"=$21,"source"=$22,"rating"=$23,"background"=$24,"mean"=$25,"rank"=$26,"popularity"=$27,"member"=$28,"voter"=$29,"user_watching"=$30,"user_completed"=$31,"user_on_hold"=$32,"user_dropped"=$33,"user_planned"=$34,"created_at"=$35,"updated_at"=$36,"deleted_at"=$37 WHERE "anime"."deleted_at" IS NULL AND "id" = $38`,
			saveQueryArgs:            []driver.Value{anime.Title, "[]", "", "", "", 0, 0, 0, 0, 0, 0, "", false, "", "", 0, 0, "", 0, "", "", "", "", "", 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, now, sqlmock.AnyArg(), nil, 1},
			saveQueryResult:          sqlmock.NewResult(0, 1),
			saveQueryError:           nil,
			deleteGenreCalled:        true,
			deleteGenreQuery:         `DELETE FROM "anime_genre" WHERE anime_id = $1`,
			deleteGenreQueryArgs:     []driver.Value{1},
			deleteGenreQueryResult:   sqlmock.NewResult(0, 1),
			deleteGenreQueryError:    nil,
			createGenreCalled:        true,
			createGenreQuery:         `INSERT INTO "anime_genre" ("anime_id","genre_id") VALUES ($1,$2)`,
			createGenreQueryArgs:     []driver.Value{1, 2},
			createGenreQueryResult:   sqlmock.NewResult(0, 1),
			createGenreQueryError:    nil,
			deletePictureCalled:      true,
			deletePictureQuery:       `DELETE FROM "anime_picture" WHERE anime_id = $1`,
			deletePictureQueryArgs:   []driver.Value{1},
			deletePictureQueryResult: sqlmock.NewResult(0, 1),
			deletePictureQueryError:  nil,
			createPictureCalled:      true,
			createPictureQuery:       `INSERT INTO "anime_picture" ("anime_id","url") VALUES ($1,$2)`,
			createPictureQueryArgs:   []driver.Value{1, "www"},
			createPictureQueryResult: sqlmock.NewResult(0, 1),
			createPictureQueryError:  nil,
			deleteRelatedCalled:      true,
			deleteRelatedQuery:       `DELETE FROM "anime_related" WHERE anime_id1 = $1`,
			deleteRelatedQueryArgs:   []driver.Value{1},
			deleteRelatedQueryResult: sqlmock.NewResult(0, 1),
			deleteRelatedQueryError:  nil,
			createRelatedCalled:      true,
			createRelatedQuery:       `INSERT INTO "anime_related" ("anime_id1","anime_id2","relation") VALUES ($1,$2,$3)`,
			createRelatedQueryArgs:   []driver.Value{1, 3, "FULL_STORY"},
			createRelatedQueryResult: sqlmock.NewResult(0, 1),
			createRelatedQueryError:  nil,
			deleteStudioCalled:       true,
			deleteStudioQuery:        `DELETE FROM "anime_studio" WHERE anime_id = $1`,
			deleteStudioQueryArgs:    []driver.Value{1},
			deleteStudioQueryResult:  sqlmock.NewResult(0, 1),
			deleteStudioQueryError:   nil,
			createStudioCalled:       true,
			createStudioQuery:        `INSERT INTO "anime_studio" ("anime_id","studio_id") VALUES ($1,$2)`,
			createStudioQueryArgs:    []driver.Value{1, 4},
			createStudioQueryResult:  sqlmock.NewResult(0, 1),
			createStudioQueryError:   errDummy,
			rollbackCalled:           true,
			expectedCode:             http.StatusInternalServerError,
			expectedError:            errors.ErrInternalDB,
		},
		{
			name:                     "error-create-history",
			param:                    anime,
			selectCalled:             true,
			selectQuery:              `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:          []driver.Value{1, 1},
			selectQueryReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"created_at"}).AddRow(&now)},
			selectQueryError:         nil,
			saveCalled:               true,
			saveQuery:                `UPDATE "anime" SET "title"=$1,"title_synonym"=$2,"title_english"=$3,"title_japanese"=$4,"picture"=$5,"start_day"=$6,"start_month"=$7,"start_year"=$8,"end_day"=$9,"end_month"=$10,"end_year"=$11,"synopsis"=$12,"nsfw"=$13,"type"=$14,"status"=$15,"episode"=$16,"episode_duration"=$17,"season"=$18,"season_year"=$19,"broadcast_day"=$20,"broadcast_time"=$21,"source"=$22,"rating"=$23,"background"=$24,"mean"=$25,"rank"=$26,"popularity"=$27,"member"=$28,"voter"=$29,"user_watching"=$30,"user_completed"=$31,"user_on_hold"=$32,"user_dropped"=$33,"user_planned"=$34,"created_at"=$35,"updated_at"=$36,"deleted_at"=$37 WHERE "anime"."deleted_at" IS NULL AND "id" = $38`,
			saveQueryArgs:            []driver.Value{anime.Title, "[]", "", "", "", 0, 0, 0, 0, 0, 0, "", false, "", "", 0, 0, "", 0, "", "", "", "", "", 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, now, sqlmock.AnyArg(), nil, 1},
			saveQueryResult:          sqlmock.NewResult(0, 1),
			saveQueryError:           nil,
			deleteGenreCalled:        true,
			deleteGenreQuery:         `DELETE FROM "anime_genre" WHERE anime_id = $1`,
			deleteGenreQueryArgs:     []driver.Value{1},
			deleteGenreQueryResult:   sqlmock.NewResult(0, 1),
			deleteGenreQueryError:    nil,
			createGenreCalled:        true,
			createGenreQuery:         `INSERT INTO "anime_genre" ("anime_id","genre_id") VALUES ($1,$2)`,
			createGenreQueryArgs:     []driver.Value{1, 2},
			createGenreQueryResult:   sqlmock.NewResult(0, 1),
			createGenreQueryError:    nil,
			deletePictureCalled:      true,
			deletePictureQuery:       `DELETE FROM "anime_picture" WHERE anime_id = $1`,
			deletePictureQueryArgs:   []driver.Value{1},
			deletePictureQueryResult: sqlmock.NewResult(0, 1),
			deletePictureQueryError:  nil,
			createPictureCalled:      true,
			createPictureQuery:       `INSERT INTO "anime_picture" ("anime_id","url") VALUES ($1,$2)`,
			createPictureQueryArgs:   []driver.Value{1, "www"},
			createPictureQueryResult: sqlmock.NewResult(0, 1),
			createPictureQueryError:  nil,
			deleteRelatedCalled:      true,
			deleteRelatedQuery:       `DELETE FROM "anime_related" WHERE anime_id1 = $1`,
			deleteRelatedQueryArgs:   []driver.Value{1},
			deleteRelatedQueryResult: sqlmock.NewResult(0, 1),
			deleteRelatedQueryError:  nil,
			createRelatedCalled:      true,
			createRelatedQuery:       `INSERT INTO "anime_related" ("anime_id1","anime_id2","relation") VALUES ($1,$2,$3)`,
			createRelatedQueryArgs:   []driver.Value{1, 3, "FULL_STORY"},
			createRelatedQueryResult: sqlmock.NewResult(0, 1),
			createRelatedQueryError:  nil,
			deleteStudioCalled:       true,
			deleteStudioQuery:        `DELETE FROM "anime_studio" WHERE anime_id = $1`,
			deleteStudioQueryArgs:    []driver.Value{1},
			deleteStudioQueryResult:  sqlmock.NewResult(0, 1),
			deleteStudioQueryError:   nil,
			createStudioCalled:       true,
			createStudioQuery:        `INSERT INTO "anime_studio" ("anime_id","studio_id") VALUES ($1,$2)`,
			createStudioQueryArgs:    []driver.Value{1, 4},
			createStudioQueryResult:  sqlmock.NewResult(0, 1),
			createStudioQueryError:   nil,
			createHistoryCalled:      true,
			createHistoryQuery:       `INSERT INTO "anime_stats_history" ("anime_id","mean","rank","popularity","member","voter","user_watching","user_completed","user_on_hold","user_dropped","user_planned","created_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id`,
			createHistoryQueryArgs:   []driver.Value{1, 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, sqlmock.AnyArg()},
			createHistoryQueryReturn: []*sqlmock.Rows{},
			createHistoryQueryError:  errDummy,
			rollbackCalled:           true,
			expectedCode:             http.StatusInternalServerError,
			expectedError:            errors.ErrInternalDB,
		},
		{
			name:                     "error-commit",
			param:                    anime,
			selectCalled:             true,
			selectQuery:              `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:          []driver.Value{1, 1},
			selectQueryReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"created_at"}).AddRow(&now)},
			selectQueryError:         nil,
			saveCalled:               true,
			saveQuery:                `UPDATE "anime" SET "title"=$1,"title_synonym"=$2,"title_english"=$3,"title_japanese"=$4,"picture"=$5,"start_day"=$6,"start_month"=$7,"start_year"=$8,"end_day"=$9,"end_month"=$10,"end_year"=$11,"synopsis"=$12,"nsfw"=$13,"type"=$14,"status"=$15,"episode"=$16,"episode_duration"=$17,"season"=$18,"season_year"=$19,"broadcast_day"=$20,"broadcast_time"=$21,"source"=$22,"rating"=$23,"background"=$24,"mean"=$25,"rank"=$26,"popularity"=$27,"member"=$28,"voter"=$29,"user_watching"=$30,"user_completed"=$31,"user_on_hold"=$32,"user_dropped"=$33,"user_planned"=$34,"created_at"=$35,"updated_at"=$36,"deleted_at"=$37 WHERE "anime"."deleted_at" IS NULL AND "id" = $38`,
			saveQueryArgs:            []driver.Value{anime.Title, "[]", "", "", "", 0, 0, 0, 0, 0, 0, "", false, "", "", 0, 0, "", 0, "", "", "", "", "", 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, now, sqlmock.AnyArg(), nil, 1},
			saveQueryResult:          sqlmock.NewResult(0, 1),
			saveQueryError:           nil,
			deleteGenreCalled:        true,
			deleteGenreQuery:         `DELETE FROM "anime_genre" WHERE anime_id = $1`,
			deleteGenreQueryArgs:     []driver.Value{1},
			deleteGenreQueryResult:   sqlmock.NewResult(0, 1),
			deleteGenreQueryError:    nil,
			createGenreCalled:        true,
			createGenreQuery:         `INSERT INTO "anime_genre" ("anime_id","genre_id") VALUES ($1,$2)`,
			createGenreQueryArgs:     []driver.Value{1, 2},
			createGenreQueryResult:   sqlmock.NewResult(0, 1),
			createGenreQueryError:    nil,
			deletePictureCalled:      true,
			deletePictureQuery:       `DELETE FROM "anime_picture" WHERE anime_id = $1`,
			deletePictureQueryArgs:   []driver.Value{1},
			deletePictureQueryResult: sqlmock.NewResult(0, 1),
			deletePictureQueryError:  nil,
			createPictureCalled:      true,
			createPictureQuery:       `INSERT INTO "anime_picture" ("anime_id","url") VALUES ($1,$2)`,
			createPictureQueryArgs:   []driver.Value{1, "www"},
			createPictureQueryResult: sqlmock.NewResult(0, 1),
			createPictureQueryError:  nil,
			deleteRelatedCalled:      true,
			deleteRelatedQuery:       `DELETE FROM "anime_related" WHERE anime_id1 = $1`,
			deleteRelatedQueryArgs:   []driver.Value{1},
			deleteRelatedQueryResult: sqlmock.NewResult(0, 1),
			deleteRelatedQueryError:  nil,
			createRelatedCalled:      true,
			createRelatedQuery:       `INSERT INTO "anime_related" ("anime_id1","anime_id2","relation") VALUES ($1,$2,$3)`,
			createRelatedQueryArgs:   []driver.Value{1, 3, "FULL_STORY"},
			createRelatedQueryResult: sqlmock.NewResult(0, 1),
			createRelatedQueryError:  nil,
			deleteStudioCalled:       true,
			deleteStudioQuery:        `DELETE FROM "anime_studio" WHERE anime_id = $1`,
			deleteStudioQueryArgs:    []driver.Value{1},
			deleteStudioQueryResult:  sqlmock.NewResult(0, 1),
			deleteStudioQueryError:   nil,
			createStudioCalled:       true,
			createStudioQuery:        `INSERT INTO "anime_studio" ("anime_id","studio_id") VALUES ($1,$2)`,
			createStudioQueryArgs:    []driver.Value{1, 4},
			createStudioQueryResult:  sqlmock.NewResult(0, 1),
			createStudioQueryError:   nil,
			createHistoryCalled:      true,
			createHistoryQuery:       `INSERT INTO "anime_stats_history" ("anime_id","mean","rank","popularity","member","voter","user_watching","user_completed","user_on_hold","user_dropped","user_planned","created_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id`,
			createHistoryQueryArgs:   []driver.Value{1, 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, sqlmock.AnyArg()},
			createHistoryQueryReturn: []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)},
			createHistoryQueryError:  nil,
			commitCalled:             true,
			commitError:              errDummy,
			expectedCode:             http.StatusInternalServerError,
			expectedError:            errors.ErrInternalDB,
		},
		{
			name:                     "ok",
			param:                    anime,
			selectCalled:             true,
			selectQuery:              `SELECT "created_at" FROM "anime" WHERE id = $1 AND "anime"."deleted_at" IS NULL ORDER BY "anime"."id" LIMIT $2`,
			selectQueryArgs:          []driver.Value{1, 1},
			selectQueryReturn:        []*sqlmock.Rows{sqlmock.NewRows([]string{"created_at"}).AddRow(&now)},
			selectQueryError:         nil,
			saveCalled:               true,
			saveQuery:                `UPDATE "anime" SET "title"=$1,"title_synonym"=$2,"title_english"=$3,"title_japanese"=$4,"picture"=$5,"start_day"=$6,"start_month"=$7,"start_year"=$8,"end_day"=$9,"end_month"=$10,"end_year"=$11,"synopsis"=$12,"nsfw"=$13,"type"=$14,"status"=$15,"episode"=$16,"episode_duration"=$17,"season"=$18,"season_year"=$19,"broadcast_day"=$20,"broadcast_time"=$21,"source"=$22,"rating"=$23,"background"=$24,"mean"=$25,"rank"=$26,"popularity"=$27,"member"=$28,"voter"=$29,"user_watching"=$30,"user_completed"=$31,"user_on_hold"=$32,"user_dropped"=$33,"user_planned"=$34,"created_at"=$35,"updated_at"=$36,"deleted_at"=$37 WHERE "anime"."deleted_at" IS NULL AND "id" = $38`,
			saveQueryArgs:            []driver.Value{anime.Title, "[]", "", "", "", 0, 0, 0, 0, 0, 0, "", false, "", "", 0, 0, "", 0, "", "", "", "", "", 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, now, sqlmock.AnyArg(), nil, 1},
			saveQueryResult:          sqlmock.NewResult(0, 1),
			saveQueryError:           nil,
			deleteGenreCalled:        true,
			deleteGenreQuery:         `DELETE FROM "anime_genre" WHERE anime_id = $1`,
			deleteGenreQueryArgs:     []driver.Value{1},
			deleteGenreQueryResult:   sqlmock.NewResult(0, 1),
			deleteGenreQueryError:    nil,
			createGenreCalled:        true,
			createGenreQuery:         `INSERT INTO "anime_genre" ("anime_id","genre_id") VALUES ($1,$2)`,
			createGenreQueryArgs:     []driver.Value{1, 2},
			createGenreQueryResult:   sqlmock.NewResult(0, 1),
			createGenreQueryError:    nil,
			deletePictureCalled:      true,
			deletePictureQuery:       `DELETE FROM "anime_picture" WHERE anime_id = $1`,
			deletePictureQueryArgs:   []driver.Value{1},
			deletePictureQueryResult: sqlmock.NewResult(0, 1),
			deletePictureQueryError:  nil,
			createPictureCalled:      true,
			createPictureQuery:       `INSERT INTO "anime_picture" ("anime_id","url") VALUES ($1,$2)`,
			createPictureQueryArgs:   []driver.Value{1, "www"},
			createPictureQueryResult: sqlmock.NewResult(0, 1),
			createPictureQueryError:  nil,
			deleteRelatedCalled:      true,
			deleteRelatedQuery:       `DELETE FROM "anime_related" WHERE anime_id1 = $1`,
			deleteRelatedQueryArgs:   []driver.Value{1},
			deleteRelatedQueryResult: sqlmock.NewResult(0, 1),
			deleteRelatedQueryError:  nil,
			createRelatedCalled:      true,
			createRelatedQuery:       `INSERT INTO "anime_related" ("anime_id1","anime_id2","relation") VALUES ($1,$2,$3)`,
			createRelatedQueryArgs:   []driver.Value{1, 3, "FULL_STORY"},
			createRelatedQueryResult: sqlmock.NewResult(0, 1),
			createRelatedQueryError:  nil,
			deleteStudioCalled:       true,
			deleteStudioQuery:        `DELETE FROM "anime_studio" WHERE anime_id = $1`,
			deleteStudioQueryArgs:    []driver.Value{1},
			deleteStudioQueryResult:  sqlmock.NewResult(0, 1),
			deleteStudioQueryError:   nil,
			createStudioCalled:       true,
			createStudioQuery:        `INSERT INTO "anime_studio" ("anime_id","studio_id") VALUES ($1,$2)`,
			createStudioQueryArgs:    []driver.Value{1, 4},
			createStudioQueryResult:  sqlmock.NewResult(0, 1),
			createStudioQueryError:   nil,
			createHistoryCalled:      true,
			createHistoryQuery:       `INSERT INTO "anime_stats_history" ("anime_id","mean","rank","popularity","member","voter","user_watching","user_completed","user_on_hold","user_dropped","user_planned","created_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id`,
			createHistoryQueryArgs:   []driver.Value{1, 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, sqlmock.AnyArg()},
			createHistoryQueryReturn: []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)},
			createHistoryQueryError:  nil,
			commitCalled:             true,
			commitError:              nil,
			expectedCode:             http.StatusOK,
			expectedError:            nil,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.dbMock.ExpectBegin().WillReturnError(test.beginError)

			if test.selectCalled {
				suite.dbMock.ExpectQuery(regexp.QuoteMeta(test.selectQuery)).
					WithArgs(test.selectQueryArgs...).
					WillReturnRows(test.selectQueryReturn...).
					WillReturnError(test.selectQueryError)
			}

			if test.saveCalled {
				suite.dbMock.ExpectExec(regexp.QuoteMeta(test.saveQuery)).
					WithArgs(test.saveQueryArgs...).
					WillReturnResult(test.saveQueryResult).
					WillReturnError(test.saveQueryError)
			}

			if test.deleteGenreCalled {
				suite.dbMock.ExpectExec(regexp.QuoteMeta(test.deleteGenreQuery)).
					WithArgs(test.deleteGenreQueryArgs...).
					WillReturnResult(test.deleteGenreQueryResult).
					WillReturnError(test.deleteGenreQueryError)
			}

			if test.createGenreCalled {
				suite.dbMock.ExpectExec(regexp.QuoteMeta(test.createGenreQuery)).
					WithArgs(test.createGenreQueryArgs...).
					WillReturnResult(test.createGenreQueryResult).
					WillReturnError(test.createGenreQueryError)
			}

			if test.deletePictureCalled {
				suite.dbMock.ExpectExec(regexp.QuoteMeta(test.deletePictureQuery)).
					WithArgs(test.deletePictureQueryArgs...).
					WillReturnResult(test.deletePictureQueryResult).
					WillReturnError(test.deletePictureQueryError)
			}

			if test.createPictureCalled {
				suite.dbMock.ExpectExec(regexp.QuoteMeta(test.createPictureQuery)).
					WithArgs(test.createPictureQueryArgs...).
					WillReturnResult(test.createPictureQueryResult).
					WillReturnError(test.createPictureQueryError)
			}

			if test.deleteRelatedCalled {
				suite.dbMock.ExpectExec(regexp.QuoteMeta(test.deleteRelatedQuery)).
					WithArgs(test.deleteRelatedQueryArgs...).
					WillReturnResult(test.deleteRelatedQueryResult).
					WillReturnError(test.deleteRelatedQueryError)
			}

			if test.createRelatedCalled {
				suite.dbMock.ExpectExec(regexp.QuoteMeta(test.createRelatedQuery)).
					WithArgs(test.createRelatedQueryArgs...).
					WillReturnResult(test.createRelatedQueryResult).
					WillReturnError(test.createRelatedQueryError)
			}

			if test.deleteStudioCalled {
				suite.dbMock.ExpectExec(regexp.QuoteMeta(test.deleteStudioQuery)).
					WithArgs(test.deleteStudioQueryArgs...).
					WillReturnResult(test.deleteStudioQueryResult).
					WillReturnError(test.deleteStudioQueryError)
			}

			if test.createStudioCalled {
				suite.dbMock.ExpectExec(regexp.QuoteMeta(test.createStudioQuery)).
					WithArgs(test.createStudioQueryArgs...).
					WillReturnResult(test.createStudioQueryResult).
					WillReturnError(test.createStudioQueryError)
			}

			if test.createHistoryCalled {
				suite.dbMock.ExpectQuery(regexp.QuoteMeta(test.createHistoryQuery)).
					WithArgs(test.createHistoryQueryArgs...).
					WillReturnRows(test.createHistoryQueryReturn...).
					WillReturnError(test.createHistoryQueryError)
			}

			if test.rollbackCalled {
				suite.dbMock.ExpectRollback()
			}

			if test.commitCalled {
				suite.dbMock.ExpectCommit().WillReturnError(test.commitError)
			}

			sql := sql.New(suite.db, 0, 0, 0)

			code, err := sql.Update(ctx, test.param)
			suite.Equal(test.expectedCode, code)
			suite.ErrorIs(test.expectedError, err)
			suite.Nil(suite.dbMock.ExpectationsWereMet())
		})
	}
}
