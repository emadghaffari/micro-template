package postgres

import (
	"micro/config"
	"sync"
)

var (
	Storage store = &psql{}
	once    sync.Once
)

// store interface is interface for store things into postgres
type store interface {
	Connect(config config.GlobalConfig) error
}

// postgres struct
type psql struct {
	// db *gorm.DB
}

// Connect method job is connect to postgres database and check migration
func (m *psql) Connect(config config.GlobalConfig) error {
	// logger := zapLogger.GetZapLogger(config.Debug())
	var err error
	once.Do(func() {})

	return err
}
