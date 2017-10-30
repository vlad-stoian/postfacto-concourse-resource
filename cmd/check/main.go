package main

import (
	"encoding/json"
	"fmt"
	"os"

	"bytes"
	"net/http"
	"net/http/httputil"

	"io/ioutil"

	"sort"

	"github.com/mitchellh/colorstring"
)

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type CheckResponse []Version

type Version struct {
	ActionItemDate string `json:"action_item_date,omitempty"`
	ActionItemID   string `json:"action_item_id,omitempty"`
}

type Source struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func Fatal(doing string, err error) {
	Sayf(colorstring.Color("[red]error %s: %s\n"), doing, err)
	os.Exit(1)
}

func Sayf(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, message, args...)
}

func GetToken(retroID, retroPassword string) string {
	url := fmt.Sprintf("https://retro-api.cfapps.io/retros/%v/login", retroID)

	buffer := bytes.NewBufferString(fmt.Sprintf(`{ "retro": { "password": "%s"} }`, retroPassword))

	req, _ := http.NewRequest("PUT", url, buffer)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, _ := http.DefaultClient.Do(req)

	if resp.StatusCode != 200 {
		response, _ := httputil.DumpResponse(resp, true)
		fmt.Fprint(os.Stderr, string(response))
	}

	serverResponse, _ := ioutil.ReadAll(resp.Body)

	var lResponse loginResponse
	json.Unmarshal(serverResponse, &lResponse)

	bearerToken := fmt.Sprintf(`Bearer %v`, lResponse.Token)

	return bearerToken
}

func GetRetroBoard(retroID, token string) RetroBoard {
	url := fmt.Sprintf("https://retro-api.cfapps.io/retros/%v/", retroID)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", token)

	resp, _ := http.DefaultClient.Do(req)

	if resp.StatusCode != 200 {
		response, _ := httputil.DumpResponse(resp, true)
		fmt.Fprint(os.Stderr, string(response))
	}

	serverResponse, _ := ioutil.ReadAll(resp.Body)

	var retroBoard RetroBoard
	_ = json.Unmarshal(serverResponse, &retroBoard)

	return retroBoard
}

type RetroBoard struct {
	Board struct {
		Slug        string       `json:"slug"`
		ActionItems []ActionItem `json:"action_items"`
		RetroItems  []RetroItem  `json:"items"`
	} `json:"retro"`
}

type ActionItem struct {
	Description string `json:"description"`
	ID          uint64 `json:"id"`
	Done        bool   `json:"done"`
	CreatedAt   string `json:"created_at"`
}

type RetroItem struct {
	Description string `json:"description"`
	Category    string `json:"category"`
	Done        bool   `json:"done,omitempty"`
	ID          uint64 `json:"id"`
}

func main() {
	var request CheckRequest

	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		Fatal("reading request from STDIN", err)
	}

	token := GetToken(request.Source.ID, request.Source.Password)
	retro := GetRetroBoard(request.Source.ID, token)

	actionItems := retro.Board.ActionItems

	var response CheckResponse

	sort.Slice(actionItems, func(i, j int) bool {
		return actionItems[i].CreatedAt < actionItems[j].CreatedAt
	})

	if request.Version.ActionItemDate != "" {
		for _, actionItem := range actionItems {
			if actionItem.CreatedAt < request.Version.ActionItemDate {
				continue
			}

			response = append(response, Version{
				ActionItemDate: actionItem.CreatedAt,
				ActionItemID:   fmt.Sprintf("%d", actionItem.ID),
			})

		}
	} else {
		lastActionItem := retro.Board.ActionItems[len(retro.Board.ActionItems)-1]
		response = append(response, Version{
			ActionItemDate: lastActionItem.CreatedAt,
			ActionItemID:   fmt.Sprintf("%d", lastActionItem.ID),
		})
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		Fatal("writing response to STDOUT", err)
	}

}
