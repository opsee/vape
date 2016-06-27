package main

import (
	"io/ioutil"
	"os"

	"github.com/keighl/mandrill"
	log "github.com/opsee/logrus"
	"github.com/opsee/vape/api"
	"github.com/opsee/vape/service"
	"github.com/opsee/vape/servicer"
	"github.com/opsee/vape/store"
	"github.com/opsee/vaper"
)

func main() {
	log.SetLevel(log.DebugLevel)
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

	intercomKey := os.Getenv("INTERCOM_KEY")
	if intercomKey == "" {
		log.Println("WARN: INTERCOM_KEY not set, we won't send user HMAC.")
	}

	closeioKey := os.Getenv("CLOSEIO_KEY")
	if closeioKey == "" {
		log.Println("WARN: CLOSEIO_KEY not set, we won't create leeeds.")
	}

	host := os.Getenv("OPSEE_HOST")
	if host == "" {
		log.Fatal("Must set the OPSEE_HOST environment variable.")
	}

	slackUrl := os.Getenv("SLACK_ENDPOINT")
	if slackUrl == "" {
		log.Println("WARN: SLACK_ENDPOINT not set, we won't post notifications.")
	}

	slackDomain := os.Getenv("VAPE_SLACK_DOMAIN")
	if slackDomain == "" {
		log.Println("WARN: VAPE_SLACK_DOMAIN not set, we won't invite users.")
	}

	slackToken := os.Getenv("VAPE_SLACK_ADMIN_TOKEN")
	if slackToken == "" {
		log.Println("WARN: VAPE_SLACK_ADMIN_TOKEN not set, we won't invite users.")
	}

	spanxHost := os.Getenv("VAPE_SPANX_HOST")
	if spanxHost == "" {
		log.Fatal("Must set the VAPE_SPANX_HOST environment variable.")
	}

	servicer.Init(host, mandrillClient, intercomKey, closeioKey, slackUrl, slackDomain, slackToken, spanxHost)

	certfile := os.Getenv("VAPE_CERT")
	certkeyfile := os.Getenv("VAPE_CERT_KEY")
	if certfile == "" || certkeyfile == "" {
		log.Fatal("VAPE_CERT and VAPE_CERT_KEY must be set, and you must have a certificate and key")
	}

	service := service.New()

	api.ListenAndServe(os.Getenv("VAPE_PUBLIC_HOST"), os.Getenv("VAPE_PRIVATE_HOST"), certfile, certkeyfile, service.Server)
}
