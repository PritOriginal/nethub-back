package deviceshandler

import "github.com/PritOriginal/nethub-back/internal/models"

type AddDeviceRequest struct {
	Hostname string `json:"hostname" binding:"required,hostname"`
	IP       string `json:"ip" binding:"required,ip"`
	Location string `json:"location" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type AddDeviceResponse struct {
	Device models.Device `json:"device"`
}

type GetDevicesResponse struct {
	Devices []models.Device `json:"devices"`
}

type GetDeviceByIdResponse struct {
	Device models.Device `json:"device"`
}

type UpdateDeviceRequest struct {
	Hostname string `json:"hostname" binding:"required,hostname"`
	IP       string `json:"ip" binding:"required,ip"`
	Location string `json:"location" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type UpdateDeviceResponse struct {
	Device models.Device `json:"device"`
}

type DeleteDeviceResponse struct {
}
