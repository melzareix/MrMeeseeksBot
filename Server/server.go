package Server

import (
	"net/http"
	"os"
)

func StartServer()  {
	mux := http.NewServeMux();
	mux.HandleFunc("/welcome", nil)
	mux.HandleFunc("/chat", nil)
	mux.HandleFunc("/", nil)
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(":" + port, mux)
}