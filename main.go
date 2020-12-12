package main

import (
	"fmt"
	"github.com/fredericobormann/go-speed/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"strconv"
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

	scheduler.StartAsync()

	router := gin.Default()
	router.StaticFile("/", "graph.png")
	if err := router.Run(":8070"); err != nil {
		log.Fatal(err)
	}
}

// measureSpeed runs a speedtest-cli command and prints its results
func measureSpeed() {
	/*output, err := exec.Command("speedtest-cli", "--secure", "--json").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	var structurizedMeasurement storage.SpeedMeasurement
	marshalErr := json.Unmarshal(output, &structurizedMeasurement)
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}
	log.Printf("%+v\n", structurizedMeasurement)
	store.SaveMeasurement(structurizedMeasurement)*/

	refreshGraph()
}

func refreshGraph() {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("Creating plot failed: %v\n", err)
	}

	xticks := plot.TimeTicks{Format: "02.01.2006\n15:04"}

	p.Title.Text = "Internet speed"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Speed"
	p.X.Tick.Marker = xticks
	p.Y.Min = 0
	p.Y.Tick.Marker = mbitTicks{}

	measurements := store.GetMeasurements()

	pts := make(plotter.XYs, len(measurements))
	for i, m := range measurements {
		pts[i].X = float64(m.Timestamp.Unix())
		pts[i].Y = m.Download
	}

	_, scatter, err := plotter.NewLinePoints(pts)
	if err != nil {
		log.Fatal(err)
	}
	scatter.Color = color.RGBA{
		R: 0,
		G: 255,
		B: 255,
	}

	p.Add(scatter)

	if err := p.Save(40*vg.Centimeter, 24*vg.Centimeter, "graph.png"); err != nil {
		log.Fatalf("Saving graph failed: %v\n", err)
	}
}

type mbitTicks struct{}

func (mbitTicks) Ticks(_, max float64) []plot.Tick {
	var ticks []plot.Tick

	for i := 0; i*1000000 < int(max); i += 10 {
		ticks = append(ticks, plot.Tick{
			Label: fmt.Sprintf("%s MBit/s", strconv.FormatInt(int64(i), 10)),
			Value: float64(i * 1000000),
		})
	}
	return ticks
}
