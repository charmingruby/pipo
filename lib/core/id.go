package core

import "github.com/oklog/ulid/v2"

func New() string {
	return ulid.Make().String()
}
