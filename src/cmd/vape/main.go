package main

import (
        "os"
        "log"
        "io/ioutil"
        "github.com/opsee/vape/store"
        "github.com/opsee/vape/api"
)

func main() {
        key, err := ioutil.ReadFile(os.Getenv("VAPE_KEYFILE"))
        if err != nil {
                log.Fatal(err)
        }

        err = store.Init(os.Getenv("POSTGRES_CONN"), key)
        if err != nil {
                log.Fatal(err)
        }

        api.ListenAndServe(":8080", os.Stdout)
}
