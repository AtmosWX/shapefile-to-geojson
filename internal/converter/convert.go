package converter

import (
	"encoding/json"
	"io/ioutil"
)

type GeoJSON struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string                 `json:"type"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

func Convert(path string, output string) error {
	reader, err := readFile(path)

	if err != nil {
		return err
	}

	shapefile := parseShapefileContent(reader)

	features := []Feature{}

	for _, point := range shapefile.Records.Points {
		feature := Feature{}
		feature.Type = "Feature"
		feature.Properties = map[string]interface{}{}
		feature.Geometry.Type = "Point"
		feature.Geometry.Coordinates = []float64{point.X, point.Y}
		features = append(features, feature)
	}

	for _, polyLine := range shapefile.Records.PolyLines {
		feature := Feature{}
		feature.Type = "Feature"
		feature.Properties = map[string]interface{}{}
		feature.Geometry.Type = "PolyLine"

		coordinates := [][][]float64{}

		for i, part := range polyLine.Parts {
			points := []Point{}

			if i < len(polyLine.Parts)-1 {
				points = polyLine.Points[part:polyLine.Parts[i+1]]
			} else {
				points = polyLine.Points[part:]
			}

			coords := [][]float64{}

			for _, point := range points {
				coords = append(coords, []float64{point.X, point.Y})
			}

			coordinates = append(coordinates, coords)
		}

		feature.Geometry.Coordinates = coordinates
		features = append(features, feature)
	}

	for _, polygon := range shapefile.Records.Polygons {
		feature := Feature{}
		feature.Type = "Feature"
		feature.Properties = map[string]interface{}{}
		feature.Geometry.Type = "Polygon"

		coordinates := [][][]float64{}

		for i, part := range polygon.Parts {
			points := []Point{}

			if i < len(polygon.Parts)-1 {
				points = polygon.Points[part:polygon.Parts[i+1]]
			} else {
				points = polygon.Points[part:]
			}

			coords := [][]float64{}

			for _, point := range points {
				coords = append(coords, []float64{point.X, point.Y})
			}

			coordinates = append(coordinates, coords)
		}

		feature.Geometry.Coordinates = coordinates
		features = append(features, feature)
	}

	for _, multiPoint := range shapefile.Records.MultiPoints {
		feature := Feature{}
		feature.Type = "Feature"
		feature.Properties = map[string]interface{}{}
		feature.Geometry.Type = "MultiPoint"

		coordinates := [][]float64{}
		for _, point := range multiPoint.Points {
			coordinates = append(coordinates, []float64{point.X, point.Y})
		}

		feature.Geometry.Coordinates = coordinates
		features = append(features, feature)
	}

	geojson := GeoJSON{}
	geojson.Type = "FeatureCollection"
	geojson.Features = features

	data, err := json.MarshalIndent(geojson, "", " ")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(output, data, 0644)

	if err != nil {
		return err
	}

	return nil
}
