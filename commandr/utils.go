package commandr

import (
	"fmt"
	markdownP "github.com/MichaelMure/go-term-markdown"
	"io"
	"strings"
	"time"
)

// RenderMarkdown render the markdown string in terminal
func RenderMarkdown(w io.Writer, markdown string) {
	result := markdownP.Render(markdown, 80, 6)
	w.Write(result)
}

// Type mimic typewriter typing content to writer
func Type(w io.Writer, format string, a ...interface{}) {

	var content string
	if len(a) > 0 {
		content = fmt.Sprintf(format, a...)
	} else {
		content = format
	}

	chars := strings.Split(content, "")

	for _, c := range chars {

		time.Sleep(20 * time.Millisecond)

		w.Write([]byte(c))
	}
}

// AddText convenience to write string to writer
func AddText(w io.Writer, content string) {
	w.Write([]byte(content))
}

// ClearTerm convenience to write clear screen to writer
func ClearTerm(w io.Writer) {
	w.Write([]byte("\033c"))
}
