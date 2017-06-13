package main

import (
	"fmt"
	"github.com/gernest/mention"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"context"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func SaveBookmarkHandler(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	raw_body := string(body[:])
	w.Write(body[:])

	f, _ := os.OpenFile("log.out", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	log.SetOutput(f)
	log.Println("msg: ", raw_body)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	err, title, tags, url := msgParser(raw_body)
	if err != nil {
		fmt.Printf("parse raw msg fail: %v\n", err)
		return
	}

	input := &github.IssueRequest{
		Title:  &title,
		Body:   &url,
		Labels: &tags,
	}

	_, _, errs := client.Issues.Create(ctx, user, repo, input)
	if errs != nil {
		fmt.Printf("Issues.Create returned error: %v\n", err)
		return
	}
}

/*
	Incoming message format
		<tweet>|<url>

	Return <title>, <tags>, <url>
*/
func msgParser(msg string) (error, string, []string, string) {
	splited_arr := strings.Split(msg, "|")
	if len(splited_arr) != 2 {
		return errors.New("Invalid msg format, length of splited array is not 2."), "", []string{}, ""
	}
	tweet, url := splited_arr[0], splited_arr[1]
	raw_tags := mention.GetTags('#', strings.NewReader(tweet))
	tags := remove(raw_tags, "bm")
	title := strings.SplitN(tweet, "#", 2)[0]

	return nil, title, tags, url
}

func remove(s []string, key string) []string {
	for i, v := range s {
		if v == key {
			s = append(s[:i], s[i+1:]...)
			break
		}
	}
	return s
}

func main() {
	fmt.Println("Start server")
	r := mux.NewRouter()
	r.HandleFunc("/save", SaveBookmarkHandler).Methods("POST")
	http.ListenAndServe(":"+PORT, r)
}
