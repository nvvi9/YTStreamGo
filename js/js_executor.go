package js

import "rogchap.com/v8go"

var iso = v8go.NewIsolate()

func ExecuteScript(script string) (string, error) {
	ctx := v8go.NewContext(iso)

	val, err := ctx.RunScript(script, "extractor.js")
	if err != nil {
		return "", err
	}

	return val.String(), nil
}
