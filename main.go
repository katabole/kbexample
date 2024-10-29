package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/katabole/kbexample/actions"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	app.Stop(context.Background())
}
