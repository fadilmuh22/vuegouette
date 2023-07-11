package main

import (
	"github.com/fadilmuh22/restskuy/config"
	"github.com/fadilmuh22/restskuy/internal"
)

func main() {
	config.Init()
	internal.StartServer()
}
