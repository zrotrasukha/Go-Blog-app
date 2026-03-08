package main

import (
	"testing"
	"time"

	"github.com/zrotrasukha/Go-Blog-app/internal/assert"
)

func TestHumanDate(t *testing.T) {
	test := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2026, 3, 6, 3, 42, 0, 0, time.UTC),
			want: "06 Mar 2026 at 03:42",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "IST",
			tm:   time.Date(2026, 3, 6, 3, 0, 0, 0, time.FixedZone("IST", 5*60*60+30*60)),
			want: "05 Mar 2026 at 21:30",
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}
