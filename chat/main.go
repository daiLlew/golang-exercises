package main

import (
	"flag"
	"fmt"
	"github.com/daiLlew/golang-exercises/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

const (
	githubRedirectURL = "http://localhost:8080/auth/callback/github"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	var clientId = flag.String("clientId", "", "Client ID to use for github OAuth")
	var clientSecret = flag.String("clientSecret", "", "Client Secret to use for github OAuth")
	var securityKey = flag.String("securityKey", "", "")
	flag.Parse()

	gomniauth.SetSecurityKey(*securityKey)
	gomniauth.WithProviders(github.New(*clientId, *clientSecret, githubRedirectURL))

	fmt.Printf("Starting server port=%s\n", *addr)
	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.Run()
	// start the web server.
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServeError:", err)
	}
}
