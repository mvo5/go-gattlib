package main

import (
	"fmt"
	"os"

	"github.com/mvo5/go-gattlib"
)
	
func run() error {
	loop := gattlib.GMainLoopNew()
	go loop.Run()
	defer loop.Quit()

	dest := os.Args[1]
	conn, err := gattlib.Connect(dest)
	if err != nil {
		return err
	}
	fmt.Println("got conn",conn)
	
	return nil
}


func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
