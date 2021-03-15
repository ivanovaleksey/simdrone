package drone

import (
	"context"
	"fmt"
	"github.com/ivanovaleksey/simdrone/pkg/distance"
	"github.com/ivanovaleksey/simdrone/pkg/types"
	"time"
)

type Drone struct {
	id      types.DroneID
	storage Storage

	lastMove     *types.DroneMove
	currentSpeed float64
}

type Storage interface {
	FindNearbyStations(ctx context.Context, p types.Position) ([]types.Station, error)
}

func New(id types.DroneID, storage Storage) *Drone {
	return &Drone{
		id:      id,
		storage: storage,
	}
}

func (d *Drone) Start(ctx context.Context, moves <-chan types.DroneMove) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case move, ok := <-moves:
			if !ok {
				return nil
			}
			res, err := d.makeMove(ctx, move)
			if err != nil {
				return err
			}
			if res.Ok {
				time.Sleep(res.Duration /10)
				fmt.Printf("%d %s %.3f %s\n", d.id, res.Timestamp.Format(time.RFC3339), res.Speed, res.Conditions)
			}
		}
	}
}

type MoveResult struct {
	Ok         bool
	Speed      float64
	Timestamp  time.Time
	Conditions types.TrafficConditions
	Duration   time.Duration
}

func (d *Drone) makeMove(ctx context.Context, move types.DroneMove) (MoveResult, error) {
	if d.lastMove == nil {
		d.lastMove = &move
		return MoveResult{}, nil
	}

	lastMove := *d.lastMove
	speed, duration := calcSpeed(lastMove, move)
	if speed == 0 {
		speed = d.currentSpeed
	}
	d.lastMove = &move
	d.currentSpeed = speed

	nearbys, err := d.storage.FindNearbyStations(ctx, move.Position)
	if err != nil {
		return MoveResult{}, err
	}
	if len(nearbys) == 0 {
		return MoveResult{}, nil
	}

	res := MoveResult{
		Ok:         true,
		Speed:      speed,
		Timestamp:  move.Timestamp,
		Duration:   duration,
		Conditions: types.RandomTrafficConditions(),
	}

	return res, nil
}

func calcSpeed(prev, next types.DroneMove) (speed float64, flightDuration time.Duration) {
	distanceBetween := distance.Between(prev.Position, next.Position)
	flightDuration = next.Timestamp.Sub(prev.Timestamp)
	timeBetween := flightDuration.Seconds()
	if timeBetween == 0 {
		return 0, 0
	}
	return distanceBetween / timeBetween, flightDuration
}
