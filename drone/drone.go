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
}

func (d *Drone) makeMove(ctx context.Context, move types.DroneMove) (MoveResult, error) {
	if d.lastMove == nil {
		d.lastMove = &move
		return MoveResult{}, nil
	}

	nearbys, err := d.storage.FindNearbyStations(ctx, move.Position)
	if err != nil {
		return MoveResult{}, err
	}
	if len(nearbys) == 0 {
		return MoveResult{}, nil
	}

	lastMove := *d.lastMove
	speed := calcSpeed(lastMove, move)
	if speed == 0 {
		speed = d.currentSpeed
	}

	res := MoveResult{
		Ok:         true,
		Speed:      speed,
		Timestamp:  move.Timestamp,
		Conditions: types.RandomTrafficConditions(),
	}

	d.lastMove = &move
	d.currentSpeed = speed
	return res, nil
}

func calcSpeed(prev, next types.DroneMove) float64 {
	distanceBetween := distance.Between(prev.Position, next.Position)
	timeBetween := next.Timestamp.Sub(prev.Timestamp).Seconds()
	if timeBetween == 0 {
		return 0
	}
	return distanceBetween / timeBetween
}
