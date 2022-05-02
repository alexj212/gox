package commandr

import (
    markdownP "github.com/MichaelMure/go-term-markdown"
    "io"
    "strings"
    "time"
)

func RenderMarkdown(w io.Writer, markdown string) {
    result := markdownP.Render(markdown, 80, 6)
    w.Write(result)
}

func Type(w io.Writer, content string) {
    chars := strings.Split(content, "")

    for _, c := range chars {

        time.Sleep(20 * time.Millisecond)

        w.Write([]byte(c))
    }
}

func AddText(w io.Writer, content string) {
    w.Write([]byte(content))
}

func ClearTerm(w io.Writer) {
    w.Write([]byte("\033c"))
}
