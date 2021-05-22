package postgres

import (
	"context"
	"micro/config"
	zapLogger "micro/pkg/logger"

	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

// Connect method job is connect to postgres database and check migration
func (p *psql) Connect(cnf config.Config) error {
	var err error

	once.Do(func() {
		p.db = pg.Connect(&pg.Options{
			User:                  cnf.POSTGRES.Username,
			Password:              cnf.POSTGRES.Password,
			Addr:                  cnf.POSTGRES.Host,
			Database:              cnf.POSTGRES.Schema,
			RetryStatementTimeout: true,
		})
		if err = p.db.Ping(context.Background()); err != nil {
			zapLogger.Prepare(zapLogger.GetZapLogger(config.Confs.GetDebug())).Development().Level(zap.ErrorLevel).Commit("init configs")
			return
		}

	})

	return err
}

func (p *psql) Get() *pg.DB {
	return p.db
}
