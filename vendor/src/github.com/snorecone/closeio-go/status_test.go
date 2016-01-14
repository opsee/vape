package closeio

import (
	"testing"
	"log"
	"os"
)

func TestStatusList(t *testing.T) {
	key := os.Getenv("CLOSEIO_KEY")
	closeAPI := New(key)
	statuses, err := closeAPI.Statuses()
	if err != nil {
		log.Println(err)
	}
	log.Println(statuses)
}
