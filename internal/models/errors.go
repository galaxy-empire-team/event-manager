package models

import "errors"

var (
	ErrCapitolAlreadyExists          = errors.New("capitol planet already exists")
	ErrPlanetCoordinatesAlreadyTaken = errors.New("planet coordinates already taken")
	ErrBuildingNotFound              = errors.New("building not found")
)
