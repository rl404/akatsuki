package cache_test

import (
	"context"
	_errors "errors"
	"net/http"
	"testing"

	"github.com/rl404/akatsuki/internal/domain/anime/entity"
	"github.com/rl404/akatsuki/internal/domain/anime/repository/cache"
	"github.com/rl404/akatsuki/internal/errors"
	"github.com/rl404/akatsuki/internal/utils"
	mockCacher "github.com/rl404/akatsuki/tests/mocks/cacher"
	mockRepository "github.com/rl404/akatsuki/tests/mocks/domain/anime"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	cacherMock *mockCacher.Cacher
	repoMock   *mockRepository.Repository
}

func TestCache(t *testing.T) {
	suite.Run(t, &testSuite{
		cacherMock: new(mockCacher.Cacher),
		repoMock:   new(mockRepository.Repository),
	})
}

func (suite *testSuite) TestGet() {
	ctx := context.Background()

	request := entity.GetRequest{
		Title: "title",
	}

	suite.repoMock.On("Get", ctx, request).Return([]*entity.Anime{{
		ID: 2,
	}}, 1, http.StatusOK, nil)

	c := cache.New(suite.cacherMock, suite.repoMock)

	anime, total, code, err := c.Get(ctx, request)
	suite.Equal(int64(2), anime[0].ID)
	suite.Equal(1, total)
	suite.Equal(http.StatusOK, code)
	suite.Nil(err)
}

func (suite *testSuite) TestGetByID() {
	ctx := context.Background()

	var emptyAnime *entity.Anime
	key := utils.GetKey("anime", 1)
	errDummy := _errors.New("dummy error")
	anime := entity.Anime{ID: 1}

	tests := []struct {
		name            string
		id              int64
		cacherGetCalled bool
		cacherGetParams []interface{}
		cacherGetReturn error
		repoCalled      bool
		repoParams      []interface{}
		repoReturn      []interface{}
		cacherSetCalled bool
		cacherSetParams []interface{}
		cacherSetReturn error
		expectedReturn  *entity.Anime
		expectedCode    int
		expectedError   error
	}{
		{
			name:            "from-cache",
			id:              1,
			cacherGetCalled: true,
			cacherGetParams: []interface{}{ctx, key, &emptyAnime},
			cacherGetReturn: nil,
			expectedReturn:  emptyAnime,
			expectedCode:    http.StatusOK,
			expectedError:   nil,
		},
		{
			name:            "error-get",
			id:              1,
			cacherGetCalled: true,
			cacherGetParams: []interface{}{ctx, key, &emptyAnime},
			cacherGetReturn: errDummy,
			repoCalled:      true,
			repoParams:      []interface{}{ctx, int64(1)},
			repoReturn:      []interface{}{nil, http.StatusInternalServerError, errDummy},
			expectedReturn:  nil,
			expectedCode:    http.StatusInternalServerError,
			expectedError:   errDummy,
		},
		{
			name:            "error-set",
			id:              1,
			cacherGetCalled: true,
			cacherGetParams: []interface{}{ctx, key, &emptyAnime},
			cacherGetReturn: errDummy,
			repoCalled:      true,
			repoParams:      []interface{}{ctx, int64(1)},
			repoReturn:      []interface{}{&anime, http.StatusOK, nil},
			cacherSetCalled: true,
			cacherSetParams: []interface{}{ctx, key, &anime},
			cacherSetReturn: errDummy,
			expectedReturn:  nil,
			expectedCode:    http.StatusInternalServerError,
			expectedError:   errors.ErrInternalCache,
		},
		{
			name:            "ok",
			id:              1,
			cacherGetCalled: true,
			cacherGetParams: []interface{}{ctx, key, &emptyAnime},
			cacherGetReturn: errDummy,
			repoCalled:      true,
			repoParams:      []interface{}{ctx, int64(1)},
			repoReturn:      []interface{}{&anime, http.StatusOK, nil},
			cacherSetCalled: true,
			cacherSetParams: []interface{}{ctx, key, &anime},
			cacherSetReturn: nil,
			expectedReturn:  &anime,
			expectedCode:    http.StatusOK,
			expectedError:   nil,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			if test.cacherGetCalled {
				suite.cacherMock.On("Get", test.cacherGetParams...).Return(test.cacherGetReturn).Once()
			}

			if test.repoCalled {
				suite.repoMock.On("GetByID", test.repoParams...).Return(test.repoReturn...).Once()
			}

			if test.cacherSetCalled {
				suite.cacherMock.On("Set", test.cacherSetParams...).Return(test.cacherSetReturn).Once()
			}

			c := cache.New(suite.cacherMock, suite.repoMock)

			anime, code, err := c.GetByID(ctx, test.id)
			suite.Equal(test.expectedReturn, anime)
			suite.Equal(test.expectedCode, code)
			suite.ErrorIs(test.expectedError, err)
		})

	}
}

