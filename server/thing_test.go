package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/snowzach/gogrpcapi/gogrpcapi"
	"github.com/snowzach/gogrpcapi/mocks"
	"github.com/snowzach/gogrpcapi/server/rpc"
)

func TestServerThingPost(t *testing.T) {

	// Mock Store and server
	ts := new(mocks.ThingStore)
	s, err := New(ts)
	assert.Nil(t, err)

	// Create Item
	i := &gogrpcapi.Thing{
		Id:   "id",
		Name: "name",
	}

	// Mock call to item store
	ts.On("ThingSave", mock.AnythingOfType("*context.emptyCtx"), i).Once().Return(i.Id, nil)

	response, err := s.ThingSave(context.Background(), i)
	assert.Nil(t, err)
	assert.Equal(t, response.Id, i.Id)

	// Check remaining expectations
	ts.AssertExpectations(t)

}

func TestServerThingGetAll(t *testing.T) {

	// Mock Store and server
	ts := new(mocks.ThingStore)
	s, err := New(ts)
	assert.Nil(t, err)

	// Create Item
	i := []*gogrpcapi.Thing{
		&gogrpcapi.Thing{
			Id:   "id1",
			Name: "name1",
		},
		&gogrpcapi.Thing{
			Id:   "id2",
			Name: "name2",
		},
	}

	// Mock call to item store
	ts.On("ThingFind", mock.AnythingOfType("*context.emptyCtx")).Once().Return(i, nil)

	response, err := s.ThingFind(context.Background(), nil)
	assert.Nil(t, err)
	assert.Equal(t, i, response.Data)

	// Check remaining expectations
	ts.AssertExpectations(t)

}

func TestServerThingGet(t *testing.T) {

	// Mock Store and server
	ts := new(mocks.ThingStore)
	s, err := New(ts)
	assert.Nil(t, err)

	// Create Item
	i := &gogrpcapi.Thing{
		Id:   "id",
		Name: "name",
	}

	// Mock call to item store
	ts.On("ThingGetById", mock.AnythingOfType("*context.emptyCtx"), "1234").Once().Return(i, nil)

	response, err := s.ThingGet(context.Background(), &rpc.ThingId{Id: "1234"})
	assert.Nil(t, err)
	assert.Equal(t, i, response)

	// Check remaining expectations
	ts.AssertExpectations(t)

}

func TestServerThingDelete(t *testing.T) {

	// Mock Store and server
	ts := new(mocks.ThingStore)
	s, err := New(ts)
	assert.Nil(t, err)

	// Mock call to item store
	ts.On("ThingDeleteById", mock.AnythingOfType("*context.emptyCtx"), "1234").Once().Return(nil)

	_, err = s.ThingDelete(context.Background(), &rpc.ThingId{Id: "1234"})
	assert.Nil(t, err)

	// Check remaining expectations
	ts.AssertExpectations(t)

}
