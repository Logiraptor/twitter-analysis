package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	. "gopkg.in/mgo.v2/bson"

	"github.com/Logiraptor/notes/engl452/project/models"

	"github.com/Jeffail/gabs"
)

type StateCounty struct {
	State, County string
}

func getStateCounty(lat, lng float32) (StateCounty, error) {
	url := fmt.Sprintf("http://data.fcc.gov/api/block/find?latitude=%f&longitude=%f&showall=false&format=json", lat, lng)
	resp, err := http.Get(url)
	if err != nil {
		return StateCounty{}, err
	}
	container, err := gabs.ParseJSONBuffer(resp.Body)
	if err != nil {
		return StateCounty{}, err
	}
	st, _ := container.Path("State.FIPS").Data().(string)
	ct, _ := container.Path("County.FIPS").Data().(string)
	if st == "" || ct == "" {
		return StateCounty{}, nil
	}
	return StateCounty{
		State:  st,
		County: ct[2:],
	}, nil
}

func getData(location StateCounty, vars ...string) (*gabs.Container, error) {
	censuskey := "834f30d25c3352951df2bd4a457e21a88f9e083f"
	url := fmt.Sprintf("http://api.census.gov/data/2014/acs1?get=%s&for=county:%s&in=state:%s&key=%s",
		strings.Join(vars, ","),
		location.County,
		location.State,
		censuskey,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	container, err := gabs.ParseJSONBuffer(resp.Body)
	if err != nil {
		return nil, err
	}
	return container, nil
}

type job struct {
	tweet models.Status
	resp  chan models.Status
}

func worker(master chan chan job, done chan struct{}) {
	incoming := make(chan job)
	for {
		select {
		case master <- incoming:
		case <-done:
			return
		}
		job := <-incoming
		sc, err := getStateCounty(job.tweet.ComputedCoords.Coordinates[1], job.tweet.ComputedCoords.Coordinates[0])
		if err != nil {
			job.resp <- models.Status{}
			continue
		}
		job.tweet.FIPSCounty = sc.County
		job.tweet.FIPSState = sc.State
		job.resp <- job.tweet
	}
}

func master(num int) chan job {
	workers := make(chan chan job)
	done := make(chan struct{})
	for i := 0; i < num; i++ {
		go worker(workers, done)
	}
	jobs := make(chan job)
	go func(jobs chan job) {
		for job := range jobs {
			worker := <-workers
			worker <- job
		}
		close(done)
	}(jobs)
	return jobs
}

func main() {
	db, err := models.NewClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	tweets := db.DB("engl452").C("tweets")

	runtime.GOMAXPROCS(runtime.NumCPU())
	queue := master(1000)
	resp := make(chan models.Status)
	go func() {
		iter := tweets.Pipe([]M{
			M{"$match": M{
				"computedcoords.coordinates": M{"$ne": nil},
				"fipsstate":                  M{"$eq": nil},
				"fipscounty":                 M{"$eq": nil},
				"place.countrycode":          "US",
			}},
			// M{"$project": M{"coords": "$computedcoords.coordinates"}},
		}).Iter()
		var tweet models.Status
		for iter.Next(&tweet) {
			queue <- job{
				tweet: tweet,
				resp:  resp,
			}
		}
		close(queue)
	}()
	for r := range resp {
		if r.FIPSCounty == "" || r.FIPSState == "" {
			fmt.Println("No Data")
			continue
		}
		// container, err := getData(r,
		// 	"B02001_008E",
		// 	"B02001_001E",
		// 	"B02001_002E",
		// 	"B02001_003E",
		// 	"B02001_005E",
		// 	"B02001_007E",
		// 	"B02001_006E",
		// 	"B02001_004E",
		// )
		// if err != nil {
		// 	fmt.Println("getData", err.Error())
		// 	continue
		// }
		fmt.Println(r.FIPSState, r.FIPSCounty)
		_, err = tweets.UpsertId(M{"id": r.ID}, r)
		if err != nil {
			fmt.Println("Mongo:", err.Error())
			continue
		}
	}
}
