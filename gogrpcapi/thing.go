package gogrpcapi

import (
	"context"
)

// ThingStore is the persistent store of things
type ThingStore interface {
	ThingGetById(context.Context, string) (*Thing, error)
	ThingSave(context.Context, *Thing) (string, error)
	ThingDeleteById(context.Context, string) error
	ThingFind(context.Context) ([]*Thing, error)
}
