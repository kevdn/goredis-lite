package cmd

import "goredis-lite/internal/server"

func main() {
	server.RunIoMultiplexingServer()
}