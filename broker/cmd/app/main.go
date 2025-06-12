package main

import (
	"broker/internal/config"
	"broker/internal/server"
)

func main() {
	server := server.NewServer(
		*config.NewConfig(),
	)
	server.Start()
}
