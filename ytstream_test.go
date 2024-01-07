package ytstream

import (
	"testing"
)

var videoIds = []string{
	"dz6YbMtj74U",
	"UqLRqzTp6Rk",
	"u0BetD0OAcs",
	"uKM9ZuQB3MA",
	"1nX0kF2UwDc",
	"kfugSz3m_zA",
	"UqLRqzTp6Rk",
	"u0BetD0OAcs",
	"uKM9ZuQB3MA",
	"1nX0kF2UwDc",
	"kfugSz3m_zA",
}

func TestVideoDataExtraction(t *testing.T) {
	for _, id := range videoIds {
		t.Run("Test VideoData extraction: "+id, func(t *testing.T) {
			videoData, err := ExtractVideoData(id)
			if err != nil {
				t.Error(err)
			}
			if len(videoData.Streams) == 0 {
				t.Error("Empty streams")
			}
		})
	}
}
