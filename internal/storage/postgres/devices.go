package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/PritOriginal/nethub-back/internal/models"
	"github.com/PritOriginal/nethub-back/internal/storage"
	"github.com/jmoiron/sqlx"
)

type DevicesStorage struct {
	DB *sqlx.DB
}

func NewDeviceStorage(db *sqlx.DB) *DevicesStorage {
	return &DevicesStorage{
		DB: db,
	}
}

func (s *DevicesStorage) AddDevice(ctx context.Context, device models.Device) (models.Device, error) {
	const op = "storage.postgres.AddDevice"

	var newDevice models.Device

	query := `
			INSERT INTO
				devices (hostname, ip, location, is_active)
			VALUES
				(:hostname, :ip, :location, :is_active)
			RETURNING
				id, hostname, ip, location, is_active, created_at 
			`

	stmt, err := s.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return newDevice, fmt.Errorf("%s: %w", op, err)
	}

	if err := stmt.GetContext(ctx, &newDevice, device); err != nil {
		return newDevice, fmt.Errorf("%s: %w", op, err)
	}

	return newDevice, nil
}

func (s *DevicesStorage) GetDevices(ctx context.Context, params models.GetDevicesParams) ([]models.Device, error) {
	const op = "storage.postgres.GetDevices"

	devices := []models.Device{}
	var conditions []string
	var args []any

	query := `
			SELECT 
				id, hostname, ip, location, is_active, created_at 
			FROM
				devices
			WHERE
				is_deleted = false
			`
	if params.Search != "" {
		conditions = append(conditions, "hostname ILIKE $?")
		pattern := "%" + params.Search + "%"
		args = append(args, pattern)
	}

	if params.IsActive != nil {
		conditions = append(conditions, "is_active = $?")
		args = append(args, *params.IsActive)
	}

	for i, condition := range conditions {
		query += " AND " + condition
		query = strings.Replace(query, "$?", fmt.Sprintf("$%d", len(args)-len(conditions)+i+1), 1)
	}

	query += " ORDER BY id"

	if err := s.DB.SelectContext(ctx, &devices, query, args...); err != nil {
		return devices, fmt.Errorf("%s: %w", op, err)
	}

	return devices, nil
}

func (s *DevicesStorage) GetDeviceById(ctx context.Context, id int) (models.Device, error) {
	const op = "storage.postgres.GetDeviceById"

	var device models.Device

	query := `
			SELECT
				id, hostname, ip, location, is_active, created_at 
			FROM
				devices
			WHERE
				id = $1 AND is_deleted = false
			`
	if err := s.DB.GetContext(ctx, &device, query, id); err != nil {
		switch err {
		case sql.ErrNoRows:
			return device, storage.ErrNotFound
		default:
			return device, fmt.Errorf("%s: %w", op, err)
		}
	}

	return device, nil
}

func (s *DevicesStorage) UpdateDevice(ctx context.Context, device models.Device) (models.Device, error) {
	const op = "storage.postgres.UpdateDevice"

	var updatedDevice models.Device

	query := `
			UPDATE 
				devices 
			SET 
				hostname = :hostname, ip = :ip, location = :location, is_active = :is_active 
			WHERE 
				is_deleted = false AND id = :id
			RETURNING
				id, hostname, ip, location, is_active, created_at 
			`
	stmt, err := s.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return updatedDevice, fmt.Errorf("%s: %w", op, err)
	}

	if err := stmt.GetContext(ctx, &updatedDevice, device); err != nil {
		switch err {
		case sql.ErrNoRows:
			return updatedDevice, storage.ErrNotFound
		default:
			return updatedDevice, fmt.Errorf("%s: %w", op, err)
		}
	}

	return updatedDevice, nil
}

func (s *DevicesStorage) DeleteDevice(ctx context.Context, id int) error {
	const op = "storage.postgres.DeleteDevice"

	query := `
			UPDATE
				devices
			SET
				is_deleted = TRUE
			WHERE
				id = $1
			`
	if _, err := s.DB.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
