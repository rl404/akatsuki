package nagato

import (
	"context"
	"net/http"
)

// GetUserInfo to get user info.
//
// Only `@me` in username param works for now.
//
// Need oauth2.
func (c *Client) GetUserInfo(username string, fields ...UserField) (*User, int, error) {
	return c.GetUserInfoWithContext(context.Background(), username, fields...)
}

// GetUserInfoWithContext to get user info with context.
//
// Only `@me` in username param works for now.
//
// Need oauth2.
func (c *Client) GetUserInfoWithContext(ctx context.Context, username string, fields ...UserField) (*User, int, error) {
	u := usernameParam{Username: username}

	if err := c.validate(&u); err != nil {
		return nil, http.StatusBadRequest, err
	}

	user, code, err := c.mal.GetUserInfoWithContext(ctx, u.Username, c.userFieldsToStrs(fields...)...)
	if err != nil {
		return nil, code, err
	}

	return c.userToUser(user), http.StatusOK, nil
}
