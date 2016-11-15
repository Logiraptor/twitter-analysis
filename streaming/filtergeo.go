package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
)

func filterGeo(out chan Status) func(Status) {
	return func(s Status) {
		s.ComputedCoords = getCoords(s)
		if s.ComputedCoords != nil {
			out <- s
		}

		// if s.User.Location == nil {
		// 	return
		// }
		// grr, err := geocode(*s.User.Location)
		// if err != nil {
		// 	return
		// }

		// switch grr.Status {
		// case "OK":
		// 	first := grr.Results[0]
		// 	s.ComputedCoords = &geo{
		// 		Coordinates: [...]float32{
		// 			first.Geometry.Location.Lng,
		// 			first.Geometry.Location.Lat,
		// 		},
		// 	}
		// 	s.Geocoded = true
		// 	out <- s
		// case "ZERO_RESULTS":
		// case "UNKNOWN_ERROR":
		// default:
		// }
	}
}

func geocode(address string) (*GeoCodeResults, error) {
	args := url.Values{
		"key":     {"AIzaSyBC0Sq-WjUS4kerSh_MEv2YozBTo7Yw23Q"},
		"address": {address},
	}
	resp, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?" + args.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result = new(GeoCodeResults)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getCoords(s Status) *geo {
	if s.Coordinates != nil {
		return s.Coordinates
	}
	if s.Geo != nil {
		return s.Geo
	}
	if s.Place == nil {
		return nil
	}
	if s.Place.BoundingBox.Type != "Polygon" {
		fmt.Println("Not polygon: ", s.Place.BoundingBox.Type)
		return nil
	}
	if len(s.Place.BoundingBox.Coordinates) == 0 {
		fmt.Println("Coords is empty: ", s.Place)
		return nil
	}
	if len(s.Place.BoundingBox.Coordinates[0]) != 4 {
		fmt.Println("Polygon is not a quad: ", s.Place)
		return nil
	}
	c := s.Place.BoundingBox.Coordinates[0]
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
