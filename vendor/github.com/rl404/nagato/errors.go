package nagato

import (
	"errors"
	"fmt"
	"strings"
)

// List of errors.
var (
	ErrInvalidDate            = errors.New("invalid date")
	ErrForumTopicMissingQuery = errors.New("at least one of BoardID, SubboardID, or Query is required")
)

func (c *Client) errRequiredField(str string) error {
	return fmt.Errorf("required %s", str)
}

func (c *Client) errGTField(str, value string) error {
	return fmt.Errorf("%s must be greater than %s", str, value)
}

func (c *Client) errGTEField(str, value string) error {
	return fmt.Errorf("%s must be greater than or equal %s", str, value)
}

func (c *Client) errLTEField(str, value string) error {
	return fmt.Errorf("%s must be lower than or equal %s", str, value)
}

func (c *Client) errOneOfField(str, value string) error {
	return fmt.Errorf("%s must be one of %s", str, strings.Join(strings.Split(value, " "), ","))
}
