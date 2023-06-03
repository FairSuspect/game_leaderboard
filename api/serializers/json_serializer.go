package serializers

import "encoding/json"

func FormatJSON(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "	")
}
