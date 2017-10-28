package main

import (
	"github.com/joho/godotenv"
	"github.com/melzareix/MrMeseeksBot/Server"
	"github.com/melzareix/MrMeseeksBot/Database"
)

func main() {
	godotenv.Load()
	Database.Init()
	Server.StartServer()
}
