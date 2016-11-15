package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()

	db, err := mgo.Dial("localhost:27017")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	m.Use(func(ctx martini.Context) {
		s := db.Clone()
		ctx.Map(s)
		ctx.Next()
		s.Close()
	})

	m.Get("/heatmap/:form", HeatMapDataHandler)
	m.Run()
}

func HeatMapDataHandler(db *mgo.Session, params martini.Params, rw http.ResponseWriter) {
	tweets := db.DB("clean_tweets").C("tweets")
	type coords struct {
		Coords [2]float32
	}
	var coordResults []coords
	err := tweets.Pipe([]M{
		M{"$match": M{
			"computedcoords": M{"$ne": nil},
			"lang":           "en",
			params["form"]:   true,
		}},
		M{"$project": M{
			"coords": "$computedcoords.coordinates",
		}},
	}).All(&coordResults)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
		return
	}

	type countType struct {
		Count  int     `json:"n"`
		AvgLen float32 `bson:"avgLen" json:"avgLen"`
	}
	var count countType
	err = tweets.Pipe([]M{
		M{"$match": M{
			"computedcoords": M{"$ne": nil},
			"lang":           "en",
			// "place.countrycode": "US",
			// "$and": []M{
			// 	{"computedcoords[0]": M{"$lt": -66.885444}},
			// 	{"computedcoords[0]": M{"$gt": -124.848974}},
			// 	{"computedcoords[1]": M{"$lt": 49.384358}},
			// 	{"computedcoords[1]": M{"$gt": 24.396308}},
			// },
			//  (-124.848974, 24.396308) - (-66.885444, 49.384358)
			params["form"]: true,
		}},
		{"$group": M{
			"_id":    nil,
			"count":  M{"$sum": 1},
			"avgLen": M{"$avg": "$text.length"},
		}},
	}).One(&count)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
		return
	}

	type outputType struct {
		Coords []coords `json:"coords"`
		countType
	}

	json.NewEncoder(rw).Encode(outputType{
		Coords:    coordResults,
		countType: count,
	})
}
