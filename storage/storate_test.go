package storage

import (
	"context"
	"github.com/ivanovaleksey/simdrone/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStorage_GetStations(t *testing.T) {
	fx := newFixture(t)
	defer fx.Finish()

	stations, err := fx.storage.GetStations(fx.ctx)

	require.NoError(t, err)
	require.Greater(t, len(stations), 0)

	expected := types.Station{
		Title: "Acton Town",
		Position: types.Position{
			Lat: 51.503071,
			Lon: -0.280303,
		},
	}
	assert.Equal(t, expected, stations[0])
}

func TestStorage_GetDroneMoves(t *testing.T) {
	t.Run("with known drone ID", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		moves, err := fx.storage.GetDroneMoves(fx.ctx, 6043)

		require.NoError(t, err)
		require.Greater(t, len(moves), 0)

		expected := types.DroneMove{
			DroneID: 6043,
			Position: types.Position{
				Lat: 51.474579,
				Lon: -0.171834,
			},
			Timestamp: time.Date(2011, 3, 22, 07, 47, 55, 0, time.UTC),
		}
		assert.Equal(t, expected, moves[0])
	})

	t.Run("with unknown drone ID", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		_, err := fx.storage.GetDroneMoves(fx.ctx, 128)

		require.Error(t, err)
	})
}

type fixture struct {
	t       *testing.T
	ctx     context.Context
	storage *Storage
}

func newFixture(t *testing.T) *fixture {
	fx := fixture{
		t:       t,
		ctx:     context.Background(),
		storage: New("data"),
	}
	return &fx
}

func (fx *fixture) Finish() {}
