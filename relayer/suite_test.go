package relayer_test

import (
	qgbtesting "github.com/celestiaorg/orchestrator-relayer/testing"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type RelayerTestSuite struct {
	suite.Suite
	Node *qgbtesting.CelestiaNetwork
}

func (s *RelayerTestSuite) SetupSuite() {
	t := s.T()
	s.Node = qgbtesting.NewCelestiaNetwork(t, time.Millisecond)
}

func (s *RelayerTestSuite) TearDownSuite() {
	s.Node.Stop()
}

func TestRelayer(t *testing.T) {
	suite.Run(t, new(RelayerTestSuite))
}
