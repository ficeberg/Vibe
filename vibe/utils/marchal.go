package utils

import (
	"encoding/json"
)

type Marshal struct{}

func (m *Marshal) S2M(source interface{}) (map[string]interface{}, error) {
	var f interface{}
	j, err := json.Marshal(source)
	if err != nil {
		return map[string]interface{}{}, err
	}
	if err = json.Unmarshal(j, &f); err != nil {
		return map[string]interface{}{}, err
	}

	return f.(map[string]interface{}), nil
}
