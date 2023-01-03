package utils

import "encoding/json"

func ToJsonString(v interface{}) string {
	marshal, _ := json.Marshal(v)
	return string(marshal)
}
