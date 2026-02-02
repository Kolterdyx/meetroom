package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go openChromium(ctx)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func openChromium(ctx context.Context) {
	cmd := exec.CommandContext(ctx,
		"chromium",
		"--kiosk",
		"--noerrdialogs",
		"--disable-infobars",
		"--disable-session-crashed-bubble",
		"--autoplay-policy=no-user-gesture-required",
		"--usenable-features=VaapiVideoDecoder",
		"--remote-debugging-port=9222",
		"http://localhost:5000/idle",
	)
	cmd.Env = append(os.Environ(), "DISPLAY=:0")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start chromium: %v", err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Chromium exited with error: %v", err)
	}
	log.Println("Chromium exited")
}

func render(w http.ResponseWriter, name string, data any) {
	err := templates[name].ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
