package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Company
// Rol ->
const (
	RolOwner string = "OWNER"
	RolStaff string = "STAFF"

	StatusPending   string = "PENDING_CONFIRMATION"
	StatusConfirmed string = "CONFIRMED"
	StatusActive    string = "ACTIVE"
	StatusInactive  string = "INACTIVE"

	DocTypeStaff   string = "Staff"
	DocTypeService string = "Service"
	DocTypeCompany string = "Company"

	SexM string = "M"
	SexW string = "W"
	SexA string = "A" //Any
)

type Company struct {
	ID          string            `json:"CompanyID"`
	UserSub     string            `json:"UserSub"`
	Name        string            `json:"Name,omitempty"`
	Description string            `json:"Description,omitempty"`
	Timestamp   string            `json:"Timestamp,omitempty"`
	Staff       []*Staff          `json:"Staff,omitempty"`
	Services    []*CompanyService `json:"Services,omitempty"`
	Rol         string            `json:"Rol,omitempty"`
	Status      string            `json:"Status,omitempty"` // ACTIVE, INACTIVE
}

type Staff struct {
	ID        string `json:"StaffID"`
	UserSub   string `json:"UserSub"`
	CompanyID string `json:"CompanyID,omitempty"`
	Name      string `json:"Name,omitempty"`
	Rol       string `json:"Rol,omitempty"`
	Status    string `json:"Status,omitempty"` // PENDING_CONFIRMATION, CONFIRMED
	CreatedBy string `json:"CreatedBy,omitempty"`
}

type CompanyService struct {
	ID          string   `json:"ServiceID"`
	UserSub     string   `json:"UserSub"`
	CompanyID   string   `json:"CompanyID,omitempty"`
	Name        string   `json:"Name,omitempty"`
	Description string   `json:"Description,omitempty"`
	ServiceType string   `json:"ServiceType,omitempty"`
	Price       int64    `json:"Price,omitempty"`
	CurrencyID  string   `json:"CurrencyID,omitempty"`
	Status      string   `json:"Status,omitempty"` // ACTIVE, INACTIVE
	Sex         string   `json:"Sex"`              // Men, Women, Any
	MinAge      int      `json:"MinAge"`
	MaxAge      int      `json:"MaxAge"`
	Tags        []string `json:"Tags,omitempty"` // For instance: HARD, SOFT, etc.
}

func NewCompanyService(companyID string, userSub string) CompanyService {
	c := CompanyService{}
	uid, _ := uuid.NewV4()
	c.ID = uid.String()
	c.MinAge = -1
	c.MaxAge = 9999
	c.CompanyID = companyID
	c.Status = StatusActive
	c.Sex = SexA
	c.UserSub = userSub
	return c
}

func NewCompany(userSub string) Company {
	c := Company{}
	uid, _ := uuid.NewV4()
	c.ID = uid.String()
	c.UserSub = userSub
	c.Rol = RolOwner
	c.Status = StatusActive
	c.Timestamp = time.Now().Format("2006-01-02T15:04:05")
	return c
}
