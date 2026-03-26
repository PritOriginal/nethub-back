package deviceshandler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/PritOriginal/nethub-back/internal/models"
	"github.com/PritOriginal/nethub-back/internal/service"
	"github.com/PritOriginal/nethub-back/pkg/logger"
	"github.com/PritOriginal/nethub-back/pkg/responses"
	"github.com/gin-gonic/gin"
)

type DevicesService interface {
	AddDevice(ctx context.Context, device models.Device) (models.Device, error)
	GetDevices(ctx context.Context, params models.GetDevicesParams) ([]models.Device, error)
	GetDeviceById(ctx context.Context, id int) (models.Device, error)
	UpdateDevice(ctx context.Context, device models.Device) (models.Device, error)
	DeleteDevice(ctx context.Context, id int) error
}

type handler struct {
	log *slog.Logger
	s   DevicesService
}

func Register(r *gin.Engine, log *slog.Logger, s DevicesService) {
	handler := handler{log: log, s: s}

	devices := r.Group("/devices")
	{
		devices.POST("", handler.AddDevice())
		devices.GET("", handler.GetDevices())
		devices.GET(":id", handler.GetDeviceById())
		devices.PUT(":id", handler.UpdateDevice())
		devices.DELETE(":id", handler.DeleteDevice())
	}
}

// AddDevice Add new device
//
//	@Summary		Add new device
//	@Description	Add new device
//	@Tags			devices
//	@Produce		json
//	@Param			request	body		deviceshandler.AddDeviceRequest	true	"device params"
//	@Success		201		{object}	responses.Response[deviceshandler.AddDeviceResponse]
//	@Failure		400		{object}	responses.Response[any]
//	@Failure		500		{object}	responses.Response[any]
//	@Router			/devices [post]
func (h *handler) AddDevice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req AddDeviceRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			h.log.Debug("failed binding request", logger.Err(err))
			responses.Fail(ctx, http.StatusBadRequest, "failed binding request")
			return
		}

		newDevice, err := h.s.AddDevice(ctx.Request.Context(), models.Device{
			Hostname: req.Hostname,
			IP:       req.IP,
			Location: req.Location,
			IsActive: req.IsActive,
		})
		if err != nil {
			h.log.Error("error add device", logger.Err(err))
			responses.Internal(ctx, "error add device")
			return
		}

		h.log.Info("add new device", slog.Any("device", newDevice))
		responses.Created(ctx, AddDeviceResponse{
			Device: newDevice,
		})
	}
}

// GetTasks returns a list of devices
//
//	@Summary		Returns a list of devices
//	@Description	Returns a list of devices using filtering
//	@Tags			devices
//	@Produce		json
//	@Param			is_active	query		boolean	false	"filter by is_active"
//	@Param			search		query		string	false	"filter by hostname"
//	@Success		200			{object}	responses.Response[deviceshandler.GetDevicesResponse]
//	@Failure		400			{object}	responses.Response[any]
//	@Failure		500			{object}	responses.Response[any]
//	@Router			/devices [get]
func (h *handler) GetDevices() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var isActive *bool
		isActiveStr := ctx.Query("is_active")
		if isActiveStr != "" {
			val, err := strconv.ParseBool(isActiveStr)
			if err != nil {
				h.log.Debug("", logger.Err(err))
				responses.Fail(ctx, http.StatusBadRequest, "failed parse is_active")
				return
			}
			isActive = &val
		}
		search := ctx.Query("search")

		devices, err := h.s.GetDevices(ctx.Request.Context(), models.GetDevicesParams{
			IsActive: isActive,
			Search:   search,
		})
		if err != nil {
			h.log.Error("error get devices", logger.Err(err))
			responses.Internal(ctx, "error get devices")
			return
		}

		responses.OK(ctx, GetDevicesResponse{
			Devices: devices,
		})
	}
}

// GetDeviceById returns the device
//
//	@Summary		Returns the device by ID
//	@Description	Returns the device by ID
//	@Tags			devices
//	@Produce		json
//	@Param			id	path		int	true	"device id"
//	@Success		200	{object}	responses.Response[deviceshandler.GetDeviceByIdResponse]
//	@Failure		400	{object}	responses.Response[any]
//	@Failure		404	{object}	responses.Response[any]
//	@Failure		500	{object}	responses.Response[any]
//	@Router			/devices/{id} [get]
func (h *handler) GetDeviceById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			h.log.Debug("failed parse id", logger.Err(err))
			responses.Fail(ctx, http.StatusBadRequest, "failed parse id")
			return
		}

		device, err := h.s.GetDeviceById(ctx.Request.Context(), id)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				h.log.Debug("device not found", slog.Int("id", id))
				responses.Fail(ctx, http.StatusNotFound, "device not found")
			} else {
				h.log.Error("error get device by id", logger.Err(err))
				responses.Internal(ctx, "error get device by id")
			}
			return
		}

		responses.OK(ctx, GetDeviceByIdResponse{
			Device: device,
		})
	}
}

// UpdateDevice update the device
//
//	@Summary		Update the device
//	@Description	Update the device by ID
//	@Tags			devices
//	@Produce		json
//	@Param			id		path		int									true	"device id"
//	@Param			request	body		deviceshandler.UpdateDeviceRequest	true	"device params"
//	@Success		200		{object}	responses.Response[deviceshandler.UpdateDeviceResponse]
//	@Failure		400		{object}	responses.Response[any]
//	@Failure		500		{object}	responses.Response[any]
//	@Router			/devices/{id} [put]
func (h *handler) UpdateDevice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			h.log.Debug("failed parse id", logger.Err(err))
			responses.Fail(ctx, http.StatusBadRequest, "failed parse id")
			return
		}

		var req UpdateDeviceRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			h.log.Debug("failed binding request", logger.Err(err))
			responses.Fail(ctx, http.StatusBadRequest, "failed binding request")
			return
		}

		updatedDevice, err := h.s.UpdateDevice(ctx.Request.Context(), models.Device{
			Id:       id,
			Hostname: req.Hostname,
			IP:       req.IP,
			Location: req.Location,
			IsActive: req.IsActive,
		})
		if err != nil {
			h.log.Error("error update device", logger.Err(err))
			responses.Internal(ctx, "error update device")
			return
		}

		h.log.Info("update device", slog.Any("device", updatedDevice))
		responses.OK(ctx, UpdateDeviceResponse{
			Device: updatedDevice,
		})
	}
}

// DeleteDevice Delete the device
//
//	@Summary		Delete the device
//	@Description	Delete the device by ID
//	@Tags			devices
//	@Produce		json
//	@Param			id	path		int	true	"device id"
//	@Success		200	{object}	responses.Response[deviceshandler.DeleteDeviceResponse]
//	@Failure		400	{object}	responses.Response[any]
//	@Failure		404	{object}	responses.Response[any]
//	@Failure		500	{object}	responses.Response[any]
//	@Router			/devices/{id} [delete]
func (h *handler) DeleteDevice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			h.log.Debug("failed parse id", logger.Err(err))
			responses.Fail(ctx, http.StatusBadRequest, "failed parse id")
			return
		}

		err = h.s.DeleteDevice(ctx.Request.Context(), id)
		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				h.log.Debug("device not found", slog.Int("id", id))
				responses.Fail(ctx, http.StatusNotFound, "device not found")
			} else {
				h.log.Error("error delete device", logger.Err(err))
				responses.Internal(ctx, "error delete device")
			}
			return
		}

		h.log.Info("delete device", slog.Int("device_id", id))
		responses.OK(ctx, DeleteDeviceResponse{})
	}
}
