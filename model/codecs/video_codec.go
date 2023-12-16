package codecs

type VideoCodec int

const (
	H263 VideoCodec = iota + 1
	H264
	MPEG4
	VP8
	VP9
)
