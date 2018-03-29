package helix

import (
	"strings"
	"time"
)

// Time is our custom time struct.
type Time struct {
	time.Time
}

// UnmarshalJSON is our custom datetime unmarshaller. Twitch sometimes
// returns datetimes as empty strings, which casuses issues with the native time
// UnmarshalJSON method when decoding the JSON string. Here we hanlde that scenario,
// by returning a zero time value for any JSON time field that is either an
// empty string or "null".
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	timeStr := strings.Trim(string(b), "\"")

	if timeStr == "" || timeStr == "null" {
		t.Time = time.Time{}
		return
	}

	t.Time, err = time.Parse(time.RFC3339, timeStr)

	return
}
