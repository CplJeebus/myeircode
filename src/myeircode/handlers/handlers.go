package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"myeircode/utils"
	. "myeircode/utils"
	"net/http"
	"os"

	uuid "github.com/nu7hatch/gouuid"
)

var c Config

func ShowCodes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, Pretty)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	c.LoadConfig()
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
	if e != nil {
		fmt.Printf("Can't Unmarshal json: %+v", e)
	}

	CurrentCodes = append(CurrentCodes, Staged)
	out, _ := json.Marshal(CurrentCodes)
	SaveCodes(c.Bucket, "codes.json", out)

	e = os.Remove(u + ".json")
	if e != nil {
		fmt.Printf("Can't clean up staged update %+v", e)
	}
}

func AddCode(w http.ResponseWriter, r *http.Request) {
	c.LoadConfig()

	switch r.Method {
	case "GET":
		fmt.Fprint(w, PrettyForm)
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
		if e != nil {
			fmt.Printf("Can't create a UUID for some god know reason %+v", e)
		}

		fn := u.String()
		fmt.Println(fn)

		e = ioutil.WriteFile(fn+".json", f, 0600)
		if e != nil {
			fmt.Printf("Unable to write updated file %+v", e)
		}

		SendMail(c, fn)

		fmt.Fprintf(w, Wait)
	default:
		fmt.Fprintf(w, "Method not supported!")
	}
}

func Challenge(function http.HandlerFunc) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		c.LoadConfig()

		switch r.Method {
		case "GET":
			q := r.URL.Query()
			cid := q.Get("cid")

			if cid != "" {
				_, e := ioutil.ReadFile(cid + ".tmp")
				if e != nil {
					fmt.Printf("Could not validate CID file %v", e)
				} else {
					function(w, r)
				}
			}
			fmt.Fprintf(w, utils.Challenge)
		case "POST":

			if r.FormValue("challenge") == "Major" {
				cid, e := uuid.NewV4()
				if e != nil {
					fmt.Printf("[CID] Can't create a UUID for some god know reason %+v", e)
				}

				e = ioutil.WriteFile(fmt.Sprint(cid)+".tmp", nil, 0600)
				if e != nil {
					fmt.Printf("Unable to write temp cid file %+v", e)
				}

				http.Redirect(w, r, "https://"+c.Host+"/?cid="+fmt.Sprint(cid), 301)
			}
		}
	}

	return f
}

func ShowJSON(w http.ResponseWriter, r *http.Request) {
	c.LoadConfig()

	b, err := DownloadFile(c.Bucket, "codes.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprint(w, string(b))
}
