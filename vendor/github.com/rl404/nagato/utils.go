package nagato

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (c *Client) dateToDate(date string) Date {
	d := Date{}

	if date == "" {
		return d
	}

	split := strings.Split(date, "-")

	if len(split) >= 1 {
		d.Year, _ = strconv.Atoi(split[0])
	}

	if len(split) >= 2 {
		d.Month, _ = strconv.Atoi(split[1])
	}

	if len(split) >= 3 {
		d.Day, _ = strconv.Atoi(split[2])
	}

	return d
}

func (c *Client) dateToStr(date Date) string {
	var d []string
	if date.Year > 0 {
		d = append(d, fmt.Sprintf("%d", date.Year))
	}

	if date.Month > 0 {
		d = append(d, fmt.Sprintf("%02d", date.Month))
	}

	if date.Day > 0 {
		d = append(d, fmt.Sprintf("%02d", date.Day))
	}

	return strings.Join(d, "-")
}

func (c *Client) validateDate(date Date) error {
	if date.Day > 0 {
		if date.Month == 0 || date.Year == 0 {
			return ErrInvalidDate
		}

		d := fmt.Sprintf("%02d-%02d-%04d", date.Day, date.Month, date.Year)
		if _, err := time.Parse("02-01-2006", d); err != nil {
			return ErrInvalidDate
		}

		return nil
	}

	if date.Month > 0 {
		if date.Year == 0 {
			return ErrInvalidDate
		}
		return nil
	}

	return nil
}

func (c *Client) initValidator() {
	_ = c.validator.RegisterValidatorError("required", c.valErrRequired)
	_ = c.validator.RegisterValidatorError("gt", c.valErrGT)
	_ = c.validator.RegisterValidatorError("gte", c.valErrGTE)
	_ = c.validator.RegisterValidatorError("lte", c.valErrLTE)
	_ = c.validator.RegisterValidatorError("oneof", c.valErrOneOf)
}

func (c *Client) validate(data interface{}) error {
	return c.validator.Validate(data)
}

func (c *Client) valErrRequired(f string, param ...string) error {
	return c.errRequiredField(f)
}

func (c *Client) valErrGT(f string, param ...string) error {
	return c.errGTField(f, param[0])
}

func (c *Client) valErrGTE(f string, param ...string) error {
	return c.errGTEField(f, param[0])
}

func (c *Client) valErrLTE(f string, param ...string) error {
	return c.errLTEField(f, param[0])
}

func (c *Client) valErrOneOf(f string, param ...string) error {
	return c.errOneOfField(f, param[0])
}

func (p *GetForumTopicsParam) validate() error {
	if p.BoardID == 0 && p.SubboardID == 0 && p.Query == "" {
		return ErrForumTopicMissingQuery
	}
	return nil
}
