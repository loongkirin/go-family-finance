package dbfactory

import (
	"github.com/loongkirin/go-family-finance/pkg/database"
	"github.com/loongkirin/go-family-finance/pkg/database/postgres"
)

func CreateDbContext(cfg database.DbConfig) (database.DbContext, error) {
	var dbcontext database.DbContext
	switch cfg.DbType {
	case "postgres":
		pgDbContext, err := postgres.NewPostgresDbContext(&cfg)
		if err != nil {
			return nil, err
		}
		dbcontext = pgDbContext
	}

	return dbcontext, nil
}
