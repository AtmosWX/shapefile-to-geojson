package converter

import (
	"io/ioutil"

	"atmoswx.com/shapefile-to-geojson/internal/reader"
)

func readFile(path string) (*reader.Reader, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return reader.NewReader(data), nil
}

func writeFile(path string, data []byte) {
}
