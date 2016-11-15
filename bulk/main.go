package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"gopkg.in/mgo.v2"

	"encoding/json"
)

type results struct {
	Meta struct {
		Next string `json:"next_results"`
	} `json:"search_metadata"`
	Statuses []Status `json:"statuses"`
}

func main() {
	client := NewClient(
		"32FiBBFYuEb7c9S2K5tTBGddb",
		"oCIjH791Bg8zkAWzXKFRoZhzeFe2UDtNWZIvNUqDMZODIzoxny",
		"27555535-JkZ8yREgZEnLddst2v9ze5v0LO5eRSn9iurOvA5xw",
		"CYX0Kwe8j4E57XemzbisiWhaOS0Y2Uoq5jp564Od7b1sP",
	)

	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	tweets := sess.DB("bulk_tweets").C("tweets")

	var query = "q=" + url.QueryEscape(os.Args[1])
	var r results
	for {
		r, err = getTweets(client, query)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if len(r.Meta.Next) == 0 {
			fmt.Println("done")
			return
		}
		query = r.Meta.Next[1:]

		for _, s := range r.Statuses {
			tweets.Insert(s)
		}

		if len(r.Statuses) == 0 {
			break
		}
	}

}

func getTweets(client *Client, query string) (results, error) {
	vals, err := url.ParseQuery(query)
	if err != nil {
		return results{}, err
	}

	var args = make(map[string]string)
	for key, val := range vals {
		args[key] = strings.Join(val, "+")
	}

	resp, err := client.Get("https://api.twitter.com/1.1/search/tweets.json", args)
	if err != nil {
		return results{}, err
	}
	defer resp.Body.Close()

	var result results
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return results{}, err
	}

	return result, nil
}
