package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"encoding/json"

	"os"

	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

func main() {
	router := mux.NewRouter()
	router.Path("/tweets").HandlerFunc(tweetHandler)
	router.Path("/").Handler(http.FileServer(http.Dir("public")))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

func tweetHandler(rw http.ResponseWriter, req *http.Request) {
	client := NewClient(
		"32FiBBFYuEb7c9S2K5tTBGddb",
		"oCIjH791Bg8zkAWzXKFRoZhzeFe2UDtNWZIvNUqDMZODIzoxny",
		"27555535-JkZ8yREgZEnLddst2v9ze5v0LO5eRSn9iurOvA5xw",
		"CYX0Kwe8j4E57XemzbisiWhaOS0Y2Uoq5jp564Od7b1sP",
	)

	websocket.Handler(func(c *websocket.Conn) {
		defer c.Close()
		resp, err := client.Post("https://stream.twitter.com/1.1/statuses/sample.json", map[string]string{})
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer resp.Body.Close()

		scan := bufio.NewScanner(resp.Body)
		for scan.Scan() {
			var s Status
			err := json.Unmarshal(scan.Bytes(), &s)
			if err != nil {
				log.Println(err.Error())
				return
			}

			s.ComputedCoords = getCoords(s)

			buf, err := json.Marshal(s)
			if err != nil {
				log.Println(err.Error())
				return
			}

			_, err = c.Write(buf)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}).ServeHTTP(rw, req)
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
