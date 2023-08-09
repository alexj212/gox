package utilx

import (
	"github.com/potakhov/loge"
	"os"
	"os/signal"
	"syscall"
)

// LoopForever on signal processing
func LoopForever(onShutdownFunc func(os.Signal)) {
	loge.Info("Entering infinite loop\n")

	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	sig := <-osSignal

	loge.Info("Exiting infinite loop received OsSignal\n")

	if onShutdownFunc != nil {
		onShutdownFunc = DefaultShutdown
	}

	onShutdownFunc(sig)
}