func (suite *testSuite) TestGetByIDs() {
	ctx := context.Background()

	request := []int64{1, 2}

	suite.repoMock.On("GetByIDs", ctx, request).Return([]*entity.Anime{{
		ID: 2,
	}}, http.StatusOK, nil)

	c := cache.New(suite.cacherMock, suite.repoMock)

	anime, code, err := c.GetByIDs(ctx, request)
	suite.Equal(int64(2), anime[0].ID)
	suite.Equal(http.StatusOK, code)
	suite.Nil(err)
}

func (suite *testSuite) TestUpdate() {
	ctx := context.Background()

	anime := entity.Anime{ID: 1}
	key := utils.GetKey("anime", 1)
	errDummy := _errors.New("dummy error")

	tests := []struct {
		name          string
		data          entity.Anime
		repoCalled    bool
		repoParams    []interface{}
		repoReturn    []interface{}
		cacherCalled  bool
		cacherParams  []interface{}
		cacherReturn  []interface{}
		expectedCode  int
		expectedError error
	}{
		{
			name:          "error-update",
			data:          anime,
			repoCalled:    true,
			repoParams:    []interface{}{ctx, anime},
			repoReturn:    []interface{}{http.StatusInternalServerError, errDummy},
			expectedCode:  http.StatusInternalServerError,
			expectedError: errDummy,
		},
		{
			name:          "error-cache",
			data:          anime,
			repoCalled:    true,
			repoParams:    []interface{}{ctx, anime},
			repoReturn:    []interface{}{http.StatusOK, nil},
			cacherCalled:  true,
			cacherParams:  []interface{}{ctx, key},
			cacherReturn:  []interface{}{errDummy},
			expectedCode:  http.StatusInternalServerError,
			expectedError: errors.ErrInternalCache,
		},
		{
			name:          "ok",
			data:          anime,
			repoCalled:    true,
			repoParams:    []interface{}{ctx, anime},
			repoReturn:    []interface{}{http.StatusOK, nil},
			cacherCalled:  true,
			cacherParams:  []interface{}{ctx, key},
			cacherReturn:  []interface{}{nil},
			expectedCode:  http.StatusOK,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			if test.repoCalled {
				suite.repoMock.On("Update", test.repoParams...).Return(test.repoReturn...).Once()
			}

			if test.cacherCalled {
				suite.cacherMock.On("Delete", test.cacherParams...).Return(test.cacherReturn...).Once()
			}

			c := cache.New(suite.cacherMock, suite.repoMock)

			code, err := c.Update(ctx, test.data)
			suite.Equal(test.expectedCode, code)
			suite.ErrorIs(test.expectedError, err)
		})
	}
}

func (suite *testSuite) TestIsOld() {
	ctx := context.Background()

	request := int64(1)

	suite.repoMock.On("IsOld", ctx, request).Return(true, http.StatusOK, nil)

	c := cache.New(suite.cacherMock, suite.repoMock)

	isOld, code, err := c.IsOld(ctx, request)
	suite.Equal(true, isOld)
	suite.Equal(http.StatusOK, code)
	suite.Nil(err)
}

