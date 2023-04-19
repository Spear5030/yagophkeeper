package main

import (
	"github.com/Spear5030/yagophkeeper/internal/client/app"
	"github.com/Spear5030/yagophkeeper/internal/client/config"
	"log"
)

var Version string
var BuildTime string

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	a, err := app.New(cfg, Version, BuildTime)
	if err != nil {
		log.Fatal(err)
	}
	err = a.Run()
	if err != nil {
		log.Println(err)
	}
}
