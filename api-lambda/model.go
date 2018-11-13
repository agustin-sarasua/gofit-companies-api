package main

type Company struct {
	UserSub   string `json:"UserSub"`
	ID        string `json:"CompanyID"`
	Name      string `json:"Name,omitempty"`
	Timestamp string `json:"Timestamp"`
}
