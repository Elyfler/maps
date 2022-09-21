package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sync"

	"github.com/glaslos/go-osm"
	"gopkg.in/fogleman/gg.v1"
)

type BoundingBox struct {
	southLatitude float64
	northLatitude float64
	westLongitude float64
	eastLongitude float64
}

type Line struct {
	x1        float64
	y1        float64
	x2        float64
	y2        float64
	lineWidth float64
}

func (l *Line) toFunc() (x1, y1, x2, y2 float64) {
	return l.x1, l.y1, l.x2, l.y2
}

type Drawing struct {
	Context *gg.Context
	Roads   Roads
	Lines   []Line
	Nodes   []osm.Node
}

func NewDrawing(width, height int, data *osm.Map) Drawing {
	var d Drawing
	d.Context = gg.NewContext(width, height)
	d.Roads.sortRoads(data)
	d.Nodes = data.Nodes
	return d
}

func (d *Drawing) DrawRoads() {
	log.Println("DrawRoads")
	d.Context.SetRGB(0, 0, 0)
	d.Context.Fill()
	d.Context.SetRGB(1, 0, 0)
	for _, line := range d.Lines {
		d.Context.SetLineWidth(line.lineWidth)
		d.Context.DrawLine(line.toFunc())
		d.Context.Stroke()
	}
	d.Context.SavePNG("out.png")
}

func (d *Drawing) ComputeLines(wg *sync.WaitGroup, ways []osm.Way, lineWidth float64) {
	defer wg.Done()
	// Ugly shit to rework later on
	for _, motorway := range ways {
		var nodes []osm.Node
		for _, ID := range motorway.Nds {
			for _, node := range d.Nodes {
				if node.ID == ID.ID {
					nodes = append(nodes, node)
				}
			}
		}
		for i := 0; i < len(nodes)-1; i++ {
			line := Line{
				x1:        nodes[i].Lat,
				y1:        nodes[i].Lng,
				x2:        nodes[i+1].Lat,
				y2:        nodes[i+1].Lng,
				lineWidth: lineWidth,
			}
			log.Println(line)
			d.Lines = append(d.Lines, line)
		}
	}
}

func (d *Drawing) Normalize(bbox BoundingBox) {
	maxWidth := math.Abs(bbox.northLatitude - bbox.southLatitude)
	maxHeight := math.Abs(bbox.westLongitude - bbox.eastLongitude)
	// log.Println(maxWidth)
	// log.Println(maxHeight)
	// log.Println(d.Context.Width())
	// log.Println(d.Context.Height())
	for k, _ := range d.Nodes {
		d.Nodes[k].Lat = ((bbox.northLatitude - d.Nodes[k].Lat) / maxWidth) * float64(d.Context.Width())
		log.Println(d.Nodes[k].Lat)
		d.Nodes[k].Lng = ((-bbox.westLongitude + d.Nodes[k].Lng) / maxHeight) * float64(d.Context.Height())
		log.Println(d.Nodes[k].Lng)
	}
}

type Roads struct {
	Motorway      []osm.Way
	MotorwayLink  []osm.Way
	Trunk         []osm.Way
	TrunkLink     []osm.Way
	Primary       []osm.Way
	PrimaryLink   []osm.Way
	Secondary     []osm.Way
	SecondaryLink []osm.Way
	Tertiary      []osm.Way
	TertiaryLink  []osm.Way
	Unclassified  []osm.Way
	Residential   []osm.Way
	Service       []osm.Way
	Pedestrian    []osm.Way
	Track         []osm.Way
	Footway       []osm.Way
	Steps         []osm.Way
	Path          []osm.Way
}

func calculateSize(b BoundingBox, scale float64) (width int, height int) {
	log.Println("calculateSize")
	baseWidth := math.Abs(b.westLongitude - b.eastLongitude)
	baseHeight := b.northLatitude - b.southLatitude
	// log.Println(baseWidth)
	// log.Println(baseHeight)
	return int(scale * baseWidth), int(scale * baseHeight)
}

