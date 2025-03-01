package postgres

import (
	"fmt"
	"sync"
	"time"

	config "github.com/loongkirin/go-family-finance/pkg/database"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type PostgresDbContext struct {
	DbConfig *config.DbConfig
	master   *gorm.DB
	slaves   []*gorm.DB
	lock     sync.RWMutex
	current  int
}

func NewPostgresDbContext(cfg *config.DbConfig) (*PostgresDbContext, error) {
	master, err := connectDB(cfg.Master)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master: %w", err)
	}

	if err := master.Use(tracing.NewPlugin()); err != nil {
		return nil, fmt.Errorf("failed to enable tracing: %w", err)
	}

	var slaves []*gorm.DB
	for _, slaveCfg := range cfg.Slaves {
		slave, err := connectDB(slaveCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to slave: %w", err)
		}
		if err := slave.Use(tracing.NewPlugin()); err != nil {
			return nil, fmt.Errorf("failed to enable tracing: %w", err)
		}
		slaves = append(slaves, slave)
	}

	return &PostgresDbContext{
		DbConfig: cfg,
		master:   master,
		slaves:   slaves,
	}, nil
}

func connectDB(cfg config.DBConnection) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable %s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.Config,
	)
	pgsqlconfig := gormpostgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false,
	}
	db, err := gorm.Open(gormpostgres.New(pgsqlconfig), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	if duration, err := time.ParseDuration(cfg.ConnMaxLifetime); err == nil {
		sqlDB.SetConnMaxLifetime(duration)
	}
	return db, nil
}

func (db *PostgresDbContext) GetMasterDb() *gorm.DB {
	return db.master
}

func (db *PostgresDbContext) GetSlaveDb() *gorm.DB {
	if len(db.slaves) == 0 {
		return db.master
	}

	db.lock.Lock()
	defer db.lock.Unlock()

	db.current = (db.current + 1) % len(db.slaves)
	return db.slaves[db.current]
}
