package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	req, _ := http.NewRequest(
		"GET",
		"http://127.0.0.1:34903/metrics",
		nil)
	res, _ := (&http.Client{}).Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("get message error: %v", err.Error())
	} else {
		fmt.Println(string(body))
	}
}
