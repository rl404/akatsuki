package mal

import (
	"context"
	"net/http"
)

// GetPeopleDetails to get people details.
//
// Undocumented.
func (c *Client) GetPeopleDetails(id int, fields ...string) (*People, int, error) {
	return c.GetPeopleDetailsWithContext(context.Background(), id, fields...)
}

// GetPeopleDetailsWithContext to get people details with context.
//
// Undocumented.
func (c *Client) GetPeopleDetailsWithContext(ctx context.Context, id int, fields ...string) (*People, int, error) {
	url := c.generateURL(map[string]interface{}{
		"fields": fields,
	}, "people", id)

	var people People
	if code, err := c.get(ctx, url, &people); err != nil {
		return nil, code, err
	}

	return &people, http.StatusOK, nil
}
