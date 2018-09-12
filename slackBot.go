/*

mybot - Illustrative Slack bot in Go

Copyright (c) 2015 RapidLoop

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type DogResp struct {
	Success string
	Message string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: mybot slack-bot-token\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	ws, id := slackConnect(os.Args[1])
	fmt.Println("Lexie is ready to play, ^C exits")

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// if there is an error returned, report it and skip this loop.
		if m.Type == "error" {
			log.Println(fmt.Errorf("error from Slack on getMessage: code %v\t%s\n", m.Error.Code, m.Error.Message))
			continue
		}

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// if so try to parse if
			parts := strings.Fields(m.Text)
			if len(parts) == 2 && parts[1] == "dogs" {
				go func(m Message) {
					m.Text = getDog()
					postMessage(ws, m)
				}(m)
			} else if len(parts) == 3 && parts[1] == "issue" && parts[2] == "count" {
				go func(m Message) {
					m.Text = issueCredits()
					postMessage(ws, m)
				}(m)
			} else if len(parts) == 2 && parts[1] == "pet" {
				go func(m Message) {
					m.Text = "tail wags"
					postMessage(ws, m)
				}(m)
			} else if len(parts) == 2 && parts[1] == "help" {
				go func(m Message) {
					// String literals with back tics capture line breaks. Neato!
					m.Text = `Commands:
@lexie *issue count*: Displays the number of issue credits for MG.
@lexie *dogs*: Fetches awesome dog images.
@lexie *pet*: love you friendly bot.`
					postMessage(ws, m)
				}(m)
			} else {
				// huh?
				m.Text = fmt.Sprintf("sorry, that does not compute\n")
				postMessage(ws, m)
			}
		}
	}
}

func issueCredits() string {
	resp, err := http.Get("https://www.drupal.org/api-d7/node/2497975.json")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	input, err := ioutil.ReadAll(resp.Body)
	var org map[string]interface{}
	json.Unmarshal(input, &org)
	return org["title"].(string) + " issue count is " + org["field_org_issue_credit_count"].(string)
}
