//go:build functional

package tests

import (
	"testing"

	"github.com/PritOriginal/nethub-back/internal/config"
	"github.com/PritOriginal/nethub-back/internal/service/devices"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	service *devices.Service
	cfg     *config.Config
}

func (st *Suite) SetupSuite() {
	st.cfg = config.MustLoadPath("../configs/config.yaml")
}

func Test(t *testing.T) {
	suite.Run(t, new(Suite))
}
