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

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func format(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)

	body, err := ioutil.ReadAll(r.Body)
	check(err)

	err = ioutil.WriteFile("./tmp/tmp.txt", body, 0644)
	check(err)

	out, err := exec.Command("clang-format", "./tmp/tmp.txt").Output()
	check(err)

	fmt.Fprintf(w, string(out))

	err = os.Remove("./tmp/tmp.txt")
	check(err)
}

func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	http.HandleFunc("/api", format)

	err := http.ListenAndServe(":8080", nil)
	check(err)
}
