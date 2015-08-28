package main

import (
	"github.com/opsee/vape/api"
	"github.com/opsee/vape/store"
	"github.com/opsee/vape/token"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	key, err := ioutil.ReadFile(os.Getenv("VAPE_KEYFILE"))
	if err != nil {
		log.Fatal(err)
	}
	token.Init(key)

	err = store.Init(os.Getenv("POSTGRES_CONN"))
	if err != nil {
		log.Fatal(err)
	}

	api.InjectLogger(os.Stdout)
	api.ListenAndServe(os.Getenv("VAPE_HOST"))
}
