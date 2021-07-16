package test

import (
	"flag"
	"fmt"
	"time"
)

var debug = flag.Bool("debug", false, "debugging explanation")
var Other = 20

const result = 10

func (*otherTestStruct) Test() {
	fmt.Println(pi*result*Other, time.Now())
}

type otherTestStruct int

const pi = 3
