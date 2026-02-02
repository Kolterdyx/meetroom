package internal

import (
	"embed"
	"html/template"
	"net/http"
)

var templates map[string]*template.Template

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

func init() {

	templates = make(map[string]*template.Template)
	templates["controller"] = template.Must(template.ParseFS(templateFS,
		"templates/layout.html",
		"templates/controller.html",
	))
	templates["idle"] = template.Must(template.ParseFS(templateFS,
		"templates/layout.html",
		"templates/idle.html",
	))
}

func render(w http.ResponseWriter, name string, data any) {
	err := templates[name].ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ControllerHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "controller", nil)
}

func IdleHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "idle", nil)
}

func NewMeetingHandler(w http.ResponseWriter, r *http.Request) {

	if err := OpenURL("https://meet.google.com/new"); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := CloseIdleTabs(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func JoinMeetingHandler(w http.ResponseWriter, r *http.Request) {
	link := r.FormValue("link")
	if link == "" {
		http.Error(w, "missing link parameter", 400)
		return
	}
	if err := OpenURL(link); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := CloseIdleTabs(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EndMeetingHandler(w http.ResponseWriter, r *http.Request) {
	if err := CloseMeetTabs(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := OpenURL("http://localhost:5000/idle"); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func StaticHandler() http.Handler {
	return http.FileServer(http.FS(staticFS))
}
