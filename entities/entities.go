package entities

import (
	"time"

	"github.com/google/uuid"
)

type Gig struct {
	ID            uuid.UUID `gorm:"primaryKey;unique;not null"`
	FreelancerId  uuid.UUID
	Title         string
	CategoryId    int
	SkillId       int
	Description   string
	PackageTypeId int
	PackageType   PackageType `gorm:"foreignKey:PackageTypeId"`
	Price         float64
	DeliveryDays  string
	IsPrivate     bool `gorm:"default:false"`
}

type ClientRequest struct {
	ID           uuid.UUID `gorm:"primaryKey;unique;not null"`
	ClientId     uuid.UUID
	Title        string
	CategoryId   int
	SkillId      int
	Description  string
	Price        float64
	DelivaryDate time.Time `gorm:"type:date"`
}

type GigImages struct {
	Id    uuid.UUID `gorm:"primaryKey;unique;not null"`
	GigId uuid.UUID
	Gig   Gig `gorm:"foreignKey:GigId"`
	Image string
}

type ClientRequestImages struct {
	Id              uuid.UUID `gorm:"primaryKey;unique;not null"`
	ClientRequestId uuid.UUID
	ClientRequest   ClientRequest `gorm:"foreignKey:ClientRequestId"`
	Image           string
}

type PackageType struct {
	Id   int
	Type string
}

type Intrest struct {
	Id              uuid.UUID `gorm:"primaryKey;unique;not null"`
	ClientRequestId uuid.UUID
	ClientRequest   ClientRequest `gorm:"foreignKey:ClientRequestId"`
	FreelancerId    uuid.UUID
	GigId           uuid.UUID
	Gig             Gig `gorm:"foreignKey:GigId"`
}

type IntrestAcknowledgment struct {
	Id        uuid.UUID `gorm:"primaryKey;unique;not null"`
	ClientId  uuid.UUID
	IntrestId uuid.UUID
	Intrest   Intrest `gorm:"foreignKey:IntrestId"`
}

type Project struct {
	Id           uuid.UUID `gorm:"primaryKey;unique;not null"`
	ClientId     uuid.UUID
	FreelancerId uuid.UUID
	GigId        uuid.UUID
	StartDate    time.Time `gorm:"type:date"`
	EndDate      time.Time `gorm:"type:date"`
	Status       string
	Price        float64
	IsManagement bool
	IsPaid       bool `gorm:"default:false"`
}

type ProjectManagement struct {
	Id            uuid.UUID `gorm:"primaryKey;unique;not null"`
	ProjectId     uuid.UUID
	Project       Project `gorm:"foreignKey:ProjectId"`
	ModulesNumber int
}
type Module struct {
	Id                uuid.UUID `gorm:"primaryKey;unique;not null"`
	ManagementId      uuid.UUID
	ProjectManagement ProjectManagement `gorm:"foreignKey:ManagementId"`
	ModuleName        string
	ModuleDescription string
	Status            int `gorm:"default:0"`
}
type ProjectFiles struct {
	Id        uuid.UUID `gorm:"primaryKey;unique;not null"`
	ProjectId uuid.UUID
	Project   Project `gorm:"foreignKey:ProjectId"`
	File      string
	IsPaid    bool `gorm:"default:false"`
}
