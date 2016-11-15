package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

func storeinMongo(dbName string) func(Status) Status {
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	tweets := sess.DB(dbName).C("tweets")

	return func(s Status) Status {
		_, err := tweets.UpsertId(s.ID, s)
		if err != nil {
			sess.DB(dbName).C("errors").Insert(bson.M{
				"error": err.Error(),
			})
		}
		return s
	}
}

func reprocess(dbName string, out chan Status) {
	defer close(out)
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	tweets := sess.DB(dbName).C("tweets")

	iter := tweets.Find(nil).Iter()
	defer iter.Close()

	count := 0.0
	lastTime := time.Now()
	var s Status
	for iter.Next(&s) {
		if time.Since(lastTime) > (time.Second * 30) {
			fmt.Printf("\r%v %v tweets / s", time.Now().Format(time.RFC3339), count/(float64(time.Since(lastTime))/float64(time.Second)))
			count = 0
			lastTime = time.Now()
		}
		if s.FipsBlock == "" {
			count++
			out <- s
		}
	}
	if err := iter.Err(); err != nil {
		fmt.Println(err)
		return
	}
}
