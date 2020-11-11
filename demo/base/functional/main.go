package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

// HandlerFunc 定义函数类型作为变量类型
type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func main() {
	hf := HandlerFunc(handler)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", bytes.NewBuffer([]byte("test")))

	hf.ServeHTTP(resp, req)

	bts, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(bts))
}

func handler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("hello world!\n"))
}
