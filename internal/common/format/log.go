package format

import "encoding/json"

func ToString(v any) string {
	data, _ := json.MarshalIndent(v, "", "	")
	return string(data)
}

func MarshalString(v any) (string, error) {
	bytes, err := json.Marshal(v)
	return string(bytes), err
}
