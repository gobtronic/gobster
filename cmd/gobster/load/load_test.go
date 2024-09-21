package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModel(t *testing.T) {
	model := NewModel()
	assert.Nil(t, model.err)
}
