package main

import (
	"log"

	"github.com/katabole/kbexample/actions"
	"github.com/katabole/kbexample/models"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	var conf actions.Config
	if err := envconfig.Process("", &conf); err != nil {
		log.Fatal(err.Error())
	}

	db, err := models.NewDB(conf.DBConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	for _, u := range []models.User{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Charlie"},
	} {
		if _, err := db.CreateUser(&u); err != nil {
			log.Fatal(err.Error())
		}
	}
}
