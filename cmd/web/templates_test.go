package main

import (
	"testing"
	"time"

	"github.com/MtLampert/alex_e_-_lets_go/internal/assert"
)

type humanDateTest struct {
	name      string
	timestamp time.Time
	expected  string
}

func TestHumanDate(t *testing.T) {
	// test cases as 'table'
	testCases := []humanDateTest{
		{
			name:      `UTC`,
			timestamp: time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			expected:  `2022-03-17 10:15`,
		},
		{
			name:      `Empty`,
			timestamp: time.Time{},
			expected:  ``,
		},
		{
			name: `CET`,
			// '1*60*60' means UTC plus 1*60*60 seconds
			// time.Date() describes the time locally to CET (Berlin),
			// so in UTC it's the expected time, 1 hour before
			timestamp: time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone(`CET`, 1*60*60)),
			expected:  `2022-03-17 09:15`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hd := humanDate(tc.timestamp)
			assert.Equal(t, hd, tc.expected)
		})
	}
}

// vim: ts=4 sw=4 fdm=indent
