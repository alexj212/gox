package utilx

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// LoopForever on signal processing
func LoopForever(onShutdownFunc func(os.Signal)) {
	log.Printf("Entering infinite loop\n")

	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	sig := <-osSignal

	log.Printf("Exiting infinite loop received OsSignal\n")

	if onShutdownFunc != nil {
		onShutdownFunc = DefaultShutdown
	}

	onShutdownFunc(sig)
}
