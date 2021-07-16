package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/sync/errgroup"
)

var (
	write   = flag.Bool("w", false, "write result to (source) file instead of stdout")
	verbose = flag.Bool("v", false, "more verbose error reporting")
)

func main() {
	flag.Parse()

	files := make(chan string)
	wg := errgroup.Group{}

	wg.Go(func() error { return Scan(files, flag.Args()...) })

	for f := range files {
		file := f
		wg.Go(func() error {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			ast, err := Sort(data)
			if err != nil {
				return err
			}
			if *write {
				w, err := os.Create(file)
				if err != nil {
					return err
				}
				_, err = w.WriteString(ast)
				return err
			} else {
				fmt.Println(ast)
			}
			return err
		})
	}

	err := wg.Wait()
	if err != nil {
		format := "%v"
		if *verbose {
			format = "%+v"
		}
		log.Fatalf(format, err)
	}
}
