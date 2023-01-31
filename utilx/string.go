package utilx

import (
	"encoding/json"
	"fmt"
)

// ToString util to convert obj to json string
func ToString(obj interface{}) string {
	b, err := json.MarshalIndent(obj, "", "")
	if err != nil {
		return fmt.Sprintf("obj: %T %v could not marshal to json err: %v", obj, obj, err)
	}
	return string(b)
}
