package commandr

import (
	"bufio"
	"bytes"
	"io"

	"github.com/fatih/color"
)

// HandleCommands handler function to execute commands
func HandleCommands(Commands *Command) (handler func(io.Writer, string)) {

	handler = func(client io.Writer, cmdLine string) {
		// log.Printf("handleMessage  - authenticated user message.Payload: [" + cmd+"]")

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		parsed, err := NewCommandArgs(cmdLine, writer)

		if err != nil {
			client.Write([]byte(color.RedString("Error parsing command: %v\n", err)))
			return
		}
		Commands.Execute(client, parsed)
		_ = writer.Flush()
		result := b.String()
		client.Write([]byte(color.WhiteString(result)))
	}
	return
}
