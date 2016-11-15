package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

func main() {
	client := NewClient(
		"32FiBBFYuEb7c9S2K5tTBGddb",
		"oCIjH791Bg8zkAWzXKFRoZhzeFe2UDtNWZIvNUqDMZODIzoxny",
		"27555535-JkZ8yREgZEnLddst2v9ze5v0LO5eRSn9iurOvA5xw",
		"CYX0Kwe8j4E57XemzbisiWhaOS0Y2Uoq5jp564Od7b1sP",
	)

	// resp, err := client.Get("https://api.twitter.com/1.1/statuses/show.json", map[string]string{
	// 	"id": "659494313726865408",
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer resp.Body.Close()

	// io.Copy(os.Stdout, resp.Body)

	// resp, err := client.Post("https://stream.twitter.com/1.1/statuses/filter.json", map[string]string{
	// "track": "bro,broo,brooo,bruh,bruhh,bruhhh,brah,brahh,brahhh",
	// })
	resp, err := client.Get("https://stream.twitter.com/1.1/statuses/sample.json", map[string]string{
		"language": "en",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	db, err := mgo.Dial("localhost:27017")
	if err != nil {
		fmt.Println(err)
		return
	}
	tweets := db.DB("engl452").C("alltweets")

	i := 0
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

		v.ComputedCoords = getCoords(v)
		v.Brah = strings.Contains(strings.ToLower(v.Text), "brah")
		v.Bruh = strings.Contains(strings.ToLower(v.Text), "bruh")
		v.Bro = strings.Contains(strings.ToLower(v.Text), "bro")

		if v.ComputedCoords != nil {
			_, err = tweets.Upsert(bson.M{"id": v.ID}, v)
			if err != nil {
				fmt.Println(err)
				return
			}

			i++
		}
		fmt.Println(i, v.Text)

		if i > 100000 {
			fmt.Println("Collected 100000 tweets! stopping")
			return
		}
	}

}

func getCoords(v Status) *geo {
	if v.Coordinates != nil {
		return v.Coordinates
	}
	if v.Geo != nil {
		return v.Geo
	}
	if v.Place == nil {
		return nil
	}
	if v.Place.BoundingBox.Type != "Polygon" {
		fmt.Println("Not polygon: ", v.Place.BoundingBox.Type)
		return nil
	}
	if len(v.Place.BoundingBox.Coordinates) == 0 {
		fmt.Println("Coords is empty: ", v.Place)
		return nil
	}
	if len(v.Place.BoundingBox.Coordinates[0]) != 4 {
		fmt.Println("Polygon is not a quad: ", v.Place)
		return nil
	}
	c := v.Place.BoundingBox.Coordinates[0]
	return &geo{
		Type: "Point",
		Coordinates: [...]float32{
			randFloat(c[0][0], c[2][0]),
			randFloat(c[0][1], c[1][1]),
		},
	}
}

func randFloat(a, b float64) float32 {
	if a > b {
		a, b = b, a
	}
	return float32(a + rand.Float64()*(b-a))
}