func main() {
	fmt.Println("Hello, world")

	brest := BoundingBox{
		westLongitude: -4.559003,
		southLatitude: 48.363099,
		eastLongitude: -4.421826,
		northLatitude: 48.416005,
	}

	// data, err := retrieveOSMData(brest)
	// if err != nil {
	// 	panic(err)
	// }
	data, err := os.ReadFile("data")
	if err != nil {
		panic(err)
	}

	realData, err := osm.DecodeString(string(data))
	if err != nil {
		panic(err)
	}
	// log.Println(realData)
	// width, height := calculateSize(brest, 10000)
	d := NewDrawing(1000, 1000, realData)
	// d := NewDrawing(width, height, data)
	d.Normalize(brest)

	wg := new(sync.WaitGroup)
	wg.Add(18)
	// wg.Add(6)
	go d.ComputeLines(wg, d.Roads.Motorway, 5)
	go d.ComputeLines(wg, d.Roads.MotorwayLink, 5)
	go d.ComputeLines(wg, d.Roads.Trunk, 5)
	go d.ComputeLines(wg, d.Roads.TrunkLink, 5)
	go d.ComputeLines(wg, d.Roads.Primary, 4)
	go d.ComputeLines(wg, d.Roads.PrimaryLink, 4)
	go d.ComputeLines(wg, d.Roads.Secondary, 3)
	go d.ComputeLines(wg, d.Roads.SecondaryLink, 3)
	go d.ComputeLines(wg, d.Roads.Tertiary, 2)
	go d.ComputeLines(wg, d.Roads.TertiaryLink, 2)
	go d.ComputeLines(wg, d.Roads.Unclassified, 1)
	go d.ComputeLines(wg, d.Roads.Residential, 1)
	go d.ComputeLines(wg, d.Roads.Service, 1)
	go d.ComputeLines(wg, d.Roads.Pedestrian, 1)
	go d.ComputeLines(wg, d.Roads.Track, 1)
	go d.ComputeLines(wg, d.Roads.Footway, 1)
	go d.ComputeLines(wg, d.Roads.Steps, 1)
	go d.ComputeLines(wg, d.Roads.Path, 1)
	wg.Wait()
	d.DrawRoads()

}

func (r *Roads) sortRoads(data *osm.Map) {
	for _, way := range data.Ways {
		for _, tag := range way.RTags {
			switch tag.Value {
			case "motorway":
				r.Motorway = append(r.Motorway, way)
			case "motorway_link":
				r.MotorwayLink = append(r.MotorwayLink, way)
			case "trunk":
				r.Trunk = append(r.Trunk, way)
			case "trunk_link":
				r.TrunkLink = append(r.TrunkLink, way)
			case "primary":
				r.Primary = append(r.Primary, way)
			case "primary_link":
				r.PrimaryLink = append(r.PrimaryLink, way)
			case "secondary":
				r.Secondary = append(r.Secondary, way)
			case "secondary_link":
				r.SecondaryLink = append(r.SecondaryLink, way)
			case "tertiary":
				r.Tertiary = append(r.Tertiary, way)
			case "tertiary_link":
				r.TertiaryLink = append(r.TertiaryLink, way)
			case "unclassified":
				r.Unclassified = append(r.Unclassified, way)
			case "residential":
				r.Residential = append(r.Residential, way)
			case "service":
				r.Service = append(r.Service, way)
			case "pedestrian":
				r.Pedestrian = append(r.Pedestrian, way)
			case "track":
				r.Track = append(r.Track, way)
			case "footway":
				r.Footway = append(r.Footway, way)
			case "steps":
				r.Steps = append(r.Steps, way)
			case "path":
				r.Path = append(r.Path, way)
			}
		}

	}
}

func makeURL(bbox BoundingBox) string {
	// TO ADJUST
	// minlat, minlon, maxlat, maxlon
	return fmt.Sprintf("https://overpass.kumi.systems/api/map?bbox=%v,%v,%v,%v", bbox.westLongitude, bbox.southLatitude, bbox.eastLongitude, bbox.northLatitude)
}

func retrieveOSMData(bbox BoundingBox) (*osm.Map, error) {

	fmt.Println(makeURL(bbox))
	resp, err := http.Get(makeURL(bbox))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// err = os.WriteFile("data", body, 0644)
	// if err != nil {
	// 	return nil, err
	// }

	return osm.DecodeString(string(body))
}
