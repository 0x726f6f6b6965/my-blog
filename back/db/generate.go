//go:generate go run ./cmd/main.go --config ./sqlboiler.toml
//go:generate sqlboiler psql --config ./sqlboiler.toml
package tool

import (
	_ "github.com/volatiletech/null/v8"
)
