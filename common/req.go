package common

import (
	"vvorker/entities"
)

type Request interface {
	*entities.DeleteWorkerRequest | *entities.LoginRequest | *entities.RegisterRequest |
		*entities.NotifyEventRequest
	Validate() bool
}
