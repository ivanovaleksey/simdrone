package storage

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/ivanovaleksey/simdrone/pkg/distance"
	"github.com/ivanovaleksey/simdrone/pkg/types"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var droneIDs = []types.DroneID{5937, 6043}

type Storage struct {
	stations     []types.Station
	stationsErr  error
	stationsOnce sync.Once

	dataDir string
}

func New(dataDir string) *Storage {
	return &Storage{
		dataDir: dataDir,
	}
}

func (s *Storage) GetDroneIDs(_ context.Context) ([]types.DroneID, error) {
	return droneIDs, nil
}

func (s *Storage) GetStations(_ context.Context) ([]types.Station, error) {
	stations, err := s.loadStationsOnce()
	if err != nil {
		return nil, errors.Wrap(err, "can't load stations")
	}
	return stations, err
}

func (s *Storage) loadStationsOnce() ([]types.Station, error) {
	s.stationsOnce.Do(func() {
		s.stations, s.stationsErr = s.loadStations()
	})
	return s.stations, s.stationsErr
}

func (s *Storage) loadStations() ([]types.Station, error) {
	const fileName = "tube.csv"

	file, err := os.Open(filepath.Join(s.dataDir, fileName))
	if err != nil {
		return nil, errors.Wrap(err, "can't open stations file")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	stations := make([]types.Station, 0, len(rows))
	for _, row := range rows {
		lat, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse station latitude")
		}
		lon, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse station longitude")
		}

		stations = append(stations, types.Station{
			Title: row[0],
			Position: types.Position{
				Lat: lat,
				Lon: lon,
			},
		})
	}
	return stations, nil
}

func (s *Storage) GetDroneMoves(_ context.Context, id types.DroneID) ([]types.DroneMove, error) {
	const dateLayout = "2006-01-02 15:04:05"

	file, err := os.Open(filepath.Join(s.dataDir, fmt.Sprintf("%d.csv", id)))
	if err != nil {
		return nil, errors.Wrap(err, "can't open moves file")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	moves := make([]types.DroneMove, 0, len(rows))
	for _, row := range rows {
		droneID, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse drone ID")
		}
		lat, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse move latitude")
		}
		lon, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse move longitude")
		}
		dateTime, err := time.Parse(dateLayout, row[3])
		if err != nil {
			return nil, errors.Wrap(err, "can't parse move time")
		}

		moves = append(moves, types.DroneMove{
			DroneID: droneID,
			Position: types.Position{
				Lat: lat,
				Lon: lon,
			},
			Timestamp: dateTime,
		})
	}
	return moves, nil
}

func (s *Storage) FindNearbyStations(ctx context.Context, pos types.Position) ([]types.Station, error) {
	stations, err := s.GetStations(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can't get stations")
	}

	var nearby []types.Station
	for _, station := range stations {
		if isNearby(station.Position, pos) {
			nearby = append(nearby, station)
		}
	}
	return nearby, nil
}

func isNearby(p1, p2 types.Position) bool {
	const nearbyDistance float64 = 350
	return distance.Between(p1, p2) < nearbyDistance
}
