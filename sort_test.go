package main

import (
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortFunctions(t *testing.T) {
	is := assert.New(t)
	data, err := ioutil.ReadFile("test/functions.go")
	is.NoError(err)

	result, err := Sort(data)
	is.NoError(err)

	// check that the functions and comments are sorted
	pos := []int{}
	for _, s := range []string{
		"PublicA", "public functions", "PublicB", "PublicC",
		"privateA", "private function", "privateB",
	} {
		pos = append(pos, strings.Index(result, s))
	}

	is.True(sort.IntsAreSorted(pos), "%v\n%s", pos, result)
}

func TestSortStructs(t *testing.T) {
	is := assert.New(t)
	data, err := ioutil.ReadFile("test/structs.go")
	is.NoError(err)

	result, err := Sort(data)
	is.NoError(err)

	pos := []int{}
	for _, s := range []string{
		"Car struct", "Drive", "accelerate",
		"Money int", "BuyCar",
	} {
		pos = append(pos, strings.Index(result, s))
	}

	is.True(sort.IntsAreSorted(pos), "%v\n%s", pos, result)
}

func TestSortConst(t *testing.T) {
	is := assert.New(t)
	data, err := ioutil.ReadFile("test/const.go")
	is.NoError(err)

	result, err := Sort(data)
	is.NoError(err)

	pos := []int{}
	for _, s := range []string{
		"first", "second", "third",
	} {
		pos = append(pos, strings.Index(result, s))
	}

	is.True(sort.IntsAreSorted(pos), "%v\n%s", pos, result)
}
func TestSortOther(t *testing.T) {
	is := assert.New(t)
	data, err := ioutil.ReadFile("test/other.go")
	is.NoError(err)

	result, err := Sort(data)
	is.NoError(err)

	pos := []int{}
	for _, s := range []string{
		"time", "pi", "result",
		"Test()",
	} {
		pos = append(pos, strings.Index(result, s))
	}

	is.True(sort.IntsAreSorted(pos), "%v\n%s", pos, result)
}
