// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated from the elasticsearch-specification DO NOT EDIT.
// https://github.com/elastic/elasticsearch-specification/tree/a0da620389f06553c0727f98f95e40dbb564fcca

package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strconv"
)

// RefreshStats type.
//
// https://github.com/elastic/elasticsearch-specification/blob/a0da620389f06553c0727f98f95e40dbb564fcca/specification/_types/Stats.ts#L168-L175
type RefreshStats struct {
	ExternalTotal             int64    `json:"external_total"`
	ExternalTotalTimeInMillis int64    `json:"external_total_time_in_millis"`
	Listeners                 int64    `json:"listeners"`
	Total                     int64    `json:"total"`
	TotalTime                 Duration `json:"total_time,omitempty"`
	TotalTimeInMillis         int64    `json:"total_time_in_millis"`
}

func (s *RefreshStats) UnmarshalJSON(data []byte) error {

	dec := json.NewDecoder(bytes.NewReader(data))

	for {
		t, err := dec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		switch t {

		case "external_total":
			var tmp interface{}
			dec.Decode(&tmp)
			switch v := tmp.(type) {
			case string:
				value, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
				s.ExternalTotal = value
			case float64:
				f := int64(v)
				s.ExternalTotal = f
			}

		case "external_total_time_in_millis":
			if err := dec.Decode(&s.ExternalTotalTimeInMillis); err != nil {
				return err
			}

		case "listeners":
			var tmp interface{}
			dec.Decode(&tmp)
			switch v := tmp.(type) {
			case string:
				value, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
				s.Listeners = value
			case float64:
				f := int64(v)
				s.Listeners = f
			}

		case "total":
			var tmp interface{}
			dec.Decode(&tmp)
			switch v := tmp.(type) {
			case string:
				value, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
				s.Total = value
			case float64:
				f := int64(v)
				s.Total = f
			}

		case "total_time":
			if err := dec.Decode(&s.TotalTime); err != nil {
				return err
			}

		case "total_time_in_millis":
			if err := dec.Decode(&s.TotalTimeInMillis); err != nil {
				return err
			}

		}
	}
	return nil
}

// NewRefreshStats returns a RefreshStats.
func NewRefreshStats() *RefreshStats {
	r := &RefreshStats{}

	return r
}
