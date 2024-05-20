package main

import (
	"github.com/Tanmai2002/takago/redirect_service/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	server.InitializeServer()
}
