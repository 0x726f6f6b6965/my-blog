package migrations

import (
	"context"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

// MigrateSchema create/migrate the schema
func MigrateSchema(ctx context.Context, uri string) error {
	resource := bindata.Resource(AssetNames(), func(name string) ([]byte, error) {
		return Asset(name)
	})
	driver, err := bindata.WithInstance(resource)
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("go-bindata", driver, uri)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Applied %d migrations!\n", 0)
			return nil
		}
		return err
	}

	v, _, _ := m.Version()
	log.Printf("Applied %d migrations!\n", v)
	return nil
}
