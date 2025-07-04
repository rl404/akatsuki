package service_test

import (
	"context"
	_errors "errors"
	"net/http"
	"testing"

	"github.com/rl404/akatsuki/internal/domain/anime/entity"
	entityGenre "github.com/rl404/akatsuki/internal/domain/genre/entity"
	entityStudio "github.com/rl404/akatsuki/internal/domain/studio/entity"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/service"
	mockAnime "github.com/rl404/akatsuki/tests/mocks/domain/anime"
	mockEmptyID "github.com/rl404/akatsuki/tests/mocks/domain/empty_id"
	mockGenre "github.com/rl404/akatsuki/tests/mocks/domain/genre"
	mockPublisher "github.com/rl404/akatsuki/tests/mocks/domain/publisher"
	mockStudio "github.com/rl404/akatsuki/tests/mocks/domain/studio"
	"github.com/stretchr/testify/suite"
)

type animeTestSuite struct {
	suite.Suite
	animeMock     *mockAnime.Repository
	emptyIDMock   *mockEmptyID.Repository
	genreMock     *mockGenre.Repository
	studioMock    *mockStudio.Repository
	publisherMock *mockPublisher.Repository
}

func TestAnime(t *testing.T) {
	suite.Run(t, &animeTestSuite{
		animeMock:     new(mockAnime.Repository),
		emptyIDMock:   new(mockEmptyID.Repository),
		genreMock:     new(mockGenre.Repository),
		studioMock:    new(mockStudio.Repository),
		publisherMock: new(mockPublisher.Repository),
	})
}

