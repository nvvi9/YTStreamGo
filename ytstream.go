package ytstream

import (
	"encoding/json"
	"fmt"
	"github.com/nvvi9/YTStreamGo/extractor"
	"github.com/nvvi9/YTStreamGo/model"
	"github.com/nvvi9/YTStreamGo/model/youtube"
	"github.com/nvvi9/YTStreamGo/network"
	"github.com/nvvi9/YTStreamGo/utils"
)

func ExtractVideoData(videoId string) (model.VideoData, error) {
	pageHtml, err := network.GetVideoPage(videoId)
	if err != nil {
		return model.VideoData{}, err
	}

	match := utils.PatternPlayerResponse.FindStringSubmatch(pageHtml)
	if match == nil {
		return model.VideoData{}, fmt.Errorf("error parsing player response")
	}

	playerResponse := match[1]

	var initialPlayerResponse = youtube.InitialPlayerResponse{}

	if err := json.Unmarshal([]byte(playerResponse), &initialPlayerResponse); err != nil {
		return model.VideoData{}, err
	}

	videoDetails := model.VideoDetailsFromInitialPlayerResponse(initialPlayerResponse)

	streams := extractor.ExtractStreams(pageHtml, initialPlayerResponse.StreamingData)

	return model.VideoData{
		VideoDetails: videoDetails,
		Streams:      streams,
	}, nil
}

func ExtractVideoDetails(videoId string) (model.VideoDetails, error) {
	pageHtml, err := network.GetVideoPage(videoId)
	if err != nil {
		return model.VideoDetails{}, err
	}

	match := utils.PatternPlayerResponse.FindStringSubmatch(pageHtml)
	if match == nil {
		return model.VideoDetails{}, fmt.Errorf("error parsing player response")
	}

	playerResponse := match[1]

	var initialPlayerResponse = youtube.InitialPlayerResponse{}

	if err := json.Unmarshal([]byte(playerResponse), &initialPlayerResponse); err != nil {
		return model.VideoDetails{}, err
	}

	return model.VideoDetailsFromInitialPlayerResponse(initialPlayerResponse), nil
}
