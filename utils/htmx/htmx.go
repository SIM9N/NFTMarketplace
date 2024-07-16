package htmx

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

func isJson(input string) bool {
	match, _ := regexp.MatchString(`^\s*(\{.*\}|\[.*\])\s*$`, input)
	return match
}

func DecodeHTMXBody(r *http.Request) (map[string]interface{}, error) {
	body, _ := io.ReadAll(r.Body)
	r.Body.Close()

	urlValues, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	jsonData := make(map[string]interface{})
	for k, v := range urlValues {
		if len(v) == 0 {
			continue
		}

		var value any

		if isJson(v[0]) {
			json.Unmarshal([]byte(v[0]), &value)
		} else {
			value = v[0]
		}

		jsonData[k] = value
	}

	return jsonData, nil
}

func DecodeHTMXQuery(r *http.Request) (map[string]interface{}, error) {
	urlValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	jsonData := make(map[string]interface{})
	for k, v := range urlValues {
		if len(v) == 0 {
			continue
		}

		var value any

		if isJson(v[0]) {
			json.Unmarshal([]byte(v[0]), &value)
		} else {
			value = v[0]
		}

		jsonData[k] = value
	}

	return jsonData, nil
}
