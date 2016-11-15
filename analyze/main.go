package main

import (
	"fmt"

	. "gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

func main() {
	db, err := mgo.Dial("localhost:27017")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	tweets := db.DB("engl452").C("tweets")

	// type count struct {
	// 	HasGeo bool `bson:"_id"`
	// 	Count  int  `bson:"count"`
	// }

	// // Count users with geo enabled vs not
	// var results []count
	// err = tweets.Pipe([]M{
	// 	{"$group": M{
	// 		"_id":   "$user.geoenabled",
	// 		"count": M{"$sum": 1},
	// 	}},
	// }).All(&results)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(results)

	// var resultsI []M
	// err = tweets.Pipe([]M{
	// 	{"$match": M{"computedcoords": M{"$ne": nil}}},
	// }).All(&resultsI)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// for _, doc := range resultsI {
	// 	fmt.Println(doc["computedcoords"])
	// }

	// fmt.Println(len(resultsI))

	type coords struct {
		Coords [2]float32
	}
	var coordResults []coords
	err = tweets.Pipe([]M{
		M{"$match": M{"place": M{"$ne": nil}, "place.countrycode": "US"}},
		M{"$project": M{"coords": "$computedcoords.coordinates"}},
	}).All(&coordResults)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, doc := range coordResults {
		fmt.Printf("new google.maps.LatLng(%v, %v),\n", doc.Coords[1], doc.Coords[0])
	}

	// json.NewEncoder(os.Stdout).Encode(coordResults)
}
