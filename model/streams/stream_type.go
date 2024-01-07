package streams

type StreamType string

const (
	Video       StreamType = "video"
	Audio       StreamType = "audio"
	Live        StreamType = "live"
	Multiplexed StreamType = "multiplexed"
)
