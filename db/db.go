package db

import (
	"github.com/Nishad4140/SkillSync_ProjectService/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(connectTo string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectTo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(
		&entities.Gig{},
		&entities.PackageType{},
		&entities.ClientRequest{},
		&entities.GigImages{},
		&entities.ClientRequestImages{},
		&entities.Intrest{},
		&entities.IntrestAcknowledgment{},
		&entities.Project{},
		&entities.ProjectFiles{},
		&entities.ProjectManagement{},
		&entities.Module{},
	)

	return db, nil
}
