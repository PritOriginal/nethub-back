package deviceshandler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http/httptest"
	"testing"

	deviceshandler "github.com/PritOriginal/nethub-back/internal/handler/devices"
	"github.com/PritOriginal/nethub-back/internal/models"
	"github.com/PritOriginal/nethub-back/internal/service"
	"github.com/PritOriginal/nethub-back/pkg/logger/slogdiscard"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DevicesSuite struct {
	suite.Suite
	log *slog.Logger
	r   *gin.Engine
	s   *deviceshandler.MockDevicesService
}

func (suite *DevicesSuite) SetupSuite() {
	suite.log = slogdiscard.NewDiscardLogger()
	gin.SetMode(gin.TestMode)
	suite.r = gin.New()
	suite.s = deviceshandler.NewMockDevicesService(suite.T())

	deviceshandler.Register(suite.r, suite.log, suite.s)
}

func TestDevices(t *testing.T) {
	suite.Run(t, new(DevicesSuite))
}

func (suite *DevicesSuite) TestAddDevice() {
	tests := []struct {
		name            string
		rawReq          string
		req             deviceshandler.AddDeviceRequest
		wantErrParseReq bool
		errAddDevice    error
		statusCode      int
	}{
		{
			name: "Ok201",
			req: deviceshandler.AddDeviceRequest{
				Hostname: "myhost",
				IP:       "0.0.0.0",
				Location: "location",
				IsActive: true,
			},
			statusCode: 201,
		},
		{
			name:            "Err400InvalidJSON",
			rawReq:          "{",
			wantErrParseReq: true,
			statusCode:      400,
		},
		{name: "Err400InvalidReq",
			req: deviceshandler.AddDeviceRequest{
				Hostname: "myhost",
				IP:       "0.0.0",
				Location: "location",
				IsActive: true,
			},
			wantErrParseReq: true,
			statusCode:      400,
		},
		{
			name: "Err500",
			req: deviceshandler.AddDeviceRequest{
				Hostname: "myhost",
				IP:       "0.0.0.0",
				Location: "location",
				IsActive: true,
			},
			errAddDevice: errors.New(""),
			statusCode:   500,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if !tt.wantErrParseReq {
				suite.s.On("AddDevice", mock.Anything, mock.Anything).Once().
					Return(models.Device{}, tt.errAddDevice)
			}

			w := httptest.NewRecorder()

			var buf *bytes.Buffer
			if tt.rawReq == "" {
				body, err := json.Marshal(tt.req)
				suite.NoError(err)
				buf = bytes.NewBuffer(body)
			} else {
				buf = bytes.NewBuffer([]byte(tt.rawReq))
			}

			req := httptest.NewRequest("POST", "/devices", buf)

			suite.r.ServeHTTP(w, req)

			suite.Equal(tt.statusCode, w.Code)
		})
	}
}

func (suite *DevicesSuite) TestGetDevices() {
	tests := []struct {
		name                 string
		query                string
		wantErrParseIsActive bool
		errGetDevices        error
		statusCode           int
	}{
		{
			name:       "Ok200-1",
			query:      "?is_active=true&search=a",
			statusCode: 200,
		},
		{
			name:       "Ok200-2",
			query:      "?is_active=1&search=a",
			statusCode: 200,
		},
		{
			name:       "Ok200-3",
			query:      "?is_active=t&search=",
			statusCode: 200,
		},
		{
			name:       "Ok200-4",
			query:      "?is_active=true",
			statusCode: 200,
		},
		{
			name:       "Ok200-5",
			query:      "",
			statusCode: 200,
		},
		{
			name:       "Ok200-6",
			query:      "?is_active=false&search=",
			statusCode: 200,
		},
		{
			name:                 "Err400",
			query:                "?is_active=abc&search=abc",
			wantErrParseIsActive: true,
			statusCode:           400,
		},
		{
			name:          "Err500",
			query:         "?is_active=true&search=",
			errGetDevices: errors.New(""),
			statusCode:    500,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if !tt.wantErrParseIsActive {
				suite.s.On("GetDevices", mock.Anything, mock.Anything).Once().
					Return([]models.Device{}, tt.errGetDevices)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/devices"+tt.query, nil)
			suite.r.ServeHTTP(w, req)

			suite.Equal(tt.statusCode, w.Code)
		})
	}
}

