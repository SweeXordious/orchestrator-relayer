package orchestrator_test

import (
	qgbtesting "github.com/celestiaorg/orchestrator-relayer/testing"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type OrchestratorTestSuite struct {
	suite.Suite
	Node *qgbtesting.CelestiaNetwork
}

func (s *OrchestratorTestSuite) SetupSuite() {
	t := s.T()
	s.Node = qgbtesting.NewCelestiaNetwork(t, time.Millisecond)
}

func (s *OrchestratorTestSuite) TearDownSuite() {
	s.Node.Stop()
}

func TestOrchestrator(t *testing.T) {
	suite.Run(t, new(OrchestratorTestSuite))
}
