package commandr

import (
    "fmt"
    "github.com/alexj212/gox"
    "github.com/alexj212/gox/term"
    "github.com/fatih/color"
    "github.com/gliderlabs/ssh"
    "github.com/go-errors/errors"
    gossh "golang.org/x/crypto/ssh"
    "io"
    "log"
    "net"
    "time"
)

type sshClient struct {
    s         ssh.Session
    user      *sshUser
    activeKey *gox.SshKey
}

//Close interface func implementation to close client down
func (s *sshClient) Close() {
    _ = s.s.Close()
}

//ExecLevel interface func implementation to return client exec level
func (s *sshClient) ExecLevel() ExecLevel {
    return s.user.level
}

//UserName interface func implementation to return client user name
func (s *sshClient) UserName() string {
    return s.s.User()
}

//History interface func implementation to return client command history
func (s *sshClient) History() []string {
    return s.user.history
}

//Write interface func implementation to write to clients stream
func (s *sshClient) Write(p []byte) (n int, err error) {
    return s.s.Write(p)
}

//WriteString interface func implementation to write string to clients stream
func (s *sshClient) WriteString(p string) {
    _, _ = s.Write([]byte(p))
}

// SshClient interface of ssh client
type SshClient interface {
    //Close interface func implementation to close client down
    Close()
    //UserName interface func implementation to return client user name
    UserName() string
    //ExecLevel interface func implementation to return client exec level
    ExecLevel() ExecLevel
    //History interface func implementation to return client command history
    History() []string
    //Write interface func implementation to write to clients stream
    Write(p []byte) (n int, err error)
    // WriteString sends text back to client
    WriteString(p string)
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
        s:    s,
        user: user,
    }

    //ptyReq, winCh, isPty := s.Pty()
    //io.WriteString(s, fmt.Sprintf("ptyReq: %v\n", ptyReq))
    //io.WriteString(s, fmt.Sprintf("winCh: %v\n", winCh))
    //io.WriteString(s, fmt.Sprintf("isPty: %v\n", isPty))
    //io.WriteString(s, fmt.Sprintf("s.User(): %v\n", s.User()))
    //

    authorizedKey := gossh.MarshalAuthorizedKey(s.PublicKey())
    c.WriteString(fmt.Sprintf("public key used by %s    key: %s\n", s.User(), string(authorizedKey)))

    activeKey, ok := user.keys[string(authorizedKey)]
    if ok {
        c.activeKey = activeKey
        c.WriteString(fmt.Sprintf("Active Key KeyType: %v Comment: %v\n", activeKey.KeyType, activeKey.Comment))
    }

    term := term.NewTerminal(s, "> ")

    for _, v := range user.history {
        term.AddHistory(v)
    }

    for {
        line, err := term.ReadLine()
        if err != nil {
            fmt.Printf("SSH ReadLine err: %v\n", err)
            break
        }

        var allowExec error
        if svc.preExecHandler != nil {
            allowExec = svc.preExecHandler(svc, c, line)
        }

        if allowExec != nil {
            msg := fmt.Sprintf("Error Pre Exec Handler disabling exec: %v\n\n", allowExec)
            term.Write([]byte(msg))
            continue
        }
        parsed, err := NewCommandArgs(line, s)
        if err != nil {
            continue
        }

        execErr := DefaultCommands.Execute(c, parsed)
        if svc.postExecHandler != nil {
            svc.postExecHandler(svc, c, line, execErr)
        }

        if execErr != nil {
            term.Write([]byte(fmt.Sprintf("Error: %v\n", execErr)))
            msg := fmt.Sprintf("Error Stack\n%v\n", errors.Wrap(execErr, 2).ErrorStack())
            term.Write([]byte(msg))
            continue
        }

        user.history = append(user.history, line)
    }
}

// ClientDecorator - func def for client initializer
type ClientDecorator func(*sshService)

type PreExecHandler func(SshService, SshClient, string) error
type PostExecHandler func(SshService, SshClient, string, error)

type sshService struct {
    s               *ssh.Server
    ln              net.Listener
    users           map[string]*sshUser
    preExecHandler  PreExecHandler
    postExecHandler PostExecHandler
}

// SetPreExecHandler - set pre exec handler
func SetPreExecHandler(val PreExecHandler) ClientDecorator {
    return func(l *sshService) {
        l.preExecHandler = val
    }
}

// SetPostExecHandler - set post exec handler
func SetPostExecHandler(val PostExecHandler) ClientDecorator {
    return func(l *sshService) {
        l.postExecHandler = val
    }
}

// NewSshService create a new instance of ssh service
func NewSshService(port int, hostKey gossh.Signer, decorators ...ClientDecorator) (SshService, error) {

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

    for _, decorator := range decorators {
        decorator(svc)
    }

    log.Printf("starting ssh server on port %s...\n", server.Addr)
    log.Printf("connections will only last %s\n", server.MaxTimeout)
    log.Printf("and timeout after %s of no activity\n", server.IdleTimeout)

    return svc, nil
}

type sshUser struct {
    name    string
    level   ExecLevel
    keys    map[string]*gox.SshKey
    history []string
}

// SshService interface of ssh service
type SshService interface {
    // Close shut down ssh service
    Close() error
    // Spawn start new go routine serving ssh
    Spawn()
    // RegisterUser register a user on the system
    RegisterUser(user string, level ExecLevel, keys []*gox.SshKey, history []string)
    // LookupUser lookup a user by name
    LookupUser(username string) (user *sshUser, ok bool)
    // AddCommand add commands to be executed
    AddCommand(cmds ...*Command)
}

// Close shut down ssh service
func (svc *sshService) Close() error {
    return svc.s.Close()
}

// Spawn start new go routine serving ssh
func (svc *sshService) Spawn() {
    go svc.s.Serve(svc.ln)
}

// LookupUser lookup a user by name
func (svc *sshService) LookupUser(username string) (user *sshUser, ok bool) {
    user, ok = svc.users[username]
    return user, ok
}

// RegisterUser register a user on the system
func (svc *sshService) RegisterUser(user string, level ExecLevel, keys []*gox.SshKey, history []string) {

    u := &sshUser{
        name:    user,
        level:   level,
        keys:    make(map[string]*gox.SshKey),
        history: history,
    }

    if u.history == nil {
        u.history = make([]string, 0)
    }

    for _, v := range keys {
        u.keys[string(v.Key)] = v
    }
    svc.users[user] = u
}

// AddCommand add commands to be executed
func (svc *sshService) AddCommand(cmds ...*Command) {
    DefaultCommands.AddCommand(cmds...)
}
