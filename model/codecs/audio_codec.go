package codecs

type AudioCodec int

const (
	MP3 AudioCodec = iota + 1
	AAC
	Vorbis
	Opus
)
