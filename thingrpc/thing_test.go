package thingrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThingString(t *testing.T) {
	thing := &Thing{
		Id:   "id1",
		Name: "name1",
	}

	assert.Equal(t, `id:"id1" name:"name1" `, thing.String())
}
