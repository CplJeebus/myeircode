package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "myeircode/utils"
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
)

var c Config

func main() {
	c.LoadConfig()
	http.HandleFunc("/", ShowCodes)
	http.HandleFunc("/api", ShowJson)
	http.HandleFunc("/new", NewCode)
	http.ListenAndServe(":8080", nil)
}

func ShowCodes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, Pretty)
}

func NewCode(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, PrettyForm)
	case "POST":
		var code = Code{
			Name: r.FormValue("name"),
			Code: r.FormValue("code"),
		}

		f, e := json.Marshal(code)
		if e != nil {
			fmt.Println(e)
		}
		u, e := uuid.NewV4()
		fn := u.String()
		e = ioutil.WriteFile(fn+".json", f, 0644)
		SendMail(c, fn)

		fmt.Fprintf(w, "Wait for authorisation")
	default:
		fmt.Fprintf(w, "Method not supported!")
	}
}

func ShowJson(w http.ResponseWriter, r *http.Request) {
	b, err := DownloadFile(c.Bucket, "codes.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintf(w, string(b))
}
