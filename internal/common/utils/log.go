package utils

import "encoding/json"

func ToString(v any) string {
	data, _ := json.MarshalIndent(v, "", "	")
	return string(data)
}