func (suite *DevicesSuite) TestGetDeviceById() {
	tests := []struct {
		name             string
		id               string
		wantErrParseId   bool
		errGetDeviceById error
		statusCode       int
	}{
		{
			name:       "Ok200",
			id:         "1",
			statusCode: 200,
		},
		{
			name:           "Err400",
			id:             "a",
			wantErrParseId: true,
			statusCode:     400,
		},
		{
			name:             "Err404",
			id:               "1",
			errGetDeviceById: service.ErrNotFound,
			statusCode:       404,
		},
		{
			name:             "Err500",
			id:               "1",
			errGetDeviceById: errors.New(""),
			statusCode:       500,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if !tt.wantErrParseId {
				suite.s.On("GetDeviceById", mock.Anything, mock.AnythingOfType("int")).Once().
					Return(models.Device{}, tt.errGetDeviceById)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/devices/"+tt.id, nil)
			suite.r.ServeHTTP(w, req)

			suite.Equal(tt.statusCode, w.Code)
		})
	}
}

func (suite *DevicesSuite) TestUpdateDevice() {
	tests := []struct {
		name            string
		id              string
		wantErrParseId  bool
		rawReq          string
		req             deviceshandler.UpdateDeviceRequest
		wantErrParseReq bool
		errUpdateDevice error
		statusCode      int
	}{
		{
			name: "Ok200",
			id:   "1",
			req: deviceshandler.UpdateDeviceRequest{
				Hostname: "myhost",
				IP:       "0.0.0.0",
				Location: "location",
				IsActive: true,
			},
			statusCode: 200,
		},
		{
			name:           "Err400InvalidId",
			id:             "a",
			wantErrParseId: true,
			statusCode:     400,
		},
		{
			name:            "Err400InvalidJSON",
			id:              "1",
			rawReq:          "{",
			wantErrParseReq: true,
			statusCode:      400,
		},
		{
			name: "Err400InvalidReq",
			id:   "1",
			req: deviceshandler.UpdateDeviceRequest{
				Hostname: "myhost",
				IP:       "0.0.0",
				Location: "location",
				IsActive: true,
			},
			wantErrParseReq: true,
			statusCode:      400,
		},
		{
			name: "Err404",
			id:   "1",
			req: deviceshandler.UpdateDeviceRequest{
				Hostname: "myhost",
				IP:       "0.0.0.0",
				Location: "location",
				IsActive: true,
			},
			errUpdateDevice: service.ErrNotFound,
			statusCode:      404,
		},
		{
			name: "Err500",
			id:   "1",
			req: deviceshandler.UpdateDeviceRequest{
				Hostname: "myhost",
				IP:       "0.0.0.0",
				Location: "location",
				IsActive: true,
			},
			errUpdateDevice: errors.New(""),
			statusCode:      500,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if !tt.wantErrParseId && !tt.wantErrParseReq {
				suite.s.On("UpdateDevice", mock.Anything, mock.Anything).Once().
					Return(models.Device{}, tt.errUpdateDevice)
			}

			w := httptest.NewRecorder()

			var buf *bytes.Buffer
			if tt.rawReq == "" {
				body, err := json.Marshal(tt.req)
				suite.NoError(err)
				buf = bytes.NewBuffer(body)
			} else {
				buf = bytes.NewBuffer([]byte(tt.rawReq))
			}

			req := httptest.NewRequest("PUT", "/devices/"+tt.id, buf)
			suite.r.ServeHTTP(w, req)

			suite.Equal(tt.statusCode, w.Code)
		})
	}
}

func (suite *DevicesSuite) TestDeleteDevice() {
	tests := []struct {
		name            string
		id              string
		wantErrParseId  bool
		errDeleteDevice error
		statusCode      int
	}{
		{
			name:       "Ok200",
			id:         "1",
			statusCode: 200,
		},
		{
			name:           "Err400",
			id:             "a",
			wantErrParseId: true,
			statusCode:     400,
		},
		{
			name:            "Err404",
			id:              "1",
			errDeleteDevice: service.ErrNotFound,
			statusCode:      404,
		},
		{
			name:            "Err500",
			id:              "1",
			errDeleteDevice: errors.New(""),
			statusCode:      500,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if !tt.wantErrParseId {
				suite.s.On("DeleteDevice", mock.Anything, mock.AnythingOfType("int")).Once().
					Return(tt.errDeleteDevice)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/devices/"+tt.id, nil)
			suite.r.ServeHTTP(w, req)

			suite.Equal(tt.statusCode, w.Code)
		})
	}
}
