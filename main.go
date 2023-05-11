package main

import (
	"github.com/fadilmuh22/restskuy/cmd"
	"github.com/fadilmuh22/restskuy/cmd/db"
	"github.com/fadilmuh22/restskuy/config"
)

func main() {
	con := db.Connect()

	config.Init()
	cmd.StartServer(con)
}
