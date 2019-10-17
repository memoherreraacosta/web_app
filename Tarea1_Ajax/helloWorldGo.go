package main

import (
	"io/ioutil"
	"net/http"
)

func main() {

	res, _ := http.Get("http://www.google.com")

	resp, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	_ = ioutil.WriteFile("body_html.txt", resp, 0644)
}
