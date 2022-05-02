package commandr

import (
    "bufio"
    "bytes"
    "github.com/fatih/color"
)

// DefaultCommands main struct for all commands
var DefaultCommands = &Command{ExecLevel: All}

// HistoryCommand command to view history of commands executed
var HistoryCommand = &Command{Use: "history", Exec: displayHistory, Short: "Show the history of commands executed", ExecLevel: All}

// UserCommand command to view user details
var UserCommand = &Command{Use: "user", Exec: displayUserInfo, Short: "Show user details about logged in user", ExecLevel: All}

// ClsCommand command to clear web screen
var ClsCommand = &Command{Use: "cls", Exec: clsCommand, Short: "send cls event to terminal client", ExecLevel: All}

// ExitCommand command to exit
var ExitCommand = &Command{Use: "exit", Exec: exitCommand, Short: "exit the session", ExecLevel: All}

func init() {
    DefaultCommands.AddCommand(HistoryCommand)
    DefaultCommands.AddCommand(UserCommand)
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
            client.Write([]byte(color.RedString("Error parsing command: %v\n", err)))
            return
        }
        Commands.Execute(client, parsed)
        writer.Flush()
        result := b.String()
        client.Write([]byte(color.WhiteString(result)))
    }
    return
}

func displayUserInfo(client Client, args *CommandArgs) (err error) {
    client.Write([]byte(color.GreenString("Username       : %v / %v\n", client.UserName(), client.ExecLevel())))
    return
}

func displayHistory(client Client, args *CommandArgs) (err error) {
    if len(client.History()) > 0 {
        for i, cmd := range client.History() {
            client.Write([]byte(color.GreenString("History[%d]: %v\n", i, cmd)))
        }
    } else {
        client.Write([]byte(color.GreenString("History is empty\n")))
    }
    return
}

func exitCommand(client Client, args *CommandArgs) (err error) {
    client.Write([]byte(color.GreenString("Bye bye 👋\n")))
    client.Close()
    return
}

func clsCommand(client Client, args *CommandArgs) (err error) {
    client.Write([]byte("\033c"))
    client.Close()
    return
}