func (suite *testSuite) TestGetOldReleasingIDs() {
	ctx := context.Background()

	suite.repoMock.On("GetOldReleasingIDs", ctx).Return([]int64{1}, http.StatusOK, nil)

	c := cache.New(suite.cacherMock, suite.repoMock)

	res, code, err := c.GetOldReleasingIDs(ctx)
	suite.Equal([]int64{1}, res)
	suite.Equal(http.StatusOK, code)
	suite.Nil(err)
}

func (suite *testSuite) TestGetOldFinishedIDs() {
	ctx := context.Background()

	suite.repoMock.On("GetOldFinishedIDs", ctx).Return([]int64{1}, http.StatusOK, nil)

	c := cache.New(suite.cacherMock, suite.repoMock)

	res, code, err := c.GetOldFinishedIDs(ctx)
	suite.Equal([]int64{1}, res)
	suite.Equal(http.StatusOK, code)
	suite.Nil(err)
}

func (suite *testSuite) TestGetOldNotYetIDs() {
	ctx := context.Background()

	suite.repoMock.On("GetOldNotYetIDs", ctx).Return([]int64{1}, http.StatusOK, nil)

	c := cache.New(suite.cacherMock, suite.repoMock)

	res, code, err := c.GetOldNotYetIDs(ctx)
	suite.Equal([]int64{1}, res)
	suite.Equal(http.StatusOK, code)
	suite.Nil(err)
}

func (suite *testSuite) TestGetMaxID() {
	ctx := context.Background()

	suite.repoMock.On("GetMaxID", ctx).Return(int64(1), http.StatusOK, nil)

	c := cache.New(suite.cacherMock, suite.repoMock)

	res, code, err := c.GetMaxID(ctx)
	suite.Equal(int64(1), res)
	suite.Equal(http.StatusOK, code)
	suite.Nil(err)
}

func (suite *testSuite) TestGetIDs() {
	ctx := context.Background()

	suite.repoMock.On("GetIDs", ctx).Return([]int64{1}, http.StatusOK, nil)

	c := cache.New(suite.cacherMock, suite.repoMock)

	res, code, err := c.GetIDs(ctx)
	suite.Equal([]int64{1}, res)
	suite.Equal(http.StatusOK, code)
	suite.Nil(err)
}

func (suite *testSuite) TestGetRelatedByIDs() {
	ctx := context.Background()

	request := []int64{1}

	suite.repoMock.On("GetRelatedByIDs", ctx, request).Return([]*entity.AnimeRelated{{AnimeID1: 1, AnimeID2: 2}}, http.StatusOK, nil)

	c := cache.New(suite.cacherMock, suite.repoMock)

	res, code, err := c.GetRelatedByIDs(ctx, request)
	suite.Equal([]*entity.AnimeRelated{{AnimeID1: 1, AnimeID2: 2}}, res)
	suite.Equal(http.StatusOK, code)
	suite.Nil(err)
}

func (suite *testSuite) TestDeleteByID() {
	ctx := context.Background()

	request := int64(1)

	suite.repoMock.On("DeleteByID", ctx, request).Return(http.StatusOK, nil)

	c := cache.New(suite.cacherMock, suite.repoMock)

	code, err := c.DeleteByID(ctx, request)
	suite.Equal(http.StatusOK, code)
	suite.Nil(err)
}

func (suite *testSuite) TestGetHistories() {
	ctx := context.Background()

	request := entity.GetHistoriesRequest{AnimeID: 1}

	suite.repoMock.On("GetHistories", ctx, request).Return([]entity.History{{Rank: 1}}, http.StatusOK, nil)

	c := cache.New(suite.cacherMock, suite.repoMock)

	res, code, err := c.GetHistories(ctx, request)
	suite.Equal([]entity.History{{Rank: 1}}, res)
	suite.Equal(http.StatusOK, code)
	suite.Nil(err)
}
