package main

import (
	"net/http"
)

func controllerHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "controller", nil)
}

func idleHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "idle", nil)
}

func newMeetingHandler(w http.ResponseWriter, r *http.Request) {

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

func joinMeetingHandler(w http.ResponseWriter, r *http.Request) {
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

func endMeetingHandler(w http.ResponseWriter, r *http.Request) {
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
