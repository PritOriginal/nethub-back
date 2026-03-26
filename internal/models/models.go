package models

import "time"

type Device struct {
	Id        int       `json:"id" db:"id"`
	Hostname  string    `json:"hostname" db:"hostname"`
	IP        string    `json:"ip" db:"ip"`
	Location  string    `json:"location" db:"location"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type GetDevicesParams struct {
	IsActive *bool
	Search   string
}
