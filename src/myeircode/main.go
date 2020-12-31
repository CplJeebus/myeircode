package main

import (
	"fmt"
	. "myeircode/utils"
	"net/http"
)

var c Config

func main() {
	c.LoadConfig()
	http.HandleFunc("/mailtest", func(w http.ResponseWriter, r *http.Request) { SendMail(w, r, c.Admin, c.MailKey) })
	http.HandleFunc("/", ShowCodes)
	http.HandleFunc("/api", ShowJson)
	http.ListenAndServe(":8080", nil)
}

func ShowCodes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, Pretty)
}

func ShowJson(w http.ResponseWriter, r *http.Request) {
	b, err := DownloadFile(c.Bucket, "codes.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintf(w, string(b))
}
