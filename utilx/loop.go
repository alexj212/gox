package utilx

import (
    "fmt"
    "os"
)

var osSignal chan os.Signal
var onShutdownFunc func(os.Signal)

func init() {
    osSignal = make(chan os.Signal, 1)

}

// DefaultShutdown default func for shutdown func for LoopForever
func DefaultShutdown(sig os.Signal) {
    fmt.Printf("caught sig: %v\n\n", sig)
    os.Exit(0)
}
