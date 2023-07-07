package main

import (
	"os"

	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/server"
)

func main() {
	if consts.Version == "dev" {
		server.NewServer("socks5", "0.0.0.0", "8888").Serve()
	} else {
		server.NewServer(os.Getenv("serverMod"), "0.0.0.0", "9000").Serve()
	}
}
