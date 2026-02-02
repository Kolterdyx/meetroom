package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates map[string]*template.Template

func main() {
	templates = make(map[string]*template.Template)
	templates["controller"] = template.Must(template.ParseFiles(
		"templates/layout.html",
		"templates/controller.html",
	))
	templates["idle"] = template.Must(template.ParseFiles(
		"templates/layout.html",
		"templates/idle.html",
	))
	http.HandleFunc("/", controllerHandler)
	http.HandleFunc("/new", newMeetingHandler)
	http.HandleFunc("/join", joinMeetingHandler)
	http.HandleFunc("/end", endMeetingHandler)
	http.HandleFunc("/idle", idleHandler)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	log.Println("Meet room control listening on :5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func render(w http.ResponseWriter, name string, data any) {
	err := templates[name].ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
