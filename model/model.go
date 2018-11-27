package model

// Company
// Rol ->
const (
	RolOwner string = "OWNER"
	RolStaff string = "STAFF"

	StatusPending   string = "PENDING_CONFIRMATION"
	StatusConfirmed string = "CONFIRMED"
	StatusActive    string = "ACTIVE"
	StatusInactive  string = "INACTIVE"
)

type Company struct {
	UserSub     string            `json:"UserSub"`
	ID          string            `json:"CompanyID"`
	SortKey     string            `json:"SortKey,omitempty"`
	Name        string            `json:"Name,omitempty"`
	Description string            `json:"Description,omitempty"`
	LogoUrl     string            `json:"LogoUrl,omitempty"`
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
	SortKey   string `json:"SortKey,omitempty"`
	Name      string `json:"Name,omitempty"`
	Rol       string `json:"Rol,omitempty"`
	Status    string `json:"Status,omitempty"` // PENDING_CONFIRMATION, CONFIRMED
	CreatedBy string `json:"CreatedBy,omitempty"`
}

type CompanyService struct {
	ID          string `json:"ServiceID"`
	UserSub     string `json:"UserSub"`
	CompanyID   string `json:"CompanyID,omitempty"`
	SortKey     string `json:"SortKey,omitempty"`
	Name        string `json:"Name,omitempty"`
	Description string `json:"Description,omitempty"`
	ServiceType string `json:"ServiceType,omitempty"`
	Price       int64  `json:"Price,omitempty"`
	CurrencyID  string `json:"CurrencyID,omitempty"`
	Status      string `json:"Status,omitempty"` // ACTIVE, INACTIVE
}
