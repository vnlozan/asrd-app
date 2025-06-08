package main

import (
	"broker/config"
	"broker/internal/server"
)

func main() {
	server := server.NewServer(
		*config.NewConfig(),
	)
	server.Start()
}
