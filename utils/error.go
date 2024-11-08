package utils

type Error struct {
	BaseError          error
	StatusCodeToReturn int
}
