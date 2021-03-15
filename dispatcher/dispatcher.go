package dispatcher

import (
	"context"
	"github.com/ivanovaleksey/simdrone/drone"
	"github.com/ivanovaleksey/simdrone/pkg/types"
	"log"
	"sync"
)

type Dispatcher struct {
	storage       Storage
	nearbyService drone.Storage

	wg     sync.WaitGroup
	cancel context.CancelFunc
	// stations []types.Station
}

type Storage interface {
	GetDroneIDs(ctx context.Context) ([]types.DroneID, error)
	GetStations(ctx context.Context) ([]types.Station, error)
	GetDroneMoves(ctx context.Context, id types.DroneID) ([]types.DroneMove, error)
}

type Drone interface {
	Start()
}

func New(storage Storage, nearby drone.Storage) Dispatcher {
	return Dispatcher{
		storage:       storage,
		nearbyService: nearby,
	}
}

func (d *Dispatcher) Start(ctx context.Context) error {
	ctx, d.cancel = context.WithCancel(ctx)

	droneIDs, err := d.storage.GetDroneIDs(ctx)
	if err != nil {
		return err
	}

	for _, droneID := range droneIDs {
		d.wg.Add(1)

		go func(id types.DroneID) {
			defer d.wg.Done()

			err := d.startDrone(ctx, id)
			if err != nil {
				log.Printf("can't start drone %d: %s\n", id, err)
				return
			}
		}(droneID)
	}
	return nil
}

func (d *Dispatcher) Close() error {
	d.cancel()
	d.wg.Wait()
	return nil
}

func (d *Dispatcher) startDrone(ctx context.Context, droneID types.DroneID) error {
	const droneCapacity = 10

	moves, err := d.storage.GetDroneMoves(ctx, droneID)
	if err != nil {
		return err
	}

	droneMoves := make(chan types.DroneMove, droneCapacity)
	go func() {
		defer close(droneMoves)

		for i := range moves {
			select {
			case <-ctx.Done():
				return
			case droneMoves <- moves[i]:
			}
		}
	}()

	log.Printf("start drone %d\n", droneID)

	dr := drone.New(droneID, d.nearbyService)
	return dr.Start(ctx, droneMoves)
}
