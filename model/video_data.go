package model

import (
	"github.com/nvvi9/YTStreamGo/model/streams"
	"github.com/nvvi9/YTStreamGo/model/youtube"
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

func VideoDetailsFromInitialPlayerResponse(response *youtube.InitialPlayerResponse) *VideoDetails {
	durationSeconds, err := strconv.ParseInt(response.VideoDetails.LengthSeconds, 10, 64)
	if err != nil {
		durationSeconds = 0
	}

	viewCount, err := strconv.ParseInt(response.VideoDetails.ViewCount, 10, 64)
	if err != nil {
		viewCount = 0
	}

	thumbnails := make([]Thumbnail, 0, len(response.VideoDetails.Thumbnail.Thumbnails))
	for _, t := range response.VideoDetails.Thumbnail.Thumbnails {
		thumbnail := Thumbnail{
			Height: t.Height,
			Width:  t.Width,
			Url:    t.Url,
		}
		thumbnails = append(thumbnails, thumbnail)
	}

	videoDetails := new(VideoDetails)
	videoDetails.Id = response.VideoDetails.VideoId
	videoDetails.Title = response.VideoDetails.Title
	videoDetails.Channel = response.VideoDetails.Author
	videoDetails.ChannelId = response.VideoDetails.ChannelId
	videoDetails.Description = response.VideoDetails.ShortDescription
	videoDetails.DurationSeconds = durationSeconds
	videoDetails.ViewCount = viewCount
	videoDetails.Thumbnails = thumbnails
	videoDetails.IsLiveStream = false

	return videoDetails
}
