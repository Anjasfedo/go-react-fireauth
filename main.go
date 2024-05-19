package main

import (
	"context"

	"github.com/Anjasfedo/go-react-fireauth/configs"
	"github.com/Anjasfedo/go-react-fireauth/server"
)

func main() {
	configs.InitFirebase(context.Background())
	defer configs.CloseFirebase()

	server.Init(":8080")
}

