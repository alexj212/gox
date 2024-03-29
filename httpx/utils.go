package httpx

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"log"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// HttpNocacheContent will set the headers for content type along with no caching.
func HttpNocacheContent(w http.ResponseWriter, content string) {
	w.Header().Set("Content-Type", content)
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

// HttpNocacheJson will set the headers on an http Response for a text/json content type along with no cache.
func HttpNocacheJson(w http.ResponseWriter) {
	HttpNocacheContent(w, "text/json")
}

// InternalServerError will return an error to the client, sending 500 error code to the client with generic string
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	log.Printf("Internal server error %s: %s", mux.CurrentRoute(r).GetName(), err.Error())
	http.Error(w, "Internal server error", 500)
}

// SendJson will return take a value and serialize it to json and return the http response.
func SendJson(w http.ResponseWriter, r *http.Request, val interface{}) {
	bytes, err := json.Marshal(val)
	if err != nil {
		log.Printf("Send json error: %v\n", err)
		InternalServerError(w, r, err)
		return
	}
	HttpNocacheContent(w, "text/json")
	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("error calling w.Write() error: %v\n", err)
	}
}

// SendText will return take a string and send to the http writer.
func SendText(w http.ResponseWriter, r *http.Request, val string) {
	HttpNocacheContent(w, "text/plain")
	_, err := w.Write([]byte(val))
	if err != nil {
		log.Printf("error calling w.Write() error: %v\n", err)
	}
}

// AddHeadersHandler will take a map of string/string and use it to set the key and value as the header name and value respectively.
func AddHeadersHandler(addHeaders map[string]string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for key, value := range addHeaders {
			w.Header().Set(key, value)
		}

		h.ServeHTTP(w, r)
	})
}

// ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
	start net.IP
	end   net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

var privateRanges = []ipRange{
	{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
	{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

// IsPrivateSubnet - check to see if this ip is in a private subnet
func IsPrivateSubnet(ipAddress net.IP) bool {
	// my use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range privateRanges {
			// check if this ip is in a private range
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

// GetIPAddress will take a http request and check headers if it has been proxied to extract what the server believes to be the client ip address.
func GetIPAddress(r *http.Request) string {
	ip := ""
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip = strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || IsPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return ip
		}
	}

	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return ip
}

// GetRequestField is a convenience member to access the http request fields and return the value or an error if it does not exist.
func GetRequestField(r *http.Request, fieldName string) (val string, err error) {
	valVar, ok := mux.Vars(r)[fieldName]
	if !ok {
		return "", errors.Errorf("Bad request: specify %v", fieldName)
	}

	return valVar, nil
}

// JSON send back json to ResponseWriter
func JSON(w http.ResponseWriter, code int, val interface{}) error {

	b, err := json.Marshal(val)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(b)
	return err
}

// ShouldBindJSON decode val into json
func ShouldBindJSON(r *http.Request, val any) error {
	err := json.NewDecoder(r.Body).Decode(val)
	if err != nil {
		return err
	}
	return nil
}
