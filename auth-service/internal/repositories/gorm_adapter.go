package repositories

import (
	"errors"

	"gorm.io/gorm"
)

// GormAdapter adapter for GORM DB
type GormAdapter struct {
	db *gorm.DB
}

// NewGormAdapter creates a new adapter for GORM
func NewGormAdapter(db *gorm.DB) IDatabase {
	return &GormAdapter{db: db}
}

func (g *GormAdapter) Create(value interface{}) IDatabase {
	if g.db == nil {
		return &GormAdapter{db: nil}
	}
	return &GormAdapter{db: g.db.Create(value)}
}

func (g *GormAdapter) Where(query interface{}, args ...interface{}) IDatabase {
	if g.db == nil {
		return &GormAdapter{db: nil}
	}
	return &GormAdapter{db: g.db.Where(query, args...)}
}

func (g *GormAdapter) First(dest interface{}, conds ...interface{}) IDatabase {
	if g.db == nil {
		return &GormAdapter{db: nil}
	}
	return &GormAdapter{db: g.db.First(dest, conds...)}
}

func (g *GormAdapter) GetError() error {
	if g.db == nil {
		return errors.New("database is nil")
	}
	return g.db.Error
}
