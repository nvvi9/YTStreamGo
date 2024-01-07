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

func ExtractVideoData(videoId string) (*model.VideoData, error) {
	pageHtml, err := network.GetVideoPage(videoId)
	if err != nil {
		return nil, err
	}

	match := utils.PatternPlayerResponse.FindStringSubmatch(pageHtml)
	if match == nil {
		return nil, fmt.Errorf("error parsing player response")
	}

	playerResponse := match[1]

	initialPlayerResponse := new(youtube.InitialPlayerResponse)

	if err := json.Unmarshal([]byte(playerResponse), initialPlayerResponse); err != nil {
		return nil, err
	}

	videoDetails := model.VideoDetailsFromInitialPlayerResponse(initialPlayerResponse)

	streams := extractor.ExtractStreams(&pageHtml, &initialPlayerResponse.StreamingData)

	videoData := new(model.VideoData)
	videoData.VideoDetails = *videoDetails
	videoData.Streams = streams
	return videoData, nil
}

func ExtractVideoDetails(videoId string) (*model.VideoDetails, error) {
	pageHtml, err := network.GetVideoPage(videoId)
	if err != nil {
		return nil, err
	}

	match := utils.PatternPlayerResponse.FindStringSubmatch(pageHtml)
	if match == nil {
		return nil, fmt.Errorf("error parsing player response")
	}

	playerResponse := match[1]

	initialPlayerResponse := new(youtube.InitialPlayerResponse)

	if err := json.Unmarshal([]byte(playerResponse), &initialPlayerResponse); err != nil {
		return nil, err
	}

	return model.VideoDetailsFromInitialPlayerResponse(initialPlayerResponse), nil
}
