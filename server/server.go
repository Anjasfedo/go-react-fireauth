package server

import (
	// "github.com/Anjasfedo/go-react-fireauth/configs"
)

func Init() {
	r := NewRouter()
	
	r.Run()
}