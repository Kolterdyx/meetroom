package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const debugPort = "9222"

type ChromeTab struct {
	ID                   string `json:"id"`
	URL                  string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

func OpenURL(u string) error {
	var client http.Client
	parsedURL, err := url.Parse(fmt.Sprintf("http://localhost:9222/json/new?%s", u))
	if err != nil {
		return err
	}
	resp, err := client.Do(&http.Request{
		Method: "PUT",
		URL:    parsedURL,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to open URL: %s", resp.Status)
	}
	return nil
}

func ListTabs() ([]ChromeTab, error) {
	resp, err := http.Get(
		fmt.Sprintf("http://localhost:%s/json", debugPort),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tabs []ChromeTab
	err = json.NewDecoder(resp.Body).Decode(&tabs)
	return tabs, err
}

func CloseMeetTabs() error {
	tabs, err := ListTabs()
	if err != nil {
		return err
	}

	for _, tab := range tabs {
		if strings.Contains(tab.URL, "meet.google.com") {
			http.Get(
				fmt.Sprintf("http://localhost:%s/json/close/%s", debugPort, tab.ID),
			)
		}
	}
	return nil
}

func CloseIdleTabs() error {
	tabs, err := ListTabs()
	if err != nil {
		return err
	}

	for _, tab := range tabs {
		if strings.Contains(tab.URL, "localhost:5000/idle") {
			http.Get(
				fmt.Sprintf("http://localhost:%s/json/close/%s", debugPort, tab.ID),
			)
		}
	}
	return nil
}

func FindMeetTab() (*ChromeTab, error) {
	resp, err := http.Get("http://localhost:9222/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tabs []ChromeTab
	if err := json.NewDecoder(resp.Body).Decode(&tabs); err != nil {
		return nil, err
	}

	for _, tab := range tabs {
		if strings.Contains(tab.URL, "meet.google.com") {
			return &tab, nil
		}
	}
	return nil, fmt.Errorf("no Meet tab found")
}

func InjectJoinClick(wsURL string) error {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	msg := map[string]interface{}{
		"id":     1,
		"method": "Runtime.evaluate",
		"params": map[string]interface{}{
			"expression": `
(() => {
  const buttons = [...document.querySelectorAll("button")];
  const join = buttons.find(b => b.innerText.match(/(join|unir)/i));
  if (join) join.click();
})();
`,
		},
	}

	return conn.WriteJSON(msg)
}
