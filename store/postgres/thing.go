package postgres

import (
	"context"
	"database/sql"

	"github.com/snowzach/gogrpcapi/gogrpcapi"
	"github.com/snowzach/gogrpcapi/store"
)

// ThingGetByID returns the the thing by ID
func (c *Client) ThingGetById(ctx context.Context, id string) (*gogrpcapi.Thing, error) {

	b := new(gogrpcapi.Thing)
	err := c.db.GetContext(ctx, b, `SELECT * FROM thing WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return b, nil

}

// ThingSave saves the thing
func (c *Client) ThingSave(ctx context.Context, i *gogrpcapi.Thing) (string, error) {

	// Generate an ID if needed
	if i.Id == "" {
		i.Id = c.newID()
	}

	_, err := c.db.ExecContext(ctx, `
		INSERT INTO thing (id, name)
		VALUES($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = $2
	`, i.Id, i.Name)
	if err != nil {
		return i.Id, err
	}
	return i.Id, nil

}

// ThingDeleteById a thing
func (c *Client) ThingDeleteById(ctx context.Context, id string) error {

	_, err := c.db.ExecContext(ctx, `DELETE FROM thing WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil

}

// ThingFind gets things
func (c *Client) ThingFind(ctx context.Context) ([]*gogrpcapi.Thing, error) {

	var bs = make([]*gogrpcapi.Thing, 0)
	err := c.db.SelectContext(ctx, &bs, `SELECT * FROM thing`)
	if err == sql.ErrNoRows {
		// No Error
	} else if err != nil {
		return bs, err
	}
	return bs, nil

}
