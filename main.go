package main

import (
	"github.com/fadilmuh22/restskuy/cmd"
	"github.com/fadilmuh22/restskuy/config"
)

func main() {
	config.Init()
	cmd.StartServer()
}
