package utilx

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode"
)

var templateFuncs = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          TrimRightSpace,
	"trimTrailingWhitespaces": TrimRightSpace,
	"appendIfNotPresent":      AppendIfNotPresent,
	"rpad":                    Rpad,
	"padSpaceAfter":           PadSpaceAfter,
	"stringInSlice":           StringInSlice,
}

func TrimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

func PadSpaceAfter(s string) string {
	if s == "" {
		return ""
	}
	return s + " "
}

// FIXME appendIfNotPresent is unused by cobra and should be removed in a version 2. It exists only for compatibility with users of cobra.

// AppendIfNotPresent will append stringToAppend to the end of s, but only if it's not yet present in s.
func AppendIfNotPresent(s, stringToAppend string) string {
	if strings.Contains(s, stringToAppend) {
		return s
	}
	return s + " " + stringToAppend
}

// Rpad adds padding to the right of a string.
func Rpad(s string, padding int) string {
	tpl := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(tpl, s)
}

// Tmpl executes the given template text on data, writing the result to w.
func Tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}

// LD compares two strings and returns the levenshtein distance between them.
func LD(s, t string, ignoreCase bool) int {
	if ignoreCase {
		s = strings.ToLower(s)
		t = strings.ToLower(t)
	}
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(s)][len(t)]
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
