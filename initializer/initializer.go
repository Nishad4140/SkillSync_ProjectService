package initializer

import (
	"github.com/Nishad4140/SkillSync_ProjectService/internal/adapter"
	"github.com/Nishad4140/SkillSync_ProjectService/internal/service"
	"github.com/Nishad4140/SkillSync_ProjectService/internal/usecase"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB) *service.ProjectService {
	adapter := adapter.NewProjectAdapter(db)
	usecase := usecase.NewProjectUsecase(adapter)
	service := service.NewProjectService(adapter, usecase)
	return service
}
