package main

import (
	"fmt"
	"net/http"
	"html/template"
	"log"
)


type HttpServer struct {
	service      LinkShortenService
	listenPort   string
	baseUrl string
	shortenForm  *template.Template
	viewShorten *template.Template
}

func createServer(storage string, port string, baseUrl string) HttpServer {
	s := HttpServer{}
	s.service = LinkShortenService{getStorage(storage)}
	s.listenPort = port
	s.baseUrl = baseUrl
	s.shortenForm = loadTemplate("shorten-form.html")
	s.viewShorten = loadTemplate("view-shorten.html")
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", s.redirectHandler)
	http.HandleFunc("/shorten", s.shortenPageHandler)
	http.HandleFunc("/status/", s.statusPageHandler)
	return s
}

func (s *HttpServer) Start() {
	err := http.ListenAndServe(":"+s.listenPort, nil)
	if err != nil {
		fmt.Println("Could not start server", err)
	}
}

func loadTemplate(name string) *template.Template {
	t, err := template.ParseFiles("src/server/templates/" + name)
	if err != nil {
		panic(fmt.Sprintln("Template ", name, " not found. ", err))
	}
	return t
}

func getStorage(name string) LinkStorage {
	switch(name) {
	case "memory":
		return NewMemoryLinkStorage()
	default:
		panic("No such storage: " + name)
	}
}

func (s *HttpServer) getExternalUrl(slug string) string {
	return s.baseUrl + "/" + slug
}

func (s *HttpServer) shortenPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		s.shortenForm.Execute(w, struct{}{})
	} else if r.Method == "POST" {
		url := r.FormValue("url")
		slug := s.service.CreateSlug(url)
		http.Redirect(w, r, "/status/"+slug, http.StatusFound)
	}
}

type StatusPage struct {
	ShortUrl  string
	TargetUrl string
}

func (s *HttpServer) statusPageHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[len("/status/"):]
	url := s.service.ExpandSlug(slug)
	statusPage := StatusPage{s.getExternalUrl(slug), url}
	s.viewShorten.Execute(w, statusPage)
}

func (s *HttpServer) redirectHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[len("/"):]

	if slug == "" {
		http.Redirect(w, r, "/shorten", 302)
		return
	}

	url := s.service.ExpandSlug(slug)

	if url == "" {
		http.Redirect(w, r, "/status/"+slug, 302)
		log.Printf("Slug [%s] is not known", slug)
		return
	}

	http.Redirect(w, r, url, 302)
	log.Printf("Redirected slug [%s] to url [%s]", slug, url)
}
