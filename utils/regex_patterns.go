package utils

import "regexp"

var (
	PatternPlayerResponse               = regexp.MustCompile(`var ytInitialPlayerResponse\s*=\s*({.+?})\s*;`)
	PatternSigEncUrl                    = regexp.MustCompile("url=(.+?)(&|$)")
	PatternSignature                    = regexp.MustCompile("s=(.+?)(&|$)")
	PatternDecryptionJsFile             = regexp.MustCompile(`\\/s\\/player\\/([^\"]+?)\\.js`)
	PatternDecryptionJsFileWithoutSlash = regexp.MustCompile(`/s/player/([^\"]+?).js`)
	PatternSignatureDecFunction         = regexp.MustCompile(`(?:\b|[^a-zA-Z0-9$])([a-zA-Z0-9$]{1,4})\s*=\s*function\s*\(\s*a\s*\)\s*\{\s*a\s*=\s*a\.split\s*\(\s*""\s*\)`)
	PatternVariableFunction             = regexp.MustCompile(`([{; =])([a-zA-Z$][a-zA-Z0-9$]{0,2})\.([a-zA-Z$][a-zA-Z0-9$]{0,2})\(`)
	PatternFunction                     = regexp.MustCompile(`([{; =])([a-zA-Z$_][a-zA-Z0-9$]{0,2})\(`)
)
