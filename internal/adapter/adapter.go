package adapter

import (
	"fmt"
	"strings"
	"time"

	"github.com/Nishad4140/SkillSync_ProjectService/entities"
	"github.com/Nishad4140/SkillSync_ProjectService/internal/helper"
	helperstruct "github.com/Nishad4140/SkillSync_ProjectService/internal/helperStruct"
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

func (project *ProjectAdapter) GetAllGigs(queryParams helperstruct.FilterQuery) ([]entities.Gig, error) {
	var res []entities.Gig
	query := "SELECT * FROM gigs"

	if queryParams.Query != "" && queryParams.Filter != "" {
		query = fmt.Sprintf("%s WHERE LOWER(%s) LIKE '%%%s%%' AND is_private = false", query, queryParams.Filter, strings.ToLower(queryParams.Query))
	} else {
		query = fmt.Sprintf("%s WHERE is_private = false", query)
	}

	if queryParams.SortBy != "" {
		if queryParams.SortDesc {
			query = fmt.Sprintf("%s ORDER BY %s DESC", query, queryParams.SortBy)
		} else {
			query = fmt.Sprintf("%s ORDER BY %s ASC", query, queryParams.SortBy)
		}
	} else {
		query = fmt.Sprintf("%s ORDER BY price ASC", query)
	}
	//to set the page number and the qty that need to display in a single responce
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		query = fmt.Sprintf("%s LIMIT 10 OFFSET 0", query)
	}

	if err := project.DB.Raw(query).Scan(&res).Error; err != nil {
		return []entities.Gig{}, err
	}
	return res, nil
}

func (project *ProjectAdapter) ClientAddRequest(req entities.ClientRequest) error {
	id := uuid.New()
	query := "INSERT INTO client_requests (id, client_id, title, category_id, skill_id, description, price, delivary_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
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

func (project *ProjectAdapter) GetAllClientRequest(clientId string, queryParams helperstruct.FilterQuery) ([]entities.ClientRequest, error) {
	var res []entities.ClientRequest
	query := "SELECT * FROM client_requests WHERE client_id = $1"

	if queryParams.Query != "" && queryParams.Filter != "" {
		query = fmt.Sprintf("%s AND LOWER(%s) LIKE '%%%s%%'", query, queryParams.Filter, strings.ToLower(queryParams.Query))
	}

	if queryParams.SortBy != "" {
		if queryParams.SortDesc {
			query = fmt.Sprintf("%s ORDER BY %s DESC", query, queryParams.SortBy)
		} else {
			query = fmt.Sprintf("%s ORDER BY %s ASC", query, queryParams.SortBy)
		}
	} else {
		query = fmt.Sprintf("%s ORDER BY price ASC", query)
	}
	//to set the page number and the qty that need to display in a single responce
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		query = fmt.Sprintf("%s LIMIT 10 OFFSET 0", query)
	}

	if err := project.DB.Raw(query, clientId).Scan(&res).Error; err != nil {
		return []entities.ClientRequest{}, err
	}
	return res, nil
}

func (project *ProjectAdapter) GetAllClientRequestForFreelancers(categoryId int, queryParams helperstruct.FilterQuery) ([]entities.ClientRequest, error) {
	var res []entities.ClientRequest

	query := "SELECT * FROM client_requests WHERE category_id = $1"

	if queryParams.Query != "" && queryParams.Filter != "" {
		query = fmt.Sprintf("%s AND LOWER(%s) LIKE '%%%s%%'", query, queryParams.Filter, strings.ToLower(queryParams.Query))
	}

	if queryParams.SortBy != "" {
		if queryParams.SortDesc {
			query = fmt.Sprintf("%s ORDER BY %s DESC", query, queryParams.SortBy)
		} else {
			query = fmt.Sprintf("%s ORDER BY %s ASC", query, queryParams.SortBy)
		}
	} else {
		query = fmt.Sprintf("%s ORDER BY price ASC", query)
	}
	//to set the page number and the qty that need to display in a single responce
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		query = fmt.Sprintf("%s LIMIT 10 OFFSET 0", query)
	}

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

