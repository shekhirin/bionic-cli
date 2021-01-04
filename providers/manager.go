package providers

import (
	"errors"
	"fmt"
	"github.com/shekhirin/bionic-cli/providers/provider"
	"github.com/shekhirin/bionic-cli/providers/twitter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var ErrProviderNotFound = errors.New("provider not found")

type Manager struct {
	db        *gorm.DB
	providers map[string]provider.Provider
}

func NewManager(dbPath string) (*Manager, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Manager{
		db: db,
		providers: map[string]provider.Provider{
			"twitter": twitter.New(db),
		},
	}, nil
}

func (m Manager) GetByName(name string) (provider.Provider, error) {
	p, ok := m.providers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, name)
	}

	return p, nil
}

func (m Manager) Migrate(p provider.Provider) error {
	return m.migrate(m.db, p)
}

func (m Manager) Reset(p provider.Provider) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Migrator().DropTable(p.Models()...); err != nil {
			return err
		}

		if err := m.migrate(tx, p); err != nil {
			return err
		}

		return nil
	})
}

func (m Manager) migrate(db *gorm.DB, p provider.Provider) error {
	return db.AutoMigrate(p.Models()...)
}
