package nagato

import (
	"context"
	"net/http"

	"github.com/rl404/nagato/mal"
)

// GetForumBoards to get forum board list.
func (c *Client) GetForumBoards() ([]ForumBoardCategory, int, error) {
	return c.GetForumBoardsWithContext(context.Background())
}

// GetForumBoardsWithContext to get forum board list with context.
func (c *Client) GetForumBoardsWithContext(ctx context.Context) ([]ForumBoardCategory, int, error) {
	boards, code, err := c.mal.GetForumBoardsWithContext(ctx)
	if err != nil {
		return nil, code, err
	}

	return c.forumCategoriesToForumCategories(boards.Categories), http.StatusOK, nil
}

// GetForumTopics to get forum topic list.
func (c *Client) GetForumTopics(param GetForumTopicsParam) ([]ForumTopic, int, error) {
	return c.GetForumTopicsWithContext(context.Background(), param)
}

// GetForumTopicsWithContext to get forum topic list with context.
func (c *Client) GetForumTopicsWithContext(ctx context.Context, param GetForumTopicsParam) ([]ForumTopic, int, error) {
	if err := c.validate(&param); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := param.validate(); err != nil {
		return nil, http.StatusBadRequest, err
	}

	topics, code, err := c.mal.GetForumTopicsWithContext(ctx, mal.GetForumTopicsParam{
		BoardID:       param.BoardID,
		SubboardID:    param.SubboardID,
		Query:         param.Query,
		TopicUsername: param.TopicUsername,
		Username:      param.Username,
		Sort:          string(param.Sort),
		Limit:         param.Limit,
		Offset:        param.Offset,
	})
	if err != nil {
		return nil, code, err
	}

	return c.forumTopicsToForumTopic(topics.Data), http.StatusOK, nil
}

// GetForumTopicDetails to get forum topic details.
func (c *Client) GetForumTopicDetails(param GetForumTopicDetailsParam) (*ForumTopicDetail, int, error) {
	return c.GetForumTopicDetailsWithContext(context.Background(), param)
}

// GetForumTopicDetailsWithContext to get forum topic details with context.
func (c *Client) GetForumTopicDetailsWithContext(ctx context.Context, param GetForumTopicDetailsParam) (*ForumTopicDetail, int, error) {
	if err := c.validate(&param); err != nil {
		return nil, http.StatusBadRequest, err
	}

	details, code, err := c.mal.GetForumTopicDetailsWithContext(ctx, mal.GetForumTopicDetailsParam{
		ID:     param.ID,
		Limit:  param.Limit,
		Offset: param.Offset,
	})
	if err != nil {
		return nil, code, err
	}

	return c.forumTopicDetailsToForumTopicDetails(details.Data), http.StatusOK, nil
}
