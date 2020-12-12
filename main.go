package main

import (
	"encoding/json"
	"github.com/fredericobormann/go-speed/storage"
	"github.com/go-co-op/gocron"
	"log"
	"os/exec"
	"time"
)

var store *storage.Store

func main() {
	store = storage.CreateDB("data.db")

	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(30).Minutes().Do(measureSpeed)
	if err != nil {
		log.Fatalf("Could not create task: %v", err)
	}

	scheduler.StartBlocking()
}

// measureSpeed runs a speedtest-cli command and prints its results
func measureSpeed() {
	output, err := exec.Command("speedtest-cli", "--secure", "--json").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	var structurizedMeasurement storage.SpeedMeasurement
	marshalErr := json.Unmarshal(output, &structurizedMeasurement)
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}
	log.Printf("%+v\n", structurizedMeasurement)
	store.SaveMeasurement(structurizedMeasurement)
}
