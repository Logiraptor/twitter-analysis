package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Jeffail/gabs"
)

type StateCounty struct {
	State, County, Tract, Block string
}

func resolveFips(out chan Status) func(Status) {
	return func(s Status) {
		sc, err := getStateCounty(s.ComputedCoords.Coordinates[1], s.ComputedCoords.Coordinates[0])
		if err != nil {
			return
		}
		s.FipsState = sc.State
		s.FipsCounty = sc.County
		s.FipsTract = sc.Tract
		s.FipsBlock = sc.Block

		out <- s
	}
}

func getStateCounty(lat, lng float32) (StateCounty, error) {
	url := fmt.Sprintf("http://data.fcc.gov/api/block/find?latitude=%f&longitude=%f&showall=false&format=json", lat, lng)
	resp, err := http.Get(url)
	if err != nil {
		return StateCounty{}, err
	}
	defer resp.Body.Close()
	container, err := gabs.ParseJSONBuffer(resp.Body)
	if err != nil {
		return StateCounty{}, err
	}
	fips, ok := container.Path("Block.FIPS").Data().(string)
	if !ok {
		return StateCounty{}, errors.New("no results")
	}
	return StateCounty{
		State:  fips[0:2],
		County: fips[2:5],
		Tract:  fips[5:11],
		Block:  fips[11:],
	}, nil
}
