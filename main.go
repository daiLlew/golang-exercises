package main

import (
	"fmt"
	"github.com/daiLlew/golang-exercises/chat"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling request.")
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	fmt.Println("Starting server...")
	r := chat.NewRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	go r.Run()
	// start the web server.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServeError:", err)
	}
}
