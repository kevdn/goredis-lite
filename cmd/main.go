package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"goredis-lite/internal/server"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	var wg sync.WaitGroup
	wg.Add(2)

	go server.RunIoMultiplexingServer(&wg)
	go server.WaitForSignal(&wg, signals)
	wg.Wait()
}
