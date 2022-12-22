package commandr

import (
	"bufio"
	"bytes"
	"github.com/fatih/color"
)

// DefaultCommands main struct for all commands
var DefaultCommands = &Command{ExecLevel: All}

// HistoryCommand command to view history of commands executed
var HistoryCommand = &Command{Use: "history", Exec: historyCmd, Short: "Show the history of commands executed", ExecLevel: All}

// WhoamiCommand command to view user details
var WhoamiCommand = &Command{Use: "whoami", Exec: whoamiCmd, Short: "Show user details about logged in user", ExecLevel: All}

// ClsCommand command to clear web screen
var ClsCommand = &Command{Use: "cls", Exec: clsCmd, Short: "send cls event to terminal client", ExecLevel: All}

// ExitCommand command to exit
var ExitCommand = &Command{Use: "exit", Exec: exitCmd, Short: "exit the session", ExecLevel: All}

func init() {
	DefaultCommands.AddCommand(HistoryCommand)
	DefaultCommands.AddCommand(WhoamiCommand)
	DefaultCommands.AddCommand(ClsCommand)
	DefaultCommands.AddCommand(ExitCommand)
	return
}

// HandleCommands handler function to execute commands
func HandleCommands(Commands *Command) (handler func(Client, string)) {

	handler = func(client Client, cmdLine string) {
		// loge.Info("handleMessage  - authenticated user message.Payload: [" + cmd+"]")

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		parsed, err := NewCommandArgs(cmdLine, writer)

		if err != nil {
			client.WriteString(color.RedString("Error parsing command: %v\n", err))
			return
		}
		Commands.Execute(client, parsed)
		_ = writer.Flush()
		result := b.String()
		client.WriteString(color.WhiteString(result))
	}
	return
}

func whoamiCmd(client Client, cmd *Command, _ *CommandArgs) (err error) {
	client.WriteString(color.GreenString("whoami  username: %v  exec level: %v\n", client.UserName(), client.ExecLevel()))
	return
}

func historyCmd(client Client, cmd *Command, _ *CommandArgs) (err error) {
	if len(client.History()) > 0 {
		for i, cmd := range client.History() {
			client.WriteString(color.GreenString("History[%d]: %v\n", i, cmd))
		}
	} else {
		client.WriteString(color.GreenString("History is empty\n"))
	}
	return
}

func exitCmd(client Client, cmd *Command, _ *CommandArgs) (err error) {
	client.WriteString(color.GreenString("Bye bye 👋\n"))
	client.Close()
	return
}

func clsCmd(client Client, cmd *Command, _ *CommandArgs) (err error) {
	client.WriteString("\033c")
	return
}
