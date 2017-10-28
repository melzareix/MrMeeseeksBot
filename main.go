package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/melzareix/MrMeseeksBot/Api/Welcome"
	"github.com/melzareix/MrMeseeksBot/Database"
	"net/http"
	"os"
)

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/welcome", Welcome.Handler)
	mux.HandleFunc("/chat", nil)
	mux.HandleFunc("/", nil)
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}
	fmt.Println("Server started at http://localhost:" + port)
	http.ListenAndServe(":"+port, mux)
}

func main() {
	godotenv.Load()
	Database.Init()
	StartServer()
}
