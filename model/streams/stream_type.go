package streams

type StreamType int

const (
	Video StreamType = iota + 1
	Audio
	Live
	Multiplexed
)
