package main

import (
	"image/color"
	"os"

	"githib.com/maps/internal"
	"github.com/glaslos/go-osm"
)

func main() {

	brest := internal.BoundingBox{
		SouthLatitude: 48.363099,
		NorthLatitude: 48.426005,
		WestLongitude: -4.559003,
		EastLongitude: -4.421826,
	}
	// data, err := retrieveOSMData(brest)
	// if err != nil {
	// 	panic(err)
	// }
	data, err := os.ReadFile("./cmd/data")
	if err != nil {
		panic(err)
	}

	realData, err := osm.DecodeString(string(data))
	if err != nil {
		panic(err)
	}
	// log.Println(realData)
	width, height := internal.CalculateSize(brest, 50000)
	d := internal.NewDrawing(width, height, realData)
	colorScheme := internal.Colors{
		BackgroundColor: color.RGBA{0, 20, 132, 255},
		RoadColor:       color.RGBA{255, 255, 255, 255},
		WaterColor:      color.RGBA{0, 0, 255, 255},
	}
	d.SetColorScheme(colorScheme)
	// d := internal.NewDrawing(width, height, data)
	d.Normalize(brest)
	d.ComputeAll()
	d.DrawRoads()
	// d.DrawWater()

	d.Context.SavePNG("out.png")
}
