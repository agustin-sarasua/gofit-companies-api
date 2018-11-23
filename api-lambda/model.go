package main

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
	UserSub     string   `json:"UserSub"`
	ID          string   `json:"CompanyID"`
	Name        string   `json:"Name,omitempty"`
	Description string   `json:"Description,omitempty"`
	LogoUrl     string   `json:"LogoUrl,omitempty"`
	Timestamp   string   `json:"Timestamp,omitempty"`
	Staff       []*Staff `json:"Staff,omitempty"`
	Rol         string   `json:"Rol,omitempty"`
	Status      string   `json:"Status,omitempty"` // ACTIVE, INACTIVE
	//DocType     string   `json:"DocType,omitempty"`
}

type Staff struct {
	UserSub   string `json:"UserSub"`
	CompanyID string `json:"CompanyID,omitempty"`
	Name      string `json:"Name,omitempty"`
	Rol       string `json:"Rol,omitempty"`
	Status    string `json:"Status,omitempty"` // PENDING_CONFIRMATION, CONFIRMED
	CreatedBy string `json:"CreatedBy,omitempty"`
	//DocType   string `json:"DocType,omitempty"`
}
