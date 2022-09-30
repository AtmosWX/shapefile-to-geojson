package converter

import (
	"atmoswx.com/shapefile-to-geojson/internal/reader"
)

type Shapefile struct {
	Header  ShapefileHeader
	Records RecordTable
}

func parseShapefileContent(rdr *reader.Reader) Shapefile {
	header := readShapfileHeader(rdr)
	records := readRecords(rdr)

	return Shapefile{header, records}
}

type ShapefileHeader struct {
	FileCode     uint32
	FileLength   uint32
	Version      uint32
	ShapeType    uint32
	BoundingRect []float64
	ZRange       []float64
	MRange       []float64
}

func readShapfileHeader(rdr *reader.Reader) ShapefileHeader {
	header := ShapefileHeader{}

	header.FileCode = rdr.ReadUint()
	rdr.StepForward(20)
	header.FileLength = rdr.ReadUint()

	rdr.SetMode(reader.LittleEndian)

	header.Version = rdr.ReadUint()
	header.ShapeType = rdr.ReadUint()
	header.BoundingRect = []float64{
		rdr.ReadDouble(),
		rdr.ReadDouble(),
		rdr.ReadDouble(),
		rdr.ReadDouble(),
	}
	header.ZRange = []float64{
		rdr.ReadDouble(),
		rdr.ReadDouble(),
	}
	header.MRange = []float64{
		rdr.ReadDouble(),
		rdr.ReadDouble(),
	}

	return header
}

type RecordTable struct {
	Nulls       []NullRecord
	Points      []PointRecord
	PolyLines   []PolyLineRecord
	Polygons    []PolygonRecord
	MultiPoints []MultiPointRecord
}

type RecordHeader struct {
	Number uint32
	Length uint32
}

type NullRecord struct {
	Header    RecordHeader
	ShapeType ShapeType
}

type PointRecord struct {
	Header    RecordHeader
	ShapeType ShapeType
	X         float64
	Y         float64
}

type PolyLineRecord struct {
	Header       RecordHeader
	ShapeType    ShapeType
	BoundingRect []float64
	PartCount    uint32
	PointCount   uint32
	Parts        []uint32
	Points       []Point
}

type PolygonRecord struct {
	Header       RecordHeader
	ShapeType    ShapeType
	BoundingRect []float64
	PartCount    uint32
	PointCount   uint32
	Parts        []uint32
	Points       []Point
}

type MultiPointRecord struct {
	Header       RecordHeader
	ShapeType    ShapeType
	BoundingRect []float64
	PointCount   uint32
	Points       []Point
}

type ShapeType uint32

const (
	NullType       ShapeType = 0
	PointType      ShapeType = 1
	PolyLineType   ShapeType = 3
	PolygonType    ShapeType = 5
	MultiPointType ShapeType = 8
)

type Point struct {
	X float64
	Y float64
}

func readRecords(rdr *reader.Reader) RecordTable {
	records := RecordTable{}

	for int(rdr.Offset) < len(rdr.Data) {
		header := RecordHeader{}

		rdr.SetMode(reader.BigEndian)
		header.Number = rdr.ReadUint()
		header.Length = rdr.ReadUint()

		rdr.SetMode(reader.LittleEndian)
		shapeType := (ShapeType)(rdr.ReadUint())

		switch shapeType {

		case NullType:
			record := NullRecord{}
			record.ShapeType = shapeType

			records.Nulls = append(records.Nulls, record)

		case PointType:
			record := PointRecord{}
			record.ShapeType = shapeType

			record.X = rdr.ReadDouble()
			record.Y = rdr.ReadDouble()

			records.Points = append(records.Points, record)

		case PolyLineType:
			record := PolyLineRecord{}
			record.ShapeType = shapeType

			record.BoundingRect = []float64{
				rdr.ReadDouble(),
				rdr.ReadDouble(),
				rdr.ReadDouble(),
				rdr.ReadDouble(),
			}

			record.PartCount = rdr.ReadUint()
			record.PointCount = rdr.ReadUint()

			for i := 0; i < int(record.PartCount); i++ {
				record.Parts = append(record.Parts, rdr.ReadUint())
			}

			for i := 0; i < int(record.PointCount); i++ {
				point := Point{}
				point.X = rdr.ReadDouble()
				point.Y = rdr.ReadDouble()
				record.Points = append(record.Points, point)
			}

			record.Header = header
			records.PolyLines = append(records.PolyLines, record)

		case PolygonType:
			record := PolygonRecord{}
			record.ShapeType = shapeType

			record.BoundingRect = []float64{
				rdr.ReadDouble(),
				rdr.ReadDouble(),
				rdr.ReadDouble(),
				rdr.ReadDouble(),
			}

			record.PartCount = rdr.ReadUint()
			record.PointCount = rdr.ReadUint()

			for i := 0; i < int(record.PartCount); i++ {
				record.Parts = append(record.Parts, rdr.ReadUint())
			}

			for i := 0; i < int(record.PointCount); i++ {
				point := Point{}
				point.X = rdr.ReadDouble()
				point.Y = rdr.ReadDouble()
				record.Points = append(record.Points, point)
			}

			record.Header = header
			records.Polygons = append(records.Polygons, record)

		case MultiPointType:
			record := MultiPointRecord{}
			record.ShapeType = shapeType

			record.BoundingRect = []float64{
				rdr.ReadDouble(),
				rdr.ReadDouble(),
				rdr.ReadDouble(),
				rdr.ReadDouble(),
			}

			record.PointCount = rdr.ReadUint()

			for i := 0; i < int(record.PointCount); i++ {
				point := Point{}
				point.X = rdr.ReadDouble()
				point.Y = rdr.ReadDouble()
				record.Points = append(record.Points, point)
			}

			records.MultiPoints = append(records.MultiPoints, record)

		default:
			return records
		}

	}

	return records
}
