package main

import (
	"fmt"
	"github.com/alexj212/gox"
	"github.com/alexj212/gox/commandr"
	"github.com/alexj212/gox/utilx"
	gossh "golang.org/x/crypto/ssh"
	"os"
	"os/user"
	"path/filepath"
)

// "github.com/tj/go-spin"
// nmap --script ssh2-enum-algos -sV -p  2022 localhost

func main() {

	usr, _ := user.Current()
	sshCertDir := filepath.Join(usr.HomeDir, ".ssh")
	authorizedKeyFile := filepath.Join(sshCertDir, "authorized_keys")

	keys, err := gox.LoadAuthorizedKeys(authorizedKeyFile)
	if err != nil {
		fmt.Printf("Unable to load authorized keys: %v\n", err)
		os.Exit(1)
	}

	appKey, err := gox.GetAppKey()
	if err != nil {
		fmt.Printf("Unable to load app key: %v\n", err)
		os.Exit(1)
	}

	hostKey, err := gossh.NewSignerFromKey(appKey)
	if err != nil {
		fmt.Printf("Unable create hostKey: %v\n", err)
		os.Exit(1)
	}
	svc, err := commandr.NewSshService(2022, hostKey)
	if err != nil {
		fmt.Printf("Unable to launch ssh server: %v\n", err)
		os.Exit(1)
	}

	svc.RegisterUser("alexj_a", commandr.Admin, keys, nil)
	svc.RegisterUser("alexj_sa", commandr.SuperAdmin, keys, nil)
	svc.RegisterUser("alexj", commandr.User, keys, nil)

	svc.AddCommand(ExitCommand)
	svc.AddCommand(EchoCommand)
	svc.AddCommand(DebugCommand)
	svc.AddCommand(TldrCmd)
	svc.AddCommand(LinesCommand)
	svc.AddCommand(AdminLevelCommand)

	svc.Spawn()
	utilx.LoopForever(nil)
}
