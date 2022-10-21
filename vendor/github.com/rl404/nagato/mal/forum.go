package mal

import (
	"context"
	"net/http"
)

// GetForumBoards to get forum board list.
func (c *Client) GetForumBoards() (*ForumBoardCategories, int, error) {
	return c.GetForumBoardsWithContext(context.Background())
}

// GetForumBoardsWithContext to get forum board list with context.
func (c *Client) GetForumBoardsWithContext(ctx context.Context) (*ForumBoardCategories, int, error) {
	url := c.generateURL(nil, "forum", "boards")

	var forumBoardCategories ForumBoardCategories
	if code, err := c.get(ctx, url, &forumBoardCategories); err != nil {
		return nil, code, err
	}

	return &forumBoardCategories, http.StatusOK, nil
}

// GetForumTopics to get forum topic list.
func (c *Client) GetForumTopics(param GetForumTopicsParam) (*ForumTopicPaging, int, error) {
	return c.GetForumTopicsWithContext(context.Background(), param)
}

// GetForumTopicsWithContext to get forum topic list with context.
func (c *Client) GetForumTopicsWithContext(ctx context.Context, param GetForumTopicsParam) (*ForumTopicPaging, int, error) {
	url := c.generateURL(map[string]interface{}{
		"board_id":        param.BoardID,
		"subboard_id":     param.SubboardID,
		"limit":           param.Limit,
		"offset":          param.Offset,
		"sort":            param.Sort,
		"q":               param.Query,
		"topic_user_name": param.TopicUsername,
		"user_name":       param.Username,
	}, "forum", "topics")

	var forumTopic ForumTopicPaging
	if code, err := c.get(ctx, url, &forumTopic); err != nil {
		return nil, code, err
	}

	return &forumTopic, http.StatusOK, nil
}

// GetForumTopicDetails to get forum topic details.
func (c *Client) GetForumTopicDetails(param GetForumTopicDetailsParam) (*ForumTopicDetailPaging, int, error) {
	return c.GetForumTopicDetailsWithContext(context.Background(), param)
}

// GetForumTopicDetailsWithContext to get forum topic details with context.
func (c *Client) GetForumTopicDetailsWithContext(ctx context.Context, param GetForumTopicDetailsParam) (*ForumTopicDetailPaging, int, error) {
	url := c.generateURL(map[string]interface{}{
		"limit":  param.Limit,
		"offset": param.Offset,
	}, "forum", "topic", param.ID)

	var forumTopic ForumTopicDetailPaging
	if code, err := c.get(ctx, url, &forumTopic); err != nil {
		return nil, code, err
	}

	return &forumTopic, http.StatusOK, nil
}