func (suite *animeTestSuite) TestGetAnime() {
	ctx := context.Background()
	errDummy := _errors.New("dummy error")

	tests := []struct {
		name               string
		param              service.GetAnimeRequest
		repoCalled         bool
		repoParams         []interface{}
		repoReturn         []interface{}
		expectedReturn     []service.Anime
		expectedPagination *service.Pagination
		expectedCode       int
		expectedError      error
	}{
		{
			name:               "error-validate",
			param:              service.GetAnimeRequest{Type: "random"},
			expectedReturn:     nil,
			expectedPagination: nil,
			expectedCode:       http.StatusBadRequest,
			expectedError:      errors.ErrOneOfField("type", "TV/OVA/ONA/MOVIE/SPECIAL/MUSIC/CM/PV/TV_SPECIAL"),
		},
		{
			name:               "error-get",
			param:              service.GetAnimeRequest{},
			repoCalled:         true,
			repoParams:         []interface{}{ctx, entity.GetRequest{Sort: "RANK", Page: 1, Limit: 20}},
			repoReturn:         []interface{}{nil, 0, http.StatusInternalServerError, errDummy},
			expectedReturn:     nil,
			expectedPagination: nil,
			expectedCode:       http.StatusInternalServerError,
			expectedError:      errDummy,
		},
		{
			name:               "ok",
			param:              service.GetAnimeRequest{},
			repoCalled:         true,
			repoParams:         []interface{}{ctx, entity.GetRequest{Sort: "RANK", Page: 1, Limit: 20}},
			repoReturn:         []interface{}{[]*entity.Anime{{ID: 1}}, 1, http.StatusOK, nil},
			expectedReturn:     []service.Anime{{ID: 1, Genres: []service.AnimeGenre{}, Related: []service.AnimeRelated{}, Studios: []service.AnimeStudio{}}},
			expectedPagination: &service.Pagination{Page: 1, Limit: 20, Total: 1},
			expectedCode:       http.StatusOK,
			expectedError:      nil,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			if test.repoCalled {
				suite.animeMock.On("Get", test.repoParams...).Return(test.repoReturn...).Once()
			}

			s := service.New(suite.animeMock, nil, nil, nil, nil, nil, nil)

			data, pagination, code, err := s.GetAnime(ctx, test.param)
			suite.Equal(test.expectedReturn, data)
			suite.Equal(test.expectedPagination, pagination)
			suite.Equal(test.expectedCode, code)

			if test.expectedError != nil {
				suite.ErrorContains(test.expectedError, err.Error())
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *animeTestSuite) TestGetAnimeByID() {
	ctx := context.Background()
	errDummy := _errors.New("dummy error")

	tests := []struct {
		name                string
		param               int64
		repoEmptyIDCalled   bool
		repoEmptyIDParams   []interface{}
		repoEmptyIDReturn   []interface{}
		repoAnimeCalled     bool
		repoAnimeParams     []interface{}
		repoAnimeReturn     []interface{}
		repoPublisherCalled bool
		repoPublisherParams []interface{}
		repoPublisherReturn []interface{}
		repoGenreCalled     bool
		repoGenreParams     []interface{}
		repoGenreReturn     []interface{}
		repoRelatedCalled   bool
		repoRelatedParams   []interface{}
		repoRelatedReturn   []interface{}
		repoStudioCalled    bool
		repoStudioParams    []interface{}
		repoStudioReturn    []interface{}
		expectedReturn      *service.Anime
		expectedCode        int
		expectedError       error
	}{
		{
			name:           "invalid-id",
			param:          -1,
			expectedReturn: nil,
			expectedCode:   http.StatusBadRequest,
			expectedError:  errors.ErrInvalidAnimeID,
		},
		{
			name:              "error-empty-id",
			param:             1,
			repoEmptyIDCalled: true,
			repoEmptyIDParams: []interface{}{ctx, int64(1)},
			repoEmptyIDReturn: []interface{}{int64(0), http.StatusInternalServerError, errDummy},
			expectedReturn:    nil,
			expectedCode:      http.StatusInternalServerError,
			expectedError:     errDummy,
		},
		{
			name:              "empty-anime",
			param:             1,
			repoEmptyIDCalled: true,
			repoEmptyIDParams: []interface{}{ctx, int64(1)},
			repoEmptyIDReturn: []interface{}{int64(0), http.StatusOK, nil},
			expectedReturn:    nil,
			expectedCode:      http.StatusNotFound,
			expectedError:     errors.ErrAnimeNotFound,
		},
		{
			name:                "error-publisher",
			param:               1,
			repoEmptyIDCalled:   true,
			repoEmptyIDParams:   []interface{}{ctx, int64(1)},
			repoEmptyIDReturn:   []interface{}{int64(0), http.StatusNotFound, errDummy},
			repoAnimeCalled:     true,
			repoAnimeParams:     []interface{}{ctx, int64(1)},
			repoAnimeReturn:     []interface{}{nil, http.StatusNotFound, errDummy},
			repoPublisherCalled: true,
			repoPublisherParams: []interface{}{ctx, int64(1), false},
			repoPublisherReturn: []interface{}{errDummy},
			expectedReturn:      nil,
			expectedCode:        http.StatusInternalServerError,
			expectedError:       errDummy,
		},
		{
			name:                "ok-publisher",
			param:               1,
			repoEmptyIDCalled:   true,
			repoEmptyIDParams:   []interface{}{ctx, int64(1)},
			repoEmptyIDReturn:   []interface{}{int64(0), http.StatusNotFound, errDummy},
			repoAnimeCalled:     true,
			repoAnimeParams:     []interface{}{ctx, int64(1)},
			repoAnimeReturn:     []interface{}{nil, http.StatusNotFound, errDummy},
			repoPublisherCalled: true,
			repoPublisherParams: []interface{}{ctx, int64(1), false},
			repoPublisherReturn: []interface{}{nil},
			expectedReturn:      nil,
			expectedCode:        http.StatusAccepted,
			expectedError:       nil,
		},
		{
			name:              "error-anime",
			param:             1,
			repoEmptyIDCalled: true,
			repoEmptyIDParams: []interface{}{ctx, int64(1)},
			repoEmptyIDReturn: []interface{}{int64(0), http.StatusNotFound, errDummy},
			repoAnimeCalled:   true,
			repoAnimeParams:   []interface{}{ctx, int64(1)},
			repoAnimeReturn:   []interface{}{nil, http.StatusInternalServerError, errDummy},
			expectedReturn:    nil,
			expectedCode:      http.StatusInternalServerError,
			expectedError:     errDummy,
		},
		{
			name:              "error-genre",
			param:             1,
			repoEmptyIDCalled: true,
			repoEmptyIDParams: []interface{}{ctx, int64(1)},
			repoEmptyIDReturn: []interface{}{int64(0), http.StatusNotFound, errDummy},
			repoAnimeCalled:   true,
			repoAnimeParams:   []interface{}{ctx, int64(1)},
			repoAnimeReturn:   []interface{}{&entity.Anime{ID: 1, GenreIDs: []int64{2}}, http.StatusOK, nil},
			repoGenreCalled:   true,
			repoGenreParams:   []interface{}{ctx, []int64{2}},
			repoGenreReturn:   []interface{}{nil, http.StatusInternalServerError, errDummy},
			expectedReturn:    nil,
			expectedCode:      http.StatusInternalServerError,
			expectedError:     errDummy,
		},
		{
			name:              "error-related",
			param:             1,
			repoEmptyIDCalled: true,
			repoEmptyIDParams: []interface{}{ctx, int64(1)},
			repoEmptyIDReturn: []interface{}{int64(0), http.StatusNotFound, errDummy},
			repoAnimeCalled:   true,
			repoAnimeParams:   []interface{}{ctx, int64(1)},
			repoAnimeReturn:   []interface{}{&entity.Anime{ID: 1, GenreIDs: []int64{2}, Related: []entity.Related{{ID: 3, Relation: entity.RelationAdaptation}}}, http.StatusOK, nil},
			repoGenreCalled:   true,
			repoGenreParams:   []interface{}{ctx, []int64{2}},
			repoGenreReturn:   []interface{}{[]*entityGenre.Genre{{ID: 2, Name: "genre"}}, http.StatusOK, nil},
			repoRelatedCalled: true,
			repoRelatedParams: []interface{}{ctx, []int64{3}},
			repoRelatedReturn: []interface{}{nil, http.StatusInternalServerError, errDummy},
			expectedReturn:    nil,
			expectedCode:      http.StatusInternalServerError,
			expectedError:     errDummy,
		},
		{
			name:              "error-studio",
			param:             1,
			repoEmptyIDCalled: true,
			repoEmptyIDParams: []interface{}{ctx, int64(1)},
			repoEmptyIDReturn: []interface{}{int64(0), http.StatusNotFound, errDummy},
			repoAnimeCalled:   true,
			repoAnimeParams:   []interface{}{ctx, int64(1)},
			repoAnimeReturn:   []interface{}{&entity.Anime{ID: 1, GenreIDs: []int64{2}, Related: []entity.Related{{ID: 3, Relation: entity.RelationAdaptation}}, StudioIDs: []int64{4}}, http.StatusOK, nil},
			repoGenreCalled:   true,
			repoGenreParams:   []interface{}{ctx, []int64{2}},
			repoGenreReturn:   []interface{}{[]*entityGenre.Genre{{ID: 2, Name: "genre"}}, http.StatusOK, nil},
			repoRelatedCalled: true,
			repoRelatedParams: []interface{}{ctx, []int64{3}},
			repoRelatedReturn: []interface{}{[]*entity.Anime{{ID: 3, Title: "title", Picture: "picture"}}, http.StatusOK, nil},
			repoStudioCalled:  true,
			repoStudioParams:  []interface{}{ctx, []int64{4}},
			repoStudioReturn:  []interface{}{nil, http.StatusInternalServerError, errDummy},
			expectedReturn:    nil,
			expectedCode:      http.StatusInternalServerError,
			expectedError:     errDummy,
		},
		{
			name:              "ok",
			param:             1,
			repoEmptyIDCalled: true,
			repoEmptyIDParams: []interface{}{ctx, int64(1)},
			repoEmptyIDReturn: []interface{}{int64(0), http.StatusNotFound, errDummy},
			repoAnimeCalled:   true,
			repoAnimeParams:   []interface{}{ctx, int64(1)},
			repoAnimeReturn:   []interface{}{&entity.Anime{ID: 1, GenreIDs: []int64{2}, Related: []entity.Related{{ID: 3, Relation: entity.RelationAdaptation}}, StudioIDs: []int64{4}}, http.StatusOK, nil},
			repoGenreCalled:   true,
			repoGenreParams:   []interface{}{ctx, []int64{2}},
			repoGenreReturn:   []interface{}{[]*entityGenre.Genre{{ID: 2, Name: "genre"}}, http.StatusOK, nil},
			repoRelatedCalled: true,
			repoRelatedParams: []interface{}{ctx, []int64{3}},
			repoRelatedReturn: []interface{}{[]*entity.Anime{{ID: 3, Title: "title", Picture: "picture"}}, http.StatusOK, nil},
			repoStudioCalled:  true,
			repoStudioParams:  []interface{}{ctx, []int64{4}},
			repoStudioReturn:  []interface{}{[]*entityStudio.Studio{{ID: 4, Name: "studio"}}, http.StatusOK, nil},
			expectedReturn: &service.Anime{
				ID:      1,
				Genres:  []service.AnimeGenre{{ID: 2, Name: "genre"}},
				Related: []service.AnimeRelated{{ID: 3, Title: "title", Picture: "picture", Relation: entity.RelationAdaptation}},
				Studios: []service.AnimeStudio{{ID: 4, Name: "studio"}},
			},
			expectedCode:  http.StatusOK,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			if test.repoEmptyIDCalled {
				suite.emptyIDMock.On("Get", test.repoEmptyIDParams...).Return(test.repoEmptyIDReturn...).Once()
			}

			if test.repoAnimeCalled {
				suite.animeMock.On("GetByID", test.repoAnimeParams...).Return(test.repoAnimeReturn...).Once()
			}

			if test.repoPublisherCalled {
				suite.publisherMock.On("PublishParseAnime", test.repoPublisherParams...).Return(test.repoPublisherReturn...).Once()
			}

			if test.repoGenreCalled {
				suite.genreMock.On("GetByIDs", test.repoGenreParams...).Return(test.repoGenreReturn...).Once()
			}

			if test.repoRelatedCalled {
				suite.animeMock.On("GetByIDs", test.repoRelatedParams...).Return(test.repoRelatedReturn...).Once()
			}

			if test.repoStudioCalled {
				suite.studioMock.On("GetByIDs", test.repoStudioParams...).Return(test.repoStudioReturn...).Once()
			}

			s := service.New(suite.animeMock, suite.genreMock, suite.studioMock, nil, suite.emptyIDMock, suite.publisherMock, nil)

			data, code, err := s.GetAnimeByID(ctx, test.param)
			suite.Equal(test.expectedReturn, data)
			suite.Equal(test.expectedCode, code)

			if test.expectedError != nil {
				suite.ErrorContains(test.expectedError, err.Error())
			} else {
				suite.Nil(err)
			}
		})
	}
}
