package main

import (
	"log"

	"github.com/streamrail/views/lib"
)

func main() {
	viewsCreator, err := lib.InitViewsCreator()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = viewsCreator.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
}
