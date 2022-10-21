package mal

import (
	"context"
	"net/http"
)

// GetUserInfo to get user info.
//
// Only `@me` in username param works for now.
//
// Need oauth2.
func (c *Client) GetUserInfo(username string, fields ...string) (*User, int, error) {
	return c.GetUserInfoWithContext(context.Background(), username, fields...)
}

// GetUserInfoWithContext to get user info with context.
//
// Only `@me` in username param works for now.
//
// Need oauth2.
func (c *Client) GetUserInfoWithContext(ctx context.Context, username string, fields ...string) (*User, int, error) {
	url := c.generateURL(map[string]interface{}{
		"fields": fields,
	}, "users", username)

	var user User
	if code, err := c.get(ctx, url, &user); err != nil {
		return nil, code, err
	}

	return &user, http.StatusOK, nil
}
