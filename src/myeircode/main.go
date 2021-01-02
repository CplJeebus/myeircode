package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "myeircode/utils"
	"net/http"
	"strings"

	uuid "github.com/nu7hatch/gouuid"
)

var c Config

func main() {
	c.LoadConfig()
	http.HandleFunc("/", CheckTLS(ShowCodes))
	http.HandleFunc("/api", ShowJson)
	http.HandleFunc("/new", AddCode)
	http.HandleFunc("/auth", Auth)
	http.ListenAndServe(":8080", nil)
}

func CheckTLS(function http.HandlerFunc) http.HandlerFunc {
	f := func(w http.ResponseWriter, req *http.Request) {
		parts := strings.Split(req.URL.String(), ":")
		if parts[0] == "http" {
			http.Redirect(w, req, "https:"+parts[1], 301)
		} else {
			function(w, req)
		}
	}

	return f
}

func ShowCodes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, Pretty)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	CurrentBytes, _ := DownloadFile(c.Bucket, "codes.json")
	var CurrentCodes []Code
	var Staged Code

	e := json.Unmarshal(CurrentBytes, &CurrentCodes)
	if e != nil {
		fmt.Println(e)
	}

	q := r.URL.Query()
	u := q.Get("id")

	StagedBytes, _ := ioutil.ReadFile(u + ".json")
	e = json.Unmarshal(StagedBytes, &Staged)
	CurrentCodes = append(CurrentCodes, Staged)
	out, _ := json.Marshal(CurrentCodes)
	SaveCodes(c.Bucket, "codes.json", out)

}

func AddCode(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println(fn)
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
