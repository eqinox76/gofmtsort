package test

import (
	"flag"
	"fmt"
	"time"
)

var debug = flag.Bool("debug", false, "debugging explanation")

const result = 10

var Other = 20

func (*otherTestStruct) Test() {
	fmt.Println(pi*result*Other, time.Now())
}

type otherTestStruct int

const pi = 3
