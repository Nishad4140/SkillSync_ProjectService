package adapter

import (
	"github.com/Nishad4140/SkillSync_ProjectService/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectAdapter struct {
	DB *gorm.DB
}

func NewProjectAdapter(db *gorm.DB) *ProjectAdapter {
	return &ProjectAdapter{
		DB: db,
	}
}

func (project *ProjectAdapter) CreateGig(req entities.Gig) error {
	id := uuid.New()
	query := "INSERT INTO gigs (id, freelancer_id, title, category_id, skill_id, description, package_type_id, price, delivery_days) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	if err := project.DB.Exec(query, id, req.FreelancerId, req.Title, req.CategoryId, req.SkillId, req.Description, req.PackageTypeId, req.Price, req.DeliveryDays).Error; err != nil {
		return err
	}
	return nil
}
