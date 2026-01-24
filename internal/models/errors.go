package models

import "errors"

var (
	ErrUserAlreadyExists             = errors.New("user already exists")
	ErrCapitolAlreadyExists          = errors.New("capitol planet already exists")
	ErrCapitolNotFound               = errors.New("capitol planet not found")
	ErrPlanetCoordinatesAlreadyTaken = errors.New("planet coordinates already taken")
)
