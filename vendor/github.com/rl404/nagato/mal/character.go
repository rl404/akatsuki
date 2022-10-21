package mal

import (
	"context"
	"net/http"
)

// GetCharacterDetails to get character details.
//
// Undocumented.
func (c *Client) GetCharacterDetails(id int, fields ...string) (*Character, int, error) {
	return c.GetCharacterDetailsWithContext(context.Background(), id, fields...)
}

// GetCharacterDetailsWithContext to get character details with context.
//
// Undocumented.
func (c *Client) GetCharacterDetailsWithContext(ctx context.Context, id int, fields ...string) (*Character, int, error) {
	url := c.generateURL(map[string]interface{}{
		"fields": fields,
	}, "characters", id)

	var character Character
	if code, err := c.get(ctx, url, &character); err != nil {
		return nil, code, err
	}

	return &character, http.StatusOK, nil
}
