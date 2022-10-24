package internal

import (
	"fmt"
	"io"
	"net/http"

	"github.com/glaslos/go-osm"
)

func makeURL(bbox BoundingBox) string {
	// TO ADJUST
	// minlat, minlon, maxlat, maxlon
	return fmt.Sprintf("https://overpass.kumi.systems/api/map?bbox=%v,%v,%v,%v", bbox.WestLongitude, bbox.SouthLatitude, bbox.EastLongitude, bbox.NorthLatitude)
}

func RetrieveOSMData(bbox BoundingBox) (*osm.Map, error) {

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
