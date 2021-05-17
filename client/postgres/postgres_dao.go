package postgres

import (
	"micro/config"
	"sync"

	"github.com/go-pg/pg/v10"
)

var (
	Storage store = &psql{}
	once    sync.Once
)

// store interface is interface for store things into postgres
type store interface {
	Connect(config config.Config) error
	Get() *pg.DB
}

// postgres struct
type psql struct {
	db *pg.DB
}