func (project *ProjectAdapter) AddIntrestToClientRequest(req entities.Intrest) error {
	id := uuid.New()

	query := "INSERT INTO intrests (id, client_request_id, freelancer_id, gig_id) VALUES ($1, $2, $3, $4)"
	if err := project.DB.Exec(query, id, req.ClientRequestId, req.FreelancerId, req.GigId).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) GetAllClientRequestIntrest(reqId string) ([]entities.Intrest, error) {
	var res []entities.Intrest

	query := "SELECT * FROM intrests WHERE client_request_id = ?"
	if err := project.DB.Raw(query, reqId).Scan(&res).Error; err != nil {
		return []entities.Intrest{}, err
	}
	return res, nil
}

func (project *ProjectAdapter) GetClientIdByRequestId(reqId string) (string, error) {
	var res string

	query := "SELECT client_id FROM client_requests WHERE id = ?"
	if err := project.DB.Raw(query, reqId).Scan(&res).Error; err != nil {
		return "", err
	}
	return res, nil
}

func (project *ProjectAdapter) ClientAddIntrestAcknowledgment(req entities.IntrestAcknowledgment) error {
	id := uuid.New()
	query := "INSERT INTO intrest_acknowledgments (id, client_id, intrest_id) VALUES ($1, $2, $3)"
	if err := project.DB.Exec(query, id, req.ClientId, req.IntrestId).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) GetIntrestById(id string) (entities.Intrest, error) {
	var res entities.Intrest
	query := "SELECT * FROM intrests WHERE id = ?"
	if err := project.DB.Raw(query, id).Scan(&res).Error; err != nil {
		return entities.Intrest{}, err
	}
	return res, nil
}

func (project *ProjectAdapter) CreateProject(req entities.Project) (entities.Project, error) {
	var proj entities.Project
	query := "INSERT INTO projects (id, client_id, freelancer_id, gig_id, start_date, end_date, status, price) VALUES ($1, $2, $3, $4, $5, $6, 'not started', $7) RETURNING *"
	if err := project.DB.Raw(query, req.Id, req.ClientId, req.FreelancerId, req.GigId, time.Now(), req.EndDate, req.Price).Scan(&proj).Error; err != nil {
		return entities.Project{}, err
	}
	id := uuid.New()
	fileQuery := "INSERT INTO project_files (id, project_id) VALUES ($1, $2)"
	if err := project.DB.Raw(fileQuery, id, req.Id).Error; err != nil {
		return entities.Project{}, err
	}
	return proj, nil
}

func (project *ProjectAdapter) UpdateProject(req entities.Project) error {
	query := "UPDATE projects SET gig_id = $1, end_date = $2 WHERE id = $3"
	if err := project.DB.Exec(query, req.GigId, req.EndDate, req.Id).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) GetProject(projectId string) (entities.Project, error) {
	var res entities.Project
	query := "SELECT * FROM projects WHERE id = ?"
	if err := project.DB.Raw(query, projectId).Scan(&res).Error; err != nil {
		return entities.Project{}, err
	}
	return res, nil
}

func (project *ProjectAdapter) GetAllProjects(req entities.Project) ([]entities.Project, error) {
	var res []entities.Project
	if req.ClientId != uuid.Nil {
		query := "SELECT * FROM projects WHERE client_id = $1"
		if err := project.DB.Raw(query, req.ClientId).Scan(&res).Error; err != nil {
			return []entities.Project{}, err
		}
	} else if req.FreelancerId != uuid.Nil {
		query := "SELECT * FROM gigs WHERE freelancer_id = $1 AND status != 'not started'"
		if err := project.DB.Raw(query, req.FreelancerId).Scan(&res).Error; err != nil {
			return []entities.Project{}, err
		}
	}

	return res, nil
}

