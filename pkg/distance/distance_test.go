package distance

import (
	"github.com/ivanovaleksey/simdrone/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBetween(t *testing.T) {
	p1 := types.Position{
		Lat: 51.476105,
		Lon: -0.100224,
	}
	p2 := types.Position{
		Lat: 51.475967,
		Lon: -0.100368,
	}

	dist := Between(p1, p2)

	assert.Equal(t, 18.30099559243436, dist)
}
