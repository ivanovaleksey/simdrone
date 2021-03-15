package types

import "math"

type Position struct {
	Lat Latitude
	Lon Longitude
}

func (p Position) ToRadians() Position {
	return Position{
		Lat: degreesToRadians(p.Lat),
		Lon: degreesToRadians(p.Lon),
	}
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

