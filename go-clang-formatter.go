package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type reqStyle struct {
	Style string `json:"style"`
	UUID  string `json:"uuid"`
}

type reqFormat struct {
	Code []byte `json:"code"`
	UUID string `json:"uuid"`
}

func get_style(s string) string {
	switch s {
	case "LLVM":
		return "-style=LLVM"
	case "Google":
		return "-style=Google"
	case "Chromium":
		return "-style=Chromium"
	case "Mozilla":
		return "-style=Mozilla"
	case "WebKit":
		return "-style=WebKit"
	default:
		return "-style=LLVM"
	}
}

func style(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)

	body, err := ioutil.ReadAll(r.Body)
	check(err)
	fmt.Println(string(body))

	var s reqStyle
	err = json.Unmarshal(body, &s)
	check(err)

	err = ioutil.WriteFile("./tmp/tmp-style."+s.UUID+".txt", []byte(get_style(s.Style)), 0644)
	check(err)

	fmt.Fprintf(w, "ok")
}

func format(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)

	body, err := ioutil.ReadAll(r.Body)
	check(err)

	var f reqFormat
	err = json.Unmarshal(body, &f)
	check(err)

	err = ioutil.WriteFile("./tmp/tmp-code."+f.UUID+"txt", f.Code, 0644)
	check(err)

	var style string
	_, err = os.Stat("./tmp/tmp-style." + f.UUID + ".txt")
	if err != nil {
		style = "-style=LLVM"
	} else {
		s, err := ioutil.ReadFile("./tmp/tmp-style." + f.UUID + ".txt")
		check(err)
		style = string(s)
		err = os.Remove("./tmp/tmp-style." + f.UUID + ".txt")
		check(err)
	}

	fmt.Println("Request format-style is", style)

	stdout, err := exec.Command("clang-format", style, "./tmp/tmp-code."+f.UUID+"txt").Output()
	check(err)

	enc := base64.StdEncoding.EncodeToString(stdout)
	fmt.Fprintf(w, enc)

	err = os.Remove("./tmp/tmp-code." + f.UUID + "txt")
	check(err)
}

func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	http.HandleFunc("/format", format)
	http.HandleFunc("/style", style)

	err := http.ListenAndServe(":8080", nil)
	check(err)
}
