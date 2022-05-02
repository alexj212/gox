package commandr

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/gliderlabs/ssh"
    gossh "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/terminal"
    "io"
    "log"
    "net"
    "time"
)

type sshClient struct {
    s       ssh.Session
    user    *sshUser
    history []string
}

func (s *sshClient) Close() {
    s.s.Close()
}

func (s *sshClient) ExecLevel() ExecLevel {
    return s.user.level
}

func (s *sshClient) UserName() string {
    return s.s.User()
}

func (s *sshClient) History() []string {
    return s.history
}

func (s *sshClient) Write(p []byte) (n int, err error) {
    return s.s.Write(p)
}

func (svc *sshService) publicKeyValidator(ctx ssh.Context, key ssh.PublicKey) bool {
    user, ok := svc.LookupUser(ctx.User())
    if !ok {
        fmt.Printf("Login Attempt %v - user not found: %v\n", ctx.RemoteAddr(), ctx.User())
        // user not found
        return false
    }

    pubKey := key.Marshal()
    _, ok = user.keys[string(pubKey)]
    return ok // allow all keys, or use ssh.KeysEqual() to compare against known keys
}

func (svc *sshService) sshSessionHandler(s ssh.Session) {

    user, ok := svc.LookupUser(s.User())
    if !ok {
        io.WriteString(s, color.RedString("\nUnknown user: %v\n\n\n", s.User()))
        s.Close()
        return
    }

    c := &sshClient{
        s:       s,
        user:    user,
        history: make([]string, 0),
    }

    //ptyReq, winCh, isPty := s.Pty()
    //io.WriteString(s, fmt.Sprintf("ptyReq: %v\n", ptyReq))
    //io.WriteString(s, fmt.Sprintf("winCh: %v\n", winCh))
    //io.WriteString(s, fmt.Sprintf("isPty: %v\n", isPty))
    //io.WriteString(s, fmt.Sprintf("s.User(): %v\n", s.User()))
    //

    authorizedKey := gossh.MarshalAuthorizedKey(s.PublicKey())
    io.WriteString(s, fmt.Sprintf("public key used by %s    key: %s\n", s.User(), string(authorizedKey)))

    term := terminal.NewTerminal(s, "> ")
    line := ""
    for {
        line, _ = term.ReadLine()
        //if line == "quit" || line == "exit" {
        //    break
        //}

        parsed, err := NewCommandArgs(line, s)
        if err != nil {
            continue
        }

        c.history = append(c.history, line)
        DefaultCommands.Execute(c, parsed)
    }
}

type sshService struct {
    s     *ssh.Server
    ln    net.Listener
    users map[string]*sshUser
}

func NewSshService(port int, hostKey gossh.Signer) (SshService, error) {

    addr := fmt.Sprintf(":%d", port)

    ln, err := net.Listen("tcp", addr)
    if err != nil {
        return nil, err
    }

    svc := &sshService{
        ln:    ln,
        users: make(map[string]*sshUser),
    }

    server := &ssh.Server{
        Addr:             addr,
        PublicKeyHandler: svc.publicKeyValidator,
        Handler:          svc.sshSessionHandler,
    }
    svc.s = server

    // server.MaxTimeout = 30 * time.Second  // absolute connection timeout, none if empty
    server.IdleTimeout = 60 * time.Second // connection timeout when no activity, none if empty
    server.AddHostKey(hostKey)

    log.Printf("starting ssh server on port %s...\n", server.Addr)
    log.Printf("connections will only last %s\n", server.MaxTimeout)
    log.Printf("and timeout after %s of no activity\n", server.IdleTimeout)

    return svc, nil
}

type sshUser struct {
    name  string
    level ExecLevel
    keys  map[string]bool
}

type SshService interface {
    Close()
    Spawn()
    RegisterUser(user string, level ExecLevel, keys []string)
    LookupUser(username string) (user *sshUser, ok bool)

    AddCommand(cmds ...*Command)
}

func (svc *sshService) Close() {
    svc.s.Close()
}

func (svc *sshService) Spawn() {
    go svc.s.Serve(svc.ln)
}

func (svc *sshService) LookupUser(username string) (user *sshUser, ok bool) {
    user, ok = svc.users[username]
    return user, ok
}

func (svc *sshService) RegisterUser(user string, level ExecLevel, keys []string) {

    u := &sshUser{
        name:  user,
        level: level,
        keys:  make(map[string]bool),
    }

    for _, v := range keys {
        u.keys[v] = true
    }
    svc.users[user] = u
}

func (svc *sshService) AddCommand(cmds ...*Command) {
    DefaultCommands.AddCommand(cmds...)
}
