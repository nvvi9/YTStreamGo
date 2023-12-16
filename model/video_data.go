package model

import (
	"github.com/nvvi9/YTStreamGo/model/streams"
	"github.com/nvvi9/YTStreamGo/model/youtube"
	"github.com/nvvi9/YTStreamGo/utils"
	"strconv"
)

type VideoData struct {
	VideoDetails VideoDetails
	Streams      []streams.Stream
}

type VideoDetails struct {
	Id              string
	Title           string
	Channel         string
	ChannelId       string
	Description     string
	DurationSeconds int64
	ViewCount       int64
	Thumbnails      []Thumbnail
	IsLiveStream    bool
}

type Thumbnail struct {
	Width  int
	Height int
	Url    string
}

func VideoDetailsFromInitialPlayerResponse(response youtube.InitialPlayerResponse) VideoDetails {
	durationSeconds, err := strconv.ParseInt(response.VideoDetails.LengthSeconds, 10, 64)
	if err != nil {
		durationSeconds = 0
	}

	viewCount, err := strconv.ParseInt(response.VideoDetails.ViewCount, 10, 64)
	if err != nil {
		viewCount = 0
	}

	thumbnails := utils.Map(response.VideoDetails.Thumbnail.Thumbnails, func(t youtube.Thumbnail) Thumbnail {
		return Thumbnail{
			Width:  t.Width,
			Height: t.Height,
			Url:    t.Url,
		}
	})

	return VideoDetails{
		Id:              response.VideoDetails.VideoId,
		Title:           response.VideoDetails.Title,
		Channel:         response.VideoDetails.Author,
		ChannelId:       response.VideoDetails.ChannelId,
		Description:     response.VideoDetails.ShortDescription,
		DurationSeconds: durationSeconds,
		ViewCount:       viewCount,
		Thumbnails:      thumbnails,
		IsLiveStream:    false,
	}
}
