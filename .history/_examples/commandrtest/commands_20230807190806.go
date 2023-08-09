package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/alexj212/gox/commandr"
	"github.com/fatih/color"
	"github.com/tj/go-spin"
)

// "github.com/tj/go-spin"

var TldrCmd = &commandr.Command{Use: "tldr", Exec: tldrCmd, Short: "echo input", ExecLevel: commandr.All}

var ExitCommand = &commandr.Command{Use: "exit", Exec: exitCmd, Short: "exit the session", ExecLevel: commandr.All}

var EchoCommand = &commandr.Command{Use: "echo", Exec: echoCmd, Short: "echo input", ExecLevel: commandr.All}

var DebugCommand = &commandr.Command{Use: "debug", Exec: debugCmd, Short: "debug", ExecLevel: commandr.All}

var LinesCommand = &commandr.Command{Use: "lines", Exec: linesCmd, Short: "lines", ExecLevel: commandr.All}

var AdminLevelCommand = &commandr.Command{Use: "admintest", Exec: adminLevelCmd, Short: "admintest", ExecLevel: commandr.Admin}

func debugCmd(client io.Writer, cmd *commandr.Command, args *commandr.CommandArgs) (err error) {
	client.Write([]byte(color.GreenString("args.CmdLine: %v\n", args.CmdLine)))
	client.Write([]byte(color.GreenString("args.Args: %v\n", strings.Join(args.Args, " | "))))
	client.Write([]byte(color.GreenString("args.PealOff: %v\n", args.PealOff(1))))
	client.Write([]byte(color.GreenString("args.Debug: %v\n", args.Debug())))
	return
}
func echoCmd(client io.Writer, cmd *Command, args *commandr.CommandArgs) (err error) {
	//text := args.PealOff(0)
	client.Write([]byte(color.GreenString("%v\n", args.PealOff(0))))
	return
}

func exitCmd(client io.Writer, cmd *Command, args *commandr.CommandArgs) (err error) {
	client.Write([]byte(color.GreenString("Bye bye 👋\n")))
	client.Close()
	return
}

func tldrCmd(client io.Writer, cmd *Command, args *commandr.CommandArgs) (err error) {
	err = args.Parse()
	if err != nil {
		return err
	}

	isLoading := true

	commandr.AddText(client, "\n")
	go func() {
		// Cool terminal loading spinner

		s := spin.New()

		for {
			if isLoading == false {
				break
			}

			text := fmt.Sprintf("\r\033[36mLoading tldr\033[m %s ", s.Next())
			commandr.AddText(client, text)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	time.Sleep(1 * time.Second)
	// Clear terminal line
	commandr.AddText(client, "\033[2K\n")

	isLoading = false

	tldrName := ""

	if len(args.Args) > 1 {
		tldrName = args.Args[1]
	}

	if tldrName == "" {
		commandr.Type(client, color.RedString("\nYou need a specify a tldr to lookup.\n"))
		return
	}

	md := fmt.Sprintf("tldrName: %v\n", tldrName)
	commandr.Type(client, md)
	//    text := RenderMarkdownTerminal(md)
	//
	//AddText(stream, text)
	return
}

func adminLevelCmd(client io.Writer, cmd *commandr.Command, args *commandr.CommandArgs) (err error) {
	//text := args.PealOff(0)
	client.Write([]byte(color.GreenString("admintest\n")))

	return
}

func linesCmd(client io.Writer, cmd *commandr.Command, args *commandr.CommandArgs) (err error) {

	cnt := args.FlagSet.Int("cnt", 5, "number of lines to print")
	err = args.Parse()

	if err != nil {
		client.Write([]byte(color.GreenString("lines err: %v\n", err)))
		return
	}

	client.Write([]byte(color.GreenString("lines cnt: %v\n", *cnt)))
	client.Write([]byte(color.GreenString("lines invoked\n")))

	for i := 0; i < *cnt; i++ {
		client.Write([]byte(color.GreenString("line[%d]\n", i)))
	}
	return
}
