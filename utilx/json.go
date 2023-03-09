package utilx

import (
	"encoding/json"
	"fmt"
	"github.com/alexj212/gox"
	"io"
	"os"
	"strconv"
	"time"
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

// FormatJson util to convert obj to json string
func FormatJson(obj interface{}) (jsonStr string, err error) {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), err
}

// UnixTimeStamp defines a timestamp encoded as epoch seconds in JSON
type UnixTimeStamp time.Time

// MarshalJSON is used to convert the timestamp to JSON
func (t UnixTimeStamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

// UnmarshalJSON is used to convert the timestamp from JSON
func (t *UnixTimeStamp) UnmarshalJSON(s []byte) (err error) {
	r := string(s)

	millis, err := strconv.ParseUint(r, 10, 64)
	if err != nil {
		return err
	}

	val := time.UnixMilli(int64(millis))

	*t = UnixTimeStamp(val)
	return nil
}

// Unix returns t as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
func (t UnixTimeStamp) Unix() int64 {
	return time.Time(t).Unix()
}

// Time returns the JSON time as a time.Time instance in UTC
func (t UnixTimeStamp) Time() time.Time {
	return time.Time(t)
}

// String returns t as a formatted string
func (t UnixTimeStamp) String() string {
	return t.Time().String()
}

// LoadJson read file marshall json to type
func LoadJson[T any](f string, v *T) (*T, error) {

	exists := gox.FileExists(f)
	if exists {
		d, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}

		fmt.Printf("json: [%s]", string(d))
		err = json.Unmarshal(d, v)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

// ReadJson json from reader -
func ReadJson[T any](r io.Reader, _ *T) (val *T, err error) {

	body, err := io.ReadAll(r)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, val)
	if err != nil {
		return
	}

	return val, nil
}

// SaveJson write json from val
func SaveJson[T any](f string, v *T) error {

	payload, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(f, payload, 0644)
}
