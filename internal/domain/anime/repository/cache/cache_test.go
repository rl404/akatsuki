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
