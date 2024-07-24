package gox

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/droundy/goopt"
	"github.com/pkg/errors"
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
		fmt.Printf("\n%s\n", string(historyBytes))
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

// GetAppKeyFileName return name of app key file
func GetAppKeyFileName() (appKeyFileName string) {
	dirname, _ := os.UserHomeDir()
	appKeyFileName = filepath.Join(dirname, fmt.Sprintf(".%s.key", GetAppName()))
	return
}

// GetAppKey read AppKey file to rsa.PrivateKey.
func GetAppKey() (*rsa.PrivateKey, error) {
	appKeyFileName := GetAppKeyFileName()
	exists := FileExists(appKeyFileName)
	if exists {
		f, err := ioutil.ReadFile(appKeyFileName)
		if err != nil {
			return nil, err
		}

		block, _ := pem.Decode(f)
		if block != nil {
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			return key, nil
		}
		return nil, errors.New("unable to extract private key")
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	// Validate Private Key
	err = key.Validate()
	if err != nil {
		return nil, err
	}
	log.Println("Private Key generated")
	payload := encodePrivateKeyToPEM(key)
	ioutil.WriteFile(appKeyFileName, payload, 0644)
	return key, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// GetAppStorageFileName return name of app storage file
func GetAppStorageFileName() (appStorageFileName string) {
	dirname, _ := os.UserHomeDir()
	appStorageFileName = filepath.Join(dirname, fmt.Sprintf(".%s.db", GetAppName()))
	return
}

// LoadAppStorage read storage, if it exists, and return the data. If it does not exist, return empty.
func LoadAppStorage[T any](filename string, base *T) (*T, bool, error) {
	if base == nil {
		return nil, false, fmt.Errorf("base config cannot be nil")
	}

	exists := FileExists(filename)
	if !exists {
		return base, false, nil
	}

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return base, false, err
	}

	err = json.Unmarshal(f, base)
	if err != nil {
		return base, false, err
	}
	return base, true, nil
}

// SaveAppStorage read AppKey file to rsa.PrivateKey.
func SaveAppStorage[T any](filename string, val *T) error {

	payload, err := json.MarshalIndent(val, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, payload, 0644)
}
