package main

import (
	"atmoswx.com/shapefile-to-geojson/internal/converter"
)

func main() {
	converter.Convert("test-data/cb_2018_us_state_20m.shp", "output.json")
}
