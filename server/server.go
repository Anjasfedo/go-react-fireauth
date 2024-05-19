package server

import (
    "log"
)

func Init(port string) {
	r := NewRouter()

	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
