package main

import (
	"github.com/Spear5030/yagophkeeper/internal/client/app"
	"github.com/Spear5030/yagophkeeper/internal/client/config"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = a.Run()
	if err != nil {
		log.Println(err)
	}
}
