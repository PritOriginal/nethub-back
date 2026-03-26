package devices

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/PritOriginal/nethub-back/internal/models"
	"github.com/PritOriginal/nethub-back/internal/service"
	"github.com/PritOriginal/nethub-back/internal/storage"
)

type Storage interface {
	AddDevice(ctx context.Context, device models.Device) (models.Device, error)
	GetDevices(ctx context.Context, params models.GetDevicesParams) ([]models.Device, error)
	GetDeviceById(ctx context.Context, id int) (models.Device, error)
	UpdateDevice(ctx context.Context, device models.Device) (models.Device, error)
	DeleteDevice(ctx context.Context, id int) error
}

type Service struct {
	log   *slog.Logger
	store Storage
}

func New(log *slog.Logger, store Storage) *Service {
	return &Service{
		log:   log,
		store: store,
	}
}

func (s *Service) AddDevice(ctx context.Context, device models.Device) (models.Device, error) {
	const op = "service.devices.AddDevice"

	newDevice, err := s.store.AddDevice(ctx, device)
	if err != nil {
		return device, fmt.Errorf("%s: %w", op, err)
	}

	return newDevice, nil
}

func (s *Service) GetDeviceById(ctx context.Context, id int) (models.Device, error) {
	const op = "service.devices.GetDeviceById"

	device, err := s.store.GetDeviceById(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return device, service.ErrNotFound
		} else {
			return device, fmt.Errorf("%s: %w", op, err)
		}
	}

	return device, nil
}

func (s *Service) GetDevices(ctx context.Context, params models.GetDevicesParams) ([]models.Device, error) {
	const op = "service.devices.GetDevices"

	devices, err := s.store.GetDevices(ctx, params)
	if err != nil {
		return devices, fmt.Errorf("%s: %w", op, err)
	}

	return devices, nil
}

func (s *Service) UpdateDevice(ctx context.Context, device models.Device) (models.Device, error) {
	const op = "service.devices.UpdateDevice"

	updatedDevice, err := s.store.UpdateDevice(ctx, device)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return updatedDevice, service.ErrNotFound
		} else {
			return updatedDevice, fmt.Errorf("%s: %w", op, err)
		}
	}

	return updatedDevice, nil
}

func (s *Service) DeleteDevice(ctx context.Context, id int) error {
	const op = "service.devices.DeleteDevice"

	_, err := s.store.GetDeviceById(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return service.ErrNotFound
		} else {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if err := s.store.DeleteDevice(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
