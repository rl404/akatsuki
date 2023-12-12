package mal

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (c *Client) handleError(body []byte) error {
	var e errorResponse
	if err := json.Unmarshal(body, &e); err != nil {
		return errors.New(string(body))
	}

	msg := e.Error
	if e.Message != "" {
		msg = e.Message
	}

	return errors.New(msg)
}

func (c *Client) get(ctx context.Context, url string, model interface{}) (int, error) {
	body, code, err := c.MakeRequest(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return code, err
	}

	if code != http.StatusOK {
		return code, c.handleError(body)
	}

	if err := json.Unmarshal(body, &model); err != nil {
		return http.StatusInternalServerError, err
	}

	return code, nil
}

func (c *Client) patch(ctx context.Context, url string, payload []byte, model interface{}) (int, error) {
	body, code, err := c.MakeRequest(ctx, http.MethodPatch, url, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}, payload)
	if err != nil {
		return code, err
	}

	if code != http.StatusOK {
		return code, c.handleError(body)
	}

	if err := json.Unmarshal(body, &model); err != nil {
		return http.StatusInternalServerError, err
	}

	return code, nil
}

func (c *Client) delete(ctx context.Context, url string) (int, error) {
	body, code, err := c.MakeRequest(ctx, http.MethodDelete, url, nil, nil)
	if err != nil {
		return code, err
	}

	if code != http.StatusOK {
		return code, c.handleError(body)
	}

	return code, nil
}

// MakeRequest to make http request.
func (c *Client) MakeRequest(ctx context.Context, method, url string, headers map[string]string, payload []byte) ([]byte, int, error) {
	c.Limiter.Take()

	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(string(payload)))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.Http.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return body, resp.StatusCode, nil
}
