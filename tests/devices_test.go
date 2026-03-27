//go:build functional

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	deviceshandler "github.com/PritOriginal/nethub-back/internal/handler/devices"
	"github.com/PritOriginal/nethub-back/pkg/responses"
	"github.com/brianvoe/gofakeit/v7"
)

var hostnames = []string{
	"localhost",
	"myhost",
	"nethub",
	"home",
}

var locations = []string{
	"Офис",
	"Дом",
	"Главный офис",
	"Здание",
	"Серверная",
	"Комната",
	"Склад",
}

func (st *Suite) TestAddDevice() {
	tests := []struct {
		name       string
		rawReq     string
		req        deviceshandler.AddDeviceRequest
		statusCode int
	}{
		{
			name: "Ok201",
			req: deviceshandler.AddDeviceRequest{
				Hostname: gofakeit.RandomString(hostnames),
				IP:       gofakeit.IPv4Address(),
				Location: gofakeit.RandomString(locations),
				IsActive: gofakeit.Bool(),
			},
			statusCode: http.StatusCreated,
		},
		{
			name:       "Err400InvalidJSON",
			rawReq:     "{",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Err400InvalidReq-Hostname",
			req: deviceshandler.AddDeviceRequest{
				Hostname: "",
				IP:       gofakeit.IPv4Address(),
				Location: gofakeit.RandomString(locations),
				IsActive: gofakeit.Bool(),
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Err400InvalidReq-Hostname",
			req: deviceshandler.AddDeviceRequest{
				Hostname: "1a",
				IP:       gofakeit.IPv4Address(),
				Location: gofakeit.RandomString(locations),
				IsActive: gofakeit.Bool(),
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Err400InvalidReq-IP",
			req: deviceshandler.AddDeviceRequest{
				Hostname: gofakeit.RandomString(hostnames),
				IP:       "1.2.3.",
				Location: gofakeit.RandomString(locations),
				IsActive: gofakeit.Bool(),
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		st.Run(tt.name, func() {
			var body *bytes.Buffer
			if tt.rawReq == "" {
				reqJSON, err := json.Marshal(tt.req)
				st.NoError(err)
				body = bytes.NewBuffer(reqJSON)
			} else {
				body = bytes.NewBuffer([]byte(tt.rawReq))
			}

			resp, err := http.Post(
				fmt.Sprintf("http://%s:%d/api/devices", st.cfg.Server.Host, st.cfg.Server.Port),
				"application/json",
				body,
			)
			st.NoError(err)
			defer resp.Body.Close()

			var response responses.Response[deviceshandler.AddDeviceResponse]
			err = json.NewDecoder(resp.Body).Decode(&response)
			st.NoError(err)

			if tt.statusCode < 300 {
				st.NotNil(response.Data.Device)
				st.Nil(response.Error)
			} else {
				st.NotNil(response.Error)
			}
			st.Equal(tt.statusCode, resp.StatusCode)
		})
	}
}

func (st *Suite) TestGetDevices() {
	tests := []struct {
		name       string
		query      string
		statusCode int
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
			name:       "Err400",
			query:      "?is_active=abc&search=abc",
			statusCode: 400,
		},
	}
	for _, tt := range tests {
		st.Run(tt.name, func() {
			response := getDevices(st, tt.query, tt.statusCode)

			if tt.statusCode < 300 {
				st.NotNil(response.Data.Devices)
				st.Nil(response.Error)
			} else {
				st.Nil(response.Data.Devices)
				st.NotNil(response.Error)
			}
		})
	}
}

func getDevices(st *Suite, query string, expectedStatusCode int) responses.Response[deviceshandler.GetDevicesResponse] {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/api/devices%s", st.cfg.Server.Host, st.cfg.Server.Port, query))
	st.NoError(err)
	defer resp.Body.Close()

	var response responses.Response[deviceshandler.GetDevicesResponse]
	err = json.NewDecoder(resp.Body).Decode(&response)
	st.NoError(err)

	st.Equal(expectedStatusCode, resp.StatusCode)

	return response
}

func (st *Suite) TestGetDeviceById() {
	getDevicesResponse := getDevices(st, "", http.StatusOK)

	tests := []struct {
		name       string
		id         string
		statusCode int
	}{
		{
			name:       "Ok200",
			id:         strconv.Itoa(getDevicesResponse.Data.Devices[0].Id),
			statusCode: http.StatusOK,
		},
		{
			name:       "Err400",
			id:         "a",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Err404",
			id:         strconv.Itoa(math.MaxInt32),
			statusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		st.Run(tt.name, func() {
			resp, err := http.Get(fmt.Sprintf("http://%s:%d/api/devices/%s", st.cfg.Server.Host, st.cfg.Server.Port, tt.id))
			st.NoError(err)
			defer resp.Body.Close()

			var response responses.Response[deviceshandler.GetDeviceByIdResponse]
			err = json.NewDecoder(resp.Body).Decode(&response)
			st.NoError(err)

			st.NotNil(response.Data.Device)
			st.Equal(tt.statusCode, resp.StatusCode)
		})
	}
}

func (st *Suite) TestUpdateDevices() {
	getDevicesResponse := getDevices(st, "", http.StatusOK)
	device := getDevicesResponse.Data.Devices[0]

	tests := []struct {
		name       string
		id         string
		rawReq     string
		req        deviceshandler.UpdateDeviceRequest
		statusCode int
	}{
		{
			name: "Ok200-update-all",
			id:   strconv.Itoa(device.Id),
			req: deviceshandler.UpdateDeviceRequest{
				Hostname: gofakeit.RandomString(hostnames),
				IP:       gofakeit.IPv4Address(),
				Location: gofakeit.RandomString(locations),
				IsActive: gofakeit.Bool(),
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Ok200-update-hostname",
			id:   strconv.Itoa(device.Id),
			req: deviceshandler.UpdateDeviceRequest{
				Hostname: gofakeit.RandomString(hostnames),
				IP:       device.IP,
				Location: device.Location,
				IsActive: device.IsActive,
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Ok200-update-ip",
			id:   strconv.Itoa(device.Id),
			req: deviceshandler.UpdateDeviceRequest{
				Hostname: device.Hostname,
				IP:       gofakeit.IPv4Address(),
				Location: device.Location,
				IsActive: device.IsActive,
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Ok200-update-location",
			id:   strconv.Itoa(device.Id),
			req: deviceshandler.UpdateDeviceRequest{
				Hostname: device.Hostname,
				IP:       device.IP,
				Location: gofakeit.RandomString(locations),
				IsActive: device.IsActive,
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Ok200-update-isActive",
			id:   strconv.Itoa(device.Id),
			req: deviceshandler.UpdateDeviceRequest{
				Hostname: device.Hostname,
				IP:       device.IP,
				Location: device.Location,
				IsActive: gofakeit.Bool(),
			},
			statusCode: http.StatusOK,
		},
		{
			name:       "Err400InvalidId",
			id:         "a",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Err400InvalidJSON",
			id:         strconv.Itoa(device.Id),
			rawReq:     "{",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Err400InvalidReq-ip",
			id:   strconv.Itoa(device.Id),
			req: deviceshandler.UpdateDeviceRequest{
				Hostname: device.Hostname,
				IP:       "0.0.0",
				Location: device.Location,
				IsActive: device.IsActive,
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Err404",
			id:   strconv.Itoa(math.MaxInt32),
			req: deviceshandler.UpdateDeviceRequest{
				Hostname: device.Hostname,
				IP:       device.IP,
				Location: device.Location,
				IsActive: device.IsActive,
			},
			statusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		st.Run(tt.name, func() {
			reqJSON, err := json.Marshal(tt.req)
			st.NoError(err)

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s:%d/api/devices/%s", st.cfg.Server.Host, st.cfg.Server.Port, tt.id), bytes.NewBuffer(reqJSON))
			st.NoError(err)

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			st.NoError(err)
			defer resp.Body.Close()

			var response responses.Response[deviceshandler.UpdateDeviceRequest]
			err = json.NewDecoder(resp.Body).Decode(&response)
			st.NoError(err)

			st.Equal(tt.statusCode, resp.StatusCode)
			if tt.statusCode < 300 {
				st.NotNil(response.Data)
				st.Nil(response.Error)
			} else {
				st.NotNil(response.Error)
			}
		})
	}
}

func (st *Suite) TestDeleteDevice() {
	getDevicesResponse := getDevices(st, "", http.StatusOK)
	device := getDevicesResponse.Data.Devices[len(getDevicesResponse.Data.Devices)-1]

	tests := []struct {
		name       string
		id         string
		statusCode int
	}{
		{
			name:       "Ok200",
			id:         strconv.Itoa(device.Id),
			statusCode: http.StatusOK,
		},
		{
			name:       "Err400",
			id:         "a",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Err404",
			id:         strconv.Itoa(math.MaxInt32),
			statusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		st.Run(tt.name, func() {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s:%d/api/devices/%s", st.cfg.Server.Host, st.cfg.Server.Port, tt.id), nil)
			st.NoError(err)

			client := &http.Client{}
			resp, err := client.Do(req)
			st.NoError(err)
			defer resp.Body.Close()

			var response responses.Response[deviceshandler.DeleteDeviceResponse]
			err = json.NewDecoder(resp.Body).Decode(&response)
			st.NoError(err)

			st.Equal(tt.statusCode, resp.StatusCode)
		})
	}
}
