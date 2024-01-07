package youtube

type InitialPlayerResponse struct {
	StreamingData StreamingData `json:"streamingData"`
	VideoDetails  VideoDetails  `json:"videoDetails"`
}

type StreamingData struct {
	ExpiresInSeconds string   `json:"expiresInSeconds"`
	Formats          []Format `json:"formats"`
	AdaptiveFormats  []Format `json:"adaptiveFormats"`
}

type Format struct {
	Itag            int     `json:"itag"`
	SignatureCipher *string `json:"signatureCipher"`
	Type            string  `json:"type"`
	Url             *string `json:"url"`
}

type VideoDetails struct {
	VideoId          string     `json:"videoId"`
	Title            string     `json:"title"`
	LengthSeconds    string     `json:"lengthSeconds"`
	ChannelId        string     `json:"channelId"`
	ShortDescription string     `json:"shortDescription"`
	Thumbnail        Thumbnails `json:"thumbnail"`
	ViewCount        string     `json:"viewCount"`
	Author           string     `json:"author"`
	IsLiveContent    bool       `json:"isLiveContent"`
}

type Thumbnails struct {
	Thumbnails []Thumbnail `json:"thumbnails"`
}

type Thumbnail struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
