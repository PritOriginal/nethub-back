package devices_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/PritOriginal/nethub-back/internal/models"
	"github.com/PritOriginal/nethub-back/internal/service/devices"
	"github.com/PritOriginal/nethub-back/internal/storage"
	"github.com/PritOriginal/nethub-back/pkg/logger/slogdiscard"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DevicesSuite struct {
	suite.Suite
	service *devices.Service
	log     *slog.Logger
	storage *devices.MockStorage
}

func (suite *DevicesSuite) SetupSuite() {
	suite.log = slogdiscard.NewDiscardLogger()
	suite.storage = devices.NewMockStorage(suite.T())
	suite.service = devices.New(suite.log, suite.storage)
}

type method[T any] struct {
	data T
	err  error
}

func TestDevices(t *testing.T) {
	suite.Run(t, new(DevicesSuite))
}

func (suite *DevicesSuite) TestAddDevice() {
	tests := []struct {
		name      string
		addDevice method[models.Device]
	}{
		{
			name: "Ok",
			addDevice: method[models.Device]{
				data: models.Device{},
			},
		},
		{
			name: "Err",
			addDevice: method[models.Device]{
				err: errors.New(""),
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			func() {
				suite.storage.On("AddDevice", mock.Anything, mock.Anything).Once().
					Return(tt.addDevice.data, tt.addDevice.err)
				if tt.addDevice.err != nil {
					return
				}
			}()

			_, gotErr := suite.service.AddDevice(context.Background(), models.Device{})

			if tt.addDevice.err == nil {
				suite.NoError(gotErr)
			} else {
				suite.NotNil(gotErr)
			}
			suite.storage.AssertExpectations(suite.T())
		})
	}
}

func (suite *DevicesSuite) TestGetDevices() {
	tests := []struct {
		name       string
		getDevices method[[]models.Device]
	}{
		{
			name: "Ok",
			getDevices: method[[]models.Device]{
				data: []models.Device{},
			},
		},
		{
			name: "Err",
			getDevices: method[[]models.Device]{
				err: errors.New(""),
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			func() {
				suite.storage.On("GetDevices", mock.Anything, mock.Anything).Once().
					Return(tt.getDevices.data, tt.getDevices.err)
				if tt.getDevices.err != nil {
					return
				}
			}()

			_, gotErr := suite.service.GetDevices(context.Background(), models.GetDevicesParams{})

			if tt.getDevices.err == nil {
				suite.NoError(gotErr)
			} else {
				suite.NotNil(gotErr)
			}
			suite.storage.AssertExpectations(suite.T())
		})
	}
}

func (suite *DevicesSuite) TestGetDeviceById() {
	tests := []struct {
		name          string
		getDeviceById method[models.Device]
	}{
		{
			name: "Ok",
			getDeviceById: method[models.Device]{
				data: models.Device{},
			},
		},
		{
			name: "ErrNotFound",
			getDeviceById: method[models.Device]{
				err: storage.ErrNotFound,
			},
		},
		{
			name: "Err",
			getDeviceById: method[models.Device]{
				err: errors.New(""),
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			func() {
				suite.storage.On("GetDeviceById", mock.Anything, mock.AnythingOfType("int")).Once().
					Return(tt.getDeviceById.data, tt.getDeviceById.err)
				if tt.getDeviceById.err != nil {
					return
				}
			}()

			_, gotErr := suite.service.GetDeviceById(context.Background(), 1)

			if tt.getDeviceById.err == nil {
				suite.NoError(gotErr)
			} else {
				suite.NotNil(gotErr)
			}
			suite.storage.AssertExpectations(suite.T())
		})
	}
}

func (suite *DevicesSuite) TestUpdateDevice() {
	tests := []struct {
		name         string
		updateDevice method[models.Device]
	}{
		{
			name: "Ok",
			updateDevice: method[models.Device]{
				data: models.Device{},
			},
		},
		{
			name: "Err",
			updateDevice: method[models.Device]{
				err: errors.New(""),
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			func() {
				suite.storage.On("UpdateDevice", mock.Anything, mock.Anything).Once().
					Return(tt.updateDevice.data, tt.updateDevice.err)
				if tt.updateDevice.err != nil {
					return
				}
			}()

			_, gotErr := suite.service.UpdateDevice(context.Background(), models.Device{})

			if tt.updateDevice.err == nil {
				suite.NoError(gotErr)
			} else {
				suite.NotNil(gotErr)
			}
			suite.storage.AssertExpectations(suite.T())
		})
	}
}

func (suite *DevicesSuite) TestDeleteDevice() {
	tests := []struct {
		name          string
		getDeviceById method[models.Device]
		deleteDevice  method[any]
	}{
		{
			name: "Ok",
			getDeviceById: method[models.Device]{
				data: models.Device{},
			},
			deleteDevice: method[any]{},
		},
		{
			name: "ErrNotFound",
			getDeviceById: method[models.Device]{
				err: storage.ErrNotFound,
			},
		},
		{
			name: "ErrGetDeviceById",
			getDeviceById: method[models.Device]{
				err: errors.New(""),
			},
		},
		{
			name: "ErrDelete",
			getDeviceById: method[models.Device]{
				data: models.Device{},
			},
			deleteDevice: method[any]{
				err: errors.New(""),
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			func() {
				suite.storage.On("GetDeviceById", mock.Anything, mock.AnythingOfType("int")).Once().
					Return(tt.getDeviceById.data, tt.getDeviceById.err)
				if tt.getDeviceById.err != nil {
					return
				}

				suite.storage.On("DeleteDevice", mock.Anything, mock.AnythingOfType("int")).Once().
					Return(tt.deleteDevice.err)
				if tt.deleteDevice.err != nil {
					return
				}
			}()

			gotErr := suite.service.DeleteDevice(context.Background(), 1)

			if tt.getDeviceById.err == nil && tt.deleteDevice.err == nil {
				suite.NoError(gotErr)
			} else {
				suite.NotNil(gotErr)
			}
			suite.storage.AssertExpectations(suite.T())
		})
	}
}
