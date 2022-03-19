package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/raulandre/anonboard/models"
	"gorm.io/gorm"
)

func sync(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "initial",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&models.Thread{},
					&models.Reply{},
				); err != nil {
					return err
				}

				if err := tx.Exec("ALTER TABLE threads ADD CONSTRAINT fk_threads FOREIGN KEY (thread_id) REFERENCES threads (id)").Error; err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					&models.Thread{},
					&models.Reply{},
				)
			},
		},
	})

	return m.Migrate()
}
