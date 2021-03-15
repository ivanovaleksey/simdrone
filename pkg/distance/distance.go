package distance

import (
	"github.com/ivanovaleksey/simdrone/pkg/types"
	"math"
)

func Between(p1, p2 types.Position) float64 {
	const earthRadius = 6371e3

	p1 = p1.ToRadians()
	p2 = p2.ToRadians()

	delta := types.Position{
		Lat: p2.Lat - p1.Lat,
		Lon: p2.Lon - p1.Lon,
	}

	a := math.Pow(math.Sin(delta.Lat/2), 2) + math.Cos(p1.Lat)*math.Cos(p2.Lat)*math.Pow(math.Sin(delta.Lon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}
