package main

import (
	"github.com/joho/godotenv"
	"github.com/melzareix/MrMeseeksBot/Server"
)

func main() {
	godotenv.Load()
	Server.StartServer()
}
