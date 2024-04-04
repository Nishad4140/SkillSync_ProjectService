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

func (project *ProjectAdapter) UpdateGig(req entities.Gig) error {
	query := "UPDATE gigs SET title = $1, category_id = $2, skill_id = $3, description = $4, package_type_id = $5, price = $6, delivery_days = $7 WHERE id = $8"
	if err := project.DB.Exec(query, req.Title, req.CategoryId, req.SkillId, req.Description, req.PackageTypeId, req.Price, req.DeliveryDays, req.ID).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) GetGigById(id string) (entities.Gig, error) {
	var res entities.Gig
	query := "SELECT * FROM gigs WHERE id = ?"
	if err := project.DB.Raw(query, id).Scan(&res).Error; err != nil {
		return entities.Gig{}, err
	}
	return res, nil
}

func (project *ProjectAdapter) GetAllFreelancerGigs(freelancerId string) ([]entities.Gig, error) {
	var res []entities.Gig
	query := "SELECT * FROM gigs WHERE freelancer_id = ?"
	if err := project.DB.Raw(query, freelancerId).Scan(&res).Error; err != nil {
		return []entities.Gig{}, err
	}
	return res, nil
}

func (project *ProjectAdapter) GetAllGigs() ([]entities.Gig, error) {
	var res []entities.Gig
	query := "SELECT * FROM gigs"
	if err := project.DB.Raw(query).Scan(&res).Error; err != nil {
		return []entities.Gig{}, err
	}
	return res, nil
}

func (project *ProjectAdapter) ClientAddRequest(req entities.ClientRequest) error {
	id := uuid.New()
	query := "INSERT INTO client_requests (id, client_id, title, category_id, skill_id, description, price, delivary_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	if err := project.DB.Exec(query, id, req.ClientId, req.Title, req.CategoryId, req.SkillId, req.Description, req.Price, req.DelivaryDate).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) ClientUpdateRequest(req entities.ClientRequest) error {
	query := "UPDATE client_requests SET title = $1, category_id = $2, skill_id = $3, description  =$4, price = $5, delivary_date = $6 WHERE id = $7"
	if err := project.DB.Exec(query, req.Title, req.CategoryId, req.SkillId, req.Description, req.Price, req.DelivaryDate, req.ID).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) GetClientRequest(reqId string) (entities.ClientRequest, error) {
	var res entities.ClientRequest
	query := "SELECT * FROM client_requests WHERE id = ?"
	if err := project.DB.Raw(query, reqId).Scan(&res).Error; err != nil {
		return entities.ClientRequest{}, err
	}
	return res, nil
}

func (project *ProjectAdapter) GetAllClientRequest(clientId string) ([]entities.ClientRequest, error) {
	var res []entities.ClientRequest
	query := "SELECT * FROM client_requests WHERE client_id = ?"
	if err := project.DB.Raw(query, clientId).Scan(&res).Error; err != nil {
		return []entities.ClientRequest{}, err
	}
	return res, nil
}

func (project *ProjectAdapter) GetAllClientRequestForFreelancers(categoryId int) ([]entities.ClientRequest, error) {
	var res []entities.ClientRequest

	query := "SELECT * FROM client_requests WHERE category_id = ?"
	if err := project.DB.Raw(query, categoryId).Scan(&res).Error; err != nil {
		return []entities.ClientRequest{}, err
	}
	return res, nil
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
	query := "SELECT * FROM package_types"
	if err := project.DB.Raw(query).Scan(&res).Error; err != nil {
		return []entities.PackageType{}, err
	}
	return res, nil
}
