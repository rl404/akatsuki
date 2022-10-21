package mal

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func (c *Client) generateURL(queries map[string]interface{}, paths ...interface{}) string {
	strPaths := []string{c.Host}
	for _, p := range paths {
		strPaths = append(strPaths, fmt.Sprintf("%v", p))
	}
	url, _ := url.Parse(strings.Join(strPaths, "/"))

	q := url.Query()
	for k, v := range queries {
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			if reflect.ValueOf(v).String() != "" {
				q.Set(k, v.(string))
			}
		case reflect.Int:
			if reflect.ValueOf(v).Int() != 0 {
				q.Set(k, fmt.Sprintf("%d", v.(int)))
			}
		case reflect.Bool:
			if reflect.ValueOf(v).Bool() {
				q.Set(k, "true")
			}
		case reflect.Slice, reflect.Array:
			vv := reflect.ValueOf(v)
			if vv.Len() > 0 {
				fields := make([]string, vv.Len())
				for i := 0; i < vv.Len(); i++ {
					fields[i] = vv.Index(i).String()
				}
				q.Set(k, strings.Join(fields, ","))
			}
		}
	}
	url.RawQuery = q.Encode()

	return url.String()
}
