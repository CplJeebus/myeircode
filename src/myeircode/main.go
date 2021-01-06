package main

import (
	h "myeircode/handlers"
	. "myeircode/utils"
	"net/http"
)

var c Config

func main() {
	c.LoadConfig()
	http.HandleFunc("/", h.Challenge(h.ShowCodes))
	http.HandleFunc("/api", h.ShowJSON)
	http.HandleFunc("/new", h.AddCode)
	http.HandleFunc("/auth", h.Auth)
	http.ListenAndServe(":8080", nil)
}
