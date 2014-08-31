package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	//"encoding/json"
)

var format_style = "LLVM"

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func set_style(s string){
	format_style = s
}

func get_style(s string) string {
	switch s{
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

	set_style(string(body))

	fmt.Fprintf(w, string("done"))
}

func format(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)

	body, err := ioutil.ReadAll(r.Body)
	check(err)

	err = ioutil.WriteFile("./tmp/tmp.txt", body, 0644)
	check(err)

	s := get_style(format_style)
	fmt.Println("current format-style is", s)

	out, err := exec.Command("clang-format", s, "./tmp/tmp.txt").Output()
	check(err)

	fmt.Fprintf(w, string(out))

	err = os.Remove("./tmp/tmp.txt")
	check(err)
}

func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	http.HandleFunc("/format", format)
	http.HandleFunc("/style", style)

	err := http.ListenAndServe(":8080", nil)
	check(err)
}
