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
	140: {140, streams.Audio, streams.M4A, codecs.AAC, 0, 0, 128, 0},
	141: {141, streams.Audio, streams.M4A, codecs.AAC, 0, 0, 256, 0},
	256: {256, streams.Audio, streams.M4A, codecs.AAC, 0, 0, 192, 0},
	258: {258, streams.Audio, streams.M4A, codecs.AAC, 0, 0, 384, 0},
	171: {171, streams.Audio, streams.WebM, codecs.Vorbis, 0, 0, 128, 0},
	249: {249, streams.Audio, streams.WebM, codecs.Opus, 0, 0, 48, 0},
	250: {250, streams.Audio, streams.WebM, codecs.Opus, 0, 0, 64, 0},
	251: {251, streams.Audio, streams.WebM, codecs.Opus, 0, 0, 160, 0},

	160: {160, streams.Video, streams.MP4, 0, codecs.H264, 144, 0, 0},
	133: {133, streams.Video, streams.MP4, 0, codecs.H264, 240, 0, 0},
	134: {134, streams.Video, streams.MP4, 0, codecs.H264, 360, 0, 0},
	135: {135, streams.Video, streams.MP4, 0, codecs.H264, 480, 0, 0},
	136: {136, streams.Video, streams.MP4, 0, codecs.H264, 720, 0, 0},
	137: {137, streams.Video, streams.MP4, 0, codecs.H264, 1080, 0, 0},
	264: {264, streams.Video, streams.MP4, 0, codecs.H264, 1440, 0, 0},
	266: {266, streams.Video, streams.MP4, 0, codecs.H264, 2160, 0, 0},
	298: {298, streams.Video, streams.MP4, 0, codecs.H264, 720, 0, 60},
	299: {299, streams.Video, streams.MP4, 0, codecs.H264, 1080, 0, 60},
	278: {278, streams.Video, streams.WebM, 0, codecs.VP9, 144, 0, 0},
	242: {242, streams.Video, streams.WebM, 0, codecs.VP9, 240, 0, 0},
	243: {243, streams.Video, streams.WebM, 0, codecs.VP9, 360, 0, 0},
	244: {244, streams.Video, streams.WebM, 0, codecs.VP9, 480, 0, 0},
	247: {247, streams.Video, streams.WebM, 0, codecs.VP9, 720, 0, 0},
	248: {248, streams.Video, streams.WebM, 0, codecs.VP9, 1080, 0, 0},
	271: {271, streams.Video, streams.WebM, 0, codecs.VP9, 1440, 0, 0},
	313: {313, streams.Video, streams.WebM, 0, codecs.VP9, 2160, 0, 0},
	302: {302, streams.Video, streams.WebM, 0, codecs.VP9, 720, 0, 60},
	308: {308, streams.Video, streams.WebM, 0, codecs.VP9, 1440, 0, 60},
	303: {303, streams.Video, streams.WebM, 0, codecs.VP9, 1080, 0, 60},
	315: {315, streams.Video, streams.WebM, 0, codecs.VP9, 2160, 0, 60},

	17: {17, streams.Multiplexed, streams.ThreeGP, codecs.AAC, codecs.MPEG4, 144, 24, 0},
	36: {36, streams.Multiplexed, streams.ThreeGP, codecs.AAC, codecs.MPEG4, 240, 32, 0},
	5:  {5, streams.Multiplexed, streams.FLV, codecs.MP3, codecs.H263, 144, 64, 0},
	43: {43, streams.Multiplexed, streams.WebM, codecs.Vorbis, codecs.VP8, 360, 128, 0},
	18: {18, streams.Multiplexed, streams.MP4, codecs.AAC, codecs.H264, 360, 96, 0},
	22: {22, streams.Multiplexed, streams.MP4, codecs.AAC, codecs.H264, 720, 192, 0},

	91: {91, streams.Live, streams.MP4, codecs.AAC, codecs.H264, 144, 48, 0},
	92: {92, streams.Live, streams.MP4, codecs.AAC, codecs.H264, 240, 48, 0},
	93: {93, streams.Live, streams.MP4, codecs.AAC, codecs.H264, 360, 128, 0},
	94: {94, streams.Live, streams.MP4, codecs.AAC, codecs.H264, 480, 128, 0},
	95: {95, streams.Live, streams.MP4, codecs.AAC, codecs.H264, 720, 256, 0},
	96: {96, streams.Live, streams.MP4, codecs.AAC, codecs.H264, 1080, 256, 0},
}

type encodedSignature struct {
	Url           string
	Signature     string
	StreamDetails streams.StreamDetails
}

