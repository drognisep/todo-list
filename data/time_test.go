package data

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTimeEntry_Deserialization(t *testing.T) {
	data := `{
  "taskID": 64,
  "start": "2022-12-24T19:49:16.4883081Z",
  "end": "2022-12-24T19:51:16.4883081Z"
}`
	var entry TimeEntry
	assert.NoError(t, json.Unmarshal([]byte(data), &entry))
	assert.Equal(t, uint64(64), entry.TaskID)
	assert.Greater(t, time.Now(), entry.Start)
	assert.Greater(t, entry.Start, time.Time{})
	assert.Equal(t, 2*time.Minute, entry.Duration())
}

func TestTimeEntry_Duration_Neg(t *testing.T) {
	data := `{
  "taskID": 64,
  "start": "2022-12-24T19:49:16.4883081Z"
}`
	var entry TimeEntry
	require.NoError(t, json.Unmarshal([]byte(data), &entry))
	assert.Equal(t, time.Duration(0), entry.Duration())
}

func TestLastWeekday(t *testing.T) {
	day := 24 * time.Hour
	fuzzSet := []time.Time{
		time.Now(),
		time.Now().Add(1 * day),
		time.Now().Add(2 * day),
		time.Now().Add(3 * day),
		time.Now().Add(4 * day),
		time.Now().Add(5 * day),
		time.Now().Add(6 * day),
	}

	tests := map[string]struct {
		weekday time.Weekday
	}{
		"Last Sunday":    {time.Sunday},
		"Last Monday":    {time.Monday},
		"Last Tuesday":   {time.Tuesday},
		"Last Wednesday": {time.Wednesday},
		"Last Thursday":  {time.Thursday},
		"Last Friday":    {time.Friday},
		"Last Saturday":  {time.Saturday},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			for _, input := range fuzzSet {
				t.Run(fmt.Sprintf("%s to %s", input.Weekday().String(), name), func(t *testing.T) {
					result := lastWeekday(tc.weekday, input)
					assert.Equal(t, tc.weekday.String(), result.Weekday().String())
				})
			}
		})
	}
}
