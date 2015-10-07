package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/keighl/mandrill"
	"github.com/opsee/vape/api"
	"github.com/opsee/vape/servicer"
	"github.com/opsee/vape/store"
	"github.com/opsee/vaper"
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
	vaper.Init(key)

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

	mandrillAPIKey := os.Getenv("MANDRILL_API_KEY")
	var mandrillClient *mandrill.Client

	if mandrillAPIKey == "" {
		log.Println("WARN: MANDRILL_API_KEY not set, we won't send emails.")
	} else {
		mandrillClient = mandrill.ClientWithKey(mandrillAPIKey)
	}

	host := os.Getenv("OPSEE_HOST")
	if host == "" {
		log.Fatal("Must set the OPSEE_HOST environment variable.")
	}

	servicer.Init(host, mandrillClient)

	api.ListenAndServe(os.Getenv("VAPE_PUBLIC_HOST"), os.Getenv("VAPE_PRIVATE_HOST"))
}