func ExtractStreams(pageHtml string, streamingData youtube.StreamingData) []streams.Stream {
	notEncodedStreams := extractNotEncodedStreams(streamingData)

	encodedStreams, err := extractEncodedStreams(pageHtml, streamingData)
	if err != nil {
		return notEncodedStreams
	}

	return append(notEncodedStreams, encodedStreams...)
}

func extractNotEncodedStreams(streamingData youtube.StreamingData) []streams.Stream {
	var result []streams.Stream
	expiresInSeconds, err := strconv.ParseInt(streamingData.ExpiresInSeconds, 10, 64)
	if err != nil {
		return result
	}

	formats := append(streamingData.Formats, streamingData.AdaptiveFormats...)

	for _, format := range formats {
		streamUrl := format.Url
		if streamUrl != "" && format.Type != "FORMAT_STREAM_TYPE_OTF" {
			itag := format.Itag
			streamDetails, ok := itags[itag]
			if ok {
				streamUrl = strings.ReplaceAll(streamUrl, "\\u0026", "&")
				result = append(result, streams.Stream{
					Url:              streamUrl,
					StreamDetails:    streamDetails,
					ExpiresInSeconds: expiresInSeconds,
				})
			}
		}
	}

	return result
}

func extractEncodedStreams(pageHtml string, streamingData youtube.StreamingData) ([]streams.Stream, error) {
	expiresInSeconds, err := strconv.ParseInt(streamingData.ExpiresInSeconds, 10, 64)
	if err != nil {
		return nil, err
	}

	formats := append(streamingData.Formats, streamingData.AdaptiveFormats...)

	encodedSignatures := getEncodedSignatures(utils.Filter(formats, func(f youtube.Format) bool {
		return f.Type != "FORMAT_STREAM_TYPE_OTF"
	}))

	if len(encodedSignatures) > 0 {
		matches := utils.PatternDecryptionJsFile.FindStringSubmatch(pageHtml)
		if matches == nil {
			matches = utils.PatternDecryptionJsFileWithoutSlash.FindStringSubmatch(pageHtml)
		}
		if matches == nil {
			return nil, fmt.Errorf("error parsing js file")
		}

		decipherJsFileName := strings.ReplaceAll(matches[0], `\/`, "/")

		signature, err := decipherSignature(encodedSignatures, decipherJsFileName)
		if err != nil {
			return nil, err
		}
		signatures := strings.Split(signature, "\n")

		var result []streams.Stream
		for i, encSignature := range encodedSignatures {
			sig := signatures[i]
			result = append(result, streams.Stream{
				Url:              fmt.Sprintf("%s&sig=%s", encSignature.Url, sig),
				StreamDetails:    encSignature.StreamDetails,
				ExpiresInSeconds: expiresInSeconds,
			})
		}

		return result, nil
	} else {
		return nil, fmt.Errorf("error getting encoded signatures")
	}
}

func decipherSignature(encSignatures []encodedSignature, decipherJsFileName string) (string, error) {
	javaScriptFile, err := network.GetJsFile(decipherJsFileName)
	if err != nil {
		return "", err
	}

	signatureDecFunctionMatches := utils.PatternSignatureDecFunction.FindStringSubmatch(javaScriptFile)
	if signatureDecFunctionMatches == nil {
		return "", fmt.Errorf("error parsing signatureDecFunction")
	}

	decipherFunctionName := signatureDecFunctionMatches[1]

	patternMainVariable, err := regexp.Compile(`(var |\s|,|;)` +
		regexp.QuoteMeta(decipherFunctionName) +
		`(=function\((.{1,3})\)\{)`)

	if err != nil {
		return "", err
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
			return "", err
		}

		mainDecipherFunctionMatches = patternMainFunction.FindStringSubmatchIndex(javaScriptFile)

		if mainDecipherFunctionMatches == nil {
			return "", fmt.Errorf("main js function not found")
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

	return decipherEncodedSignatures(utils.Map(encSignatures, func(s encodedSignature) string {
		return s.Signature
	}), decipherFunctions, decipherFunctionName)
}

func decipherEncodedSignatures(encSignatures []string, decipherFunctions string, decipherFunctionName string) (string, error) {
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

func getEncodedSignatures(formats []youtube.Format) []encodedSignature {
	var encodedSignatures []encodedSignature

	for _, format := range formats {
		streamDetails, ok := itags[format.Itag]
		if ok && format.SignatureCipher != "" {
			sigEncUrlMatches := utils.PatternSigEncUrl.FindStringSubmatch(format.SignatureCipher)
			signatureMatches := utils.PatternSignature.FindStringSubmatch(format.SignatureCipher)

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

			encodedSignatures = append(encodedSignatures, encodedSignature{signatureUrl, signature, streamDetails})
		}
	}

	return encodedSignatures
}
