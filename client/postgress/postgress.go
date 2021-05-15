package postgress

import (
	"micro/config"
	"sync"
)

var (
	Storage store = &psql{}
	once    sync.Once
)

// store interface is interface for store things into postgress
type store interface {
	Connect(config config.GlobalConfig) error
}

// postgress struct
type psql struct {
	// db *gorm.DB
}

// Connect method job is connect to postgress database and check migration
func (m *psql) Connect(config config.GlobalConfig) error {
	// logger := zapLogger.GetZapLogger(config.Debug())
	var err error
	once.Do(func() {})

	return err
}
