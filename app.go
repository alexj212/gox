package gox

import (
	"fmt"
	"github.com/droundy/goopt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	showHistory  = goopt.Flag([]string{"-H"}, nil, "Show history", "")
	clearHistory = goopt.Flag([]string{"--XH"}, nil, "Clear history", "")
)

// HandleHistory will save the command line history to file in the home dir. If no arguments are passed to the app, it will
// display history of past commands executed. If an error occurs it will be returned.
func HandleHistory() (exitApp bool, err error) {

	if *clearHistory {
		historyFileName := GetHistoryFileName()
		fmt.Printf("Removing history file %s\n", historyFileName)
		os.Remove(historyFileName)
	}

	if *showHistory {
		historyBytes, err := GetHistoryFile()
		if err != nil {
			return true, err
		}

		if len(historyBytes) == 0 {
			usage := goopt.Usage()
			fmt.Printf("Command has never been run\nUsage:\n%s\n", usage)
			return true, nil
		}
		fmt.Printf("Showing history for command: %s\n", GetAppName())
		fmt.Printf("\n%s\b", string(historyBytes))
		return true, nil
	}

	appendHistory()
	return false, nil
}

// GetAppName return name of executable sans path
func GetAppName() (appName string) {
	_, appName = filepath.Split(os.Args[0])
	return
}

// GetHistoryFileName return name of history file
func GetHistoryFileName() (historyFileName string) {
	dirname, _ := os.UserHomeDir()
	historyFileName = filepath.Join(dirname, fmt.Sprintf(".%s.history", GetAppName()))
	return
}

// GetHistoryFile read history file to []byte.
func GetHistoryFile() ([]byte, error) {
	historyFileName := GetHistoryFileName()
	exists := FileExists(historyFileName)
	if exists {
		f, err := ioutil.ReadFile(historyFileName)
		if err != nil {
			return nil, err
		}
		return f, err
	}

	payload := make([]byte, 0)
	ioutil.WriteFile(historyFileName, payload, 0644)
	return payload, nil
}

func appendHistory() {
	historyFileName := GetHistoryFileName()
	var fullCmdLine string
	if strings.Contains(os.Args[0], "go-build") {
		fullCmdLine = "go run . "
		fullCmdLine = fullCmdLine + strings.Join(os.Args[1:], " ")
	} else {
		fullCmdLine = strings.Join(os.Args, " ")
	}

	historyBytes, err := GetHistoryFile()
	if err != nil {
		return
	}

	history := string(historyBytes)
	if strings.Contains(history, fullCmdLine) {
		log.Printf("cmd line: %s already in history file\n", fullCmdLine)
		return
	}

	f, err := os.OpenFile(historyFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error opening history file: %s error: %v\n", historyFileName, err)
		return
	}

	defer f.Close()

	if _, err = f.WriteString(fullCmdLine + "\n"); err != nil {
		log.Printf("Error writing to  history file: %s error: %v\n", historyFileName, err)
	}
}

// GetDefaultLogDir return default log dir for the app
func GetDefaultLogDir() (logDir string) {
	logDir = filepath.Join(os.TempDir(), GetAppName())
	return
}
