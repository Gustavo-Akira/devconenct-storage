package uuidgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUIDGenerator_ShouldGenerateARandomUUID(t *testing.T) {
	generator := &UUIDGenerator{}
	id := generator.NewID()
	assert.NotEmpty(t, id)
}
