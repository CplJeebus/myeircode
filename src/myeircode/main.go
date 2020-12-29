package main

import (
	"fmt"
	. "myeircode/utils"
	"net/http"
)

var c Config

func main() {
	c.LoadConfig()
	http.HandleFunc("/", ShowCodes)
	http.ListenAndServe(":8080", nil)
}

func ShowCodes(w http.ResponseWriter, r *http.Request) {
	b, err := DownloadFile(c.Bucket)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintf(w, string(b))
}
