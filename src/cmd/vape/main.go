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
	keyPath := os.Getenv("VAPE_KEYFILE")
	if keyPath == "" {
		log.Fatal("Must set VAPE_KEYFILE environment variable.")
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Println("Unable to read keyfile:", keyPath)
		log.Fatal(err)
	}
	token.Init(key)

	pgConn := os.Getenv("POSTGRES_CONN")
	if pgConn == "" {
		log.Fatal("Must set POSTGRES_CONN environment variable.")
	}
	err = store.Init(pgConn)
	if err != nil {
		log.Println("Unable to open postgres store:", pgConn)
		log.Fatal(err)
	}

	api.InjectLogger(os.Stdout)
	publicHost := os.Getenv("VAPE_PUBLIC_HOST")
	privateHost := os.Getenv("VAPE_PRIVATE_HOST")
	if publicHost == "" {
		log.Fatal("Must set VAPE_PUBLIC_HOST environment variable.")
	}
	if privateHost == "" {
		log.Fatal("Must set VAPE_PRIVATE_HOST environment variable.")
	}

	api.ListenAndServe(os.Getenv("VAPE_PUBLIC_HOST"), os.Getenv("VAPE_PRIVATE_HOST"))
}
