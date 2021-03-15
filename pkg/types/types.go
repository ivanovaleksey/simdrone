package types

import (
	"math/rand"
	"time"
)

type (
	DroneID   = int64
	Latitude  = float64
	Longitude = float64
)

type Station struct {
	Title string
	Position
}

type DroneMove struct {
	DroneID DroneID
	Position
	Timestamp time.Time
}

type TrafficConditions int

const (
	TrafficConditionsLight TrafficConditions = iota + 1
	TrafficConditionsModerate
	TrafficConditionsHeavy
)

func (tc TrafficConditions) String() string {
	switch tc {
	case TrafficConditionsLight:
		return "light"
	case TrafficConditionsModerate:
		return "moderate"
	case TrafficConditionsHeavy:
		return "heavy"
	}
	return "unknown"
}

var allConditions = []TrafficConditions{TrafficConditionsLight, TrafficConditionsModerate, TrafficConditionsHeavy}

func RandomTrafficConditions() TrafficConditions {
	idx := rand.Intn(len(allConditions))
	return allConditions[idx]
}
