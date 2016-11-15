package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type censusData struct {
	Tweets     float64
	BroTweets  float64
	BruhTweets float64
	BrahTweets float64

	TotalPop          float64
	TwoPop            float64
	WhitePop          float64
	BlackPop          float64
	AsianPop          float64
	OtherPop          float64
	NativeHawaiianPop float64
	AmericanIndianPop float64
}

type CensusStage struct {
	data    map[string]map[string]censusData
	metrics chan metric
	sync.RWMutex
}

func (c *CensusStage) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c.RLock()
	defer c.RUnlock()
	err := json.NewEncoder(rw).Encode(c.data)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
}

func (c *CensusStage) getData(state, county string) (censusData, bool) {
	c.RLock()
	defer c.RUnlock()
	stateMap, ok := c.data[state]
	if !ok {
		return censusData{}, false
	}
	countyMap, ok := stateMap[county]
	return countyMap, ok
}

func (c *CensusStage) putData(state, county string, data censusData) {
	c.Lock()
	defer c.Unlock()
	stateMap, ok := c.data[state]
	if !ok {
		stateMap = make(map[string]censusData)
	}
	stateMap[county] = data
	c.data[state] = stateMap
}

func (c *CensusStage) process(s Status) Status {
	curr, ok := c.getData(s.FipsState, s.FipsCounty)
	if ok {
		curr.Tweets++
		if s.Bro {
			curr.BroTweets++
		}
		if s.Bruh {
			curr.BruhTweets++
		}
		if s.Brah {
			curr.BrahTweets++
		}
		c.putData(s.FipsState, s.FipsCounty, curr)
		return s
	}

	data, err := getData(s.FipsState, s.FipsCounty,
		"B02001_001E",
		"B02001_008E",
		"B02001_002E",
		"B02001_003E",
		"B02001_005E",
		"B02001_007E",
		"B02001_006E",
		"B02001_004E",
	)
	if err != nil {
		return s
	}
	if len(data) != 2 {
		return s
	}

	var bro, brah, bruh float64
	if s.Bro {
		bro = 1
	}
	if s.Bruh {
		bruh = 1
	}
	if s.Brah {
		brah = 1
	}
	if (bro == bruh) && (bro == brah) && bro == 0 {
		fmt.Println(s.Truncated, s.Text)
	}
	row := data[1]
	censusdata := censusData{
		Tweets:            1,
		BroTweets:         bro,
		BruhTweets:        bruh,
		BrahTweets:        brah,
		TotalPop:          orZero(strconv.ParseFloat(row[0], 64)),
		TwoPop:            orZero(strconv.ParseFloat(row[1], 64)),
		WhitePop:          orZero(strconv.ParseFloat(row[2], 64)),
		BlackPop:          orZero(strconv.ParseFloat(row[3], 64)),
		AsianPop:          orZero(strconv.ParseFloat(row[4], 64)),
		OtherPop:          orZero(strconv.ParseFloat(row[5], 64)),
		NativeHawaiianPop: orZero(strconv.ParseFloat(row[6], 64)),
		AmericanIndianPop: orZero(strconv.ParseFloat(row[7], 64)),
	}
	c.putData(s.FipsState, s.FipsCounty, censusdata)
	return s
}

func (c *CensusStage) computeStats(rw http.ResponseWriter, req *http.Request) {
	c.RLock()
	var all []censusData
	for _, counties := range c.data {
		for _, data := range counties {
			all = append(all, data)
		}
	}
	c.RUnlock()

	tweets := make([]float64, len(all))
	broTweets := make([]float64, len(all))
	bruhTweets := make([]float64, len(all))
	brahTweets := make([]float64, len(all))
	totalPop := make([]float64, len(all))
	twoPop := make([]float64, len(all))
	whitePop := make([]float64, len(all))
	blackPop := make([]float64, len(all))
	asianPop := make([]float64, len(all))
	otherPop := make([]float64, len(all))
	nativeHawaiianPop := make([]float64, len(all))
	americanIndianPop := make([]float64, len(all))
	totalTweets := 0.0
	broTotalTweets := 0.0
	bruhTotalTweets := 0.0
	brahTotalTweets := 0.0
	counties := 0.0
	for i, data := range all {
		counties++
		totalTweets += data.Tweets
		broTotalTweets += data.BroTweets
		bruhTotalTweets += data.BruhTweets
		brahTotalTweets += data.BrahTweets
		tweets[i] = data.Tweets
		broTweets[i] = data.BroTweets
		bruhTweets[i] = data.BruhTweets
		brahTweets[i] = data.BrahTweets
		totalPop[i] = data.TotalPop / 1e6
		twoPop[i] = data.TwoPop / 1e6
		whitePop[i] = data.WhitePop / 1e6
		blackPop[i] = data.BlackPop / 1e6
		asianPop[i] = data.AsianPop / 1e6
		otherPop[i] = data.OtherPop / 1e6
		nativeHawaiianPop[i] = data.NativeHawaiianPop / 1e6
		americanIndianPop[i] = data.AmericanIndianPop / 1e6
	}

	out := map[string]interface{}{
		"Counties":   counties,
		"Tweets":     totalTweets,
		"BroTweets":  broTotalTweets,
		"BruhTweets": bruhTotalTweets,
		"BrahTweets": brahTotalTweets,
		"All":        computeCorrelations(tweets, totalPop, twoPop, whitePop, blackPop, asianPop, otherPop, nativeHawaiianPop, americanIndianPop),
		"Bro":        computeCorrelations(broTweets, totalPop, twoPop, whitePop, blackPop, asianPop, otherPop, nativeHawaiianPop, americanIndianPop),
		"Bruh":       computeCorrelations(bruhTweets, totalPop, twoPop, whitePop, blackPop, asianPop, otherPop, nativeHawaiianPop, americanIndianPop),
		"Brah":       computeCorrelations(brahTweets, totalPop, twoPop, whitePop, blackPop, asianPop, otherPop, nativeHawaiianPop, americanIndianPop),
	}
	err := json.NewEncoder(rw).Encode(out)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
}

func computeCorrelations(tweets, totalPop,
	twoPop, whitePop, blackPop, asianPop,
	otherPop, nativeHawaiianPop, americanIndianPop []float64) map[string]float64 {
	return map[string]float64{
		"TotalPop":          coorelation(tweets, totalPop),
		"TwoPop":            coorelation(tweets, twoPop),
		"WhitePop":          coorelation(tweets, whitePop),
		"BlackPop":          coorelation(tweets, blackPop),
		"AsianPop":          coorelation(tweets, asianPop),
		"OtherPop":          coorelation(tweets, otherPop),
		"NativeHawaiianPop": coorelation(tweets, nativeHawaiianPop),
		"AmericanIndianPop": coorelation(tweets, americanIndianPop),
	}
}

func orZero(i float64, e error) float64 {
	if e != nil {
		return 0
	}
	return i
}

func getData(state, county string, vars ...string) ([][]string, error) {
	censuskey := "834f30d25c3352951df2bd4a457e21a88f9e083f"
	url := fmt.Sprintf("http://api.census.gov/data/2014/acs1?get=%s&for=county:%s&in=state:%s&key=%s",
		strings.Join(vars, ","),
		county,
		state,
		censuskey,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data [][]string
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
