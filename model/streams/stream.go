package streams

import (
	"github.com/nvvi9/YTStreamGo/model/codecs"
)

type Stream struct {
	Url              string
	StreamDetails    StreamDetails
	ExpiresInSeconds int64
}

type StreamDetails struct {
	Itag       int
	Type       StreamType
	Extension  Extension
	AudioCodec codecs.AudioCodec
	VideoCodec codecs.VideoCodec
	Quality    int
	Bitrate    int
	Fps        int
}
