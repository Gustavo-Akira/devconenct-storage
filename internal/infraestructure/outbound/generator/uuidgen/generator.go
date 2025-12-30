package uuidgen

import "github.com/google/uuid"

type UUIDGenerator struct{}

func (UUIDGenerator) NewID() string {
	return uuid.NewString()
}
