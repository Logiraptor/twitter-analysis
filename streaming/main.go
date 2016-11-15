package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var broRegexp = regexp.MustCompile(`(?i)\bb+r+o+\b`)
var bruhRegexp = regexp.MustCompile(`(?i)\bb+r+u+h+\b`)
var brahRegexp = regexp.MustCompile(`(?i)\bb+r+a+h+\b`)

type Stage func(chan Status) func(Status)

type Conf struct {
	Workers int
	Op      Stage
}

type ChainConf []Conf

type Chain struct {
	in, out chan Status
}

func (c *Chain) stage(out chan Status) func(s Status) {
	go func() {
		for st := range c.out {
			out <- st
		}
	}()
	return func(s Status) {
		c.in <- s
	}
}

func NewChain(confs ChainConf) Chain {
	const buffer = 100
	in := make(chan Status, buffer)

	currentIn := in
	currentOut := in
	for _, conf := range confs {
		currentIn, currentOut = currentOut, make(chan Status, buffer)
		go func(conf Conf, in, out chan Status) {
			var wg sync.WaitGroup
			wg.Add(conf.Workers)
			for i := 0; i < conf.Workers; i++ {
				go func(f func(Status)) {
					for st := range in {
						f(st)
					}
					wg.Done()
				}(conf.Op(out))
			}
			wg.Wait()
			close(out)
		}(conf, currentIn, currentOut)
	}

	return Chain{in: in, out: currentOut}
}

func transform(f func(Status) Status) Stage {
	return func(out chan Status) func(Status) {
		return func(s Status) {
			out <- f(s)
		}
	}
}

func caps(s Status) Status {
	return Status{
		Text: strings.ToUpper(s.Text),
	}
}

func replace(s Status) Status {
	return Status{
		Text: strings.Replace(s.Text, " ", "_", -1),
	}
}

func stream(out chan Status) {
	for {
		client := NewClient(
			"32FiBBFYuEb7c9S2K5tTBGddb",
			"oCIjH791Bg8zkAWzXKFRoZhzeFe2UDtNWZIvNUqDMZODIzoxny",
			"27555535-JkZ8yREgZEnLddst2v9ze5v0LO5eRSn9iurOvA5xw",
			"CYX0Kwe8j4E57XemzbisiWhaOS0Y2Uoq5jp564Od7b1sP",
		)
		resp, err := client.Post("https://stream.twitter.com/1.1/statuses/filter.json", map[string]string{
			"track": "bro,broo,brooo,bruh,bruhh,bruhhh,brah,brahh,brahhh",
		})
		// resp, err := client.Get("https://stream.twitter.com/1.1/statuses/sample.json", map[string]string{
		// 	"language": "en",
		// })
		if err != nil {
			if resp.StatusCode == 503 {
				continue
			}
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		scan := bufio.NewScanner(resp.Body)
		for scan.Scan() {
			var v Status
			err := json.Unmarshal(scan.Bytes(), &v)
			if err != nil {
				fmt.Println(err)
				if err, ok := err.(*json.UnmarshalTypeError); ok {
					fmt.Println(scan.Text()[err.Offset])
				}
				return
			}
			v.Bro = broRegexp.MatchString(v.Text)
			v.Bruh = bruhRegexp.MatchString(v.Text)
			v.Brah = brahRegexp.MatchString(v.Text)
			out <- v
		}
	}
}

func truncate(s Status) Status {
	if len(s.Text) > 15 {
		s.Text = s.Text[:15]
	}
	return s
}

func dropWhenFull(out chan Status) func(Status) {
	return func(s Status) {
		select {
		case out <- s:
		default:
			fmt.Println("dropping")
		}
	}
}

func reportSize(prefix string, out chan metric) Stage {
	return func(next chan Status) func(Status) {
		return func(s Status) {
			out <- metric{
				Name:  prefix,
				Value: len(next),
				Max:   cap(next),
			}
			next <- s
		}
	}
}

type metric struct {
	Name  string
	Value int
	Max   int
	Sum   bool
}

func main() {
	metrics := make(chan metric)
	census := &CensusStage{
		metrics: metrics,
		data:    make(map[string]map[string]censusData),
	}
	chain := NewChain(
		ChainConf{
			{1, reportSize("raw", metrics)},
			{1, dropWhenFull},
			{1, reportSize("geocode", metrics)},
			{5, filterGeo}, // google geocode, drop those without coords
			{1, reportSize("fips", metrics)},
			{100, resolveFips},
			{1, reportSize("postfips", metrics)},
			{1, transform(storeinMongo("clean_tweets"))},
		},
	)

	go stream(chain.in)
	// go reprocess("clean_tweets", chain.in)

	var mets = make(map[string]metric)

	var metslock sync.Mutex
	go func() {
		for {
			select {
			case <-chain.out:
			case m := <-metrics:
				metslock.Lock()
				mets[m.Name] = m
				metslock.Unlock()
			}
		}
	}()

	http.HandleFunc("/metrics", func(rw http.ResponseWriter, req *http.Request) {
		metslock.Lock()
		json.NewEncoder(rw).Encode(mets)
		metslock.Unlock()
	})
	http.Handle("/data", census)
	http.HandleFunc("/coorelation", census.computeStats)
	http.ListenAndServe(":8080", nil)
}
