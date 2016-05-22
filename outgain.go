package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

func main() {
	// Entry point for the application, set up the database, server here
	http.Handle("/", appHandler(hello))
	http.ListenAndServe(":8080", nil)
}

type httpError struct {
	Error   error
	Message string
	Code    int
}

type page struct {
	Title,
	Message string
}

type appHandler func(http.ResponseWriter, *http.Request) *httpError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		http.Error(w, e.Message, e.Code)
	}
}

func hello(w http.ResponseWriter, r *http.Request) *httpError {
	p := page{Title: "Hello", Message: "Hello, World!"}
	html, err := ioutil.ReadFile("templates/base.tmpl")
	if err != nil {
		return &httpError{Error: err, Message: "Page Not Found", Code: 404}
	}
	tmpl, err := template.New("Hello").Parse(string(html))
	if err != nil {
		return &httpError{Error: err, Message: "Internal Server Error", Code: 500}
	}
	if err = tmpl.Execute(w, p); err != nil {
		return &httpError{Error: err, Message: "Internal Server Error", Code: 500}
	}

	return nil
}
