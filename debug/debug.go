package debug

import "encoding/json"

// GetJSON ...
func GetJSON(data interface{}) string {
	j, _ := json.Marshal(data)
	return string(j)
}
