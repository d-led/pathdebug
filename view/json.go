package view

import (
	"encoding/json"
)

func RenderJson() ([]byte, error) {
	results := getResults()
	return json.MarshalIndent(results, "", " ")
}
