package commandr

import (
	"bufio"
	"bytes"
	"io"

	"github.com/fatih/color"
)

// DefaultCommands main struct for all commands
var DefaultCommands = &Command{ExecLevel: All}

// ClsCommand command to clear web screen
var ClsCommand = &Command{Use: "cls", Exec: clsCmd, Short: "send cls event to terminal client", ExecLevel: All}

// ExitCommand command to exit
var ExitCommand = &Command{Use: "exit", Exec: exitCmd, Short: "exit the session", ExecLevel: All}

func init() {
	DefaultCommands.AddCommand(ClsCommand)
	DefaultCommands.AddCommand(ExitCommand)
	return
}

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

// HandleCommands handler function to execute commands
func HandleCommandsA(Commands *Command) (handler func(io.Writer, string)) {

	handler = func(client io.Writer, cmdLine string) {
		// loge.Info("handleMessage  - authenticated user message.Payload: [" + cmd+"]")

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

func exitCmd(client io.Writer, cmd *Command, _ *CommandArgs) (err error) {
	client.Write([]byte((color.GreenString("Bye bye 👋\n"))))
	return
}

func clsCmd(client io.Writer, cmd *Command, _ *CommandArgs) (err error) {
	client.Write([]byte(("\033c")))
	return
}
