package main

import (
	"github.com/Anjasfedo/go-react-fireauth/configs"
	"github.com/Anjasfedo/go-react-fireauth/server"
)

func main() {
	configs.InitFirebase()

	server.Init()
}

