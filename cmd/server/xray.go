package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"

	"github.com/DVKunion/SeaMoon/plugins/xray"
)

// Serve do common serve
func Serve(ctx context.Context, url string) error {
	server, err := xray.StartServer()
	if err != nil {
		fmt.Println("Failed to start:", err)
		// Configuration error. Exit with a special value to prevent systemd from restarting.
		os.Exit(23)
	}
	if err := server.Start(); err != nil {
		fmt.Println("Failed to start:", err)
		os.Exit(-1)
	}
	defer server.Close()
	runtime.GC()
	debug.FreeOSMemory()

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
		<-osSignals
	}
	return nil
}
