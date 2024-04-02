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

func (project *ProjectAdapter) GetPackgeTypeByName(name string) (entities.PackageType, error) {
	var res entities.PackageType
	query := "SELECT * FROM package_types WHERE type = ?"
	if err := project.DB.Raw(query, name).Scan(&res).Error; err != nil {
		return entities.PackageType{}, err
	}
	return res, nil
}

func (project *ProjectAdapter) GetPackgeTypeById(id int32) (entities.PackageType, error) {
	var res entities.PackageType
	query := "SELECT * FROM package_types WHERE id = ?"
	if err := project.DB.Raw(query, id).Scan(&res).Error; err != nil {
		return entities.PackageType{}, err
	}
	return res, nil
}

func (project *ProjectAdapter) AddPackageType(name string) error {
	query := "INSERT INTO package_types (type) VALUES ($1)"
	if err := project.DB.Exec(query, name).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) EditPackgeType(req entities.PackageType) error {
	query := "UPDATE package_types SET type = $1 WHERE id = $2"
	if err := project.DB.Exec(query, req.Type, req.Id).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) GetAllPAckageType() ([]entities.PackageType, error) {
	var res []entities.PackageType
	query := "SELECT * FROM pakcage_types"
	if err := project.DB.Raw(query).Scan(&res).Error; err != nil {
		return []entities.PackageType{}, err
	}
	return res, nil
}
