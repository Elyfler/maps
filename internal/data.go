package internal

import "github.com/glaslos/go-osm"

// Roads ...
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

type Waters struct {
	Bay       []osm.Way
	Coastline []osm.Way
	Water     []osm.Way
	// spring, wetland, peninsula, mud, cape, beach
}

func (w *Waters) sortWaters(data *osm.Map) {
	// bay, coastline, water
	for _, way := range data.Ways {
		for _, tag := range way.RTags {
			switch tag.Value {
			case "bay":
				w.Bay = append(w.Bay, way)
			case "coastline":
				w.Coastline = append(w.Coastline, way)
			case "water":
				w.Water = append(w.Water, way)
			}
		}
	}
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
