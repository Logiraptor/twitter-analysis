package main

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestStateCounty(t *testing.T) {
	sc, err := getStateCounty(30.216452099999998, -92.0599479)
	assert.Nil(t, err)
	assert.EqualValues(t, &StateCounty{
		State:  "22",
		County: "055",
	}, sc)
}
