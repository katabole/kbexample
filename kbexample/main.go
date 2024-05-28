package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/dankinder/katabole/kbexample/actions"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	var conf actions.Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	app, err := actions.NewApp(conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	app.Start()

	// Wait for an interrupt (ctrl-c)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	app.Stop(context.Background())
}