func (project *ProjectAdapter) RemoveProjects(projectId string) error {
	query := "DELETE FROM projects WHERE id = $1 AND status = 'not started'"
	if err := project.DB.Exec(query, projectId).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) ProjectManagement(req helperstruct.ProjectManagement) (string, error) {
	if !req.IsManagementNeeded {
		query := "UPDATE projects SET is_management = $1 WHERE id = $2"
		if err := project.DB.Exec(query, false, req.ProjectId).Error; err != nil {
			return "", err
		}
		return "", nil
	}
	query := "UPDATE projects SET is_management = $1 WHERE id = $2"
	if err := project.DB.Exec(query, true, req.ProjectId).Error; err != nil {
		return "", err
	}
	managementId := uuid.New()
	managementQuery := "INSERT INTO project_managements (id, project_id, modules_number) VALUES ($1, $2, $3)"
	if err := project.DB.Exec(managementQuery, managementId, req.ProjectId, req.ModuleNumber).Error; err != nil {
		return "", err
	}
	return managementId.String(), nil
}

func (project *ProjectAdapter) ModuleManagement(req helperstruct.ModuleManagement) error {
	moduleId := uuid.New()
	moduleQuery := "INSERT INTO modules (id, management_id, module_name, module_description) VALUES ($1, $2, $3, $4)"
	if err := project.DB.Exec(moduleQuery, moduleId, req.ManagementId, req.ModuleDetails[0], req.ModuleDetails[1]).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) GetManagementByProjectId(projectId string) (entities.ProjectManagement, error) {
	var managementData entities.ProjectManagement
	query := "SELECT * FROM project_managements WHERE project_id = $1"
	if err := project.DB.Raw(query, projectId).Scan(&managementData).Error; err != nil {
		return entities.ProjectManagement{}, err
	}
	return managementData, nil
}

func (project *ProjectAdapter) GetModuleByManagementId(projectId string) ([]entities.Module, error) {
	var moduleData []entities.Module
	query := "SELECT * FROM modules WHERE management_id = $1"
	if err := project.DB.Raw(query, projectId).Scan(&moduleData).Error; err != nil {
		return []entities.Module{}, err
	}
	return moduleData, nil
}

func (project *ProjectAdapter) ModuleStatusUpdate(req entities.Module) error {
	query := "UPDATE modules SET status = $1 WHERE id = $2"
	if err := project.DB.Exec(query, req.Status, req.Id).Error; err != nil {
		return err
	}
	var management entities.ProjectManagement
	managementQuery := "SELECT * FROM project_managements WHERE id = $1"
	if err := project.DB.Raw(managementQuery, req.ManagementId).Scan(&management).Error; err != nil {
		return err
	}

	var values []int
	moduleQuery := "SELECT status FROM modules WHERE management_id = $1"
	if err := project.DB.Raw(moduleQuery, req.ManagementId).Scan(&values).Error; err != nil {
		return err
	}
	percentage := helper.PercentageCaluculation(management.ModulesNumber, values)

	projectQuery := "UPDATE projects SET status = $1 WHERE id = $2"
	if err := project.DB.Exec(projectQuery, percentage, management.ProjectId).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) FreelancerUpdateStatus(req entities.Project) error {

	if req.Status == "100" {
		req.Status = "completed"
	}

	query := "UPDATE projects SET status = $1 WHERE id = $2"
	if err := project.DB.Exec(query, req.Status, req.Id).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) UpdatePaymentStatus(projectId string) error {
	query := "UPDATE projects SET is_paid = true WHERE id = $1"
	if err := project.DB.Exec(query, projectId).Error; err != nil {
		return err
	}
	fileQuery := "UPDATE project_files SET is_paid = true WHERE project_id = $1"
	if err := project.DB.Exec(fileQuery, projectId).Error; err != nil {
		return err
	}
	return nil
}

func (project *ProjectAdapter) UploadFile(file, projectId string) (string, error) {
	query := "UPDATE project_files SET file = $1 WHERE project_id = $2"
	if err := project.DB.Exec(query, file, projectId).Error; err != nil {
		return "", err
	}
	return file, nil
}

func (project *ProjectAdapter) GetFile(projectId string) (entities.ProjectFiles, error) {
	var res entities.ProjectFiles
	query := "SELECT * FROM project_files WHERE project_id = $1"
	if err := project.DB.Raw(query, projectId).Scan(&res).Error; err != nil {
		return entities.ProjectFiles{}, err
	}
	return res, nil
}
