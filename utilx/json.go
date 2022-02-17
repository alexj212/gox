package utilx

import (
    "encoding/json"
)

// ConvertMapTo converts a map to a struct and returns the json string as well extracted from the map. Uses struct tags to populate struct
func ConvertMapTo(m map[string]interface{}, target interface{}) (result interface{}, jsonString string, err error) {
    jsonBytes, err := json.Marshal(m)
    if err != nil {
        return nil, jsonString, err
    }

    // convert json to struct
    err = json.Unmarshal(jsonBytes, target)
    if err != nil {
        return nil, jsonString, err
    }
    jsonString = string(jsonBytes)
    return target, jsonString, nil
}

func FormatJson(obj interface{}) (jsonStr string, err error) {
    b, err := json.MarshalIndent(obj, "", "  ")
    if err != nil {
        return "", err
    }
    return string(b), err
}
