package main

import (
	"context"
	"github.com/Kolterdyx/meetroom/internal"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	http.HandleFunc("/", internal.ControllerHandler)
	http.HandleFunc("/new", internal.NewMeetingHandler)
	http.HandleFunc("/join", internal.JoinMeetingHandler)
	http.HandleFunc("/end", internal.EndMeetingHandler)
	http.HandleFunc("/idle", internal.IdleHandler)

	http.Handle("/static/", internal.StaticHandler())

	log.Println("Meet room control listening on :5000")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go openChromium(ctx)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func openChromium(ctx context.Context) {
	cmd := exec.CommandContext(ctx,
		"chromium",
		//"--kiosk",
		"--noerrdialogs",
		"--disable-infobars",
		"--use-fake-ui-for-media-stream",
		"--disable-features=TranslateUI",
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
