package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("%s %s \n %s\n", r.Method, r.RequestURI, string(body))
		fmt.Fprintf(w, "%s %s \n %s", r.Method, r.RequestURI, string(body))
	})

	http.ListenAndServe(":8889", nil)
}
