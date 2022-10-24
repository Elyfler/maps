package internal

import (
	"image/color"
	"log"
	"math"
	"sync"

	"github.com/glaslos/go-osm"
	"gopkg.in/fogleman/gg.v1"
)

type BoundingBox struct {
	SouthLatitude float64
	NorthLatitude float64
	WestLongitude float64
	EastLongitude float64
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

func (l *Line) toPoints() (point1 gg.Point, point2 gg.Point) {
	return gg.Point{X: l.x1, Y: l.y1}, gg.Point{X: l.x2, Y: l.y2}
}

type Colors struct {
	BackgroundColor color.Color
	RoadColor       color.Color
	WaterColor      color.Color
}

type Drawing struct {
	Context     *gg.Context
	Roads       Roads
	Waters      Waters
	Lines       []Line
	PointsWater map[int64][]gg.Point
	Nodes       []osm.Node
	Colors      Colors
}

func NewDrawing(width, height int, data *osm.Map) Drawing {
	var d Drawing
	d.Context = gg.NewContext(width, height)
	d.PointsWater = make(map[int64][]gg.Point)
	d.Roads.sortRoads(data)
	d.Waters.sortWaters(data)
	d.Nodes = data.Nodes
	return d
}

func (d *Drawing) SetColorScheme(colors Colors) {
	d.Colors = colors
}

func (d *Drawing) DrawRoads() {
	log.Println("DrawRoads")
	d.Context.SetColor(d.Colors.BackgroundColor)
	d.Context.Clear()
	d.Context.SetColor(d.Colors.RoadColor)
	for _, line := range d.Lines {
		d.Context.SetLineWidth(line.lineWidth)
		d.Context.DrawLine(line.toFunc())
		d.Context.Stroke()
	}
}

func (d *Drawing) DrawWater() {
	log.Println("DrawWater")
	d.Context.SetColor(d.Colors.WaterColor)
	d.Context.SetLineWidth(5)
	d.Context.NewSubPath()
	for id := range d.PointsWater {
		for i, point := range d.PointsWater[id] {
			log.Println(point)
			if point.X > float64(d.Context.Width())-1 {
				log.Println("Reached the end")
			}
			if point.Y > float64(d.Context.Height())-1 {
				log.Println("Reached the Y end")
			}
			if i == 0 {
				d.Context.MoveTo(point.X, point.Y)
			} else {
				d.Context.LineTo(point.X, point.Y)
				d.Context.FillPreserve()
			}
		}
	}
}

func (d *Drawing) ComputeLines(wg *sync.WaitGroup, ways []osm.Way, lineWidth float64) {
	defer wg.Done()
	// Ugly shit to rework later on
	for _, way := range ways {
		var nodes []osm.Node
		for _, ID := range way.Nds {
			for _, node := range d.Nodes {
				if node.ID == ID.ID {
					nodes = append(nodes, node)
				}
			}
		}
		for i := 0; i < len(nodes)-1; i++ {
			line := Line{
				x1:        nodes[i].Lng,
				y1:        nodes[i].Lat,
				x2:        nodes[i+1].Lng,
				y2:        nodes[i+1].Lat,
				lineWidth: lineWidth,
			}
			log.Println(line)
			d.Lines = append(d.Lines, line)
		}
	}
}

func (d *Drawing) ComputeWaters(ways []osm.Way, lineWidth float64) {
	for _, water := range ways {
		for _, ID := range water.Nds {
			for _, node := range d.Nodes {
				if node.ID == ID.ID {
					d.PointsWater[water.ID] = append(d.PointsWater[water.ID], gg.Point{X: node.Lng, Y: node.Lat})
				}
			}
		}
	}
}

func (d *Drawing) Normalize(bbox BoundingBox) {
	maxHeight := math.Abs(bbox.NorthLatitude - bbox.SouthLatitude)
	maxWidth := math.Abs(bbox.WestLongitude - bbox.EastLongitude)
	for k := range d.Nodes {
		d.Nodes[k].Lat = ((bbox.NorthLatitude - d.Nodes[k].Lat) / maxHeight) * float64(d.Context.Height())
		d.Nodes[k].Lng = ((-bbox.WestLongitude + d.Nodes[k].Lng) / maxWidth) * float64(d.Context.Width())
	}
}

func (d *Drawing) ComputeAll() {
	wg := new(sync.WaitGroup)
	wg.Add(18)
	go d.ComputeLines(wg, d.Roads.Motorway, 6)
	go d.ComputeLines(wg, d.Roads.MotorwayLink, 6)
	go d.ComputeLines(wg, d.Roads.Trunk, 6)
	go d.ComputeLines(wg, d.Roads.TrunkLink, 6)
	go d.ComputeLines(wg, d.Roads.Primary, 5)
	go d.ComputeLines(wg, d.Roads.PrimaryLink, 5)
	go d.ComputeLines(wg, d.Roads.Secondary, 4)
	go d.ComputeLines(wg, d.Roads.SecondaryLink, 4)
	go d.ComputeLines(wg, d.Roads.Tertiary, 3)
	go d.ComputeLines(wg, d.Roads.TertiaryLink, 3)
	go d.ComputeLines(wg, d.Roads.Unclassified, 2)
	go d.ComputeLines(wg, d.Roads.Residential, 2)
	go d.ComputeLines(wg, d.Roads.Service, 2)
	go d.ComputeLines(wg, d.Roads.Pedestrian, 2)
	go d.ComputeLines(wg, d.Roads.Track, 2)
	go d.ComputeLines(wg, d.Roads.Footway, 2)
	go d.ComputeLines(wg, d.Roads.Steps, 2)
	go d.ComputeLines(wg, d.Roads.Path, 2)
	wg.Wait()
	d.ComputeWaters(d.Waters.Water, 6)
	d.ComputeWaters(d.Waters.Bay, 6)
	d.ComputeWaters(d.Waters.Coastline, 6)
}

func CalculateSize(b BoundingBox, scale float64) (width int, height int) {
	log.Println("calculateSize")
	baseWidth := math.Abs(b.WestLongitude - b.EastLongitude)
	baseHeight := b.NorthLatitude - b.SouthLatitude
	// log.Println(baseWidth)
	// log.Println(baseHeight)
	return int(scale * baseWidth), int(scale * baseHeight)
}
