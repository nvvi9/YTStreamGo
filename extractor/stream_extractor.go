package extractor

import (
	"bytes"
	"fmt"
	"github.com/nvvi9/YTStreamGo/js"
	"github.com/nvvi9/YTStreamGo/model/codecs"
	"github.com/nvvi9/YTStreamGo/model/streams"
	"github.com/nvvi9/YTStreamGo/model/youtube"
	"github.com/nvvi9/YTStreamGo/network"
	"github.com/nvvi9/YTStreamGo/utils"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var itags = map[int]streams.StreamDetails{
	140: {Itag: 140, Type: streams.Audio, Extension: streams.M4A, AudioCodec: codecs.AAC, VideoCodec: "", Quality: 0, Bitrate: 128, Fps: 0},
	141: {Itag: 141, Type: streams.Audio, Extension: streams.M4A, AudioCodec: codecs.AAC, VideoCodec: "", Quality: 0, Bitrate: 256, Fps: 0},
	256: {Itag: 256, Type: streams.Audio, Extension: streams.M4A, AudioCodec: codecs.AAC, VideoCodec: "", Quality: 0, Bitrate: 192, Fps: 0},
	258: {Itag: 258, Type: streams.Audio, Extension: streams.M4A, AudioCodec: codecs.AAC, VideoCodec: "", Quality: 0, Bitrate: 384, Fps: 0},
	171: {Itag: 171, Type: streams.Audio, Extension: streams.WebM, AudioCodec: codecs.Vorbis, VideoCodec: "", Quality: 0, Bitrate: 128, Fps: 0},
	249: {Itag: 249, Type: streams.Audio, Extension: streams.WebM, AudioCodec: codecs.Opus, VideoCodec: "", Quality: 0, Bitrate: 48, Fps: 0},
	250: {Itag: 250, Type: streams.Audio, Extension: streams.WebM, AudioCodec: codecs.Opus, VideoCodec: "", Quality: 0, Bitrate: 64, Fps: 0},
	251: {Itag: 251, Type: streams.Audio, Extension: streams.WebM, AudioCodec: codecs.Opus, VideoCodec: "", Quality: 0, Bitrate: 160, Fps: 0},

	160: {Itag: 160, Type: streams.Video, Extension: streams.MP4, AudioCodec: "", VideoCodec: codecs.H264, Quality: 144, Bitrate: 0, Fps: 0},
	133: {Itag: 133, Type: streams.Video, Extension: streams.MP4, AudioCodec: "", VideoCodec: codecs.H264, Quality: 240, Bitrate: 0, Fps: 0},
	134: {Itag: 134, Type: streams.Video, Extension: streams.MP4, AudioCodec: "", VideoCodec: codecs.H264, Quality: 360, Bitrate: 0, Fps: 0},
	135: {Itag: 135, Type: streams.Video, Extension: streams.MP4, AudioCodec: "", VideoCodec: codecs.H264, Quality: 480, Bitrate: 0, Fps: 0},
	136: {Itag: 136, Type: streams.Video, Extension: streams.MP4, AudioCodec: "", VideoCodec: codecs.H264, Quality: 720, Bitrate: 0, Fps: 0},
	137: {Itag: 137, Type: streams.Video, Extension: streams.MP4, AudioCodec: "", VideoCodec: codecs.H264, Quality: 1080, Bitrate: 0, Fps: 0},
	264: {Itag: 264, Type: streams.Video, Extension: streams.MP4, AudioCodec: "", VideoCodec: codecs.H264, Quality: 1440, Bitrate: 0, Fps: 0},
	266: {Itag: 266, Type: streams.Video, Extension: streams.MP4, AudioCodec: "", VideoCodec: codecs.H264, Quality: 2160, Bitrate: 0, Fps: 0},
	298: {Itag: 298, Type: streams.Video, Extension: streams.MP4, AudioCodec: "", VideoCodec: codecs.H264, Quality: 720, Bitrate: 0, Fps: 60},
	299: {Itag: 299, Type: streams.Video, Extension: streams.MP4, AudioCodec: "", VideoCodec: codecs.H264, Quality: 1080, Bitrate: 0, Fps: 60},
	278: {Itag: 278, Type: streams.Video, Extension: streams.WebM, AudioCodec: "", VideoCodec: codecs.VP9, Quality: 144, Bitrate: 0, Fps: 0},
	242: {Itag: 242, Type: streams.Video, Extension: streams.WebM, AudioCodec: "", VideoCodec: codecs.VP9, Quality: 240, Bitrate: 0, Fps: 0},
	243: {Itag: 243, Type: streams.Video, Extension: streams.WebM, AudioCodec: "", VideoCodec: codecs.VP9, Quality: 360, Bitrate: 0, Fps: 0},
	244: {Itag: 244, Type: streams.Video, Extension: streams.WebM, AudioCodec: "", VideoCodec: codecs.VP9, Quality: 480, Bitrate: 0, Fps: 0},
	247: {Itag: 247, Type: streams.Video, Extension: streams.WebM, AudioCodec: "", VideoCodec: codecs.VP9, Quality: 720, Bitrate: 0, Fps: 0},
	248: {Itag: 248, Type: streams.Video, Extension: streams.WebM, AudioCodec: "", VideoCodec: codecs.VP9, Quality: 1080, Bitrate: 0, Fps: 0},
	271: {Itag: 271, Type: streams.Video, Extension: streams.WebM, AudioCodec: "", VideoCodec: codecs.VP9, Quality: 1440, Bitrate: 0, Fps: 0},
	313: {Itag: 313, Type: streams.Video, Extension: streams.WebM, AudioCodec: "", VideoCodec: codecs.VP9, Quality: 2160, Bitrate: 0, Fps: 0},
	302: {Itag: 302, Type: streams.Video, Extension: streams.WebM, AudioCodec: "", VideoCodec: codecs.VP9, Quality: 720, Bitrate: 0, Fps: 60},
	308: {Itag: 308, Type: streams.Video, Extension: streams.WebM, AudioCodec: "", VideoCodec: codecs.VP9, Quality: 1440, Bitrate: 0, Fps: 60},
	303: {Itag: 303, Type: streams.Video, Extension: streams.WebM, AudioCodec: "", VideoCodec: codecs.VP9, Quality: 1080, Bitrate: 0, Fps: 60},
	315: {Itag: 315, Type: streams.Video, Extension: streams.WebM, AudioCodec: "", VideoCodec: codecs.VP9, Quality: 2160, Bitrate: 0, Fps: 60},

	17: {Itag: 17, Type: streams.Multiplexed, Extension: streams.ThreeGP, AudioCodec: codecs.AAC, VideoCodec: codecs.MPEG4, Quality: 144, Bitrate: 24, Fps: 0},
	36: {Itag: 36, Type: streams.Multiplexed, Extension: streams.ThreeGP, AudioCodec: codecs.AAC, VideoCodec: codecs.MPEG4, Quality: 240, Bitrate: 32, Fps: 0},
	5:  {Itag: 5, Type: streams.Multiplexed, Extension: streams.FLV, AudioCodec: codecs.MP3, VideoCodec: codecs.H263, Quality: 144, Bitrate: 64, Fps: 0},
	43: {Itag: 43, Type: streams.Multiplexed, Extension: streams.WebM, AudioCodec: codecs.Vorbis, VideoCodec: codecs.VP8, Quality: 360, Bitrate: 128, Fps: 0},
	18: {Itag: 18, Type: streams.Multiplexed, Extension: streams.MP4, AudioCodec: codecs.AAC, VideoCodec: codecs.H264, Quality: 360, Bitrate: 96, Fps: 0},
	22: {Itag: 22, Type: streams.Multiplexed, Extension: streams.MP4, AudioCodec: codecs.AAC, VideoCodec: codecs.H264, Quality: 720, Bitrate: 192, Fps: 0},

	91: {Itag: 91, Type: streams.Live, Extension: streams.MP4, AudioCodec: codecs.AAC, VideoCodec: codecs.H264, Quality: 144, Bitrate: 48, Fps: 0},
	92: {Itag: 92, Type: streams.Live, Extension: streams.MP4, AudioCodec: codecs.AAC, VideoCodec: codecs.H264, Quality: 240, Bitrate: 48, Fps: 0},
	93: {Itag: 93, Type: streams.Live, Extension: streams.MP4, AudioCodec: codecs.AAC, VideoCodec: codecs.H264, Quality: 360, Bitrate: 128, Fps: 0},
	94: {Itag: 94, Type: streams.Live, Extension: streams.MP4, AudioCodec: codecs.AAC, VideoCodec: codecs.H264, Quality: 480, Bitrate: 128, Fps: 0},
	95: {Itag: 95, Type: streams.Live, Extension: streams.MP4, AudioCodec: codecs.AAC, VideoCodec: codecs.H264, Quality: 720, Bitrate: 256, Fps: 0},
	96: {Itag: 96, Type: streams.Live, Extension: streams.MP4, AudioCodec: codecs.AAC, VideoCodec: codecs.H264, Quality: 1080, Bitrate: 256, Fps: 0},
}

type encodedSignature struct {
	Url           string
	Signature     string
	StreamDetails streams.StreamDetails
}

func ExtractStreams(pageHtml *string, streamingData *youtube.StreamingData) []streams.Stream {
	var streams []streams.Stream

	if notEncodedStreams, _ := extractNotEncodedStreams(streamingData); notEncodedStreams != nil {
		streams = append(streams, notEncodedStreams...)
	}

	if encodedStreams, _ := extractEncodedStreams(pageHtml, streamingData); encodedStreams != nil {
		streams = append(streams, encodedStreams...)
	}

	return streams
}

func extractNotEncodedStreams(streamingData *youtube.StreamingData) ([]streams.Stream, error) {
	expiresInSeconds, err := strconv.ParseInt(streamingData.ExpiresInSeconds, 10, 64)
	if err != nil {
		return nil, err
	}

	var result []streams.Stream
	formats := append(streamingData.Formats, streamingData.AdaptiveFormats...)

	for _, format := range formats {
		if format.Type != "FORMAT_STREAM_TYPE_OTF" {
			if formatStreamUrl := format.Url; formatStreamUrl != nil {
				itag := format.Itag
				streamDetails, ok := itags[itag]
				if ok {
					streamUrl := strings.ReplaceAll(*formatStreamUrl, "\\u0026", "&")
					stream := streams.Stream{
						Url:              streamUrl,
						StreamDetails:    streamDetails,
						ExpiresInSeconds: expiresInSeconds,
					}
					result = append(result, stream)
				}
			}
		}
	}

	return result, nil
}

func extractEncodedStreams(pageHtml *string, streamingData *youtube.StreamingData) ([]streams.Stream, error) {
	expiresInSeconds, err := strconv.ParseInt(streamingData.ExpiresInSeconds, 10, 64)
	if err != nil {
		return nil, err
	}

	var formats []youtube.Format
	for _, f := range append(streamingData.Formats, streamingData.AdaptiveFormats...) {
		if f.Type != "FORMAT_STREAM_TYPE_OTF" {
			formats = append(formats, f)
		}
	}

	encodedSignatures := getEncodedSignatures(formats)

	if len(encodedSignatures) > 0 {
		matches := utils.PatternDecryptionJsFile.FindStringSubmatch(*pageHtml)
		if matches == nil {
			matches = utils.PatternDecryptionJsFileWithoutSlash.FindStringSubmatch(*pageHtml)
		}
		if matches == nil {
			return nil, fmt.Errorf("error parsing js file")
		}

		decipherJsFileName := strings.ReplaceAll(matches[0], `\/`, "/")

		signature, err := decipherSignature(encodedSignatures, decipherJsFileName)
		if err != nil {
			return nil, err
		}
		signatures := strings.Split(*signature, "\n")

		result := make([]streams.Stream, 0, len(encodedSignatures))
		for i, encSignature := range encodedSignatures {
			sig := signatures[i]
			stream := streams.Stream{
				Url:              fmt.Sprintf("%s&sig=%s", encSignature.Url, sig),
				StreamDetails:    encSignature.StreamDetails,
				ExpiresInSeconds: expiresInSeconds,
			}
			result = append(result, stream)
		}

		return result, nil
	} else {
		return nil, fmt.Errorf("error getting encoded signatures")
	}
}

func decipherSignature(encSignatures []*encodedSignature, decipherJsFileName string) (*string, error) {
	javaScriptFile, err := network.GetJsFile(decipherJsFileName)
	if err != nil {
		return nil, err
	}

	signatureDecFunctionMatches := utils.PatternSignatureDecFunction.FindStringSubmatch(javaScriptFile)
	if signatureDecFunctionMatches == nil {
		return nil, fmt.Errorf("error parsing signatureDecFunction")
	}

	decipherFunctionName := signatureDecFunctionMatches[1]

	patternMainVariable, err := regexp.Compile(`(var |\s|,|;)` +
		regexp.QuoteMeta(decipherFunctionName) +
		`(=function\((.{1,3})\)\{)`)

	if err != nil {
		return nil, err
	}

	var mainDecipherFunct string

	mainDecipherFunctionMatches := patternMainVariable.FindStringSubmatchIndex(javaScriptFile)
	if mainDecipherFunctionMatches != nil {
		mainDecipherFunct = "var " + decipherFunctionName + javaScriptFile[mainDecipherFunctionMatches[4]:mainDecipherFunctionMatches[5]]
	} else {
		patternMainFunction, err := regexp.Compile(`function ` +
			regexp.QuoteMeta(decipherFunctionName) +
			`(\((.{1,3})\)\{)`)
		if err != nil {
			return nil, err
		}

		mainDecipherFunctionMatches = patternMainFunction.FindStringSubmatchIndex(javaScriptFile)

		if mainDecipherFunctionMatches == nil {
			return nil, fmt.Errorf("main js function not found")
		}

		mainDecipherFunct = "function " + decipherFunctionName + javaScriptFile[mainDecipherFunctionMatches[4]:mainDecipherFunct[5]]
	}

	startIndex := mainDecipherFunctionMatches[1]
	i := startIndex
	braces := 1

	for i < len(javaScriptFile) {
		if braces == 0 && startIndex+5 < i {
			mainDecipherFunct += javaScriptFile[startIndex:i] + ";"
			break
		}

		if javaScriptFile[i] == '{' {
			braces++
		} else if javaScriptFile[i] == '}' {
			braces--
		}
		i++
	}

	decipherFunctions := mainDecipherFunct
	variableFunctionMatches := utils.PatternVariableFunction.FindAllStringSubmatchIndex(mainDecipherFunct, -1)

	for _, match := range variableFunctionMatches {
		variableDef := fmt.Sprintf("var %s={", mainDecipherFunct[match[4]:match[5]])

		if strings.Contains(decipherFunctions, variableDef) {
			continue
		}

		startIndex := strings.Index(javaScriptFile, variableDef) + len(variableDef)
		i = startIndex
		braces = 1

		for i < len(javaScriptFile) {
			if braces == 0 {
				decipherFunctions += fmt.Sprintf("%s%s;", variableDef, javaScriptFile[startIndex:i])
				break
			}
			if javaScriptFile[i] == '{' {
				braces++
			} else if javaScriptFile[i] == '}' {
				braces--
			}
			i++
		}
	}

	functionMatches := utils.PatternFunction.FindAllStringSubmatchIndex(mainDecipherFunct, -1)
	for _, match := range functionMatches {
		functionDef := fmt.Sprintf("function %s", mainDecipherFunct[match[4]:match[5]])

		if strings.Contains(decipherFunctions, functionDef) {
			continue
		}

		startIndex = strings.Index(javaScriptFile, functionDef) + len(functionDef)
		i = 0
		braces = 0

		for i < len(javaScriptFile) {
			if braces == 0 && startIndex+5 < i {
				decipherFunctions += fmt.Sprintf("%s%s;", functionDef, javaScriptFile[startIndex:i])
				break
			}

			if javaScriptFile[i] == '{' {
				braces++
			} else if javaScriptFile[i] == '}' {
				braces--
			}
			i++
		}
	}

	signatures := make([]string, 0, len(encSignatures))
	for _, encSignature := range encSignatures {
		signatures = append(signatures, encSignature.Signature)
	}

	return decipherEncodedSignatures(signatures, decipherFunctions, decipherFunctionName)
}

func decipherEncodedSignatures(encSignatures []string, decipherFunctions string, decipherFunctionName string) (*string, error) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s function decipher(){return ", decipherFunctions))

	for i, s := range encSignatures {
		buffer.WriteString(fmt.Sprintf("%s('%s", decipherFunctionName, s))
		if i < len(encSignatures)-1 {
			buffer.WriteString("')+\"\\n\"+")
		} else {
			buffer.WriteString("')")
		}
	}

	buffer.WriteString("};decipher();")

	script := buffer.String()
	return js.ExecuteScript(script)
}

func getEncodedSignatures(formats []youtube.Format) []*encodedSignature {
	var encodedSignatures []*encodedSignature

	for _, format := range formats {
		streamDetails, ok := itags[format.Itag]
		if ok && format.SignatureCipher != nil {
			sigEncUrlMatches := utils.PatternSigEncUrl.FindStringSubmatch(*format.SignatureCipher)
			signatureMatches := utils.PatternSignature.FindStringSubmatch(*format.SignatureCipher)

			if sigEncUrlMatches == nil || signatureMatches == nil {
				continue
			}

			signatureUrl, err := url.QueryUnescape(sigEncUrlMatches[1])
			if err != nil {
				continue
			}

			signature, err := url.QueryUnescape(signatureMatches[1])
			if err != nil {
				continue
			}

			encSignature := new(encodedSignature)
			encSignature.Url = signatureUrl
			encSignature.Signature = signature
			encSignature.StreamDetails = streamDetails

			encodedSignatures = append(encodedSignatures, encSignature)
		}
	}

	return encodedSignatures
}
